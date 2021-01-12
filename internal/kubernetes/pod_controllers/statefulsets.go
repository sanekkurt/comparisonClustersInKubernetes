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

func addItemsToStatefulSetList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v1.StatefulSetList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToStatefulSetList started")
	defer log.Debugf("addItemsToStatefulSetList completed")

	var (
		batch        *v1.StatefulSetList
		statefulSets = &v1.StatefulSetList{
			Items: make([]v1.StatefulSet, 0),
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
			batch, err = clientSet.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{
				Limit:    limit,
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("addItemsToStatefulSetList: %d objects received", len(batch.Items))

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

func CompareStateFulSets(ctx context.Context, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", "statefulset"))

		clientSet1, clientSet2, namespace = config.FromContext(ctx)

		isClustersDiffer bool
	)

	ctx = logging.WithLogger(ctx, log)

	statefulSet1, err := addItemsToStatefulSetList(ctx, clientSet1, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain statefulsets list from 1st cluster: %w", err)
	}

	statefulSet2, err := addItemsToStatefulSetList(ctx, clientSet2, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain statefulsets list from 2st cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := prepareStatefulSetMaps(ctx, statefulSet1, statefulSet2, skipEntityList.GetByKind("statefulsets"))

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

// prepareStatefulSetMaps prepares StatefulSet maps for comparison
func prepareStatefulSetMaps(ctx context.Context, obj1, obj2 *v1.StatefulSetList, skipEntities skipper.SkipComponentNames) ([]AbstractPodController, map[string]types.IsAlreadyComparedFlag, []AbstractPodController, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
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
			log.Debugf("statefulset/%s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, AbstractPodController{
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
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("statefulset/%s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, AbstractPodController{
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
