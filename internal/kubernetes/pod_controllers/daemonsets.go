package pod_controllers

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

func CompareDaemonSets(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntities config.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	daemonSets1, err := clientSet1.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain daemonsets list from 1st cluster: %w", err)
	}

	daemonSets2, err := clientSet2.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain daemonsets list from 2nd cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := PrepareDaemonSetMaps(daemonSets1, daemonSets2)

	isClustersDiffer = ComparePodControllerSpecs(map1, map2, apc1List, apc2List, namespace)

	return isClustersDiffer, nil
}

// PrepareDaemonSetMaps prepares DaemonSet maps for comparison
func PrepareDaemonSetMaps(obj1, obj2 *v1.DaemonSetList) ([]AbstractPodController, map[string]types.IsAlreadyComparedFlag, []AbstractPodController, map[string]IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		map1     = make(map[string]types.IsAlreadyComparedFlag)
		apc1List = make([]AbstractPodController, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if _, ok := ToSkipEntities[ObjectKindWrapper(value.Kind)][value.Name]; ok {
			logging.Log.Debugf("daemonset %s is skipped from comparison due to its name", value.Name)
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
		if _, ok := ToSkipEntities[ObjectKindWrapper(value.Kind)][value.Name]; ok {
			logging.Log.Debugf("daemonset %s is skipped from comparison due to its name", value.Name)
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
