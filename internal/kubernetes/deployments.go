package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

// AddValueDeploymentsInMap add value deployments in map
func AddValueDeploymentsInMap(deployments1, deployments2 *v1.DeploymentList) (map[string]CheckerFlag, map[string]CheckerFlag) { //nolint:gocritic,unused
	mapDeployments1 := make(map[string]CheckerFlag)
	mapDeployments2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range deployments1.Items {
		if _, ok := Entities["deployments"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapDeployments1[value.Name] = indexCheck
	}
	for index, value := range deployments2.Items {
		if _, ok := Entities["deployments"][value.Name]; ok {
			continue
		}
		indexCheck.Index = index
		mapDeployments2[value.Name] = indexCheck
	}
	return mapDeployments1, mapDeployments2
}

// SetInformationAboutDeployments set information about deployments
func SetInformationAboutDeployments(map1, map2 map[string]CheckerFlag, deployments1, deployments2 *v1.DeploymentList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		logging.Log.Infof("deployment counts are different")
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

				logging.Log.Debugf("----- Start checking deployment: '%s' -----", name)
				if *deployments1.Items[index1.Index].Spec.Replicas != *deployments2.Items[index2.Index].Spec.Replicas {
					logging.Log.Infof("deployment '%s':  number of replicas is different: %d and %d", deployments1.Items[index1.Index].Name, *deployments1.Items[index1.Index].Spec.Replicas, *deployments2.Items[index2.Index].Spec.Replicas)
					flag = true
				} else {
					// fill in the information that will be used for comparison
					object1 := InformationAboutObject{
						Template: deployments1.Items[index1.Index].Spec.Template,
						Selector: deployments1.Items[index1.Index].Spec.Selector,
					}
					object2 := InformationAboutObject{
						Template: deployments2.Items[index2.Index].Spec.Template,
						Selector: deployments2.Items[index2.Index].Spec.Selector,
					}
					err := CompareContainers(object1, object2, namespace, Client1, Client2)
					if err != nil {
						logging.Log.Infof("Deployment %s: %s", name, err.Error())
						flag = true
					}
				}
				logging.Log.Debugf("----- End checking deployment: '%s' -----", name)
			} else {
				logging.Log.Infof("Deployment '%s' - 1 cluster. Does not exist on another cluster", name)
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
			logging.Log.Infof("Deployment '%s' - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}
	return flag
}
