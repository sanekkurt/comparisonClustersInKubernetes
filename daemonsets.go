package main

import (
	v1 "k8s.io/api/apps/v1"
	"sync"
)

//Добавление значений DaemonSets в карту для дальнейшего сравнения
func AddValueDaemonSetsMap(daemonSets1, daemonSets2 *v1.DaemonSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapDaemonSets1 := make(map[string]CheckerFlag)
	mapDaemonSets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range daemonSets1.Items {
		indexCheck.index = index
		mapDaemonSets1[value.Name] = indexCheck
	}
	for index, value := range daemonSets2.Items {
		indexCheck.index = index
		mapDaemonSets2[value.Name] = indexCheck
	}
	return mapDaemonSets1, mapDaemonSets2
}

//Получение информации о DaemonSets
func SetInformationAboutDaemonSets(map1, map2 map[string]CheckerFlag, daemonSets1, daemonSets2 *v1.DaemonSetList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		log.Infof("DaemonSet count are different")
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
				log.Debugf("----- Start checking daemonset: '%s' -----", name)

				//заполняем информацию, которая будет использоваться при сравнении
				object1 := InformationAboutObject{
					Template: daemonSets1.Items[index1.index].Spec.Template,
					Selector: daemonSets1.Items[index1.index].Spec.Selector,
				}
				object2 := InformationAboutObject{
					Template: daemonSets2.Items[index2.index].Spec.Template,
					Selector: daemonSets2.Items[index2.index].Spec.Selector,
				}

				err := CompareContainers(object1, object2, namespace, client1, client2)
				if err != nil {
					log.Infof("DaemonSet %s: %s", name, err.Error())
					flag = true
				}
				log.Debugf("----- End checking daemonset: '%s' -----", name)
			} else {
				log.Infof("DaemonSet '%s' does not exist in 2nd cluster", name)
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
		if !index.check {
			log.Infof("DaemonSet '%s' does not exist in 1s cluster", name)
			flag = true
		}
	}
	return flag
}
