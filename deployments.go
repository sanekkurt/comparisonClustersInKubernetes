package main

import (
	v1 "k8s.io/api/apps/v1"
)

func AddValueDeploymentsInMap(deployments1 *v1.DeploymentList, deployments2 *v1.DeploymentList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapDeployments1 := make(map[string]CheckerFlag)
	mapDeployments2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range deployments1.Items {
		indexCheck.index = index
		mapDeployments1[value.Name] = indexCheck
	}
	for index, value := range deployments2.Items {
		indexCheck.index = index
		mapDeployments2[value.Name] = indexCheck
	}
	return mapDeployments1, mapDeployments2
}

func SetInformationAboutDeployments(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, deployments1 *v1.DeploymentList, deployments2 *v1.DeploymentList, namespace string) bool {
	var flag bool
	if len(map1) != len(map2) {
		log.Infof("deployment counts are different")
		flag = true
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2

			log.Debugf("----- Start checking deployment: '%s' -----", name)
			//fmt.Printf("----- Start checking deployment: '%s' -----\n", name)
			if *deployments1.Items[index1.index].Spec.Replicas != *deployments2.Items[index2.index].Spec.Replicas {
				log.Infof("deployment '%s':  number of replicas is different: %d and %d", deployments1.Items[index1.index].Name, *deployments1.Items[index1.index].Spec.Replicas, *deployments2.Items[index2.index].Spec.Replicas)
				flag = true
			} else {
				//заполняем информацию, которая будет использоваться при сравнении
				object1 := InformationAboutObject{
					Template: deployments1.Items[index1.index].Spec.Template,
					Selector: deployments1.Items[index1.index].Spec.Selector,
				}
				object2 := InformationAboutObject{
					Template: deployments2.Items[index2.index].Spec.Template,
					Selector: deployments2.Items[index2.index].Spec.Selector,
				}
				err := CompareContainers(object1, object2, namespace, client1, client2)
				if err != nil {
					log.Infof("Deployment %s: %s", name, err.Error())
					flag = true
				}
			}
			log.Debugf("----- End checking deployment: '%s' -----", name)
		} else {
			log.Infof("Deployment '%s' - 1 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}
	for name, index := range map2 {
		if index.check == false {
			log.Infof("Deployment '%s' - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}
	return flag
}
