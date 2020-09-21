package main

import (
	v12 "k8s.io/api/core/v1"
)

// Добавление значений Secrets в карту для дальнейшего сравнения
func AddValueSecretsInMap(secrets1, secrets2 *v12.SecretList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapSecrets1 := make(map[string]CheckerFlag)
	mapSecrets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range secrets1.Items {
		if checkContinueTypes(value.Type) {
			continue
		} else {
			indexCheck.index = index
			mapSecrets1[value.Name] = indexCheck
		}

	}
	for index, value := range secrets2.Items {
		if checkContinueTypes(value.Type) {
			continue
		} else {
			indexCheck.index = index
			mapSecrets2[value.Name] = indexCheck
		}
	}
	return mapSecrets1, mapSecrets2
}

// Получение информации о Secrets
func SetInformationAboutSecrets(map1, map2 map[string]CheckerFlag, secrets1, secrets2 *v12.SecretList) bool {
	var flag bool
	if len(map1) != len(map2) {
		log.Infof("secret counts are different")
		flag = true
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			//проверка на тип секрета, который проверять не нужно

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

			log.Debugf("----- End checking secret: '%s' -----", name)
		} else {

			log.Infof("secret '%s' does not exist in 2nd cluster", name)
			flag = true

		}
	}
	for name, index := range map2 {
		if !index.check {

			log.Infof("secret '%s' does not exist in 1st cluster", name)
			flag = true

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
