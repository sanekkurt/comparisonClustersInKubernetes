package pod_controllers

import (
	"fmt"

	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

// CompareContainers main function for compare containers
func CompareContainers(deploymentSpec1, deploymentSpec2 types.InformationAboutObject, namespace string, clientSet1, clientSet2 kubernetes.Interface) error {
	logging.Log.Debug("Start checking containers")

	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		return ErrorDiffersTemplatesNumber
	}
	matchLabelsString1 := common.ConvertMatchLabelsToString(deploymentSpec1.Selector.MatchLabels)
	matchLabelsString2 := common.ConvertMatchLabelsToString(deploymentSpec2.Selector.MatchLabels)
	if matchLabelsString1 != matchLabelsString2 {
		return ErrorMatchlabelsNotEqual
	}
	pods1, pods2 := common.GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
	for i := 0; i < len(containersDeploymentTemplate1); i++ {
		if containersDeploymentTemplate1[i].Name != containersDeploymentTemplate2[i].Name {
			return ErrorContainerNamesTemplate
		}
		if containersDeploymentTemplate1[i].Image != containersDeploymentTemplate2[i].Image {
			return ErrorContainerImagesTemplate
		}
		if err := CompareEnvInContainers(containersDeploymentTemplate1[i].Env, containersDeploymentTemplate2[i].Env, namespace, clientSet1, clientSet2); err != nil {
			return err
		}
		if len(pods1.Items) != len(pods2.Items) {
			return ErrorPodsCount
		}
		for j := 0; j < len(pods1.Items); j++ {
			containersStatusesInPod1 := GetContainerStatusesInPod(pods1.Items[j].Status.ContainerStatuses)
			containersStatusesInPod2 := GetContainerStatusesInPod(pods2.Items[j].Status.ContainerStatuses)
			if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
				return ErrorContainersCountInPod
			}
			var flag int
			var containerWithSameNameFound bool
			for f := 0; f < len(containersStatusesInPod1); f++ {
				if containersDeploymentTemplate1[i].Name == containersStatusesInPod1[f].Name && containersDeploymentTemplate1[i].Name == containersStatusesInPod2[f].Name { //nolint:gocritic,unused
					flag++
					if containersDeploymentTemplate1[i].Image != containersStatusesInPod1[f].Image || containersDeploymentTemplate1[i].Image != containersStatusesInPod2[f].Image {
						return ErrorContainerImageTemplatePod
					}
					for _, value := range containersStatusesInPod2 {
						if containersStatusesInPod1[f].Name == value.Name {
							containerWithSameNameFound = true
							if containersStatusesInPod1[f].Image != value.Image {
								return fmt.Errorf("%w. \nPods name: '%s'. Image name on pod1: '%s'. Image name on pod2: '%s'", ErrorDifferentImageInPods, value.Name, containersStatusesInPod1[j].Image, value.Image)
							}
							if containersStatusesInPod1[f].ImageID != value.ImageID {
								return fmt.Errorf("%w. Pods name: '%s'. ImageID on pod1: '%s'. ImageID on pod2: '%s'", ErrorDifferentImageIDInPods, value.Name, containersStatusesInPod1[j].ImageID, value.ImageID)
							}
						}
					}
					if !containerWithSameNameFound {
						return fmt.Errorf("%w. Name container: %s", ErrorContainerNotFound, containersStatusesInPod1[j].Name)
					}

				}
			}

		}

	}
	return nil
}

// GetContainerStatusesInPod get statuses containers in Pod
func GetContainerStatusesInPod(containerStatuses []v12.ContainerStatus) map[int]types.Container {
	infoAboutContainer := make(map[int]types.Container)
	var container types.Container
	for index, value := range containerStatuses {
		container.Name = value.Name
		container.Image = value.Image
		container.ImageID = value.ImageID
		infoAboutContainer[index] = container
	}
	return infoAboutContainer
}

// CompareEnvInContainers compare environment in containers
func CompareEnvInContainers(env1, env2 []v12.EnvVar, namespace string, clientSet1, clientSet2 kubernetes.Interface) error {
	if len(env1) != len(env2) {
		return ErrorNumberVariables
	}
	for i := 0; i < len(env1); i++ {
		if env1[i].ValueFrom != nil && env2[i].ValueFrom != nil {
			if env1[i].ValueFrom.ConfigMapKeyRef != nil && env2[i].ValueFrom.ConfigMapKeyRef != nil {
				if env1[i].ValueFrom.ConfigMapKeyRef.Key != env2[i].ValueFrom.ConfigMapKeyRef.Key || env1[i].ValueFrom.ConfigMapKeyRef.Name != env2[i].ValueFrom.ConfigMapKeyRef.Name {
					return fmt.Errorf("%w. Different ValueFrom: ValueFrom ConfigMapKeyRef in container 1 - %s:%s. ValueFrom ConfigMapKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[i].ValueFrom.ConfigMapKeyRef.Name, env1[i].ValueFrom.ConfigMapKeyRef.Key, env2[i].ValueFrom.ConfigMapKeyRef.Name, env2[i].ValueFrom.ConfigMapKeyRef.Key)
				}
				// logic check on configMapKey
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
				if env1[i].ValueFrom.SecretKeyRef.Key != env2[i].ValueFrom.SecretKeyRef.Key || env1[i].ValueFrom.SecretKeyRef.Name != env2[i].ValueFrom.SecretKeyRef.Name {
					return fmt.Errorf("%w. Different ValueFrom: ValueFrom SecretKeyRef in container 1 - %s:%s. ValueFrom SecretKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[i].ValueFrom.SecretKeyRef.Name, env1[i].ValueFrom.SecretKeyRef.Key, env2[i].ValueFrom.SecretKeyRef.Name, env2[i].ValueFrom.SecretKeyRef.Key)
				}
				// logic check on secretKey
				secret1, err := clientSet1.CoreV1().Secrets(namespace).Get(env1[i].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				secret2, err := clientSet2.CoreV1().Secrets(namespace).Get(env2[i].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					panic(err.Error())
				}
				if string(secret1.Data[env1[i].ValueFrom.SecretKeyRef.Key]) != string(secret2.Data[env2[i].ValueFrom.SecretKeyRef.Key]) {
					return fmt.Errorf("%w. Environment in container 1: SecretKeyRef.Key = %s, value = %s. Environment in container 2: SecretKeyRef.Key = %s, value = %s", ErrorDifferentValueSecretKey, env1[i].ValueFrom.SecretKeyRef.Key, string(secret1.Data[env1[i].ValueFrom.SecretKeyRef.Key]), env2[i].ValueFrom.SecretKeyRef.Key, string(secret2.Data[env2[i].ValueFrom.SecretKeyRef.Key]))
				}
			}

		} else if env1[i].ValueFrom != nil || env2[i].ValueFrom != nil {
			return fmt.Errorf("%w. Different ValueFrom: ValueFrom in container 1 - %s. ValueFrom in container 2 - %s", ErrorEnvironmentNotEqual, env1[i].ValueFrom, env2[i].ValueFrom)
		}
		if env1[i].Name != env2[i].Name || env1[i].Value != env2[i].Value {
			return fmt.Errorf("%w. Environment in container 1: name - '%s', value - '%s'. Environment in container 2: name - '%s', value - '%s'", ErrorEnvironmentNotEqual, env1[i].Name, env1[i].Value, env2[i].Name, env2[i].Value)
		}

	}
	return nil
}
