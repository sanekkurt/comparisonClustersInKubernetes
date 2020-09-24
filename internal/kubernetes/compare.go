package kubernetes

import (
	"fmt"
	v12 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
	"strings"
	"sync"
)

// CompareClusters main compare function, runs functions for comparing clusters by different parameters one at a time: Deployments, StatefulSets, DaemonSets, ConfigMaps
func CompareClusters(clientSet1, clientSet2 kubernetes.Interface, namespaces []string) (bool, error) {
	type ResStr struct {
		IsClusterDiffer bool
		Err             error
	}

	var (
		wg = &sync.WaitGroup{}

		resCh = make(chan ResStr, len(namespaces))
	)

	for _, namespace := range namespaces {
		wg.Add(1)

		go func(wg *sync.WaitGroup, resCh chan ResStr, namespace string) {
			var (
				isClustersDiffer bool
			)

			defer func() {
				wg.Done()
			}()

			depl1, err := clientSet1.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain deployments list from 1st cluster: %w", err),
				}
				return
			}

			depl2, err := clientSet2.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain deployments list from 2nd cluster: %w", err),
				}
				return
			}

			apc1List, map1, apc2List, map2 := PrepareDeploymentMaps(depl1, depl2)
			if ComparePodControllerSpecs(map1, map2, apc1List, apc2List, namespace) {
				isClustersDiffer = true
			}

			statefulSet1, err := clientSet1.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain statefulsets list from 1st cluster: %w", err),
				}
				return
			}

			statefulSet2, err := clientSet2.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain statefulsets list from 2nd cluster: %w", err),
				}
				return
			}

			apc1List, map1, apc2List, map2 = PrepareStatefulSetMaps(statefulSet1, statefulSet2)
			if ComparePodControllerSpecs(map1, map2, apc1List, apc2List, namespace) {
				isClustersDiffer = true
			}

			daemonSets1, err := clientSet1.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain daemonsets list from 1st cluster: %w", err),
				}
				return
			}
			daemonSets2, err := clientSet2.AppsV1().DaemonSets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain daemonsets list from 2nd cluster: %w", err),
				}
				return
			}

			apc1List, map1, apc2List, map2 = PrepareDaemonSetMaps(daemonSets1, daemonSets2)
			if ComparePodControllerSpecs(map1, map2, apc1List, apc2List, namespace) {
				isClustersDiffer = true
			}

			configMaps1, err := clientSet1.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain configmaps list from 1st cluster: %w", err),
				}
				return
			}

			configMaps2, err := clientSet2.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain configmaps list from 2nd cluster: %w", err),
				}
				return
			}

			mapConfigMaps1, mapConfigMaps2 := AddValueConfigMapsInMap(configMaps1, configMaps2)
			if SetInformationAboutConfigMaps(mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2) {
				isClustersDiffer = true
			}

			secrets1, err := clientSet1.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain secrets list from 1st cluster: %w", err),
				}
				return
			}

			secrets2, err := clientSet2.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain secrets list from 2nd cluster: %w", err),
				}
				return
			}

			mapSecrets1, mapSecrets2 := AddValueSecretsInMap(secrets1, secrets2)
			if SetInformationAboutSecrets(mapSecrets1, mapSecrets2, secrets1, secrets2) {
				isClustersDiffer = true
			}

			services1, err := clientSet1.CoreV1().Services(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain services list from 1st cluster: %w", err),
				}
				return
			}
			services2, err := clientSet2.CoreV1().Services(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain services list from 2nd cluster: %w", err),
				}
				return
			}
			mapServices1, mapServices2 := AddValueServicesInMap(services1, services2)
			if SetInformationAboutServices(mapServices1, mapServices2, services1, services2) {
				isClustersDiffer = true
			}

			ingresses1, err := clientSet1.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain ingresses list from 1st cluster: %w", err),
				}
				return
			}

			ingresses2, err := clientSet2.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
			if err != nil {
				resCh <- ResStr{
					IsClusterDiffer: false,
					Err:             fmt.Errorf("cannot obtain ingresses list from 2nd cluster: %w", err),
				}
				return
			}

			mapIngresses1, mapIngresses2 := AddValueIngressesInMap(ingresses1, ingresses2)
			SetInformationAboutIngresses(mapIngresses1, mapIngresses2, ingresses1, ingresses2)
			fmt.Println(mapIngresses1, mapIngresses2)

			resCh <- ResStr{
				IsClusterDiffer: isClustersDiffer,
				Err:             nil,
			}
		}(wg, resCh, namespace)
	}

	wg.Wait()

	close(resCh)

	for res := range resCh {
		if res.Err != nil {
			return false, res.Err
		}
		if res.IsClusterDiffer {
			return res.IsClusterDiffer, nil
		}
	}

	return false, nil
}

