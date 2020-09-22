package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

// AddValueDaemonSetsMap add value daemonSets in map
func AddValueDaemonSetsMap(daemonSets1, daemonSets2 *v1.DaemonSetList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapDaemonSets1 := make(map[string]CheckerFlag)
	mapDaemonSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range daemonSets1.Items {
		if _, ok := Entities["daemonsets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapDaemonSets1[value.Name] = indexCheck
	}
	for index, value := range daemonSets2.Items {
		if _, ok := Entities["daemonsets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapDaemonSets2[value.Name] = indexCheck
	}
	return mapDaemonSets1, mapDaemonSets2
}

// SetInformationAboutDaemonSets set information about daemonSets
func SetInformationAboutDaemonSets(map1, map2 map[string]CheckerFlag, daemonSets1, daemonSets2 *v1.DaemonSetList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("DaemonSet count are different")
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
				logging.Log.Debugf("----- Start checking daemonset: '%s' -----", name)

				// fill in the information that will be used for comparison
				object1 := InformationAboutObject{
					Template: daemonSets1.Items[index1.Index].Spec.Template,
					Selector: daemonSets1.Items[index1.Index].Spec.Selector,
				}
				object2 := InformationAboutObject{
					Template: daemonSets2.Items[index2.Index].Spec.Template,
					Selector: daemonSets2.Items[index2.Index].Spec.Selector,
				}

				err := CompareContainers(object1, object2, namespace, Client1, Client2)
				if err != nil {
					logging.Log.Infof("DaemonSet %s: %s", name, err.Error())
					flag = true
				}
				logging.Log.Debugf("----- End checking daemonset: '%s' -----", name)
			} else {
				logging.Log.Infof("DaemonSet '%s' does not exist in 2nd cluster", name)
				flag = true
			}
		channel <- flag
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
			logging.Log.Infof("DaemonSet '%s' does not exist in 1s cluster", name)
			flag = true
		}
	}
	return flag
}
