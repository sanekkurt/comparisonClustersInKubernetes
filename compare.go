package main

import (
	"errors"
	"fmt"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

//основная сравнивающая функция, поочередно запускает функции для сравнения кластеров по разным параметрам: Deployments, StatefulSets, DaemonSets, ConfigMaps
func Compare(clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface, namespaces []string) {
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

func CompareContainers(deploymentSpec1 InformationAboutObject, deploymentSpec2 InformationAboutObject, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) error {
	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		fmt.Printf("!!!The number templates of containers differs!!!\n")
		return ErrorDiffersTemplatesNumber
	}
	matchLabelsString1 := ConvertMatchlabelsToString(deploymentSpec1.Selector.MatchLabels)
	matchLabelsString2 := ConvertMatchlabelsToString(deploymentSpec2.Selector.MatchLabels)
	if matchLabelsString1 != matchLabelsString2 {
		fmt.Printf("!!!MatchLabels are not equal!!!\n")
		return ErrorMatchlabelsNotEqual
	}
	pods1, pods2 := GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
	for i := 0; i < len(containersDeploymentTemplate1); i++ {
		if containersDeploymentTemplate1[i].Name != containersDeploymentTemplate2[i].Name {
			fmt.Printf("!!!Container names in template are not equal!!!\n")
			return ErrorContainerNamesTemplate
		}
		if containersDeploymentTemplate1[i].Image != containersDeploymentTemplate2[i].Image {
			fmt.Printf("!!!Container name images in template are not equal!!!\n")
			return ErrorContainerImagesTemplate
		}
		if err := CompareEnvInContainers(containersDeploymentTemplate1[i].Env, containersDeploymentTemplate2[i].Env, namespace, clientSet1, clientSet2); err != nil {
			return err
		}
		if len(pods1.Items) != len(pods2.Items) {
			fmt.Printf("!!!The pods count are different!!!\n")
			return ErrorPodsCount
		}
		for j := 0; j < len(pods1.Items); j++ {
			containersStatusesInPod1 := GetContainerStatusesInPod(pods1.Items[j].Status.ContainerStatuses)
			containersStatusesInPod2 := GetContainerStatusesInPod(pods2.Items[j].Status.ContainerStatuses)
			if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
				fmt.Printf("!!!The containers count in pod are different!!!\n")
				return ErrorContainersCountInPod
			}
			var flag int
			var containerWithSameNameFound bool
			for f := 0; f < len(containersStatusesInPod1); f++ {
				if containersDeploymentTemplate1[i].Name == containersStatusesInPod1[f].name && containersDeploymentTemplate1[i].Name == containersStatusesInPod2[f].name {
					flag++
					if containersDeploymentTemplate1[i].Image != containersStatusesInPod1[f].image || containersDeploymentTemplate1[i].Image != containersStatusesInPod2[f].image {
						fmt.Printf("!!!The container image in the template does not match the actual image in the Pod!!!\n")
						return ErrorContainerImageTemplatePod
					}
					for _, value := range containersStatusesInPod2 {
						if containersStatusesInPod1[f].name == value.name {
							containerWithSameNameFound = true
							if containersStatusesInPod1[f].image != value.image {
								textForError := fmt.Sprintf("!!!The Image in Pods is different!!!\nPods name: '%s'\nImage name on pod1: '%s'\nImage name on pod2: '%s'\n\n", value.name, containersStatusesInPod1[j].image, value.image)
								fmt.Printf(textForError)
								ErrorDifferentImageInPods = errors.New(textForError)
								return ErrorDifferentImageInPods
							}
							if containersStatusesInPod1[f].imageID != value.imageID {
								textForError := fmt.Sprintf("!!!The ImageID in Pods is different!!!\nPods name: '%s'\nImageID on pod1: '%s'\nImageID on pod2: '%s'\n\n", value.name, containersStatusesInPod1[j].imageID, value.imageID)
								fmt.Printf(textForError)
								ErrorDifferentImageIdInPods = errors.New(textForError)
								return ErrorDifferentImageIdInPods
							}
						}
					}
					if !containerWithSameNameFound {
						textForError := fmt.Sprintf("Container '%s' not found on other pod", containersStatusesInPod1[j].name)
						ErrorContainerNotFound = errors.New(textForError)
						fmt.Printf("!!!Container '%s' not found on other pod!!!\n", containersStatusesInPod1[j].name)
						return ErrorContainerNotFound
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

func CompareEnvInContainers(env1 []v12.EnvVar, env2 []v12.EnvVar, namespace string, clientSet1 kubernetes.Interface, clientSet2 kubernetes.Interface) error {
	if len(env1) != len(env2) {
		fmt.Printf("!!!The number of variables in containers differs!!!\n")
		return ErrorNumberVariables
	}
	for i := 0; i < len(env1); i++ {
		if *env1[i].ValueFrom == *env2[i].ValueFrom {
			if env1[i].ValueFrom != nil {
				if env1[i].ValueFrom.ConfigMapKeyRef != nil && env2[i].ValueFrom.ConfigMapKeyRef != nil {
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
						stringError := fmt.Sprintf("The value for the ConfigMapKey is different.\nEnvironment in container 1: ConfigMapKeyRef.Key = %s, value = %s\nEnvironment in container 2: ConfigMapKeyRef.Key = %s, value = %s\n", env1[i].ValueFrom.ConfigMapKeyRef.Key, configMap1.Data[env1[i].ValueFrom.ConfigMapKeyRef.Key], env2[i].ValueFrom.ConfigMapKeyRef.Key, configMap2.Data[env2[i].ValueFrom.ConfigMapKeyRef.Key])
						fmt.Println(stringError)
						ErrorDifferentValueConfigMapKey = errors.New(stringError)
						return ErrorDifferentValueConfigMapKey
					}
				} else if env1[i].ValueFrom.SecretKeyRef != nil && env2[i].ValueFrom.SecretKeyRef != nil {
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
						stringError := fmt.Sprintf("The value for the SecretKey is different.\nEnvironment in container 1: SecretKeyRef.Key = %s, value = %s\nEnvironment in container 2: SecretKeyRef.Key = %s, value = %s\n", env1[i].ValueFrom.SecretKeyRef.Key, string(secret1.Data[env1[i].ValueFrom.SecretKeyRef.Key]), env2[i].ValueFrom.SecretKeyRef.Key, string(secret2.Data[env2[i].ValueFrom.SecretKeyRef.Key]))
						fmt.Println(stringError)
						ErrorDifferentValueSecretKey = errors.New(stringError)
						return ErrorDifferentValueSecretKey
					}
				}
			} else {
				if env1[i].Name != env2[i].Name || env1[i].Value != env2[i].Value {
					stringError := fmt.Sprintf("Environment in container 1 not equal environment in container 2:\nEnvironment in container 1: name - '%s', value - '%s'\nEnvironment in container 2: name - '%s', value - '%s'\n", env1[i].Name, env1[i].Value, env2[i].Name, env2[i].Value)
					fmt.Println(stringError)
					ErrorEnvironmentNotEqual = errors.New(stringError)
					return ErrorEnvironmentNotEqual
				}
			}
		} else {
			stringError := fmt.Sprintf("Environment in container 1 not equal environment in container 2. Different ValueFrom:\nValueFrom in container 1 - %s\nValueFrom in container 2 - %s\n", env1[i].ValueFrom, env2[i].ValueFrom)
			fmt.Println(stringError)
			ErrorEnvironmentNotEqual = errors.New(stringError)
			return ErrorEnvironmentNotEqual
		}
	}
	return nil
}
