package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v12 "k8s.io/api/core/v1"
	"sync"
)

// AddValueServicesInMap add value secrets in map
func AddValueServicesInMap(services1, services2 *v12.ServiceList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapServices1 := make(map[string]CheckerFlag)
	mapServices2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range services1.Items {
		if _, ok := Entities["services"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapServices1[value.Name] = indexCheck


	}
	for index, value := range services2.Items {
		if _, ok := Entities["services"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapServices2[value.Name] = indexCheck

	}
	return mapServices1, mapServices2
}

// SetInformationAboutServices set information about services
func SetInformationAboutServices(map1, map2 map[string]CheckerFlag, services1, services2 *v12.ServiceList) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("service counts are different")
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
				logging.Log.Debugf("----- Start checking service: '%s' -----", name)
				if len(services1.Items[index1.Index].Labels) != len(services2.Items[index2.Index].Labels) {
					logging.Log.Infof("the number of labels is not equal in services. Name service: '%s'. In first cluster %d, in second cluster %d", services1.Items[index1.Index].Name, len(services1.Items[index1.Index].Labels), len(services2.Items[index2.Index].Labels))
					flag = true
				} else {
					for key, value := range services1.Items[index1.Index].Labels {
						if services2.Items[index2.Index].Labels[key] != value {
							logging.Log.Infof("labels in services don't match. Name service: '%s'. In first cluster: '%s'-'%s', in second cluster value = '%s'", services1.Items[index1.Index].Name, key, value, services2.Items[index2.Index].Labels[key] )
							flag = true
						}
					}
				}
				err := CompareSpecInServices(services1.Items[index1.Index], services2.Items[index2.Index])
				if err != nil {
					logging.Log.Infof("Service %s: %s", name, err.Error())
					flag = true
				}
				logging.Log.Debugf("----- End checking service: '%s' -----", name)
			} else {
				logging.Log.Infof("service '%s' does not exist in 2nd cluster", name)
				flag = true
				channel <- flag
			}
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

			logging.Log.Infof("service '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}
	return flag
}