// CompareContainers main function for compare containers
func CompareContainers(deploymentSpec1, deploymentSpec2 InformationAboutObject, namespace string, clientSet1, clientSet2 kubernetes.Interface) error {
	containersDeploymentTemplate1 := deploymentSpec1.Template.Spec.Containers
	containersDeploymentTemplate2 := deploymentSpec2.Template.Spec.Containers
	if len(containersDeploymentTemplate1) != len(containersDeploymentTemplate2) {
		return ErrorDiffersTemplatesNumber
	}
	matchLabelsString1 := ConvertMatchLabelsToString(deploymentSpec1.Selector.MatchLabels)
	matchLabelsString2 := ConvertMatchLabelsToString(deploymentSpec2.Selector.MatchLabels)
	if matchLabelsString1 != matchLabelsString2 {
		return ErrorMatchlabelsNotEqual
	}
	pods1, pods2 := GetPodsListOnMatchLabels(deploymentSpec1.Selector.MatchLabels, namespace, clientSet1, clientSet2)
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
func GetContainerStatusesInPod(containerStatuses []v12.ContainerStatus) map[int]Container {
	infoAboutContainer := make(map[int]Container)
	var container Container
	for index, value := range containerStatuses {
		container.Name = value.Name
		container.Image = value.Image
		container.ImageID = value.ImageID
		infoAboutContainer[index] = container
	}
	return infoAboutContainer
}

// GetPodsListOnMatchLabels get pods list
func GetPodsListOnMatchLabels(matchLabels map[string]string, namespace string, clientSet1, clientSet2 kubernetes.Interface) (*v12.PodList, *v12.PodList) { //nolint:gocritic,unused
	matchLabelsString := ConvertMatchLabelsToString(matchLabels)
	pods1, err := clientSet1.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	pods2, err := clientSet2.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}
	return pods1, pods2
}

