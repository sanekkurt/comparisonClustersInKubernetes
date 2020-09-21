package main

import (
	v1 "k8s.io/api/apps/v1"
	"sync"
)

// Добавление значений StatefulSets в карту для дальнейшего сравнения
func AddValueStatefulSetsInMap(stateFulSets1, stateFulSets2 *v1.StatefulSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapStatefulSets1 := make(map[string]CheckerFlag)
	mapStatefulSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range stateFulSets1.Items {
		indexCheck.index = index
		mapStatefulSets1[value.Name] = indexCheck
	}
	for index, value := range stateFulSets2.Items {
		indexCheck.index = index
		mapStatefulSets2[value.Name] = indexCheck
	}
	return mapStatefulSets1, mapStatefulSets2
}

// Получение информации о StatefulSets
func SetInformationAboutStatefulSets(map1, map2 map[string]CheckerFlag, statefulSets1, statefulSets2 *v1.StatefulSetList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		log.Infof("StatefulSets count are different")
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
				index1.check = true
				map1[name] = index1
				index2.check = true
				map2[name] = index2
				log.Debugf("----- Start checking statefulset: '%s' -----", name)
				if *statefulSets1.Items[index1.index].Spec.Replicas != *statefulSets2.Items[index2.index].Spec.Replicas {
					log.Infof("statefulset '%s':  number of replicas is different: %d and %d", statefulSets1.Items[index1.index].Name, *statefulSets1.Items[index1.index].Spec.Replicas, *statefulSets2.Items[index2.index].Spec.Replicas)
					flag = true
				} else {
					//заполняем информацию, которая будет использоваться при сравнении
					object1 := InformationAboutObject{
						Template: statefulSets1.Items[index1.index].Spec.Template,
						Selector: statefulSets1.Items[index1.index].Spec.Selector,
					}
					object2 := InformationAboutObject{
						Template: statefulSets2.Items[index2.index].Spec.Template,
						Selector: statefulSets2.Items[index2.index].Spec.Selector,
					}

					err := CompareContainers(object1, object2, namespace, client1, client2)
					if err != nil {
						log.Infof("StatefulSet %s: %s", name, err.Error())
						flag = true
					}

				}
				log.Debugf("----- End checking statefulset: '%s' -----", name)
			} else {
				log.Infof("StatefulSet '%s' does not exist in 2nd cluster", name)
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
		if !index.check {
			log.Infof("StatefulSet '%s' does not exist in 1st cluster", name)
			flag = true
		}
	}
	return flag
}
