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

	//clientset1 := GetClientSet(kubeconfig1)
	//clientset2 := GetClientSet(kubeconfig2)
	client1 = GetClientSet(kubeconfig1)
	client2 = GetClientSet(kubeconfig2)

	//распарсинг yaml файлов в глобальные переменные, чтобы в будущем получить из них URL
	YamlToStruct("kubeconfig1.yaml", &kubeconfig1YamlStruct)
	YamlToStruct("kubeconfig2.yaml", &kubeconfig2YamlStruct)

	Compare(client1, client2, "default")
}

//переводит yaml в структуру
func YamlToStruct(nameYamlFile string, nameStruct *KubeconfigYaml) {
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
func GetClientSet(kubeconfig *string) *kubernetes.Clientset {
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

//основная сравнивающая функция, поочередно запускает функции для сравнения кластеров по разным параметрам: Deployments, StatefulSets, DaemonSets, ConfigMaps
func Compare(clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface, namespaces ...string) {
	for _, namespace := range namespaces {
		depl1, err := clientSet1.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		depl2, err := clientSet2.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapDeployments1, mapDeployments2 := AddValueDeploymentsInMap(depl1, depl2)
		SetInformationAboutDeployments(mapDeployments1, mapDeployments2, depl1, depl2, namespace)

		statefulSet1, err := clientSet1.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		statefulSet2, err := clientSet2.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapStatefulSets1, mapStatefulSets2 := AddValueStatefulSetsInMap(statefulSet1, statefulSet2)
		SetInformationAboutStatefulSets(mapStatefulSets1, mapStatefulSets2, statefulSet1, statefulSet2, namespace)

		daemonSets1, err := clientSet1.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		daemonSets2, err := clientSet2.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapDaemonSets1, mapDaemonSets2 := AddValueDaemonSetsMap(daemonSets1, daemonSets2)
		SetInformationAboutDaemonSets(mapDaemonSets1, mapDaemonSets2, daemonSets1, daemonSets2, namespace)

		configMaps1, err := clientSet1.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		configMaps2, err := clientSet2.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapConfigMaps1, mapConfigMaps2 := AddValueConfigMapsInMap(configMaps1, configMaps2)
		SetInformationAboutConfigMaps(mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2)

		secrets1, err := clientSet1.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		secrets2, err := clientSet2.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapSecrets1, mapSecrets2 := AddValueSecretsInMap(secrets1, secrets2)
		SetInformationAboutSecrets(mapSecrets1, mapSecrets2, secrets1, secrets2)
		//compareDeployments(depl1, depl2)
		//compareReplicasInDeployments(depl1, depl2)
		//compareImagesInDeployments(depl1, depl2)
	}
}

func SetInformationAboutSecrets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, secrets1 *v12.SecretList, secrets2 *v12.SecretList) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The secrets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		// type=kubernetes.io/service-account-token -> skip
		// type=kubernetes.io/dockercfg -> skip
		//secrets1.Items[0].Type
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking secret: '%s' -----\n", name)
			if len(secrets1.Items[index1.index].Data) != len(secrets2.Items[index2.index].Data){
				fmt.Printf("!!!Config map '%s' in 1 kluster have '%d' key value pair but 2 kluster have '%d' key value pair!!!\n", name, len(secrets1.Items[index1.index].Data), len(secrets2.Items[index2.index].Data))
			} else {
				for key, value := range secrets1.Items[index1.index].Data {
					if string(value) != string(secrets2.Items[index2.index].Data[key]){
						fmt.Printf("!!!The key value pair does not match. In 1 kluster %s: %s. In 2 kluster %s: %s.!!!\n", key, string(value), key, string(secrets2.Items[index2.index].Data[key]))
					}
				}
			}
			fmt.Printf("----- End checking secret: '%s' -----\n\n", name)
		} else {
			fmt.Printf("Secret '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("Secret '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}

func SetInformationAboutConfigMaps(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, configMaps1 *v12.ConfigMapList, configMaps2 *v12.ConfigMapList) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The configmaps count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking configmap: '%s' -----\n", name)
			if len(configMaps1.Items[index1.index].Data) != len(configMaps2.Items[index2.index].Data){
				fmt.Printf("!!!Config map '%s' in 1 cluster have '%d' key value pair but 2 kluster have '%d' key value pair!!!\n", name, len(configMaps1.Items[index1.index].Data), len(configMaps2.Items[index2.index].Data))
			} else {
				for key, value := range configMaps1.Items[index1.index].Data {
					if configMaps2.Items[index2.index].Data[key] != value{
						fmt.Printf("!!!The key value pair does not match. In 1 kluster %s: %s. In 2 kluster %s: %s.!!!\n", key, value, key, configMaps2.Items[index2.index].Data[key])
					}
				}
			}
			fmt.Printf("----- End checking configmap: '%s' -----\n\n", name)
		} else {
			fmt.Printf("ConfigMap '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("ConfigMap '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}

func SetInformationAboutDaemonSets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, daemonSets1 *v1.DaemonSetList, daemonSets2 *v1.DaemonSetList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The daemonsets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking daemonset: '%s' -----\n", name)

			//заполняем информацию, которая будет использоваться при сравнении
			object1 := InformationAboutObject{
				Template: daemonSets1.Items[index1.index].Spec.Template,
				Selector: daemonSets1.Items[index1.index].Spec.Selector,
			}
			object2 := InformationAboutObject{
				Template: daemonSets2.Items[index2.index].Spec.Template,
				Selector: daemonSets2.Items[index2.index].Spec.Selector,
			}
			//CompareContainers(deployment1.Items[index1.index].Spec, deployment2.Items[index2.index].Spec, namespace)
			CompareContainers(object1, object2, namespace, client1, client2)

			fmt.Printf("----- End checking daemonset: '%s' -----\n\n", name)
		} else {
			fmt.Printf("DaemonSet '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("DaemonSet '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
}

func SetInformationAboutStatefulSets(map1 map[string]CheckerFlag, map2 map[string]CheckerFlag, statefulSets1 *v1.StatefulSetList, statefulSets2 *v1.StatefulSetList, namespace string) {
	if len(map1) != len(map2) {
		fmt.Printf("!!!The statefulsets count are different!!!\n\n")
	}
	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok == true {
			index1.check = true
			map1[name] = index1
			index2.check = true
			map2[name] = index2
			fmt.Printf("----- Start checking statefulset: '%s' -----\n", name)
			if *statefulSets1.Items[index1.index].Spec.Replicas != *statefulSets2.Items[index2.index].Spec.Replicas {
				fmt.Printf("!!!The replicas count are different!!!\n%s '%s' replicas: %d\n%s '%s' replicas: %d\n", kubeconfig1YamlStruct.Clusters[0].Cluster.Server, statefulSets1.Items[index1.index].Name, *statefulSets1.Items[index1.index].Spec.Replicas, kubeconfig2YamlStruct.Clusters[0].Cluster.Server, statefulSets2.Items[index2.index].Name, *statefulSets2.Items[index2.index].Spec.Replicas)
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

				//CompareContainers(deployment1.Items[index1.index].Spec, deployment2.Items[index2.index].Spec, namespace)
				CompareContainers(object1, object2, namespace, client1, client2)
			}
			fmt.Printf("----- End checking statefulset: '%s' -----\n\n", name)
		} else {
			fmt.Printf("StatefulSet '%s' - 1 cluster. Does not exist on another cluster\n\n", name)
		}
	}
	for name, index := range map2 {
		if index.check == false {
			fmt.Printf("StatefulSet '%s' - 2 cluster. Does not exist on another cluster\n\n", name)
		}
	}
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
				CompareContainers(object1, object2, namespace, client1, client2)
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

//func CreateFile(object1, object2 InformationAboutObject, job string) error {
//	objectJsonByte1, err := json.Marshal(object1)
//	if err != nil {
//		return err
//	}
//	objectJsonByte2, err := json.Marshal(object2)
//	if err != nil {
//		return err
//	}
//	filename := fmt.Sprintf("object'%s'.txt", job)
//	err = ioutil.WriteFile(filename, objectJsonByte1,0777)
//	if err != nil {
//		return err
//	}
//	err = ioutil.WriteFile("file.txt", objectJsonByte2,0777)
//	if err != nil {
//		return err
//	}
//	return nil
//	//data, _ := ioutil.ReadFile("file.txt")
//}

func CompareContainers(deploymentSpec1 InformationAboutObject, deploymentSpec2 InformationAboutObject, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) error {
	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		fmt.Printf("!!!The number templates of containers differs!!!\n")
		return errors.New("The number templates of containers differs")
	} else {
		matchLabelsString1 := ConvertMatchlabelsToString(deploymentSpec1.Selector.MatchLabels)
		matchLabelsString2 := ConvertMatchlabelsToString(deploymentSpec2.Selector.MatchLabels)
		if matchLabelsString1 != matchLabelsString2 {
			fmt.Printf("!!!MatchLabels are not equal!!!\n")
			return errors.New("MatchLabels are not equal")
		}
		pods1, pods2 := GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
		for i := 0; i < len(containersDeploymentTemplate1); i++ {
			if containersDeploymentTemplate1[i].Name != containersDeploymentTemplate2[i].Name {
				fmt.Printf("!!!Container names are not equal!!!\n")
				return errors.New("Container names are not equal")
			} else if containersDeploymentTemplate1[i].Image != containersDeploymentTemplate2[i].Image {
				fmt.Printf("!!!Container name images are not equal!!!\n")
				return errors.New("Container name images are not equal")
			} else {
				if len(pods1.Items) != len(pods2.Items) {
					fmt.Printf("!!!The pods count are different!!!\n")
					return errors.New("The pods count are different")
				} else {
					for j := 0; j < len(pods1.Items); j++ {
						containersStatusesInPod1 := GetContainerStatusesInPod(pods1.Items[j].Status.ContainerStatuses)
						containersStatusesInPod2 := GetContainerStatusesInPod(pods2.Items[j].Status.ContainerStatuses)
						if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
							fmt.Printf("!!!The containers count in pod are different!!!\n")
							return errors.New("The containers count in pod are different")
						} else {
							var flag int
							var containerWithSameNameFound bool
							for f := 0; f < len(containersStatusesInPod1); f++ {
								if containersDeploymentTemplate1[i].Name == containersStatusesInPod1[f].name {
									flag++
									if containersDeploymentTemplate1[i].Image != containersStatusesInPod1[f].image {
										fmt.Printf("!!!The container image in the template does not match the actual image in the Pod!!!\n")
										return errors.New("The container image in the template does not match the actual image in the Pod")
									} else {
										for _, value := range containersStatusesInPod2 {
											if containersStatusesInPod1[f].name == value.name {
												containerWithSameNameFound = true
												if containersStatusesInPod1[f].image != value.image {
													textForError := fmt.Sprintf("!!!The Image in Pods is different!!!\nPods name: '%s'\nImage name on pod1: '%s'\nImage name on pod2: '%s'\n\n", value.name, containersStatusesInPod1[j].image, value.image)
													fmt.Printf(textForError)
													return errors.New(textForError)
												} else if containersStatusesInPod1[f].imageID != value.imageID {
													textForError := fmt.Sprintf("!!!The ImageID in Pods is different!!!\nPods name: '%s'\nImageID on pod1: '%s'\nImageID on pod2: '%s'\n\n", value.name, containersStatusesInPod1[j].imageID, value.imageID)
													fmt.Printf(textForError)
													return errors.New(textForError)
												}
											}
										}
										if !containerWithSameNameFound {
											textForError := fmt.Sprintf("Container '%s' not found on other pod", containersStatusesInPod1[j].name)
											fmt.Printf("!!!Container '%s' not found on other pod!!!\n", containersStatusesInPod1[j].name)
											return errors.New(textForError)
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return nil
	}
}

func GetContainerStatusesInPod(containerStatuses []v12.ContainerStatus) map[int]Container {
	infoAboutContainer := make(map[int]Container)
	var container Container
	for index, value := range containerStatuses {
		container.name = value.Name
		container.image = value.Image
		container.imageID = value.ImageID
		infoAboutContainer[index] = container
	}
	return infoAboutContainer
}

//получает айдишник раскатанного образа на контейнерах
func GetPodsListOnMatchLabels(matchLabels map[string]string, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) (*v12.PodList, *v12.PodList) {
	matchLabelsString := ConvertMatchlabelsToString(matchLabels)
	pods1, err := clientSet1.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	pods2, err := clientSet2.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	//return pods1.Items[0].Status.ContainerStatuses[0].ImageID, pods2.Items[0].Status.ContainerStatuses[0].ImageID
	return pods1, pods2
}

func ConvertMatchlabelsToString(matchLabels map[string]string) string {
	values := []string{}
	for key, value := range matchLabels {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}
	//супермегафича склеивания строчек
	return strings.Join(values, ",")
}

func AddValueSecretsInMap(secrets1 *v12.SecretList, secrets2 *v12.SecretList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapSecrets1 := make(map[string]CheckerFlag)
	mapSecrets2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range secrets1.Items {
		indexCheck.index = index
		mapSecrets1[value.Name] = indexCheck
	}
	for index, value := range secrets2.Items {
		indexCheck.index = index
		mapSecrets2[value.Name] = indexCheck
	}
	return mapSecrets1, mapSecrets2
}

func AddValueConfigMapsInMap(configMaps1 *v12.ConfigMapList, configMaps2 *v12.ConfigMapList) (map[string]CheckerFlag, map[string]CheckerFlag) {
	mapConfigMap1 := make(map[string]CheckerFlag)
	mapConfigMap2 := make(map[string]CheckerFlag)
	var indexCheck CheckerFlag

	for index, value := range configMaps1.Items {
		indexCheck.index = index
		mapConfigMap1[value.Name] = indexCheck
	}
	for index, value := range configMaps2.Items {
		indexCheck.index = index
		mapConfigMap2[value.Name] = indexCheck
	}
	return mapConfigMap1, mapConfigMap2
}

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

func AddValueStatefulSetsInMap(stateFulSets1 *v1.StatefulSetList, stateFulSets2 *v1.StatefulSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
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

func AddValueDaemonSetsMap(daemonSets1 *v1.DaemonSetList, daemonSets2 *v1.DaemonSetList) (map[string]CheckerFlag, map[string]CheckerFlag) {
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