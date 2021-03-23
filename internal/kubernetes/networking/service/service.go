package service

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/consts"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"

	"sync"

	"go.uber.org/zap"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	objectKind = "service"
)

type Comparator struct {
	Kind      string
	Namespace string
	BatchSize int64
}

func NewServicesComparator(ctx context.Context, namespace string) *Comparator {
	return &Comparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),
	}
}

func (cmp *Comparator) fieldSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) labelSelectorProvider(ctx context.Context) string {
	return ""
}

//func servicesRetrieveBatchLimit(ctx context.Context) int64 {
//	cfg := config.FromContext(ctx)
//
//	if limit := cfg.Networking.Services.BatchSize; limit != 0 {
//		return limit
//	}
//
//	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
//		return limit
//	}
//
//	return 25
//}

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]v12.Service, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]v12.Service)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.CoreV1().Services(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.CoreV1().Services(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	return objects, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]v12.Service, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *v12.ServiceList
		objects = make(map[string]v12.Service)

		continueToken string

		err error
	)

	log.Debugf("%T: collectFromClusterWithoutExcludes started", cmp)
	defer log.Debugf("%T: collectFromClusterWithoutExcludes completed", cmp)

forOuterLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.CoreV1().Services(cmp.Namespace).List(metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.fieldSelectorProvider(ctx),
				LabelSelector: cmp.labelSelectorProvider(ctx),
				Continue:      continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %ss retrieved", len(batch.Items), cmp.Kind)

		forInnerLoop:
			for _, obj := range batch.Items {
				if _, ok := objects[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoop
				}

				objects[obj.Name] = obj
			}

			if batch.Continue == "" {
				break forOuterLoop
			}

			continueToken = batch.Continue
		}
	}

	return objects, nil
}

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]v12.Service, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludes(ctx)
	} else {
		return cmp.collectIncludedFromCluster(ctx)
	}
}

//type ServicesComparator struct {
//}
//
//func NewServicesComparator(ctx context.Context, namespace string) ServicesComparator {
//	return ServicesComparator{}
//}

// Compare compares list of services objects in two given k8s-clusters
func (cmp *Comparator) Compare(ctx context.Context) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Networking.Enabled ||
		!cfg.Networking.Services.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	objects, err := cmp.collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve objects for comparision: %w", err)
	}

	diff := cmp.compare(ctx, objects[0], objects[1])

	return diff, nil

	//services1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, servicesRetrieveBatchLimit(ctx))
	//if err != nil {
	//	return nil, fmt.Errorf("cannot obtain services list from 1st cluster: %w", err)
	//}
	//
	//services2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, servicesRetrieveBatchLimit(ctx))
	//if err != nil {
	//	return nil, fmt.Errorf("cannot obtain services list from 2st cluster: %w", err)
	//}
	//
	//mapServices1, mapServices2 := prepareServiceMaps(ctx, services1, services2)
	//
	//_ = compareServicesSpecs(ctx, mapServices1, mapServices2, services1, services2)
	//
	//return nil, nil
}

