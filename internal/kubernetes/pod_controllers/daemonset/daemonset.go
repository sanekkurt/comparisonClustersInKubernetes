package daemonset

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
	daemonSetKind = "daemonset"
)

func daemonSetsRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Workloads.PodControllers.DaemonSets.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v1.DaemonSetList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch      *v1.DaemonSetList
		daemonSets = &v1.DaemonSetList{
			Items: make([]v1.DaemonSet, 0),
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
			batch, err = clientSet.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{
				Limit:    limit,
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

			daemonSets.Items = append(daemonSets.Items, batch.Items...)

			daemonSets.TypeMeta = batch.TypeMeta
			daemonSets.ListMeta = batch.ListMeta

			if batch.Continue == "" {
				break forLoop
			}

			continueToken = batch.Continue
		}
	}

	daemonSets.Continue = ""

	return daemonSets, err
}

type DaemonSetsComparator struct {
}

func NewDaemonSetsComparator(ctx context.Context, namespace string) DaemonSetsComparator {
	return DaemonSetsComparator{}
}

// Compare compares list of DaemonSets in two given k8s-clusters
func (cmp DaemonSetsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", daemonSetKind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Workloads.PodControllers.Enabled ||
		!cfg.Workloads.PodControllers.StatefulSets.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", daemonSetKind)
		return nil, nil
	}

	daemonSets1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, daemonSetsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain daemonsets list from 1st cluster: %w", err)
	}

	daemonSets2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, daemonSetsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain daemonsets list from 2st cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := prepareDaemonSetMaps(ctx, daemonSets1, daemonSets2)

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

// prepareDaemonSetMaps prepares DaemonSet maps for comparison
func prepareDaemonSetMaps(ctx context.Context, obj1, obj2 *v1.DaemonSetList) ([]common.AbstractPodController, map[string]types.IsAlreadyComparedFlag, []common.AbstractPodController, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
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
		if cfg.ExcludesIncludes.IsSkippedEntity(daemonSetKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("daemonset/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "daemonsets",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         nil,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})
	}

	for index, value := range obj2.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(daemonSetKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("daemonset/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "daemonsets",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         nil,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})
	}

	return apc1List, map1, apc2List, map2
}
