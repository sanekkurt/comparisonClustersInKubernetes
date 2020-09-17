package main

import (
	"fmt"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
	"strings"
)

//основная сравнивающая функция, поочередно запускает функции для сравнения кластеров по разным параметрам: Deployments, StatefulSets, DaemonSets, ConfigMaps
func Compare(clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface, namespaces []string) bool {
	var flag bool
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
		if SetInformationAboutDeployments(mapDeployments1, mapDeployments2, depl1, depl2, namespace) {
			flag = true
		}

		statefulSet1, err := clientSet1.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		statefulSet2, err := clientSet2.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapStatefulSets1, mapStatefulSets2 := AddValueStatefulSetsInMap(statefulSet1, statefulSet2)
		if SetInformationAboutStatefulSets(mapStatefulSets1, mapStatefulSets2, statefulSet1, statefulSet2, namespace) {
			flag = true
		}

		daemonSets1, err := clientSet1.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		daemonSets2, err := clientSet2.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapDaemonSets1, mapDaemonSets2 := AddValueDaemonSetsMap(daemonSets1, daemonSets2)
		if SetInformationAboutDaemonSets(mapDaemonSets1, mapDaemonSets2, daemonSets1, daemonSets2, namespace) {
			flag = true
		}

		configMaps1, err := clientSet1.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		configMaps2, err := clientSet2.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapConfigMaps1, mapConfigMaps2 := AddValueConfigMapsInMap(configMaps1, configMaps2)
		if SetInformationAboutConfigMaps(mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2) {
			flag = true
		}

		secrets1, err := clientSet1.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		secrets2, err := clientSet2.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		mapSecrets1, mapSecrets2 := AddValueSecretsInMap(secrets1, secrets2)
		if SetInformationAboutSecrets(mapSecrets1, mapSecrets2, secrets1, secrets2) {
			flag = true
		}
	}
	return flag
}

func CompareContainers(deploymentSpec1 InformationAboutObject, deploymentSpec2 InformationAboutObject, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) error {
	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		//fmt.Printf("!!!The number templates of containers differs!!!\n")
		return ErrorDiffersTemplatesNumber
	}
	matchLabelsString1 := ConvertMatchLabelsToString(deploymentSpec1.Selector.MatchLabels)
	matchLabelsString2 := ConvertMatchLabelsToString(deploymentSpec2.Selector.MatchLabels)
	if matchLabelsString1 != matchLabelsString2 {
		//fmt.Printf("!!!MatchLabels are not equal!!!\n")
		return ErrorMatchlabelsNotEqual
	}
	pods1, pods2 := GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
	for i := 0; i < len(containersDeploymentTemplate1); i++ {
		if containersDeploymentTemplate1[i].Name != containersDeploymentTemplate2[i].Name {
			//fmt.Printf("!!!Container names in template are not equal!!!\n")
			return ErrorContainerNamesTemplate
		}
		if containersDeploymentTemplate1[i].Image != containersDeploymentTemplate2[i].Image {
			//fmt.Printf("!!!Container name images in template are not equal!!!\n")
			return ErrorContainerImagesTemplate
		}
		if err := CompareEnvInContainers(containersDeploymentTemplate1[i].Env, containersDeploymentTemplate2[i].Env, namespace, clientSet1, clientSet2); err != nil {
			return err
		}
		if len(pods1.Items) != len(pods2.Items) {
			//fmt.Printf("!!!The pods count are different!!!\n")
			return ErrorPodsCount
		}
		for j := 0; j < len(pods1.Items); j++ {
			containersStatusesInPod1 := GetContainerStatusesInPod(pods1.Items[j].Status.ContainerStatuses)
			containersStatusesInPod2 := GetContainerStatusesInPod(pods2.Items[j].Status.ContainerStatuses)
			if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
				//fmt.Printf("!!!The containers count in pod are different!!!\n")
				return ErrorContainersCountInPod
			}
			var flag int
			var containerWithSameNameFound bool
			for f := 0; f < len(containersStatusesInPod1); f++ {
				if containersDeploymentTemplate1[i].Name == containersStatusesInPod1[f].name && containersDeploymentTemplate1[i].Name == containersStatusesInPod2[f].name {
					flag++
					if containersDeploymentTemplate1[i].Image != containersStatusesInPod1[f].image || containersDeploymentTemplate1[i].Image != containersStatusesInPod2[f].image {
						//fmt.Printf("!!!The container image in the template does not match the actual image in the Pod!!!\n")
						return ErrorContainerImageTemplatePod
					}
					for _, value := range containersStatusesInPod2 {
						if containersStatusesInPod1[f].name == value.name {
							containerWithSameNameFound = true
							if containersStatusesInPod1[f].image != value.image {
								return fmt.Errorf("%w. \nPods name: '%s'. Image name on pod1: '%s'. Image name on pod2: '%s'.",ErrorDifferentImageInPods, value.name, containersStatusesInPod1[j].image, value.image)
							}
							if containersStatusesInPod1[f].imageID != value.imageID {
								return fmt.Errorf("%w. Pods name: '%s'. ImageID on pod1: '%s'. ImageID on pod2: '%s'.", ErrorDifferentImageIdInPods, value.name, containersStatusesInPod1[j].imageID, value.imageID)
							}
						}
					}
					if !containerWithSameNameFound {
						return fmt.Errorf("%w. Name container: %s", ErrorContainerNotFound, containersStatusesInPod1[j].name)
					}

				}
			}

		}

	}
	return nil
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
	matchLabelsString := ConvertMatchLabelsToString(matchLabels)
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

