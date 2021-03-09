package deployment

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/consts"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	pccommon "k8s-cluster-comparator/internal/kubernetes/pod_controllers/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	objectKind = "deployment"
)

type DeploymentsComparator struct {
	Kind      string
	Namespace string
	BatchSize int64
}

func NewDeploymentsComparator(ctx context.Context, namespace string) *DeploymentsComparator {
	return &DeploymentsComparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),
	}
}

func (cmp *DeploymentsComparator) fieldSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *DeploymentsComparator) labelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *DeploymentsComparator) collectIncludedFromCluster(ctx context.Context) (map[string]appsv1.Deployment, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]appsv1.Deployment)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.AppsV1().Deployments(cmp.Namespace).Get(string(name), metav1.GetOptions{})
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
		obj, err := clientSet.AppsV1().Deployments(cmp.Namespace).Get(string(name), metav1.GetOptions{})
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

func (cmp *DeploymentsComparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]appsv1.Deployment, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *appsv1.DeploymentList
		objects = make(map[string]appsv1.Deployment)

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
			batch, err = clientSet.AppsV1().Deployments(cmp.Namespace).List(metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.fieldSelectorProvider(ctx),
				LabelSelector: cmp.labelSelectorProvider(ctx),
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

				if *obj.Spec.Replicas != obj.Status.ReadyReplicas {
					log.With(zap.String("objectName", obj.Name)).Warnf("%s/%s is progressing now, comparison might be inaccurate", cmp.Kind, obj.Name)
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

func (cmp *DeploymentsComparator) collectFromCluster(ctx context.Context) (map[string]appsv1.Deployment, error) {
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

// Compare compares list of Deployment objects in two given k8s-clusters
func (cmp *DeploymentsComparator) Compare(ctx context.Context) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)

		err error
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Workloads.PodControllers.Enabled ||
		!cfg.Workloads.PodControllers.Deployments.Enabled {
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

func (cmp *DeploymentsComparator) collect(ctx context.Context) ([]map[string]appsv1.Deployment, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]appsv1.Deployment, 2, 2)
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

func (cmp *DeploymentsComparator) compare(ctx context.Context, map1, map2 map[string]appsv1.Deployment) []types.KubeObjectsDifference {
	var (
		apcs = make([]map[string]*pccommon.AbstractPodController, 2, 2)
	)

	for idx, objs := range []map[string]appsv1.Deployment{map1, map2} {
		apcs[idx] = cmp.prepareAPCMap(ctx, objs)
	}

	diffs := pccommon.CompareAbstractPodControllerMaps(ctx, cmp.Kind, apcs[0], apcs[1])

	return diffs
}

func (cmp *DeploymentsComparator) prepareAPCMap(ctx context.Context, objs map[string]appsv1.Deployment) map[string]*pccommon.AbstractPodController {
	var (
		apcs = make(map[string]*pccommon.AbstractPodController)
	)

	for name, obj := range objs {
		apcs[name] = &pccommon.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       cmp.Kind,
					APIVersion: "apps/v1",
				},
				Meta: obj.ObjectMeta,
			},
			Name:             obj.Name,
			Labels:           obj.Labels,
			Annotations:      obj.Annotations,
			Replicas:         obj.Spec.Replicas,
			PodLabelSelector: obj.Spec.Selector,
			PodTemplateSpec:  obj.Spec.Template,
		}
	}

	return apcs
}
