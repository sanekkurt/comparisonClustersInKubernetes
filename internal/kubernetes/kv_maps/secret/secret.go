package secret

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/consts"
	"k8s-cluster-comparator/internal/kubernetes/common"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	objectKind = "secret"
)


type Comparator struct {
	Kind string
	Namespace string
	BatchSize int64
}

func NewComparator(ctx context.Context, namespace string) *Comparator {
	return &Comparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),

	}
}

func (cmp *Comparator) FieldSelectorProvider(ctx context.Context) string {
	var (
		cfg = config.FromContext(ctx)
		fieldSelector string
	)

	for t := range cfg.Configs.Secrets.SkipTypesMap {
		fieldSelector += fmt.Sprintf("type!=%s,", t)
	}

	return fieldSelector
}

func (cmp *Comparator) LabelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]corev1.Secret, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]corev1.Secret)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.CoreV1().Secrets(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.CoreV1().Secrets(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	return objects, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]corev1.Secret, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *corev1.SecretList
		objects = make(map[string]corev1.Secret)

		continueToken string

		err error
	)

	log.Debugf("%T: collectFromClusterWithoutExcludes started", cmp)
	defer log.Debugf("%T: collectFromClusterWithoutExcludes completed", cmp)


forOuterLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.CoreV1().Secrets(cmp.Namespace).List(metav1.ListOptions{
				Limit:    cmp.BatchSize,
				FieldSelector: cmp.FieldSelectorProvider(ctx),
				LabelSelector: cmp.LabelSelectorProvider(ctx),
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %ss retrieved", len(batch.Items), cmp.Kind)

		forInnerLoop:
			for _, obj := range batch.Items {
				if _, ok := objects[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoop
				}

				objects[obj.Name] = obj
			}

			if batch.Continue == "" {
				break forOuterLoop
			}

			continueToken = batch.Continue
		}
	}

	return objects, nil
}

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]corev1.Secret, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludes(ctx)
	} else {
		return cmp.collectIncludedFromCluster(ctx)
	}
}

// Compare compares list of Secret objects in two given k8s-clusters
func (cmp *Comparator) Compare(ctx context.Context) (*diff.DiffsStorage, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)

		err error
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Configs.Enabled ||
		!cfg.Configs.Secrets.Enabled {
		log.Debugf("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	objects, err := cmp.collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve objects for comparision: %w", err)
	}

	diff := cmp.compare(ctx, objects[0], objects[1])

	return diff, nil
}

func  (cmp *Comparator) collect(ctx context.Context) ([]map[string]corev1.Secret, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]corev1.Secret, 2, 2)
		wg = &sync.WaitGroup{}

		err error
	)

	wg.Add(2)

	for idx, clientSet := range []kubernetes.Interface{
		cfg.Connections.Cluster1.ClientSet,
		cfg.Connections.Cluster2.ClientSet,
	} {
		go func(idx int, clientSet kubernetes.Interface) {
			defer wg.Done()

			objects[idx], err = cmp.collectFromCluster(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Fatalf("cannot obtain %ss from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]corev1.Secret) []types.ObjectsDiff {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.ObjectsDiff, 0)
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
	}

	for name, obj1 := range map1 {
		ctx := logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return nil
		default:
			if obj2, ok := map2[name]; ok {
				diff := compareSecretSpecs(ctx, name, &obj1, &obj2)

				diffs = append(diffs, diff...)

				delete(map1, name)
				delete(map2, name)
			} else {
				log.With(zap.String("objectName", name)).Warnf("%s does not exist in 2nd cluster", cmp.Kind)
			}
		}
	}

	for name, _ := range map2 {
		log.With(zap.String("objectName", name)).Warnf("%s does not exist in 1st cluster", cmp.Kind)
	}

	return diffs
}

func compareSecretSpecs(ctx context.Context, name string, cm1, cm2 *corev1.Secret) []types.ObjectsDiff {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))
		cfg = config.FromContext(ctx)
	)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("Secret/%s compare started", name)
	defer func() {
		log.Debugf("Secret/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, cm1.ObjectMeta, cm2.ObjectMeta)
	common.AreKVBytesMapsEqual(ctx, cm1.Data, cm2.Data, nil, cfg.Configs.Secrets.DumpDifferentValues)

	return nil
}