// ConvertMatchLabelsToString convert MatchLabels to string
func ConvertMatchLabelsToString(matchLabels map[string]string) string {
	keys := []string{}
	for key, _ := range matchLabels { //nolint
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := []string{}
	for i := 0; i < len(keys); i++ {
		values = append(values, fmt.Sprintf("%s=%s", keys[i], matchLabels[keys[i]]))
	}
	return strings.Join(values, ",")
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

// CompareSpecInServices compares spec in services
func CompareSpecInServices(service1, service2 v12.Service) error {
	if len(service1.Spec.Ports) != len(service2.Spec.Ports) {
		return fmt.Errorf("%w. Name service: '%s'. In first service - %d ports, in second service - '%d' ports", ErrorPortsCountDifferent, service1.Name, len(service1.Spec.Ports), len(service2.Spec.Ports))
	}
	for index, value := range service1.Spec.Ports {
		if value != service2.Spec.Ports[index] {
			return fmt.Errorf("%w. Name service: '%s'. First service: %s-%d-%s. Second service: %s-%d-%s", ErrorPortInServicesDifferent, service1.Name, value.Name, value.Port, value.Protocol, service2.Spec.Ports[index].Name, service2.Spec.Ports[index].Port, service2.Spec.Ports[index].Protocol)
		}
	}
	if len(service1.Spec.Selector) != len(service2.Spec.Selector) {
		return fmt.Errorf("%w. Name service: '%s'. In first service - %d selectors, in second service - '%d' selectors", ErrorSelectorsCountDifferent, service1.Name, len(service1.Spec.Selector), len(service2.Spec.Selector))
	}
	for key, value := range service1.Spec.Selector {
		if service2.Spec.Selector[key] != value {
			return fmt.Errorf("%w. Name service: '%s'. First service: %s-%s. Second service: %s-%s", ErrorSelectorInServicesDifferent, service1.Name, key, value, key, service2.Spec.Selector[key])
		}
	}
	if service1.Spec.Type != service2.Spec.Type {
		return fmt.Errorf("%w. Name service: '%s'. First service type: %s. Second service type: %s", ErrorTypeInServicesDifferent, service1.Name, service1.Spec.Type, service2.Spec.Type)
	}
	return nil
}


// CompareSpecInIngresses compare spec in the ingresses
func CompareSpecInIngresses(ingress1, ingress2 v1beta12.Ingress) error { //nolint
	if ingress1.Spec.TLS != nil && ingress2.Spec.TLS != nil {
		if len(ingress1.Spec.TLS) != len(ingress2.Spec.TLS) {
			return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d TLS. In second ingress - %d TLS", ErrorTLSCountDifferent, ingress1.Name, len(ingress1.Spec.TLS), len(ingress2.Spec.TLS) )
		}
		for index, value := range ingress1.Spec.TLS {
			if value.SecretName != ingress2.Spec.TLS[index].SecretName{
				return fmt.Errorf("%w. Name ingress: '%s'. First ingress: '%s'. Second ingress: '%s'", ErrorSecretNameInTLSDifferent, ingress1.Name, value.SecretName, ingress2.Spec.TLS[index].SecretName)
			}
			if value.Hosts != nil && ingress2.Spec.TLS[index].Hosts != nil {
				if len(value.Hosts) != len(ingress2.Spec.TLS[index].Hosts) {
					return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d hosts. In second ingress - %d hosts", ErrorHostsCountDifferent, ingress1.Name, len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
				}
				for i:=0; i<len(value.Hosts); i++ {
					if value.Hosts[i] != ingress2.Spec.TLS[index].Hosts[i] {
						return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorNameHostDifferent, ingress1.Name, value.Hosts[i], ingress2.Spec.TLS[index].Hosts[i])
					}
				}
			} else if value.Hosts != nil || ingress2.Spec.TLS[index].Hosts != nil {
				return fmt.Errorf("%w", ErrorHostsInIngressesDifferent)
			}
		}
	} else if ingress1.Spec.TLS != nil || ingress2.Spec.TLS != nil {
		return fmt.Errorf("%w", ErrorTLSInIngressesDifferent)
	}
	if ingress1.Spec.Backend != nil && ingress2.Spec.Backend != nil {
		err := CompareIngressesBackend(*ingress1.Spec.Backend, *ingress2.Spec.Backend, ingress1.Name)
		if err != nil {
			return err
		}
	} else if ingress1.Spec.Backend != nil || ingress2.Spec.Backend != nil {
		return fmt.Errorf("%w", ErrorBackendInIngressesDifferent)
	}
	if ingress1.Spec.Rules != nil && ingress2.Spec.Rules != nil {
		if len(ingress1.Spec.Rules) != len(ingress2.Spec.Rules){
			return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' rules. In second ingress - '%d' rules", ErrorRulesCountDifferent, ingress1.Name, len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
		}
		for index, value := range ingress1.Spec.Rules {
			if value.Host != ingress2.Spec.Rules[index].Host {
				return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorHostNameInRuleDifferent, ingress1.Name, value.Host, ingress2.Spec.Rules[index].Host)
			}
			if value.HTTP != nil && ingress2.Spec.Rules[index].HTTP != nil {
				err := CompareIngressesHTTP(*value.HTTP, *ingress2.Spec.Rules[index].HTTP, ingress1.Name)
				if err != nil {
					return err
				}
			} else if value.HTTP != nil || ingress2.Spec.Rules[index].HTTP != nil {
				return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
			}
			if value.IngressRuleValue.HTTP != nil && ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
				err := CompareIngressesHTTP(*value.IngressRuleValue.HTTP, *ingress2.Spec.Rules[index].IngressRuleValue.HTTP, ingress1.Name)
				if err != nil {
					return err
				}
			} else if value.IngressRuleValue.HTTP != nil || ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
				return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
			}

		}
	} else if ingress1.Spec.Rules != nil || ingress2.Spec.Rules != nil {
		return fmt.Errorf("%w",ErrorRulesInIngressesDifferent)
	}
	return nil
}

// CompareIngressesBackend compare backend in ingresses
func CompareIngressesBackend(backend1, backend2 v1beta12.IngressBackend, name string) error {
	if backend1.ServiceName != backend2.ServiceName {
		return fmt.Errorf("%w. Name ingress: '%s'. Service name in first ingress: '%s'. Service name in second ingress: '%s'", ErrorServiceNameInBackendDifferent, name, backend1.ServiceName, backend2.ServiceName )
	}
	if backend1.ServicePort.Type != backend2.ServicePort.Type || backend1.ServicePort.IntVal != backend2.ServicePort.IntVal || backend1.ServicePort.StrVal != backend2.ServicePort.StrVal {
		return fmt.Errorf("%w. Name ingress: '%s'", ErrorBackendServicePortDifferent, name )
	}
	return nil
}

// CompareIngressesHTTP compare http in ingresses
func CompareIngressesHTTP(http1, http2 v1beta12.HTTPIngressRuleValue, name string) error {
	if len(http1.Paths) != len(http2.Paths) {
		return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' paths. In second ingress - '%d' paths", ErrorPathsCountDifferent, name, len(http1.Paths), len(http2.Paths))
	}
	for i:=0; i < len(http1.Paths); i++ {
		if http1.Paths[i].Path != http2.Paths[i].Path {
			return fmt.Errorf("%w. Name ingress: '%s'. Name path in first ingress - '%s'. Name path in second ingress - '%s'", ErrorPathValueDifferent, name, http1.Paths[i].Path, http2.Paths[i].Path)
		}
		err := CompareIngressesBackend(http1.Paths[i].Backend, http2.Paths[i].Backend, name)
		if err != nil {
			return err
		}
	}
	return nil
}