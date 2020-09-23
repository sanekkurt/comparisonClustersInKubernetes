package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v12 "k8s.io/api/core/v1"
	"sync"
)

// AddValueSecretsInMap add value secrets in map
func AddValueSecretsInMap(secrets1, secrets2 *v12.SecretList) (map[string]IsAlreadyComparedFlag, map[string]IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapSecrets1 := make(map[string]IsAlreadyComparedFlag)
	mapSecrets2 := make(map[string]IsAlreadyComparedFlag)
	var indexCheck IsAlreadyComparedFlag

	for index, value := range secrets1.Items {
		if checkContinueTypes(value.Type) {
			continue
		}
		if _, ok := ToSkipEntities["secrets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapSecrets1[value.Name] = indexCheck

	}
	for index, value := range secrets2.Items {
		if checkContinueTypes(value.Type) {
			continue
		}
		if _, ok := ToSkipEntities["secrets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapSecrets2[value.Name] = indexCheck

	}
	return mapSecrets1, mapSecrets2
}

// SetInformationAboutSecrets set information about secrets
func SetInformationAboutSecrets(map1, map2 map[string]IsAlreadyComparedFlag, secrets1, secrets2 *v12.SecretList) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("secret counts are different")
		flag = true
	}
	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))
	for name, index1 := range map1 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, channel chan bool, name string, index1 IsAlreadyComparedFlag, map1, map2 map[string]IsAlreadyComparedFlag) {
			defer func() {
				wg.Done()
			}()
			if index2, ok := map2[name]; ok {
				index1.Check = true
				map1[name] = index1
				index2.Check = true
				map2[name] = index2
				// проверка на тип секрета, который проверять не нужно

				logging.Log.Debugf("----- Start checking secret: '%s' -----", name)
				if len(secrets1.Items[index1.Index].Data) != len(secrets2.Items[index2.Index].Data) {
					logging.Log.Infof("secret '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(secrets1.Items[index1.Index].Data), len(secrets2.Items[index2.Index].Data))
					flag = true
				} else {
					for key, value := range secrets1.Items[index1.Index].Data {
						v1 := string(value)
						v2 := string(secrets2.Items[index2.Index].Data[key])

						if v1 != v2 {
							logging.Log.Infof("secret '%s', values by key '%s' do not match: '%s' and %s", name, key, v1, v2)
							flag = true
						}
					}
				}

				logging.Log.Debugf("----- End checking secret: '%s' -----", name)
			} else {

				logging.Log.Infof("secret '%s' does not exist in 2nd cluster", name)
				flag = true
				channel <- flag
			}
		}(wg, channel, name, index1, map1, map2)
	}
	wg.Wait()
	close(channel)
	for ch := range channel {
		if ch {
			flag = true
		}
	}
	for name, index := range map2 {
		if !index.Check {

			logging.Log.Infof("secret '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}
	return flag
}

func checkContinueTypes(secretType v12.SecretType) bool {
	var skip bool
	for _, skipType := range SkipTypes {
		if secretType == skipType {
			skip = true
		}
	}
	return skip
}
