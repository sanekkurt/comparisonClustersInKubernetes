package pod_controllers

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/logging"
)

func CompareDeployments(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntities config.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	depl1, err := clientSet1.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain deployments list from 1st cluster: %w", err)
	}

	depl2, err := clientSet2.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain deployments list from 2nd cluster: %w", err)
	}

	apc1List, map1, apc2List, map2 := PrepareDeploymentMaps(depl1, depl2)

	isClustersDiffer = ComparePodControllerSpecs(map1, map2, apc1List, apc2List, namespace)

	return isClustersDiffer, nil
}

// PrepareDeploymentMaps prepare deployment maps for comparison
func PrepareDeploymentMaps(obj1, obj2 *v1.DeploymentList) ([]AbstractPodController, map[string]IsAlreadyComparedFlag, []AbstractPodController, map[string]IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		map1     = make(map[string]IsAlreadyComparedFlag)
		apc1List = make([]AbstractPodController, 0)

		map2     = make(map[string]IsAlreadyComparedFlag)
		apc2List = make([]AbstractPodController, 0)

		indexCheck IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if _, ok := ToSkipEntities[ObjectKindWrapper(value.Kind)][value.Name]; ok {
			logging.Log.Debugf("deployment %s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, AbstractPodController{
			Metadata: AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "deployments",
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
		if _, ok := ToSkipEntities[ObjectKindWrapper(value.Kind)][value.Name]; ok {
			logging.Log.Debugf("deployment %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, AbstractPodController{
			Metadata: AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "deployments",
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
