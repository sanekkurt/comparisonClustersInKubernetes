package pod_controllers

import (
	"errors"
	"fmt"
	"strings"

	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

// CompareContainers main function for compare containers
func CompareContainers(deploymentSpec1, deploymentSpec2 types.InformationAboutObject, namespace string, simplifiedVerification, switchFatalDifferentTag bool, clientSet1, clientSet2 kubernetes.Interface) error {
	log.Debug("Start checking containers")

	var (
		containersDeploymentTemplate1 = deploymentSpec1.Template.Spec.Containers
		containersDeploymentTemplate2 = deploymentSpec2.Template.Spec.Containers

		pods1 *v12.PodList
		pods2 *v12.PodList

		matchLabelsString1 string
		matchLabelsString2 string
	)

	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		return ErrorDiffersTemplatesNumber
	}

	if !simplifiedVerification {
		matchLabelsString1 = common.ConvertMatchLabelsToString(deploymentSpec1.Selector.MatchLabels)
		matchLabelsString2 = common.ConvertMatchLabelsToString(deploymentSpec2.Selector.MatchLabels)

		if matchLabelsString1 != matchLabelsString2 {
			return ErrorMatchlabelsNotEqual
		}

		pods1, pods2 = common.GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
	}

	for podTemplate1ContainerIdx := 0; podTemplate1ContainerIdx < len(containersDeploymentTemplate1); podTemplate1ContainerIdx++ {

		if containersDeploymentTemplate1[podTemplate1ContainerIdx].Name != containersDeploymentTemplate2[podTemplate1ContainerIdx].Name {
			return ErrorContainerNamesTemplate
		}

		if containersDeploymentTemplate1[podTemplate1ContainerIdx].Image != containersDeploymentTemplate2[podTemplate1ContainerIdx].Image {
			return ErrorContainerImagesTemplate
		}

		if err := CompareEnvInContainers(containersDeploymentTemplate1[podTemplate1ContainerIdx].Env, containersDeploymentTemplate2[podTemplate1ContainerIdx].Env, namespace, simplifiedVerification, clientSet1, clientSet2); err != nil {
			return err
		}

		if err := CompareMassStringsInContainers(containersDeploymentTemplate1[podTemplate1ContainerIdx].Command, containersDeploymentTemplate2[podTemplate1ContainerIdx].Command); err != nil {
			return fmt.Errorf("%w. Name container: %s. %s", ErrorContainerCommandsDifferent, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name, err)
		}

		if err := CompareMassStringsInContainers(containersDeploymentTemplate1[podTemplate1ContainerIdx].Args, containersDeploymentTemplate2[podTemplate1ContainerIdx].Args); err != nil {
			return fmt.Errorf("%w. Name container: %s. %s", ErrorContainerArgumentsDifferent, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name, err)
		}

		if containersDeploymentTemplate1[podTemplate1ContainerIdx].LivenessProbe != nil && containersDeploymentTemplate2[podTemplate1ContainerIdx].LivenessProbe != nil {
			if err := CompareProbeInContainers(*containersDeploymentTemplate1[podTemplate1ContainerIdx].LivenessProbe, *containersDeploymentTemplate1[podTemplate1ContainerIdx].LivenessProbe, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name, ErrorContainerLivenessProbeDifferent); err != nil {
				return err
			}
		} else if containersDeploymentTemplate1[podTemplate1ContainerIdx].LivenessProbe != nil || containersDeploymentTemplate2[podTemplate1ContainerIdx].LivenessProbe != nil {
			return fmt.Errorf("%w. Name container: %s. One of the containers is missing Liveness probe", ErrorContainerLivenessProbeDifferent, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name)
		}

		if containersDeploymentTemplate1[podTemplate1ContainerIdx].ReadinessProbe != nil && containersDeploymentTemplate2[podTemplate1ContainerIdx].ReadinessProbe != nil {
			if err := CompareProbeInContainers(*containersDeploymentTemplate1[podTemplate1ContainerIdx].ReadinessProbe, *containersDeploymentTemplate1[podTemplate1ContainerIdx].ReadinessProbe, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name, ErrorContainerReadinessProbeDifferent); err != nil {
				return err
			}
		} else if containersDeploymentTemplate1[podTemplate1ContainerIdx].ReadinessProbe != nil || containersDeploymentTemplate2[podTemplate1ContainerIdx].ReadinessProbe != nil {
			return fmt.Errorf("%w. Name container: %s. One of the containers is missing Readiness probe", ErrorContainerReadinessProbeDifferent, containersDeploymentTemplate1[podTemplate1ContainerIdx].Name)
		}

		if !simplifiedVerification {

			if len(pods1.Items) != len(pods2.Items) {
				return ErrorPodsCount
			}

			for controlledPod1Idx := range pods1.Items {
				var (
					flag                       int
					containerWithSameNameFound bool
					templateHasAbsolutePath    bool

					containersStatusesInPod1 = GetContainerStatusesInPod(pods1.Items[controlledPod1Idx].Status.ContainerStatuses)
					containersStatusesInPod2 = GetContainerStatusesInPod(pods2.Items[controlledPod1Idx].Status.ContainerStatuses)

					containersDeploymentTemplateSplitLabel = strings.Split(containersDeploymentTemplate1[podTemplate1ContainerIdx].Image, ":")
				)

				if strings.Contains(containersDeploymentTemplate1[podTemplate1ContainerIdx].Image, "/") {
					templateHasAbsolutePath = true
				}

				if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
					return ErrorContainersCountInPod
				}

				for controlledPod1ContainerStatusIdx := range containersStatusesInPod1 {
					if containersDeploymentTemplate1[podTemplate1ContainerIdx].Name == containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name && containersDeploymentTemplate1[podTemplate1ContainerIdx].Name == containersStatusesInPod2[controlledPod1ContainerStatusIdx].Name { //nolint:gocritic,unused

						flag++

						var containersStatusesInPod1SplitLabel []string
						var containersStatusesInPod2SplitLabel []string

						if templateHasAbsolutePath {
							containersStatusesInPod1SplitLabel = strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, ":")
							containersStatusesInPod2SplitLabel = strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, ":")
						} else {
							if strings.Contains(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, "/") || strings.Contains(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, "/"){
								pathImage1 := strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, "/")
								pathImage2 := strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, "/")
								containersStatusesInPod1SplitLabel = strings.Split(pathImage1[len(pathImage1)-1], ":")
								containersStatusesInPod2SplitLabel = strings.Split(pathImage2[len(pathImage2)-1], ":")
							} else {
								containersStatusesInPod1SplitLabel = strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, ":")
								containersStatusesInPod2SplitLabel = strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, ":")
							}
						}

						if containersDeploymentTemplateSplitLabel[0] != containersStatusesInPod1SplitLabel[0] || containersDeploymentTemplateSplitLabel[0] != containersStatusesInPod2SplitLabel[0] { //nolint:gocritic,unused
							return ErrorContainerImageTemplatePod
						}

						if len(containersDeploymentTemplateSplitLabel) > 1 {
							if containersDeploymentTemplateSplitLabel[1] != containersStatusesInPod1SplitLabel[1] || containersDeploymentTemplateSplitLabel[1] != containersStatusesInPod2SplitLabel[1] {
								log.Infof("the container image tag in the template does not match the actual image tag in the pod: template image tag - %s, pod1 image tag - %s, pod2 image tag - %s", containersDeploymentTemplateSplitLabel[1], containersStatusesInPod1SplitLabel[1], containersStatusesInPod2SplitLabel[1]) // !!!!!

								if switchFatalDifferentTag {
									return ErrorContainerImageTagTemplatePod
								}
							}
						}

						for _, value := range containersStatusesInPod2 {
							if containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name == value.Name {
								containerWithSameNameFound = true
								if containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image != value.Image {
									return fmt.Errorf("%w. \nPods name: '%s'. Image name on pod1: '%s'. Image name on pod2: '%s'", ErrorDifferentImageInPods, value.Name, containersStatusesInPod1[controlledPod1Idx].Image, value.Image)
								}
								if containersStatusesInPod1[controlledPod1ContainerStatusIdx].ImageID != value.ImageID {
									return fmt.Errorf("%w. Pods name: '%s'. ImageID on pod1: '%s'. ImageID on pod2: '%s'", ErrorDifferentImageIDInPods, value.Name, containersStatusesInPod1[controlledPod1Idx].ImageID, value.ImageID)
								}
							}
						}
						if !containerWithSameNameFound {
							return fmt.Errorf("%w. Name container: %s", ErrorContainerNotFound, containersStatusesInPod1[controlledPod1Idx].Name)
						}
					}
				}
			}
		}

	}

	log.Debug("Stop checking containers")

	return nil
}

