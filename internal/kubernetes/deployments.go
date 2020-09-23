package kubernetes

import (
	v1 "k8s.io/api/apps/v1"
)

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
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, AbstractPodController{
			Kind:             ObjectKindWrapper(value.Kind),
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
			continue
		}
		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, AbstractPodController{
			Kind:             ObjectKindWrapper(value.Kind),
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
