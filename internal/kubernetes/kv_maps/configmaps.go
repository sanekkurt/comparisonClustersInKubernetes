package kv_maps

import (
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"fmt"
	"sync"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

func addItemsToConfigMapList(clientSet kubernetes.Interface, namespace string, limit int64) (*v12.ConfigMapList, error) {
	log.Debugf("addItemsToConfigMapList started")
	defer log.Debugf("addItemsToConfigMapList completed")

	var (
		batch      *v12.ConfigMapList
		configMaps = &v12.ConfigMapList{
			Items: make([]v12.ConfigMap, 0),
		}

		continueToken string

		err error
	)

	for {
		batch, err = clientSet.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("addItemsToConfigMapList: %d objects received", len(batch.Items))

		configMaps.Items = append(configMaps.Items, batch.Items...)

		configMaps.TypeMeta = batch.TypeMeta
		configMaps.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	configMaps.Continue = ""

	return configMaps, err
}

// CompareConfigMaps compares list of configmap objects in two given k8s-clusters
func CompareConfigMaps(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	configMaps1, err := addItemsToConfigMapList(clientSet1, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 1st cluster: %w", err)
	}

	configMaps2, err := addItemsToConfigMapList(clientSet2, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 2st cluster: %w", err)
	}

	mapConfigMaps1, mapConfigMaps2 := prepareConfigMapMaps(configMaps1, configMaps2, skipEntityList.GetByKind("configmaps"))

	isClustersDiffer = compareConfigMapsSpecs(mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2)

	return isClustersDiffer, nil
}

// prepareConfigMapMaps add value ConfigMaps in map
func prepareConfigMapMaps(configMaps1, configMaps2 *v12.ConfigMapList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapConfigMap1 := make(map[string]types.IsAlreadyComparedFlag)
	mapConfigMap2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range configMaps1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("configmap %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapConfigMap1[value.Name] = indexCheck
	}

	for index, value := range configMaps2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("configmap %s is skipped from comparison due to its name", value.Name)
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

	log.Debugf("----- Start checking configmap: '%s' -----", name)

	if !AreKVMapsEqual(cm1.ObjectMeta.Labels, cm2.ObjectMeta.Labels, common.SkippedKubeLabels) {
		log.Infof("metadata of configmap '%s' differs: different labels", cm1.Name)
		channel <- true
		return
	}

	if !AreKVMapsEqual(cm1.ObjectMeta.Annotations, cm2.ObjectMeta.Annotations, nil) {
		log.Infof("metadata of configmap '%s' differs: different annotations", cm1.Name)
		channel <- true
		return
	}

	if len(cm1.Data) != len(cm2.Data) {
		log.Infof("config map '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(cm1.Data), len(cm2.Data))
		flag = true
	} else {
		for key, value := range cm1.Data {
			if cm2.Data[key] != value {
				log.Infof("configmap '%s', values by key '%s' do not match: '%s' and %s", name, key, value, cm2.Data[key])
				flag = true
			}
		}
	}
	log.Debugf("----- End checking configmap: '%s' -----", name)

	channel <- flag
}

// compareConfigMapsSpecs set information about config maps
func compareConfigMapsSpecs(map1, map2 map[string]types.IsAlreadyComparedFlag, configMaps1, configMaps2 *v12.ConfigMapList) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		log.Infof("configmaps count are different")
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
			log.Infof("ConfigMap '%s' - 1 cluster. Does not exist on another cluster", name)
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
			log.Infof("ConfigMap '%s' - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}

	return flag
}
