package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

var (
	kubeconfig            = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	kubeconfig1YamlStruct KubeconfigYaml
	kubeconfig2YamlStruct KubeconfigYaml
	client1               *kubernetes.Clientset
	client2               *kubernetes.Clientset
)

func main() {
	flag.Parse()
	home, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	kubeconfig1 := flag.String("kubeconfig1", filepath.Join(home, "kubeconfig1.yaml"), "(optional) absolute path to the kubeconfig file")
	kubeconfig2 := flag.String("kubeconfig2", filepath.Join(home, "kubeconfig2.yaml"), "(optional) absolute path to the kubeconfig file")
	fmt.Println(*kubeconfig1, *kubeconfig2)
	/*if *kubeconfig == "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		*kubeconfig = os.Getenv("KUBECONFIG")
	}*/

	//clientset1 := getClientSet(kubeconfig1)
	//clientset2 := getClientSet(kubeconfig2)
	client1 = getClientSet(kubeconfig1)
	client2 = getClientSet(kubeconfig2)

	//распарсинг yaml файлов в глобальные переменные, чтобы в будущем получить из них URL
	yamlToStruct("kubeconfig1.yaml", &kubeconfig1YamlStruct)
	yamlToStruct("kubeconfig2.yaml", &kubeconfig2YamlStruct)

	compare(client1, client2, "default")
	//fmt.Println("[INFO] Connecting to ", *kubeconfig)

	// use the current context in kubeconfig
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//if err != nil {
	//	panic(err.Error())
	//}

	// create the clientset
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}
	//pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("%d", len(pods.Items))
	//nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{/*FieldSelector: "metadata.name=minikube"*/})
	//if err != nil {
	//	panic(err.Error())
	//}
	////fmt.Printf("There are %d pods in the cluster\n", len(nodes.Items))
	//fmt.Printf(nodes.Items[0].Name)

	//depl, err := clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//for _, d := range depl.Items {
	//	fmt.Printf("%#v\n", d)
	//}
	//
	//fmt.Println("[INFO] Finished")
}

//переводит yaml в структуру
func yamlToStruct(nameYamlFile string, nameStruct *KubeconfigYaml) {
	data, err := ioutil.ReadFile(nameYamlFile)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(data, nameStruct)
	if err != nil {
		panic(err.Error())
	}
}

