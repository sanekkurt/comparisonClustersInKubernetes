package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

// AddValueStatefulSetsInMap add value StatefulSets in map
func AddValueStatefulSetsInMap(stateFulSets1, stateFulSets2 *v1.StatefulSetList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapStatefulSets1 := make(map[string]CheckerFlag)
	mapStatefulSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range stateFulSets1.Items {
		if _, ok := Entities["statefulsets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapStatefulSets1[value.Name] = indexCheck
	}
	for index, value := range stateFulSets2.Items {
		if _, ok := Entities["statefulsets"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapStatefulSets2[value.Name] = indexCheck
	}
	return mapStatefulSets1, mapStatefulSets2
}

// SetInformationAboutStatefulSets set information about StatefulSets
func SetInformationAboutStatefulSets(map1, map2 map[string]CheckerFlag, statefulSets1, statefulSets2 *v1.StatefulSetList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("StatefulSets count are different")
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
				logging.Log.Debugf("----- Start checking statefulset: '%s' -----", name)
				if *statefulSets1.Items[index1.Index].Spec.Replicas != *statefulSets2.Items[index2.Index].Spec.Replicas {
					logging.Log.Infof("statefulset '%s':  number of replicas is different: %d and %d", statefulSets1.Items[index1.Index].Name, *statefulSets1.Items[index1.Index].Spec.Replicas, *statefulSets2.Items[index2.Index].Spec.Replicas)
					flag = true
				} else {
					// fill in the information that will be used for comparison
					object1 := InformationAboutObject{
						Template: statefulSets1.Items[index1.Index].Spec.Template,
						Selector: statefulSets1.Items[index1.Index].Spec.Selector,
					}
					object2 := InformationAboutObject{
						Template: statefulSets2.Items[index2.Index].Spec.Template,
						Selector: statefulSets2.Items[index2.Index].Spec.Selector,
					}

					err := CompareContainers(object1, object2, namespace, Client1, Client2)
					if err != nil {
						logging.Log.Infof("StatefulSet %s: %s", name, err.Error())
						flag = true
					}

				}
				logging.Log.Debugf("----- End checking statefulset: '%s' -----", name)
			} else {
				logging.Log.Infof("StatefulSet '%s' does not exist in 2nd cluster", name)
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
			logging.Log.Infof("StatefulSet '%s' does not exist in 1st cluster", name)
			flag = true
		}
	}
	return flag
}
