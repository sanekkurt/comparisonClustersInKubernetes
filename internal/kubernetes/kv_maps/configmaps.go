package kv_maps

import (
	"context"

	"go.uber.org/zap"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"fmt"
	"sync"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

func GetConfigMapByName(ctx context.Context, clientSet kubernetes.Interface, namespace, configMapName string) (*v12.ConfigMap, error) {
	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	return configMap, err
}

func addItemsToConfigMapList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v12.ConfigMapList, error) {
	log := logging.FromContext(ctx)

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
func CompareConfigMaps(ctx context.Context, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", "configmap"))

		clientSet1, clientSet2, namespace = config.FromContext(ctx)

		isClustersDiffer bool
	)
	ctx = logging.WithLogger(ctx, log)

	configMaps1, err := addItemsToConfigMapList(ctx, clientSet1, namespace, secretObjectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 1st cluster: %w", err)
	}

	configMaps2, err := addItemsToConfigMapList(ctx, clientSet2, namespace, secretObjectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain configmaps list from 2st cluster: %w", err)
	}

	mapConfigMaps1, mapConfigMaps2 := prepareConfigMapMaps(ctx, configMaps1, configMaps2, skipEntityList.GetByKind("configmaps"))

	isClustersDiffer = compareConfigMapsSpecs(ctx, mapConfigMaps1, mapConfigMaps2, configMaps1, configMaps2)

	return isClustersDiffer, nil
}

// prepareConfigMapMaps add value ConfigMaps in map
func prepareConfigMapMaps(ctx context.Context, configMaps1, configMaps2 *v12.ConfigMapList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	log := logging.FromContext(ctx)

	mapConfigMap1 := make(map[string]types.IsAlreadyComparedFlag)
	mapConfigMap2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range configMaps1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("configmap/%s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapConfigMap1[value.Name] = indexCheck
	}

	for index, value := range configMaps2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("configmap/%s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapConfigMap2[value.Name] = indexCheck
	}

	return mapConfigMap1, mapConfigMap2
}

func compareConfigMapSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, cm1, cm2 *v12.ConfigMap) {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))

		flag bool
	)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking configmap/%s -----", name)

	if !metadata.IsMetadataDiffers(ctx, cm1.ObjectMeta, cm2.ObjectMeta) {
		channel <- true
		return
	}

	if !common.AreKVMapsEqual(ctx, cm1.Data, cm2.Data, nil) {
		flag = true
	}

	log.Debugf("----- End checking configmap/%s -----", name)

	channel <- flag
}

// compareConfigMapsSpecs set information about config maps
func compareConfigMapsSpecs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, configMaps1, configMaps2 *v12.ConfigMapList) bool {
	var (
		log = logging.FromContext(ctx)

		flag bool
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return true

		default:
			if index2, ok := map2[name]; ok {
				wg.Add(1)

				index1.Check = true
				map1[name] = index1
				index2.Check = true
				map2[name] = index2

				compareConfigMapSpecInternals(ctx, wg, channel, name, &configMaps1.Items[index1.Index], &configMaps2.Items[index2.Index])
			} else {
				log.Infof("configmap/%s - 1 cluster. Does not exist on another cluster", name)
				flag = true
			}
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
			log.Infof("configmap/%s - 2 cluster. Does not exist on another cluster", name)
			flag = true
		}
	}

	return flag
}
