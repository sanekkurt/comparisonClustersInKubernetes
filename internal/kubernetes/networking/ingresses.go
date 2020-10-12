package networking

import (
	"fmt"

	v1beta12 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"sync"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

// CompareIngresses compares list of ingresses objects in two given k8s-clusters
func CompareIngresses(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	ingresses1, err := clientSet1.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain ingresses list from 1st cluster: %w", err)
	}

	ingresses2, err := clientSet2.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain ingresses list from 2nd cluster: %w", err)
	}

	mapIngresses1, mapIngresses2 := prepareIngressMaps(ingresses1, ingresses2, skipEntityList.GetByKind("services"))

	isClustersDiffer = setInformationAboutIngresses(mapIngresses1, mapIngresses2, ingresses1, ingresses2)

	return isClustersDiffer, nil
}

// prepareIngressMaps add value secrets in map
func prepareIngressMaps(ingresses1, ingresses2 *v1beta12.IngressList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapIngresses1 := make(map[string]types.IsAlreadyComparedFlag)
	mapIngresses2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range ingresses1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("ingress %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapIngresses1[value.Name] = indexCheck

	}
	for index, value := range ingresses2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("ingress %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapIngresses2[value.Name] = indexCheck

	}
	return mapIngresses1, mapIngresses2
}

func compareIngressSpecInternals(wg *sync.WaitGroup, channel chan bool, name string, ing1, ing2 *v1beta12.Ingress) {
	var (
		flag bool
	)
	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking ingress: '%s' -----", name)

	if len(ing1.Labels) != len(ing2.Labels) {
		log.Infof("the number of labels is not equal in ingresses. Name ingress: '%s'. In first cluster %d, in second cluster %d", ing1.Name, len(ing1.Labels), len(ing2.Labels))
		flag = true
	} else {
		for key, value := range ing1.Labels {
			if ing2.Labels[key] != value {
				log.Infof("labels in ingresses don't match. Name ingress: '%s'. In first cluster: '%s'-'%s', in second cluster value = '%s'", ing1.Name, key, value, ing2.Labels[key])
				flag = true
			}
		}
	}

	err := compareSpecInIngresses(*ing1, *ing2)
	if err != nil {
		log.Infof("Ingress %s: %s", name, err.Error())
		flag = true
	}

	log.Debugf("----- End checking ingress: '%s' -----", name)
	channel <- flag
}

// setInformationAboutIngresses set information about services
func setInformationAboutIngresses(map1, map2 map[string]types.IsAlreadyComparedFlag, ingresses1, ingresses2 *v1beta12.IngressList) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		log.Infof("ingress counts are different")
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareIngressSpecInternals(wg, channel, name, &ingresses1.Items[index1.Index], &ingresses2.Items[index2.Index])
		} else {
			log.Infof("ingress '%s' does not exist in 2nd cluster", name)
			flag = true
			channel <- flag
		}
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

			log.Infof("ingress '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}

	return flag
}