func ConvertMatchLabelsToString(matchLabels map[string]string) string {
	keys := []string{}
	for key, _ := range matchLabels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := []string{}
	for i:=0; i<len(keys); i++ {
		values = append(values, fmt.Sprintf("%s=%s", keys[i], matchLabels[keys[i]]))
	}
	//for key, value := range matchLabels {
	//	values = append(values, fmt.Sprintf("%s=%s", key, value))
	//}
	//супермегафича склеивания строчек
	return strings.Join(values, ",")
}

func CompareEnvInContainers(env1 []v12.EnvVar, env2 []v12.EnvVar, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) error {
	if len(env1) != len(env2) {
		return ErrorNumberVariables
	}
	for i := 0; i < len(env1); i++ {
		if env1[i].ValueFrom != nil && env2[i].ValueFrom != nil {
			if env1[i].ValueFrom.ConfigMapKeyRef != nil && env2[i].ValueFrom.ConfigMapKeyRef != nil {
				if env1[i].ValueFrom.ConfigMapKeyRef.Key != env2[i].ValueFrom.ConfigMapKeyRef.Key || env1[i].ValueFrom.ConfigMapKeyRef.Name != env2[i].ValueFrom.ConfigMapKeyRef.Name{
					return fmt.Errorf("%w. Different ValueFrom: ValueFrom ConfigMapKeyRef in container 1 - %s:%s. ValueFrom ConfigMapKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[i].ValueFrom.ConfigMapKeyRef.Name, env1[i].ValueFrom.ConfigMapKeyRef.Key , env2[i].ValueFrom.ConfigMapKeyRef.Name, env2[i].ValueFrom.ConfigMapKeyRef.Key)
				}
				//ЛОГИКА ПРОВЕРКИ НА КОНФИГМАП КЕЙ
				configMap1, err := clientSet1.CoreV1().ConfigMaps(namespace).Get(env1[i].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				configMap2, err := clientSet2.CoreV1().ConfigMaps(namespace).Get(env2[i].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				if configMap1.Data[env1[i].ValueFrom.ConfigMapKeyRef.Key] != configMap2.Data[env2[i].ValueFrom.ConfigMapKeyRef.Key] {
					return fmt.Errorf("%w. Environment in container 1: ConfigMapKeyRef.Key = %s, value = %s. Environment in container 2: ConfigMapKeyRef.Key = %s, value = %s", ErrorDifferentValueConfigMapKey, env1[i].ValueFrom.ConfigMapKeyRef.Key, configMap1.Data[env1[i].ValueFrom.ConfigMapKeyRef.Key], env2[i].ValueFrom.ConfigMapKeyRef.Key, configMap2.Data[env2[i].ValueFrom.ConfigMapKeyRef.Key])
				}
			} else if env1[i].ValueFrom.SecretKeyRef != nil && env2[i].ValueFrom.SecretKeyRef != nil {
				if env1[i].ValueFrom.SecretKeyRef.Key != env2[i].ValueFrom.SecretKeyRef.Key || env1[i].ValueFrom.SecretKeyRef.Name != env2[i].ValueFrom.SecretKeyRef.Name{
					return fmt.Errorf("%w. Different ValueFrom: ValueFrom SecretKeyRef in container 1 - %s:%s. ValueFrom SecretKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[i].ValueFrom.SecretKeyRef.Name, env1[i].ValueFrom.SecretKeyRef.Key , env2[i].ValueFrom.SecretKeyRef.Name, env2[i].ValueFrom.SecretKeyRef.Key)
				}
				//ЛОГИКА ПРОВЕРКИ НА СЕКРЕТ КЕЙ
				secret1, err := clientSet1.CoreV1().Secrets(namespace).Get(env1[i].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				secret2, err := clientSet2.CoreV1().Secrets(namespace).Get(env2[i].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				if string(secret1.Data[env1[i].ValueFrom.SecretKeyRef.Key]) != string(secret2.Data[env2[i].ValueFrom.SecretKeyRef.Key]) {
					return fmt.Errorf("%w. Environment in container 1: SecretKeyRef.Key = %s, value = %s. Environment in container 2: SecretKeyRef.Key = %s, value = %s", ErrorDifferentValueSecretKey,env1[i].ValueFrom.SecretKeyRef.Key, string(secret1.Data[env1[i].ValueFrom.SecretKeyRef.Key]), env2[i].ValueFrom.SecretKeyRef.Key, string(secret2.Data[env2[i].ValueFrom.SecretKeyRef.Key]))
				}
			}

		} else if env1[i].ValueFrom != nil || env2[i].ValueFrom != nil {
			return fmt.Errorf("%w. Different ValueFrom: ValueFrom in container 1 - %s. ValueFrom in container 2 - %s",ErrorEnvironmentNotEqual, env1[i].ValueFrom, env2[i].ValueFrom)
		}
		if env1[i].Name != env2[i].Name || env1[i].Value != env2[i].Value {
			return fmt.Errorf("%w. Environment in container 1: name - '%s', value - '%s'. Environment in container 2: name - '%s', value - '%s'", ErrorEnvironmentNotEqual, env1[i].Name, env1[i].Value, env2[i].Name, env2[i].Value )
		}

	}
	return nil
}