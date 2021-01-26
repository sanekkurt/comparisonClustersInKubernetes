package networking

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"sync"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	objectBatchLimit = 25
)

func addItemsToServiceList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v12.ServiceList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToServiceList started")
	defer log.Debugf("addItemsToServiceList completed")

	var (
		batch    *v12.ServiceList
		services = &v12.ServiceList{
			Items: make([]v12.Service, 0),
		}

		continueToken string

		err error
	)

	for {
		batch, err = clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("addItemsToServiceList: %d objects received", len(batch.Items))

		services.Items = append(services.Items, batch.Items...)

		services.TypeMeta = batch.TypeMeta
		services.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	services.Continue = ""

	return services, err
}

// CompareServices compares list of services objects in two given k8s-clusters
func CompareServices(ctx context.Context, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", "service"))

		clientSet1, clientSet2, namespace = config.FromContext(ctx)

		isClustersDiffer bool
	)
	ctx = logging.WithLogger(ctx, log)

	services1, err := addItemsToServiceList(ctx, clientSet1, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain services list from 1st cluster: %w", err)
	}

	services2, err := addItemsToServiceList(ctx, clientSet2, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain services list from 2st cluster: %w", err)
	}

	mapServices1, mapServices2 := prepareServiceMaps(ctx, services1, services2, skipEntityList.GetByKind("secrets"))

	isClustersDiffer = compareServicesSpecs(ctx, mapServices1, mapServices2, services1, services2)

	return isClustersDiffer, nil
}

// prepareServiceMaps add value secrets in map
func prepareServiceMaps(ctx context.Context, services1, services2 *v12.ServiceList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	log := logging.FromContext(ctx)

	mapServices1 := make(map[string]types.IsAlreadyComparedFlag)
	mapServices2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range services1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("service %s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		mapServices1[value.Name] = indexCheck
	}

	for index, value := range services2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("service %s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		mapServices2[value.Name] = indexCheck

	}

	return mapServices1, mapServices2
}

func compareServiceSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, svc1, svc2 *v12.Service) {
	var (
		log  = logging.FromContext(ctx).With(zap.String("objectName", name))
		flag bool
	)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking service: '%s' -----", name)

	if !metadata.IsMetadataDiffers(ctx, svc1.ObjectMeta, svc2.ObjectMeta) {
		channel <- true
		return
	}

	err := compareSpecInServices(ctx, *svc1, *svc2)
	if err != nil {
		log.Infof("Service %s: %s", name, err.Error())
		flag = true
	}

	log.Debugf("----- End checking service: '%s' -----", name)
	channel <- flag
}

// compareServicesSpecs set information about services
func compareServicesSpecs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, services1, services2 *v12.ServiceList) bool {
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
		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareServiceSpecInternals(ctx, wg, channel, name, &services1.Items[index1.Index], &services2.Items[index2.Index]) // тут была горутина

		} else {
			log.Infof("service '%s' does not exist in 2nd cluster", name)
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

			log.Infof("service '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}

	return flag
}

// compareSpecInServices compares spec in services
func compareSpecInServices(ctx context.Context, service1, service2 v12.Service) error {
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
