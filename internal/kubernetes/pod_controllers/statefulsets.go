package pod_controllers

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

func CompareStateFulSets(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)
	statefulSet1, err := clientSet1.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain statefulsets list from 1st cluster: %w", err)
	}

	statefulSet2, err := clientSet2.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain statefulsets list from 2nd cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := prepareStatefulSetMaps(statefulSet1, statefulSet2, skipEntityList.GetByKind("statefulsets"))

	isClustersDiffer = comparePodControllerSpecs(&clusterCompareTask{
		Client:                   clientSet1,
		APCList:                  apc1List,
		IsAlreadyCheckedFlagsMap: map1,
	}, &clusterCompareTask{
		Client:                   clientSet2,
		APCList:                  apc2List,
		IsAlreadyCheckedFlagsMap: map2,
	}, namespace)

	return isClustersDiffer, nil
}

// prepareStatefulSetMaps prepares StatefulSet maps for comparison
func prepareStatefulSetMaps(obj1, obj2 *v1.StatefulSetList, skipEntities skipper.SkipComponentNames) ([]AbstractPodController, map[string]types.IsAlreadyComparedFlag, []AbstractPodController, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		map1     = make(map[string]types.IsAlreadyComparedFlag)
		apc1List = make([]AbstractPodController, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("statefulset %s is skipped from comparison due to its name", value.Name)
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
			log.Debugf("statefulset %s is skipped from comparison due to its name", value.Name)
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