// compareSpecInIngresses compare spec in the ingresses
func compareSpecInIngresses(ingress1, ingress2 v1beta12.Ingress) error { //nolint
	if ingress1.Spec.TLS != nil && ingress2.Spec.TLS != nil {
		if len(ingress1.Spec.TLS) != len(ingress2.Spec.TLS) {
			return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d TLS. In second ingress - %d TLS", ErrorTLSCountDifferent, ingress1.Name, len(ingress1.Spec.TLS), len(ingress2.Spec.TLS))
		}
		for index, value := range ingress1.Spec.TLS {
			if value.SecretName != ingress2.Spec.TLS[index].SecretName {
				return fmt.Errorf("%w. Name ingress: '%s'. First ingress: '%s'. Second ingress: '%s'", ErrorSecretNameInTLSDifferent, ingress1.Name, value.SecretName, ingress2.Spec.TLS[index].SecretName)
			}
			if value.Hosts != nil && ingress2.Spec.TLS[index].Hosts != nil {
				if len(value.Hosts) != len(ingress2.Spec.TLS[index].Hosts) {
					return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d hosts. In second ingress - %d hosts", ErrorHostsCountDifferent, ingress1.Name, len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
				}
				for i := 0; i < len(value.Hosts); i++ {
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
		err := compareIngressesBackend(*ingress1.Spec.Backend, *ingress2.Spec.Backend, ingress1.Name)
		if err != nil {
			return err
		}
	} else if ingress1.Spec.Backend != nil || ingress2.Spec.Backend != nil {
		return fmt.Errorf("%w", ErrorBackendInIngressesDifferent)
	}
	if ingress1.Spec.Rules != nil && ingress2.Spec.Rules != nil {
		if len(ingress1.Spec.Rules) != len(ingress2.Spec.Rules) {
			return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' rules. In second ingress - '%d' rules", ErrorRulesCountDifferent, ingress1.Name, len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
		}
		for index, value := range ingress1.Spec.Rules {
			if value.Host != ingress2.Spec.Rules[index].Host {
				return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorHostNameInRuleDifferent, ingress1.Name, value.Host, ingress2.Spec.Rules[index].Host)
			}
			if value.HTTP != nil && ingress2.Spec.Rules[index].HTTP != nil {
				err := compareIngressesHTTP(*value.HTTP, *ingress2.Spec.Rules[index].HTTP, ingress1.Name)
				if err != nil {
					return err
				}
			} else if value.HTTP != nil || ingress2.Spec.Rules[index].HTTP != nil {
				return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
			}
			if value.IngressRuleValue.HTTP != nil && ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
				err := compareIngressesHTTP(*value.IngressRuleValue.HTTP, *ingress2.Spec.Rules[index].IngressRuleValue.HTTP, ingress1.Name)
				if err != nil {
					return err
				}
			} else if value.IngressRuleValue.HTTP != nil || ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
				return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
			}

		}
	} else if ingress1.Spec.Rules != nil || ingress2.Spec.Rules != nil {
		return fmt.Errorf("%w", ErrorRulesInIngressesDifferent)
	}
	return nil
}

// compareIngressesBackend compare backend in ingresses
func compareIngressesBackend(backend1, backend2 v1beta12.IngressBackend, name string) error {
	if backend1.ServiceName != backend2.ServiceName {
		return fmt.Errorf("%w. Name ingress: '%s'. Service name in first ingress: '%s'. Service name in second ingress: '%s'", ErrorServiceNameInBackendDifferent, name, backend1.ServiceName, backend2.ServiceName)
	}
	if backend1.ServicePort.Type != backend2.ServicePort.Type || backend1.ServicePort.IntVal != backend2.ServicePort.IntVal || backend1.ServicePort.StrVal != backend2.ServicePort.StrVal {
		return fmt.Errorf("%w. Name ingress: '%s'", ErrorBackendServicePortDifferent, name)
	}
	return nil
}

// compareIngressesHTTP compare http in ingresses
func compareIngressesHTTP(http1, http2 v1beta12.HTTPIngressRuleValue, name string) error {
	if len(http1.Paths) != len(http2.Paths) {
		return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' paths. In second ingress - '%d' paths", ErrorPathsCountDifferent, name, len(http1.Paths), len(http2.Paths))
	}
	for i := 0; i < len(http1.Paths); i++ {
		if http1.Paths[i].Path != http2.Paths[i].Path {
			return fmt.Errorf("%w. Name ingress: '%s'. Name path in first ingress - '%s'. Name path in second ingress - '%s'", ErrorPathValueDifferent, name, http1.Paths[i].Path, http2.Paths[i].Path)
		}
		err := compareIngressesBackend(http1.Paths[i].Backend, http2.Paths[i].Backend, name)
		if err != nil {
			return err
		}
	}
	return nil
}
