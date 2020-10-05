package kv_maps

import (
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	"fmt"
	"sync"
)

func CompareConfigMaps(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntities config.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)
	configMaps1, err := clientSet1.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 1st cluster: %w", err)
	}

	configMaps2, err := clientSet2.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 2nd cluster: %w", err)
	}

	mapConfigMaps1, mapConfigMaps2 := AddValueConfigMapsInMap(configMaps1, configMaps2)

	isClustersDiffer = CompareConfigMapsSpecs(mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2)

	return isClustersDiffer, nil
}

// AddValueConfigMapsInMap add value ConfigMaps in map
func AddValueConfigMapsInMap(configMaps1, configMaps2 *v12.ConfigMapList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapConfigMap1 := make(map[string]types.IsAlreadyComparedFlag)
	mapConfigMap2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range configMaps1.Items {
		if _, ok := ToSkipEntities["configmaps"][value.Name]; ok {
			logging.Log.Debugf("configmap %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapConfigMap1[value.Name] = indexCheck
	}
	for index, value := range configMaps2.Items {
		if _, ok := ToSkipEntities["configmaps"][value.Name]; ok {
			logging.Log.Debugf("configmap %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapConfigMap2[value.Name] = indexCheck
	}
	return mapConfigMap1, mapConfigMap2
}

func compareConfigMapSpecInternals(wg *sync.WaitGroup, channel chan bool, name string, cm1, cm2 *v12.ConfigMap) {
	var (
		flag bool
	)
	defer func() {
		wg.Done()
	}()

	logging.Log.Debugf("----- Start checking configmap: '%s' -----", name)
	if len(cm1.Data) != len(cm2.Data) {
		logging.Log.Infof("config map '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(cm1.Data), len(cm2.Data))
		flag = true
	} else {
		for key, value := range cm1.Data {
			if cm2.Data[key] != value {
				logging.Log.Infof("configmap '%s', values by key '%s' do not match: '%s' and %s", name, key, value, cm2.Data[key])
				flag = true
			}
		}
	}
	logging.Log.Debugf("----- End checking configmap: '%s' -----", name)

	channel <- flag
}

// CompareConfigMapsSpecs set information about config maps
func CompareConfigMapsSpecs(map1, map2 map[string]IsAlreadyComparedFlag, configMaps1, configMaps2 *v12.ConfigMapList) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		logging.Log.Infof("configmaps count are different")
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

			compareConfigMapSpecInternals(wg, channel, name, &configMaps1.Items[index1.Index], &configMaps2.Items[index2.Index])
		} else {
			logging.Log.Infof("ConfigMap '%s' - 1 cluster. Does not exist on another cluster", name)
			flag = true
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
			logging.Log.Infof("ConfigMap '%s' - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}
	return flag
}
