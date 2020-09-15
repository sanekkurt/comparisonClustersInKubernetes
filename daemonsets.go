package main

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
)

func AddValueDaemonSetsMap(daemonSets1 *v1.DaemonSetList, daemonSets2 *v1.DaemonSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapDaemonSets1 := make(map[string]CheckerFlag)
	mapDaemonSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range daemonSets1.Items {
		indexCheck.index = index
		mapDaemonSets1[value.Name] = indexCheck
	}
	for index, value := range daemonSets2.Items {
		indexCheck.index = index
		mapDaemonSets2[value.Name] = indexCheck
	}
	return mapDaemonSets1, mapDaemonSets2
}

func SetInformationAboutDaemonSets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, daemonSets1 *v1.DaemonSetList, daemonSets2 *v1.DaemonSetList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The daemonsets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking daemonset: '%s' -----\n", name)

			//заполняем информацию, которая будет использоваться при сравнении
			object1 := InformationAboutObject{
				Template: daemonSets1.Items[index1.index].Spec.Template,
				Selector: daemonSets1.Items[index1.index].Spec.Selector,
			}
			object2 := InformationAboutObject{
				Template: daemonSets2.Items[index2.index].Spec.Template,
				Selector: daemonSets2.Items[index2.index].Spec.Selector,
			}
			//CompareContainers(deployment1.Items[index1.index].Spec, deployment2.Items[index2.index].Spec, namespace)
			CompareContainers(object1, object2, namespace, client1, client2)

			fmt.Printf("----- End checking daemonset: '%s' -----\n\n", name)
		} else {
			fmt.Printf("DaemonSet '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("DaemonSet '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}
