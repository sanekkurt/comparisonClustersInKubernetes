package main

import (
	"fmt"
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

func SetInformationAboutDeployments(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, deployments1 *v1.DeploymentList, deployments2 *v1.DeploymentList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The deployments count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking deployment: '%s' -----\n", name)
			if *deployments1.Items[index1.index].Spec.Replicas != *deployments2.Items[index2.index].Spec.Replicas {
				fmt.Printf("!!!The replicas count are different!!!\n%s '%s' replicas: %d\n%s '%s' replicas: %d\n", kubeconfig1YamlStruct.Clusters[0].Cluster.Server, deployments1.Items[index1.index].Name, *deployments1.Items[index1.index].Spec.Replicas, kubeconfig2YamlStruct.Clusters[0].Cluster.Server, deployments2.Items[index2.index].Name, *deployments2.Items[index2.index].Spec.Replicas)
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
				log.Infof("Deployment %s: %w", name, err)
			}
			fmt.Printf("----- End checking deployment: '%s' -----\n\n", name)
		} else {
			fmt.Printf("Deployment '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("Deployment '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}