// GetContainerStatusesInPod get statuses containers in Pod
func GetContainerStatusesInPod(containerStatuses []v12.ContainerStatus) map[int]types.Container {
	var (
		container types.Container

		infoAboutContainer = make(map[int]types.Container)
	)

	for index, value := range containerStatuses {
		container.Name = value.Name
		container.Image = value.Image
		container.ImageID = value.ImageID
		infoAboutContainer[index] = container
	}

	return infoAboutContainer
}

// CompareEnvInContainers compare environment in containers
func CompareEnvInContainers(env1, env2 []v12.EnvVar, namespace string, simplifiedVerification bool, clientSet1, clientSet2 kubernetes.Interface) error {
	log.Debug("Start compare environments in containers")
	if len(env1) != len(env2) {
		return ErrorNumberVariables
	}

	for pod1EnvIdx := range env1 {
		if !simplifiedVerification {
			if env1[pod1EnvIdx].ValueFrom != nil && env2[pod1EnvIdx].ValueFrom != nil {
				if env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef != nil && env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef != nil {
					if env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key != env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key || env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name != env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name {
						return fmt.Errorf("%w. Different ValueFrom: ValueFrom ConfigMapKeyRef in container 1 - %s:%s. ValueFrom ConfigMapKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key)
					}

					// logic check on configMapKey
					configMap1, err := clientSet1.CoreV1().ConfigMaps(namespace).Get(env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						panic(err.Error())
					}

					configMap2, err := clientSet2.CoreV1().ConfigMaps(namespace).Get(env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						panic(err.Error())
					}

					if configMap1.Data[env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key] != configMap2.Data[env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key] {
						return fmt.Errorf("%w. Environment in container 1: ConfigMapKeyRef.Key = %s, value = %s. Environment in container 2: ConfigMapKeyRef.Key = %s, value = %s", ErrorDifferentValueConfigMapKey, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, configMap1.Data[env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key], env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, configMap2.Data[env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key])
					}
				} else if env1[pod1EnvIdx].ValueFrom.SecretKeyRef != nil && env2[pod1EnvIdx].ValueFrom.SecretKeyRef != nil {
					if env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key != env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key || env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name != env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name {
						return fmt.Errorf("%w. Different ValueFrom: ValueFrom SecretKeyRef in container 1 - %s:%s. ValueFrom SecretKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key)
					}
					// logic check on secretKey
					secret1, err := clientSet1.CoreV1().Secrets(namespace).Get(env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						panic(err.Error())
					}
					secret2, err := clientSet2.CoreV1().Secrets(namespace).Get(env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
					if err != nil {
						panic(err.Error())
					}
					if string(secret1.Data[env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]) != string(secret2.Data[env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]) {
						return fmt.Errorf("%w. Environment in container 1: SecretKeyRef.Key = %s, value = %s. Environment in container 2: SecretKeyRef.Key = %s, value = %s", ErrorDifferentValueSecretKey, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, string(secret1.Data[env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]), env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, string(secret2.Data[env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]))
					}
				}
			} else if env1[pod1EnvIdx].ValueFrom != nil || env2[pod1EnvIdx].ValueFrom != nil {
				return fmt.Errorf("%w. Different ValueFrom: ValueFrom in container 1 - %s. ValueFrom in container 2 - %s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom, env2[pod1EnvIdx].ValueFrom)
			}
		}

		if env1[pod1EnvIdx].Name != env2[pod1EnvIdx].Name || env1[pod1EnvIdx].Value != env2[pod1EnvIdx].Value {
			return fmt.Errorf("%w. Environment in container 1: name - '%s', value - '%s'. Environment in container 2: name - '%s', value - '%s'", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].Name, env1[pod1EnvIdx].Value, env2[pod1EnvIdx].Name, env2[pod1EnvIdx].Value)
		}
	}
	return nil
}

// CompareCommandsOrArgsInContainer compares commands or args in containers
func CompareMassStringsInContainers(mass1, mass2 []string) error {
	for index, value := range mass1 {
		if value != mass2[index] {
			return errors.New(fmt.Sprintf("value in container 1  - %s, value in container 2 - %s", value, mass2[index]))
		}
	}
	return nil
}

func CompareProbeInContainers(probe1, probe2 v12.Probe, nameContainer string, er error) error {

	err := CompareMassStringsInContainers(probe1.Exec.Command, probe2.Exec.Command)
	if err != nil {
		return fmt.Errorf("%w. Containers name: %s. Different exec command: %s", er, nameContainer, err)
	}

	if probe1.TCPSocket.Host != probe2.TCPSocket.Host {
		return fmt.Errorf("%w. Containers name: %s. Different TCPSocket.Host: container 1 host - %s, container 2 host - %s", er, nameContainer, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
	}

	if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal || probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal || probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
		return fmt.Errorf("%w. Containers name: %s. Different TCPSocket.Port: container 1 port - %s, container 2 port - %s", er, nameContainer, fmt.Sprintln(probe1.TCPSocket.Port), fmt.Sprintln(probe2.TCPSocket.Port))
	}

	if probe1.HTTPGet.Host != probe2.HTTPGet.Host {
		return fmt.Errorf("%w. Containers name: %s. Different HTTPGet.Host: container 1 host - %s, container 2 host - %s", er, nameContainer, probe1.HTTPGet.Host, probe2.HTTPGet.Host)
	}

	if probe1.HTTPGet.HTTPHeaders != nil && probe2.HTTPGet.HTTPHeaders != nil {
		if len(probe1.HTTPGet.HTTPHeaders) != len(probe2.HTTPGet.HTTPHeaders) {
			return fmt.Errorf("%w. Containers name: %s. Different count HTTPGet.HTTPHeaders: container 1 count - %d, container 2 count - %d", er, nameContainer, len(probe1.HTTPGet.HTTPHeaders), len(probe2.HTTPGet.HTTPHeaders))
		}

		for index, value := range probe1.HTTPGet.HTTPHeaders {
			if value.Name != probe2.HTTPGet.HTTPHeaders[index].Name {
				return fmt.Errorf("%w. Containers name: %s. Different name header: container 1 header name - %s, container 2 header name - %s", er, nameContainer, value.Name, probe2.HTTPGet.HTTPHeaders[index].Name)
			}

			if value.Value != probe2.HTTPGet.HTTPHeaders[index].Value {
				return fmt.Errorf("%w. Containers name: %s. Different value header. Name header - %s. Container 1 header value - %s, container 2 header value - %s", er, nameContainer, value.Name, value.Value ,probe2.HTTPGet.HTTPHeaders[index].Value)
			}
		}

	} else if probe1.HTTPGet.HTTPHeaders != nil || probe2.HTTPGet.HTTPHeaders != nil {
		return fmt.Errorf("%w. Containers name: %s. One of the containers is missing headers", er, nameContainer)
	}

	if probe1.HTTPGet.Path != probe2.HTTPGet.Path {
		return fmt.Errorf("%w. Containers name: %s. Different HTTPGet.Path: container 1 path - %s, container 2 path - %s", er, nameContainer, probe1.HTTPGet.Path, probe2.HTTPGet.Path)
	}

	if probe1.HTTPGet.Port.IntVal != probe2.HTTPGet.Port.IntVal || probe1.HTTPGet.Port.StrVal != probe2.HTTPGet.Port.StrVal || probe1.HTTPGet.Port.Type != probe2.HTTPGet.Port.Type {
		return fmt.Errorf("%w. Containers name: %s. Different HTTPGet.Port: container 1 port - %s, container 2 port - %s", er, nameContainer, fmt.Sprintln(probe1.HTTPGet.Port), fmt.Sprintln(probe2.HTTPGet.Port))
	}

	if probe1.FailureThreshold != probe2.FailureThreshold {
		return fmt.Errorf("%w. Containers name: %s. Different FailureThreshold: container 1 - %d, container 2 - %d", er, nameContainer, probe1.FailureThreshold, probe2.FailureThreshold)
	}

	if probe1.InitialDelaySeconds != probe2.InitialDelaySeconds {
		return fmt.Errorf("%w. Containers name: %s. Different InitialDelaySeconds: container 1 - %d, container 2 - %d", er, nameContainer, probe1.InitialDelaySeconds, probe2.InitialDelaySeconds)
	}

	if probe1.PeriodSeconds != probe2.PeriodSeconds {
		return fmt.Errorf("%w. Containers name: %s. Different PeriodSeconds: container 1 - %d, container 2 - %d", er, nameContainer, probe1.PeriodSeconds, probe2.PeriodSeconds)
	}

	if probe1.SuccessThreshold != probe2.SuccessThreshold {
		return fmt.Errorf("%w. Containers name: %s. Different SuccessThreshold: container 1 - %d, container 2 - %d", er, nameContainer, probe1.SuccessThreshold, probe2.SuccessThreshold)
	}

	if probe1.TimeoutSeconds != probe2.TimeoutSeconds {
		return fmt.Errorf("%w. Containers name: %s. Different TimeoutSeconds: container 1 - %d, container 2 - %d", er, nameContainer, probe1.TimeoutSeconds, probe2.TimeoutSeconds)
	}

	return nil
}