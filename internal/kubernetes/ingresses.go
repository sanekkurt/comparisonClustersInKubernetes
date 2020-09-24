package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v1beta12 "k8s.io/api/networking/v1beta1"
	"sync"
)

// AddValueIngressesInMap add value secrets in map
func AddValueIngressesInMap(ingresses1, ingresses2 *v1beta12.IngressList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapIngresses1 := make(map[string]CheckerFlag)
	mapIngresses2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range ingresses1.Items {
		if _, ok := Entities["ingresses"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapIngresses1[value.Name] = indexCheck


	}
	for index, value := range ingresses2.Items {
		if _, ok := Entities["ingresses"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapIngresses2[value.Name] = indexCheck

	}
	return mapIngresses1, mapIngresses2
}

// SetInformationAboutIngresses set information about services
func SetInformationAboutIngresses(map1, map2 map[string]CheckerFlag, ingresses1, ingresses2 *v1beta12.IngressList) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("ingress counts are different")
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
				logging.Log.Debugf("----- Start checking ingress: '%s' -----", name)
				if len(ingresses1.Items[index1.Index].Labels) != len(ingresses2.Items[index2.Index].Labels) {
					logging.Log.Infof("the number of labels is not equal in ingresses. Name ingress: '%s'. In first cluster %d, in second cluster %d", ingresses1.Items[index1.Index].Name, len(ingresses1.Items[index1.Index].Labels), len(ingresses2.Items[index2.Index].Labels))
					flag = true
				} else {
					for key, value := range ingresses1.Items[index1.Index].Labels {
						if ingresses2.Items[index2.Index].Labels[key] != value {
							logging.Log.Infof("labels in ingresses don't match. Name ingress: '%s'. In first cluster: '%s'-'%s', in second cluster value = '%s'", ingresses1.Items[index1.Index].Name, key, value, ingresses2.Items[index2.Index].Labels[key] )
							flag = true
						}
					}
				}
				err := CompareSpecInIngresses(ingresses1.Items[index1.Index], ingresses2.Items[index2.Index])
				if err != nil {
					logging.Log.Infof("Ingress %s: %s", name, err.Error())
					flag = true
				}
				logging.Log.Debugf("----- End checking ingress: '%s' -----", name)
			} else {
				logging.Log.Infof("ingress '%s' does not exist in 2nd cluster", name)
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

			logging.Log.Infof("ingress '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}
	return flag
}