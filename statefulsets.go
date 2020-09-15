package main

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
)

func AddValueStatefulSetsInMap(stateFulSets1 *v1.StatefulSetList, stateFulSets2 *v1.StatefulSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapStatefulSets1 := make(map[string]CheckerFlag)
	mapStatefulSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range stateFulSets1.Items {
		indexCheck.index = index
		mapStatefulSets1[value.Name] = indexCheck
	}
	for index, value := range stateFulSets2.Items {
		indexCheck.index = index
		mapStatefulSets2[value.Name] = indexCheck
	}
	return mapStatefulSets1, mapStatefulSets2
}

func SetInformationAboutStatefulSets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, statefulSets1 *v1.StatefulSetList, statefulSets2 *v1.StatefulSetList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The statefulsets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking statefulset: '%s' -----\n", name)
			if *statefulSets1.Items[index1.index].Spec.Replicas != *statefulSets2.Items[index2.index].Spec.Replicas {
				fmt.Printf("!!!The replicas count are different!!!\n%s '%s' replicas: %d\n%s '%s' replicas: %d\n", kubeconfig1YamlStruct.Clusters[0].Cluster.Server, statefulSets1.Items[index1.index].Name, *statefulSets1.Items[index1.index].Spec.Replicas, kubeconfig2YamlStruct.Clusters[0].Cluster.Server, statefulSets2.Items[index2.index].Name, *statefulSets2.Items[index2.index].Spec.Replicas)
			} else {
				//заполняем информацию, которая будет использоваться при сравнении
				object1 := InformationAboutObject{
					Template: statefulSets1.Items[index1.index].Spec.Template,
					Selector: statefulSets1.Items[index1.index].Spec.Selector,
				}
				object2 := InformationAboutObject{
					Template: statefulSets2.Items[index2.index].Spec.Template,
					Selector: statefulSets2.Items[index2.index].Spec.Selector,
				}

				//CompareContainers(deployment1.Items[index1.index].Spec, deployment2.Items[index2.index].Spec, namespace)
				CompareContainers(object1, object2, namespace, client1, client2)
			}
			fmt.Printf("----- End checking statefulset: '%s' -----\n\n", name)
		} else {
			fmt.Printf("StatefulSet '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("StatefulSet '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}
