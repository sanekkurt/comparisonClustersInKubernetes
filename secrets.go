package main

import (
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

func SetInformationAboutSecrets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, secrets1 *v12.SecretList, secrets2 *v12.SecretList) bool {
	var flag bool
	if len(map1) != len(map2) {
		log.Infof("secret counts are different")
		flag = true
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			//проверка на тип секрета, который проверять не нужно
			if checkContinueTypes(secrets1.Items[index1.index].Type) == true {
				continue
			} else {
				log.Debugf("----- Start checking secret: '%s' -----", name)
				if len(secrets1.Items[index1.index].Data) != len(secrets2.Items[index2.index].Data) {
					log.Infof("secret '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(secrets1.Items[index1.index].Data), len(secrets2.Items[index2.index].Data))
					flag = true
				} else {
					for key, value := range secrets1.Items[index1.index].Data {
						v1 := string(value)
						v2 := string(secrets2.Items[index2.index].Data[key])

						if v1 != v2 {
							log.Infof("secret '%s', values by key '%s' do not match: '%s' and %s", name, key, v1, v2)
							flag = true
						}
					}
				}
			}
			log.Debugf("----- End checking secret: '%s' -----", name)
		} else {
			if checkContinueTypes(secrets1.Items[index1.index].Type) == true {
				continue
			} else {
				log.Infof("secret '%s' does not exist in 2nd cluster", name)
				flag = true
			}
		}
	}
	for name, index := range map2 {
		if index.check == false {
			if checkContinueTypes(secrets2.Items[index.index].Type) == true { //secrets2.Items[index.index].Type == skipType1 || secrets2.Items[index.index].Type == skipType2 || secrets2.Items[index.index].Type == skipType3
				continue
			} else {
				log.Infof("secret '%s' does not exist in 1st cluster", name)
				flag = true
			}
		}
	}
	return flag
}

func checkContinueTypes(secretType v12.SecretType) bool {
	var skip bool
	for _, skipType := range skipTypes {
		if secretType == skipType {
			skip = true
		}
	}
	return skip
}
