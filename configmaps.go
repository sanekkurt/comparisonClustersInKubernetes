package main

import (
	"fmt"
	v12 "k8s.io/api/core/v1"
)

func AddValueConfigMapsInMap(configMaps1 *v12.ConfigMapList, configMaps2 *v12.ConfigMapList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapConfigMap1 := make(map[string]CheckerFlag)
	mapConfigMap2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range configMaps1.Items {
		indexCheck.index = index
		mapConfigMap1[value.Name] = indexCheck
	}
	for index, value := range configMaps2.Items {
		indexCheck.index = index
		mapConfigMap2[value.Name] = indexCheck
	}
	return mapConfigMap1, mapConfigMap2
}

func SetInformationAboutConfigMaps(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, configMaps1 *v12.ConfigMapList, configMaps2 *v12.ConfigMapList) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The configmaps count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking configmap: '%s' -----\n", name)
			if len(configMaps1.Items[index1.index].Data) != len(configMaps2.Items[index2.index].Data) {
				fmt.Printf("!!!Config map '%s' in 1 cluster have '%d' key value pair but 2 kluster have '%d' key value pair!!!\n", name, len(configMaps1.Items[index1.index].Data), len(configMaps2.Items[index2.index].Data))
			} else {
				for key, value := range configMaps1.Items[index1.index].Data {
					if configMaps2.Items[index2.index].Data[key] != value {
						fmt.Printf("!!!The key value pair does not match. In 1 kluster %s: %s. In 2 kluster %s: %s.!!!\n", key, value, key, configMaps2.Items[index2.index].Data[key])
					}
				}
			}
			fmt.Printf("----- End checking configmap: '%s' -----\n\n", name)
		} else {
			fmt.Printf("ConfigMap '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("ConfigMap '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}
