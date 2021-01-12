package pod_controllers

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	objectBatchLimit = 25
)

func addItemsToDaemonSetsList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v1.DaemonSetList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToDaemonSetsList started")
	defer log.Debugf("addItemsToDaemonSetsList completed")
	var (
		batch      *v1.DaemonSetList
		daemonSets = &v1.DaemonSetList{
			Items: make([]v1.DaemonSet, 0),
		}

		continueToken string

		err error
	)

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

			log.Debugf("addItemsToDaemonSetsList: %d objects received", len(batch.Items))

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

// CompareDaemonSets compares list of daemonsets objects in two given k8s-clusters
func CompareDaemonSets(ctx context.Context, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", "daemonset"))

		clientSet1, clientSet2, namespace = config.FromContext(ctx)

		isClustersDiffer bool
	)
	ctx = logging.WithLogger(ctx, log)

	daemonSets1, err := addItemsToDaemonSetsList(ctx, clientSet1, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain daemonsets list from 1st cluster: %w", err)
	}

	daemonSets2, err := addItemsToDaemonSetsList(ctx, clientSet2, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain daemonsets list from 2st cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := prepareDaemonSetMaps(ctx, daemonSets1, daemonSets2, skipEntityList.GetByKind("daemonsets"))

	isClustersDiffer, err = ComparePodControllers(ctx, &clusterCompareTask{
		Client:                   clientSet1,
		APCList:                  apc1List,
		IsAlreadyCheckedFlagsMap: map1,
	}, &clusterCompareTask{
		Client:                   clientSet2,
		APCList:                  apc2List,
		IsAlreadyCheckedFlagsMap: map2,
	}, namespace)

	return isClustersDiffer, err
}

// prepareDaemonSetMaps prepares DaemonSet maps for comparison
func prepareDaemonSetMaps(ctx context.Context, obj1, obj2 *v1.DaemonSetList, skipEntities skipper.SkipComponentNames) ([]AbstractPodController, map[string]types.IsAlreadyComparedFlag, []AbstractPodController, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)

		map1     = make(map[string]types.IsAlreadyComparedFlag)
		apc1List = make([]AbstractPodController, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {

		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("daemonset/%s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, AbstractPodController{
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
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("daemonset/%s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, AbstractPodController{
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
