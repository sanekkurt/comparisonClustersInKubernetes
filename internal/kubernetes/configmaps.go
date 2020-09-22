package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v12 "k8s.io/api/core/v1"
	"sync"
)

// AddValueConfigMapsInMap add value ConfigMaps in map
func AddValueConfigMapsInMap(configMaps1, configMaps2 *v12.ConfigMapList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapConfigMap1 := make(map[string]CheckerFlag)
	mapConfigMap2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range configMaps1.Items {
		if _, ok := Entities["configmaps"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapConfigMap1[value.Name] = indexCheck
	}
	for index, value := range configMaps2.Items {
		if _, ok := Entities["configmaps"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapConfigMap2[value.Name] = indexCheck
	}
	return mapConfigMap1, mapConfigMap2
}

// SetInformationAboutConfigMaps set information about config maps
func SetInformationAboutConfigMaps(map1, map2 map[string]CheckerFlag, configMaps1, configMaps2 *v12.ConfigMapList) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("configmaps count are different")
		flag = true
	}
	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))
	for name, index1 := range map1 {
		wg.Add(1)
		go func(wg *sync.WaitGroup, channel chan bool, name string, index1 CheckerFlag, map1, map2 map[string]CheckerFlag) {
			defer func() {
				wg.Done()
			}()
			if index2, ok := map2[name]; ok {
				index1.Check = true
				map1[name] = index1
				index2.Check = true
				map2[name] = index2
				logging.Log.Debugf("----- Start checking configmap: '%s' -----", name)
				if len(configMaps1.Items[index1.Index].Data) != len(configMaps2.Items[index2.Index].Data) {
					logging.Log.Infof("config map '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(configMaps1.Items[index1.Index].Data), len(configMaps2.Items[index2.Index].Data))
					flag = true
				} else {
					for key, value := range configMaps1.Items[index1.Index].Data {
						if configMaps2.Items[index2.Index].Data[key] != value {
							logging.Log.Infof("configmap '%s', values by key '%s' do not match: '%s' and %s", name, key, value, configMaps2.Items[index2.Index].Data[key])
							flag = true
						}
					}
				}
				logging.Log.Debugf("----- End checking configmap: '%s' -----", name)
			} else {
				logging.Log.Infof("ConfigMap '%s' - 1 cluster. Does not exist on another cluster", name)
				flag = true
			}
			channel <- flag
		}(wg, channel, name,index1, map1, map2)
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
			logging.Log.Infof("ConfigMap '%s' - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}
	return flag
}