func (cmp *Comparator) collect(ctx context.Context) ([]map[string]v12.Service, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]v12.Service, 2, 2)
		wg      = &sync.WaitGroup{}

		err error
	)

	wg.Add(2)

	for idx, clientSet := range []kubernetes.Interface{
		cfg.Connections.Cluster1.ClientSet,
		cfg.Connections.Cluster2.ClientSet,
	} {
		go func(idx int, clientSet kubernetes.Interface) {
			defer wg.Done()

			objects[idx], err = cmp.collectFromCluster(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Fatalf("cannot obtain %ss from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]v12.Service) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.KubeObjectsDifference, 0)
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
	}

	for name, obj1 := range map1 {
		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return nil
		default:
			if obj2, ok := map2[name]; ok {
				diff := compareServicesSpecs(ctx, name, &obj1, &obj2)

				diffs = append(diffs, diff...)

				delete(map1, name)
				delete(map2, name)
			} else {
				log.With(zap.String("objectName", name)).Warnf("%s does not exist in 2nd cluster", cmp.Kind)
			}
		}
	}

	for name, _ := range map2 {
		log.With(zap.String("objectName", name)).Warnf("%s does not exist in 1st cluster", cmp.Kind)
	}

	return diffs
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v12.ServiceList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch    *v12.ServiceList
		services = &v12.ServiceList{
			Items: make([]v12.Service, 0),
		}

		continueToken string

		err error
	)

	log.Debugf("fillInComparisonMap started")
	defer log.Debugf("fillInComparisonMap completed")

	for {
		batch, err = clientSet.CoreV1().Services(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

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

//// prepareServiceMaps add value secrets in map
//func prepareServiceMaps(ctx context.Context, services1, services2 *v12.ServiceList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
//	var (
//		log = logging.FromContext(ctx)
//		cfg = config.FromContext(ctx)
//
//		mapServices1 = make(map[string]types.IsAlreadyComparedFlag)
//		mapServices2 = make(map[string]types.IsAlreadyComparedFlag)
//
//		indexCheck types.IsAlreadyComparedFlag
//	)
//
//	for index, value := range services1.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(serviceKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("service/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		indexCheck.Index = index
//		mapServices1[value.Name] = indexCheck
//	}
//
//	for index, value := range services2.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(serviceKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("service/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		indexCheck.Index = index
//		mapServices2[value.Name] = indexCheck
//
//	}
//
//	return mapServices1, mapServices2
//}

//// compareServicesSpecs set information about services
//func compareServicesSpecs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, services1, services2 *v12.ServiceList) bool {
//	var (
//		log = logging.FromContext(ctx)
//
//		flag bool
//	)
//
//	if len(map1) != len(map2) {
//		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
//		flag = true
//	}
//
//	wg := &sync.WaitGroup{}
//	channel := make(chan bool, len(map1))
//
//	for name, index1 := range map1 {
//		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))
//
//		if index2, ok := map2[name]; ok {
//			wg.Add(1)
//
//			index1.Check = true
//			map1[name] = index1
//			index2.Check = true
//			map2[name] = index2
//
//			compareServiceSpecInternals(ctx, wg, channel, name, &services1.Items[index1.Index], &services2.Items[index2.Index]) // тут была горутина
//
//		} else {
//			log.With(zap.String("objectName", name)).Warn("service does not exist in 2nd cluster")
//			flag = true
//			channel <- flag
//		}
//	}
//
//	wg.Wait()
//
//	close(channel)
//
//	for ch := range channel {
//		if ch {
//			flag = true
//		}
//	}
//
//	for name, index := range map2 {
//		if !index.Check {
//			log.With(zap.String("objectName", name)).Warn("service does not exist in 1st cluster")
//			flag = true
//		}
//	}
//
//	return flag
//}

//compareServicesSpecs set information about services
func compareServicesSpecs(ctx context.Context, name string, svc1, svc2 *v12.Service) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))
	)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("service/%s compare started", name)
	defer func() {
		log.Debugf("service/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, svc1.ObjectMeta, svc2.ObjectMeta)
	compareSpecInServices(ctx, *svc1, *svc2)

	return nil
}

//func compareServiceSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, svc1, svc2 *v12.Service) {
//	var (
//		log  = logging.FromContext(ctx).With(zap.String("objectName", name))
//		flag bool
//	)
//	ctx = logging.WithLogger(ctx, log)
//
//	defer func() {
//		wg.Done()
//	}()
//
//	log.Debugf("----- Start checking service: '%s' -----", name)
//
//	if !metadata.IsMetadataDiffers(ctx, svc1.ObjectMeta, svc2.ObjectMeta) {
//		channel <- true
//		return
//	}
//
//	err := compareSpecInServices(ctx, *svc1, *svc2)
//	if err != nil {
//		log.Infof("Service %s: %s", name, err.Error())
//		flag = true
//	}
//
//	log.Debugf("----- End checking service: '%s' -----", name)
//	channel <- flag
//}

// compareSpecInServices compares spec in services
func compareSpecInServices(ctx context.Context, service1, service2 v12.Service) error {
	var (
		log = logging.FromContext(ctx)
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return nil
	default:
		if len(service1.Spec.Ports) != len(service2.Spec.Ports) {
			log.With(zap.String("objectName", service1.Name)).Warnf("%s. %d vs %d", ErrorPortsCountDifferent.Error(), len(service1.Spec.Ports), len(service2.Spec.Ports))
			return nil
		}

		for index, value := range service1.Spec.Ports {
			if value != service2.Spec.Ports[index] {
				//return fmt.Errorf("%w. Name service: '%s'. First service: %s-%d-%s. Second service: %s-%d-%s", ErrorPortInServicesDifferent, service1.Name, value.Name, value.Port, value.Protocol, service2.Spec.Ports[index].Name, service2.Spec.Ports[index].Port, service2.Spec.Ports[index].Protocol)
				log.With(zap.String("objectName", service1.Name)).Warnf("%s. %s-%d-%s vs %s-%d-%s", ErrorPortInServicesDifferent.Error(), value.Name, value.Port, value.Protocol, service2.Spec.Ports[index].Name, service2.Spec.Ports[index].Port, service2.Spec.Ports[index].Protocol)
				return nil
			}
		}

		if len(service1.Spec.Selector) != len(service2.Spec.Selector) {
			//return fmt.Errorf("%w. Name service: '%s'. In first service - %d selectors, in second service - '%d' selectors", ErrorSelectorsCountDifferent, service1.Name, len(service1.Spec.Selector), len(service2.Spec.Selector))
			log.With(zap.String("objectName", service1.Name)).Warnf("%s. %d vs %d", ErrorSelectorsCountDifferent.Error(), len(service1.Spec.Selector), len(service2.Spec.Selector))
			return nil
		}

		for key, value := range service1.Spec.Selector {
			if service2.Spec.Selector[key] != value {
				//return fmt.Errorf("%w. Name service: '%s'. First service: %s-%s. Second service: %s-%s", ErrorSelectorInServicesDifferent, service1.Name, key, value, key, service2.Spec.Selector[key])
				log.With(zap.String("objectName", service1.Name)).Warnf("%s. %s-%s vs %s-%s", ErrorSelectorInServicesDifferent.Error(), key, value, key, service2.Spec.Selector[key])
				return nil

			}
		}

		if service1.Spec.Type != service2.Spec.Type {
			//return fmt.Errorf("%w. Name service: '%s'. First service type: %s. Second service type: %s", ErrorTypeInServicesDifferent, service1.Name, service1.Spec.Type, service2.Spec.Type)
			log.With(zap.String("objectName", service1.Name)).Warnf("%s. %s vs %s", ErrorTypeInServicesDifferent.Error(), service1.Spec.Type, service2.Spec.Type)
		}
		return nil
	}
}
