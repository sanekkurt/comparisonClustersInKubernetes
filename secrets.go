package main

import (
	"fmt"
	v12 "k8s.io/api/core/v1"
)

func AddValueSecretsInMap(secrets1 *v12.SecretList, secrets2 *v12.SecretList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapSecrets1 := make(map[string]CheckerFlag)
	mapSecrets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range secrets1.Items {
		indexCheck.index = index
		mapSecrets1[value.Name] = indexCheck
	}
	for index, value := range secrets2.Items {
		indexCheck.index = index
		mapSecrets2[value.Name] = indexCheck
	}
	return mapSecrets1, mapSecrets2
}

func SetInformationAboutSecrets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, secrets1 *v12.SecretList, secrets2 *v12.SecretList) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The secrets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			//проверка на тип секрета, который проверять не нужно
			if secrets1.Items[index1.index].Type == skipType1 || secrets1.Items[index1.index].Type == skipType2 {
				continue
			} else {
				fmt.Printf("----- Start checking secret: '%s' -----\n", name)
				if len(secrets1.Items[index1.index].Data) != len(secrets2.Items[index2.index].Data) {
					fmt.Printf("!!!Config map '%s' in 1 kluster have '%d' key value pair but 2 kluster have '%d' key value pair!!!\n", name, len(secrets1.Items[index1.index].Data), len(secrets2.Items[index2.index].Data))
				} else {
					for key, value := range secrets1.Items[index1.index].Data {
						if string(value) != string(secrets2.Items[index2.index].Data[key]) {
							fmt.Printf("!!!The key value pair does not match. In 1 kluster %s: %s. In 2 kluster %s: %s.!!!\n", key, string(value), key, string(secrets2.Items[index2.index].Data[key]))
						}
					}
				}
			}
			fmt.Printf("----- End checking secret: '%s' -----\n\n", name)
		} else {
			if secrets1.Items[index1.index].Type == skipType1 || secrets1.Items[index1.index].Type == skipType2 {
				continue
			} else {
				fmt.Printf("Secret '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
			}
		}
	}
	for name, index := range map2 {
		if index.check == false {
			if secrets2.Items[index.index].Type == skipType1 || secrets2.Items[index.index].Type == skipType2 {
				continue
			} else {
				fmt.Printf("Secret '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
			}
		}
	}
}