//читает конфигурацию из yaml файла по переданному пути
func getClientSet(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

//основная сравнивающая функция, поочередно запускает функции для сравнения по разным параметрам
func compare(clientSet1 *kubernetes.Clientset, clientSet2 *kubernetes.Clientset, namespaces ...string) {
	for _, namespace := range namespaces {
		depl1, err := clientSet1.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		depl2, err := clientSet2.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		mapDeployments1, mapDeployments2 := addValueDeploymentsInMap(depl1, depl2)
		setInformationAboutDeployments(mapDeployments1, mapDeployments2, depl1, depl2, namespace)

		//compareDeployments(depl1, depl2)
		//compareReplicasInDeployments(depl1, depl2)
		//compareImagesInDeployments(depl1, depl2)
	}
}

func setInformationAboutDeployments(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, deployment1 *v1.DeploymentList, deployment2 *v1.DeploymentList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The deployments count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking deployments: '%s' -----\n", name)
			if *deployment1.Items[index1.index].Spec.Replicas != *deployment2.Items[index2.index].Spec.Replicas {
				fmt.Printf("!!!The replicas count are different!!!\n%s '%s' replicas: %d\n%s '%s' replicas: %d\n", kubeconfig1YamlStruct.Clusters[0].Cluster.Server, deployment1.Items[index1.index].Name, *deployment1.Items[index1.index].Spec.Replicas, kubeconfig2YamlStruct.Clusters[0].Cluster.Server, deployment2.Items[index2.index].Name, *deployment2.Items[index2.index].Spec.Replicas)
			} else {
				compareContainers(deployment1.Items[index1.index].Spec, deployment2.Items[index2.index].Spec, namespace)
			}
			fmt.Printf("----- End checking deployments: '%s' -----\n\n", name)
		} else {
			fmt.Printf("'%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("'%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}

func compareContainers(deploymentSpec1 v1.DeploymentSpec, deploymentSpec2 v1.DeploymentSpec, namespace string) error {
	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		fmt.Printf("!!!The number of containers differs!!!")
		return errors.New("The number of containers differs")
	} else {
		matchLabelsString1 := convertMatchlabelsToString(deploymentSpec1.Selector.MatchLabels)
		matchLabelsString2 := convertMatchlabelsToString(deploymentSpec2.Selector.MatchLabels)
		if matchLabelsString1 != matchLabelsString2 {
			fmt.Printf("!!!MatchLabels are not equal!!!")
			return errors.New("MatchLabels are not equal")
		}
		pods1, pods2 := getPodsListOnmatchLabels(deploymentSpec1.Selector.MatchLabels, namespace)
		for i := 0; i < len(containersDeploymentTemplate1); i++ {
			if containersDeploymentTemplate1[i].Name != containersDeploymentTemplate2[i].Name {
				fmt.Printf("!!!Container names are not equal!!!")
				return errors.New("Container names are not equal")
			} else if containersDeploymentTemplate1[i].Image != containersDeploymentTemplate2[i].Image {
				fmt.Printf("!!!Container name images are not equal!!!")
				return errors.New("Container name images are not equal")
			} else {
				if len(pods1.Items) != len(pods2.Items) {
					fmt.Printf("!!!The replicas count are different!!!")
					return errors.New("The replicas count are different")
				} else {
					for j := 0; j < len(pods1.Items); j++ {
						containersStatusesInPod1 := getContainerStatusesInPod(pods1.Items[j].Status.ContainerStatuses)
						containersStatusesInPod2 := getContainerStatusesInPod(pods2.Items[j].Status.ContainerStatuses)
						if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
							fmt.Printf("!!!The containers count in pod are different!!!")
							return errors.New("The containers count in pod are different")
						} else {
							var flag int
							var containerWithSameNameFound bool
							for f:=0; f<len(containersStatusesInPod1); f++{
								if containersDeploymentTemplate1[i].Name == containersStatusesInPod1[f].name {
									flag++
									if containersDeploymentTemplate1[i].Image != containersStatusesInPod1[f].image {
										fmt.Printf("!!!The container image in the template does not match the actual image in the Pod!!!")
										return errors.New("The container image in the template does not match the actual image in the Pod")
									} else {
										for _, value := range containersStatusesInPod2{
											if containersStatusesInPod1[f].name == value.name{
												containerWithSameNameFound = true
												if containersStatusesInPod1[f].image != value.image{
													textForError:=fmt.Sprintf("!!!The Image in Pods is different!!!\nPods name: '%s'\nImage name on pod1: '%s'\nImage name on pod2: '%s'\n\n",value.name, containersStatusesInPod1[j].image,value.image)
													fmt.Printf(textForError)
													return errors.New(textForError)
												} else if containersStatusesInPod1[f].imageID != value.imageID{
													textForError:=fmt.Sprintf("!!!The ImageID in Pods is different!!!\nPods name: '%s'\nImageID on pod1: '%s'\nImageID on pod2: '%s'\n\n",value.name, containersStatusesInPod1[j].imageID,value.imageID)
													fmt.Printf(textForError)
													return errors.New(textForError)
												}
											}
										}
										if !containerWithSameNameFound{
											textForError:=fmt.Sprintf("Container '%s' not found on other pod", containersStatusesInPod1[j].name)
											fmt.Printf("!!!Container '%s' not found on other pod!!!", containersStatusesInPod1[j].name)
											return errors.New(textForError)
										}
									}
								}
							}

						}
						/*if containersDeploymentTemplate1[i].Name != containersStatusesInPod1["Name"] || containersDeploymentTemplate1[i].Name != containersStatusesInPod2["Name"]{
							fmt.Printf("!!!The container name in the template does not match the actual name in the Pod!!!")
							return errors.New("The container name in the template does not match the actual name in the Pod")
						} else if  containersDeploymentTemplate1[i].Image != containersStatusesInPod1["Image"] || containersDeploymentTemplate1[i].Image != containersStatusesInPod2["Image"] {
							fmt.Printf("!!!The container image in the template does not match the actual image in the Pod!!!")
							return errors.New("The container image in the template does not match the actual image in the Pod")
						} else if  containersStatusesInPod1["ImageID"] != containersStatusesInPod2["ImageID"] {
							fmt.Printf("!!!The ImageID in Pods is different!!!")
							return errors.New("The ImageID in Pods is different")
						}*/
					}
				}
				//if pods1, pods2 := getPodsListOnmatchLabels(deploymentSpec1.Selector.MatchLabels, namespace); pods1.Items[0].Status.ContainerStatuses[0].ImageID != pods2.Items[0].Status.ContainerStatuses[0].ImageID {
				//	fmt.Printf("!!!Container imageID are not equal!!!")
				//	return errors.New("Container imageID are not equal")
				//}
			}
		}
		return nil
	}
}

func getContainerStatusesInPod(containerStatuses []v12.ContainerStatus) map[int]Container {
	infoAboutContainer := make(map[int]Container)
	var container Container
	for index, value := range containerStatuses{
		container.name = value.Name
		container.image = value.Image
		container.imageID = value.ImageID
		infoAboutContainer[index] = container
	}
	return infoAboutContainer
}

//получает айдишник раскатанного образа на контейнерах
func getPodsListOnmatchLabels(matchLabels map[string]string, namespace string) (*v12.PodList, *v12.PodList) {
	matchLabelsString := convertMatchlabelsToString(matchLabels)
	pods1, err := client1.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	pods2, err := client2.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	//return pods1.Items[0].Status.ContainerStatuses[0].ImageID, pods2.Items[0].Status.ContainerStatuses[0].ImageID
	return pods1, pods2
}

func convertMatchlabelsToString(matchLabels map[string]string) string {
	values := []string{}
	for key, value := range matchLabels {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}
	//супермегафича склеивания строчек
	return strings.Join(values, ",")
}

func addValueDeploymentsInMap(deployment1 *v1.DeploymentList, deployment2 *v1.DeploymentList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapDeployments1 := make(map[string]CheckerFlag)
	mapDeployments2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range deployment1.Items {
		indexCheck.index = index
		mapDeployments1[value.Name] = indexCheck
	}
	for index, value := range deployment2.Items {
		indexCheck.index = index
		mapDeployments2[value.Name] = indexCheck
	}
	return mapDeployments1, mapDeployments2
}

//func compareDeployments(deployment1 *v1.DeploymentList, deployment2 *v1.DeploymentList) {
//	countDeployments1 := len(deployment1.Items)
//	countDeployments2 := len(deployment2.Items)
//	if countDeployments1 != countDeployments2 {
//		fmt.Printf("!!!The deployments count are different!!!\n")
//	}
//	//вызов функции для проверки несовпадающих деплойментов первого кластера
//	badNames1 := badDeploymentsInCluster(deployment1, deployment2)
//	if badNames1 != nil {
//		fmt.Printf("\nBad deployments in 1 cluster:\n")
//		for _, value := range badNames1 {
//			fmt.Printf("%s\n", value)
//		}
//	}
//	//вызов функции для проверки несовпадающих деплойментов второго кластера
//	badNames2 := badDeploymentsInCluster(deployment2, deployment1)
//	if badNames2 != nil {
//		fmt.Printf("\nBad deployments in 2 cluster:\n")
//		for _, value := range badNames2 {
//			fmt.Printf("%s\n", value)
//		}
//	}
//}
//
////сравнивает реплики в кластерах
//func compareReplicasInDeployments(deployment1 *v1.DeploymentList, deployment2 *v1.DeploymentList) {
//	countDeployments1 := len(deployment1.Items)
//	countDeployments2 := len(deployment2.Items)
//	for i := 0; i < countDeployments1; i++ {
//		for j := 0; j < countDeployments2; j++ {
//			if deployment1.Items[i].Name == deployment2.Items[j].Name {
//				if *deployment1.Items[i].Spec.Replicas != *deployment2.Items[j].Spec.Replicas {
//					fmt.Printf("!!!The replicas count are different!!!\n%s '%s' replicas: %d\n%s '%s' replicas: %d\n\n", kubeconfig1YamlStruct.Clusters[0].Cluster.Server, deployment1.Items[i].Name, *deployment1.Items[i].Spec.Replicas, kubeconfig2YamlStruct.Clusters[0].Cluster.Server, deployment2.Items[j].Name, *deployment2.Items[j].Spec.Replicas)
//				}
//			}
//		}
//	}
//}
//
//func compareImagesInDeployments(deployment1 *v1.DeploymentList, deployment2 *v1.DeploymentList) {
//	countDeployments1 := len(deployment1.Items)
//	countDeployments2 := len(deployment2.Items)
//	//пробегаемся по деплойментам
//	for i := 0; i < countDeployments1; i++ {
//		for j := 0; j < countDeployments2; j++ {
//			//			compareContainers(deployment1.Items[i].Spec.Template.Spec.Containers)
//			//если их имена равны то пробегаемся по их контейнерам, чтобы сравнить в них image
//			if deployment1.Items[i].Name == deployment2.Items[j].Name {
//				if len(deployment1.Items[i].Spec.Template.Spec.Containers) != len(deployment2.Items[j].Spec.Template.Spec.Containers) {
//					fmt.Printf("!!!In deployments '%s' different number of containers", deployment1.Items[i].Name)
//					badContainers1 := badContainersInCluster(deployment1.Items[i].Spec.Template.Spec.Containers, deployment2.Items[j].Spec.Template.Spec.Containers)
//					if badContainers1 != nil {
//						fmt.Printf("\nBad containers in deployments '%s':\n", deployment1.Items[i].Name)
//						for _, value := range badContainers1 {
//							fmt.Printf("%s\n", value)
//						}
//					}
//					badContainers2 := badContainersInCluster(deployment2.Items[j].Spec.Template.Spec.Containers, deployment1.Items[i].Spec.Template.Spec.Containers)
//					if badContainers2 != nil {
//						fmt.Printf("\nBad containers in deployments '%s':\n", deployment2.Items[j].Name)
//						for _, value := range badContainers2 {
//							fmt.Printf("%s\n", value)
//						}
//					}
//				}
//				for a := 0; a < len(deployment1.Items[i].Spec.Template.Spec.Containers); a++ {
//					for b := 0; b < len(deployment2.Items[j].Spec.Template.Spec.Containers); b++ {
//
//						if deployment1.Items[i].Spec.Template.Spec.Containers[a].Image == deployment2.Items[j].Spec.Template.Spec.Containers[b].Image {
//
//						}
//
//					}
//				}
//			}
//		}
//	}
//}
//
//func badContainersInCluster(containers1 []v12.Container, containers2 []v12.Container) []string {
//	var containersNames []string
//	flag := 0
//
//	for i := 0; i < len(containers1); i++ {
//		for j := 0; j < len(containers2); j++ {
//			if containers1[i].Name != containers2[j].Name {
//				flag++
//			}
//		}
//		if flag == len(containers2) {
//			containersNames = append(containersNames, containers1[i].Name)
//		}
//		flag = 0
//	}
//	return containersNames
//}
//
//func badDeploymentsInCluster(depl1 *v1.DeploymentList, depl2 *v1.DeploymentList) []string {
//	var names1 []string
//	flag := 0
//
//	for i := 0; i < len(depl1.Items); i++ {
//		for j := 0; j < len(depl2.Items); j++ {
//			if depl1.Items[i].Name != depl2.Items[j].Name {
//				flag++
//			}
//		}
//		if flag == len(depl2.Items) {
//			names1 = append(names1, depl1.Items[i].Name)
//		}
//		flag = 0
//	}
//	return names1
//}
