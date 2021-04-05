package service

import (
	"context"
	"fmt"

	"k8s-cluster-comparator/internal/consts"
	"k8s-cluster-comparator/internal/kubernetes/diff"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"

	"sync"

	kubectx "k8s-cluster-comparator/internal/kubernetes/context"

	"go.uber.org/zap"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	objectKind = "service"
)

type Comparator struct {
	Kind      string
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
	return ""
}

func (cmp *Comparator) LabelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]v12.Service, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]v12.Service)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.CoreV1().Services(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
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
		obj, err := clientSet.CoreV1().Services(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
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

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]v12.Service, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *v12.ServiceList
		objects = make(map[string]v12.Service)

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
			batch, err = clientSet.CoreV1().Services(cmp.Namespace).List(ctx, metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.FieldSelectorProvider(ctx),
				LabelSelector: cmp.LabelSelectorProvider(ctx),
				Continue:      continueToken,
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

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]v12.Service, error) {
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

// Compare compares list of services objects in two given k8s-clusters
func (cmp *Comparator) Compare(ctx context.Context) (*diff.DiffsStorage, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Networking.Enabled ||
		!cfg.Networking.Services.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	objects, err := cmp.collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve objects for comparision: %w", err)
	}

	_ = cmp.compare(ctx, objects[0], objects[1])

	return nil, nil

}

func (cmp *Comparator) collect(ctx context.Context) ([]map[string]v12.Service, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]v12.Service, 2, 2)
		wg      = &sync.WaitGroup{}

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

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]v12.Service) error {
	var (
		log = logging.FromContext(ctx)
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
	}

	for name, obj1 := range map1 {
		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return nil
		default:
			if obj2, ok := map2[name]; ok {
				compareServicesSpecs(ctx, name, &obj1, &obj2)

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

	return nil
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v12.ServiceList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch    *v12.ServiceList
		services = &v12.ServiceList{
			Items: make([]v12.Service, 0),
		}

		continueToken string

		err error
	)

	log.Debugf("fillInComparisonMap started")
	defer log.Debugf("fillInComparisonMap completed")

	for {
		batch, err = clientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

		services.Items = append(services.Items, batch.Items...)

		services.TypeMeta = batch.TypeMeta
		services.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	services.Continue = ""

	return services, err
}

//compareServicesSpecs set information about services
func compareServicesSpecs(ctx context.Context, name string, svc1, svc2 *v12.Service) {
	var (
		log = logging.FromContext(ctx)

		diffStorage = diff.StorageFromContext(ctx)
		diffsBatch  = diffStorage.NewLazyBatch(svc1.TypeMeta, svc1.ObjectMeta)
	)

	ctx = diff.WithDiffBatch(ctx, diffsBatch)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("service/%s compare started", name)
	defer func() {
		log.Debugf("service/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, svc1.ObjectMeta, svc2.ObjectMeta)
	compareSpecInServices(ctx, *svc1, *svc2)

	return
}

// compareSpecInServices compares spec in services
func compareSpecInServices(ctx context.Context, service1, service2 v12.Service) {
	var (
		log = logging.FromContext(ctx)

		diffsBatch = diff.BatchFromContext(ctx)
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return
	default:
		if len(service1.Spec.Ports) != len(service2.Spec.Ports) {
			//log.Warnf("%s. %d vs %d", ErrorPortsCountDifferent.Error(), len(service1.Spec.Ports), len(service2.Spec.Ports))
			diffsBatch.Add(ctx, true, zap.WarnLevel, "%s. %d vs %d", ErrorPortsCountDifferent.Error(), len(service1.Spec.Ports), len(service2.Spec.Ports))
			return
		}

		for index, value := range service1.Spec.Ports {
			if value != service2.Spec.Ports[index] {
				diffsBatch.Add(ctx, true, zap.WarnLevel, "%s. %s-%d-%s vs %s-%d-%s", ErrorPortInServicesDifferent.Error(), value.Name, value.Port, value.Protocol, service2.Spec.Ports[index].Name, service2.Spec.Ports[index].Port, service2.Spec.Ports[index].Protocol)
				return
			}
		}

		if len(service1.Spec.Selector) != len(service2.Spec.Selector) {
			diffsBatch.Add(ctx, true, zap.WarnLevel, "%s. %d vs %d", ErrorSelectorsCountDifferent.Error(), len(service1.Spec.Selector), len(service2.Spec.Selector))
			return
		}

		for key, value := range service1.Spec.Selector {
			if service2.Spec.Selector[key] != value {
				diffsBatch.Add(ctx, true, zap.WarnLevel, "%s. %s-%s vs %s-%s", ErrorSelectorInServicesDifferent.Error(), key, value, key, service2.Spec.Selector[key])
				return

			}
		}

		if service1.Spec.Type != service2.Spec.Type {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorTypeInServicesDifferent.Error(), service1.Spec.Type, service2.Spec.Type)
		}
		return
	}
}
