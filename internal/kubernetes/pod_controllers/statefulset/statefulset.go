package statefulset

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	statefulSetKind = "statefulset"
)

func statefulSetsRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Workloads.PodControllers.Deployments.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v1.StatefulSetList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch        *v1.StatefulSetList
		statefulSets = &v1.StatefulSetList{
			Items: make([]v1.StatefulSet, 0),
		}

		continueToken string

		err error
	)

	log.Debugf("fillInComparisonMap started")
	defer log.Debugf("fillInComparisonMap completed")

forLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{
				Limit:    limit,
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

			statefulSets.Items = append(statefulSets.Items, batch.Items...)

			statefulSets.TypeMeta = batch.TypeMeta
			statefulSets.ListMeta = batch.ListMeta

			if batch.Continue == "" {
				break forLoop
			}

			continueToken = batch.Continue
		}
	}

	statefulSets.Continue = ""

	return statefulSets, err
}

type StatefulSetsComparator struct {
}

func NewStatefulSetsComparator(ctx context.Context, namespace string) StatefulSetsComparator {
	return StatefulSetsComparator{}
}

func (cmp StatefulSetsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", statefulSetKind))
		cfg = config.FromContext(ctx)
	)

	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Workloads.PodControllers.Enabled ||
		!cfg.Workloads.PodControllers.StatefulSets.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", statefulSetKind)
		return nil, nil
	}

	statefulSet1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, statefulSetsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain statefulsets list from 1st cluster: %w", err)
	}

	statefulSet2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, statefulSetsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain statefulsets list from 2st cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := prepareStatefulSetMaps(ctx, statefulSet1, statefulSet2)

	_, err = common.ComparePodControllers(ctx, &common.ClusterCompareTask{
		Client:                   cfg.Connections.Cluster1.ClientSet,
		APCList:                  apc1List,
		IsAlreadyCheckedFlagsMap: map1,
	}, &common.ClusterCompareTask{
		Client:                   cfg.Connections.Cluster2.ClientSet,
		APCList:                  apc2List,
		IsAlreadyCheckedFlagsMap: map2,
	}, namespace)

	return nil, err
}

// prepareStatefulSetMaps prepares StatefulSet maps for comparison
func prepareStatefulSetMaps(ctx context.Context, obj1, obj2 *v1.StatefulSetList) ([]common.AbstractPodController, map[string]types.IsAlreadyComparedFlag, []common.AbstractPodController, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		map1     = make(map[string]types.IsAlreadyComparedFlag)
		apc1List = make([]common.AbstractPodController, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]common.AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(statefulSetKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("statefulset/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "statefulsets",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         value.Spec.Replicas,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})
	}

	for index, value := range obj2.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(statefulSetKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("statefulset/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "statefulsets",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         value.Spec.Replicas,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})
	}

	return apc1List, map1, apc2List, map2
}
