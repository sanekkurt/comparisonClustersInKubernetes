package ingress

import (
	"context"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/consts"
	"k8s-cluster-comparator/internal/kubernetes/diff"

	"go.uber.org/zap"
	v1 "k8s.io/api/networking/v1"

	"sync"

	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	v1beta12 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//const (
//	ingressKind = "ingress"
//)
//
//func ingressesRetrieveBatchLimit(ctx context.Context) int64 {
//	cfg := config.FromContext(ctx)
//
//	if limit := cfg.Networking.Ingresses.BatchSize; limit != 0 {
//		return limit
//	}
//
//	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
//		return limit
//	}
//
//	return 25
//}
//
//func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v1beta12.IngressList, error) {
//	var (
//		log       = logging.FromContext(ctx)
//		clientSet = kubectx.ClientSetFromContext(ctx)
//
//		batch     *v1beta12.IngressList
//		ingresses = &v1beta12.IngressList{
//			Items: make([]v1beta12.Ingress, 0),
//		}
//
//		continueToken string
//
//		err error
//	)
//
//	log.Debugf("fillInComparisonMap started")
//	defer log.Debugf("fillInComparisonMap completed")
//
//	for {
//		batch, err = clientSet.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{
//			Limit:    limit,
//			Continue: continueToken,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))
//
//		ingresses.Items = append(ingresses.Items, batch.Items...)
//
//		ingresses.TypeMeta = batch.TypeMeta
//		ingresses.ListMeta = batch.ListMeta
//
//		if batch.Continue == "" {
//			break
//		}
//
//		continueToken = batch.Continue
//	}
//
//	ingresses.Continue = ""
//
//	return ingresses, err
//}
//
//type IngressesComparator struct {
//}
//
//func NewIngressesComparator(ctx context.Context, namespace string) IngressesComparator {
//	return IngressesComparator{}
//}
//
//// Compare compares list of ingresses objects in two given k8s-clusters
//func (cmd IngressesComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
//	var (
//		log = logging.FromContext(ctx).With(zap.String("kind", ingressKind))
//		cfg = config.FromContext(ctx)
//	)
//	ctx = logging.WithLogger(ctx, log)
//
//	if !cfg.Networking.Enabled ||
//		!cfg.Networking.Ingresses.Enabled {
//		log.Infof("'%s' kind skipped from comparison due to configuration", ingressKind)
//		return nil, nil
//	}
//
//	ingresses1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, ingressesRetrieveBatchLimit(ctx))
//	if err != nil {
//		return nil, fmt.Errorf("cannot obtain ingresses list from 1st cluster: %w", err)
//	}
//
//	ingresses2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, ingressesRetrieveBatchLimit(ctx))
//	if err != nil {
//		return nil, fmt.Errorf("cannot obtain ingresses list from 2st cluster: %w", err)
//	}
//
//	mapIngresses1, mapIngresses2 := prepareIngressMaps(ctx, ingresses1, ingresses2)
//
//	_ = setInformationAboutIngresses(ctx, mapIngresses1, mapIngresses2, ingresses1, ingresses2)
//
//	return nil, nil
//}
//
//// prepareIngressMaps add value secrets in map
//func prepareIngressMaps(ctx context.Context, ingresses1, ingresses2 *v1beta12.IngressList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
//	var (
//		log = logging.FromContext(ctx)
//		cfg = config.FromContext(ctx)
//
//		mapIngresses1 = make(map[string]types.IsAlreadyComparedFlag)
//		mapIngresses2 = make(map[string]types.IsAlreadyComparedFlag)
//
//		indexCheck types.IsAlreadyComparedFlag
//	)
//
//	for index, value := range ingresses1.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(ingressKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("ingress/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		indexCheck.Index = index
//		mapIngresses1[value.Name] = indexCheck
//
//	}
//	for index, value := range ingresses2.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(ingressKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("ingress/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		indexCheck.Index = index
//		mapIngresses2[value.Name] = indexCheck
//
//	}
//	return mapIngresses1, mapIngresses2
//}
//
//func compareIngressSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, ing1, ing2 *v1beta12.Ingress) {
//	var (
//		log = logging.FromContext(ctx).With(zap.String("objectName", name))
//
//		flag bool
//	)
//	ctx = logging.WithLogger(ctx, log)
//
//	defer func() {
//		wg.Done()
//	}()
//
//	log.Debugf("----- Start checking ingress/%s -----", name)
//
//	if !metadata.IsMetadataDiffers(ctx, ing1.ObjectMeta, ing2.ObjectMeta) {
//		channel <- true
//		return
//	}
//
//	err := compareSpecInIngresses(ctx, *ing1, *ing2)
//	if err != nil {
//		log.Warnw(err.Error())
//		flag = true
//	}
//
//	log.Debugf("----- End checking ingress/%s -----", name)
//	channel <- flag
//}
//
//// setInformationAboutIngresses set information about ingresses
//func setInformationAboutIngresses(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, ingresses1, ingresses2 *v1beta12.IngressList) bool {
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
//			compareIngressSpecInternals(ctx, wg, channel, name, &ingresses1.Items[index1.Index], &ingresses2.Items[index2.Index])
//		} else {
//			log.With(zap.String("objectName", name)).Warn("ingress does not exist in 2nd cluster")
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
//	for name, index := range map2 {
//		if !index.Check {
//			log.With(zap.String("objectName", name)).Warn("ingress does not exist in 1st cluster")
//			flag = true
//		}
//	}
//
//	return flag
//}

const (
	objectKind = "ingress"
)

type Comparator struct {
	Kind      string
	Namespace string
	BatchSize int64
}

func NewComparator(ctx context.Context, namespace string) *Comparator {
	return &Comparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),
	}
}

func (cmp *Comparator) FieldSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) LabelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) collectIncludedFromClusterV1Beta1(ctx context.Context) (map[string]v1beta12.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objectsV1beta1 = make(map[string]v1beta12.Ingress)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		objV1beta1, err := clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objectsV1beta1[objV1beta1.Name] = *objV1beta1
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objectsV1beta1[obj.Name] = *obj
	}

	return objectsV1beta1, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludesV1Beta1(ctx context.Context) (map[string]v1beta12.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batchV1Beta1   *v1beta12.IngressList
		objectsV1Beta1 = make(map[string]v1beta12.Ingress)

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
			batchV1Beta1, err = clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).List(ctx, metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.FieldSelectorProvider(ctx),
				LabelSelector: cmp.LabelSelectorProvider(ctx),
				Continue:      continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %ss retrieved", len(batchV1Beta1.Items), cmp.Kind)

		forInnerLoopV1Beta1:
			for _, obj := range batchV1Beta1.Items {
				if _, ok := objectsV1Beta1[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoopV1Beta1
				}

				objectsV1Beta1[obj.Name] = obj
			}

			if batchV1Beta1.Continue == "" {
				break forOuterLoop
			}

			continueToken = batchV1Beta1.Continue

		}

	}

	return objectsV1Beta1, nil
}

func (cmp *Comparator) collectFromClusterV1Beta1(ctx context.Context) (map[string]v1beta12.Ingress, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludesV1Beta1(ctx)
	} else {
		return cmp.collectIncludedFromClusterV1Beta1(ctx)
	}
	//return cmp.collectIncludedFromClusterV1Beta1(ctx)
}

func (cmp *Comparator) collectIncludedFromClusterV1(ctx context.Context) (map[string]v1.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objectsV1 = make(map[string]v1.Ingress)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		objV1, err := clientSet.NetworkingV1().Ingresses(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objectsV1[objV1.Name] = *objV1
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.NetworkingV1().Ingresses(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objectsV1[obj.Name] = *obj
	}

	return objectsV1, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludesV1(ctx context.Context) (map[string]v1.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batchV1   *v1.IngressList
		objectsV1 = make(map[string]v1.Ingress)

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
			batchV1, err = clientSet.NetworkingV1().Ingresses(cmp.Namespace).List(ctx, metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.FieldSelectorProvider(ctx),
				LabelSelector: cmp.LabelSelectorProvider(ctx),
				Continue:      continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %s retrieved", len(batchV1.Items), cmp.Kind)

		forInnerLoopV1Beta1:
			for _, obj := range batchV1.Items {
				if _, ok := objectsV1[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoopV1Beta1
				}

				objectsV1[obj.Name] = obj
			}

			if batchV1.Continue == "" {
				break forOuterLoop
			}

			continueToken = batchV1.Continue

		}

	}

	return objectsV1, nil
}

func (cmp *Comparator) collectFromClusterV1(ctx context.Context) (map[string]v1.Ingress, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludesV1(ctx)
	} else {
		return cmp.collectIncludedFromClusterV1(ctx)
	}

}

// Compare compares list of services objects in two given k8s-clusters
func (cmp *Comparator) Compare(ctx context.Context) (*diff.DiffsStorage, error) {
	var (
		log          = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg          = config.FromContext(ctx)
		resourceName = "ingresses"
		resourceKind = "Ingress"
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Networking.Enabled ||
		!cfg.Networking.Ingresses.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	if checkServerResourcesForGroupVersion(cfg.Connections.Cluster1.ClientSet, "networking.k8s.io/v1beta1", resourceName, resourceKind) {

		objectsV1Beta1, err := cmp.collectV1Beta1(ctx)

		if err != nil {
			log.With(zap.String("apiVersion", "v1Beta1")).Errorf("cannot retrieve objects for comparision: %s", err.Error())
		} else {

			objectsV1Map1, objectsV1Map2 := ingressesV1Beta1ToV1(ctx, objectsV1Beta1[0], objectsV1Beta1[1])

			func(ctx context.Context) {
				var (
					log = logging.FromContext(ctx).With(zap.String("apiVersion", "v1Beta1"))
				)

				ctx = logging.WithLogger(ctx, log)
				cmp.compare(ctx, objectsV1Map1, objectsV1Map2)
			}(ctx)

		}
	} else {
		log.With(zap.String("apiVersion", "v1Beta1")).Warnf("%s", ErrorApiV1Beta1NotSupported.Error())
	}

	if checkServerResourcesForGroupVersion(cfg.Connections.Cluster1.ClientSet, "networking.k8s.io/v1", resourceName, resourceKind) {
		objectsV1, err := cmp.collectV1(ctx)

		if err != nil {
			log.With(zap.String("apiVersion", "v1")).Errorf("cannot retrieve objects for comparision: %s", err.Error())
		} else {

			func(ctx context.Context) {
				var (
					log = logging.FromContext(ctx).With(zap.String("apiVersion", "v1"))
				)

				ctx = logging.WithLogger(ctx, log)
				cmp.compare(ctx, objectsV1[0], objectsV1[0])
			}(ctx)

		}
	} else {
		log.With(zap.String("apiVersion", "v1")).Warnf("%s", ErrorApiV1NotSupported.Error())
	}

	return nil, nil
}

func checkServerResourcesForGroupVersion(clientSet kubernetes.Interface, groupVersion, resourceName, resourceKind string) bool {
	resourceList, _ := clientSet.Discovery().ServerResourcesForGroupVersion(groupVersion)
	if resourceList != nil {
		if resourceList.APIResources != nil {
			for _, resource := range resourceList.APIResources {
				if resource.Name == resourceName && resource.Kind == resourceKind {
					return true
				}
			}
		}
	}
	return false
}

func (cmp *Comparator) collectV1Beta1(ctx context.Context) ([]map[string]v1beta12.Ingress, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("apiVersion", "v1Beta1"))
		cfg = config.FromContext(ctx)

		objects = make([]map[string]v1beta12.Ingress, 2, 2)
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

			objects[idx], err = cmp.collectFromClusterV1Beta1(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Errorf("cannot obtain %s from cluster #%d: %s", cmp.Kind, idx+1, err.Error())

			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) collectV1(ctx context.Context) ([]map[string]v1.Ingress, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("apiVersion", "v1Beta1"))
		cfg = config.FromContext(ctx)

		objects = make([]map[string]v1.Ingress, 2, 2)
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

			objects[idx], err = cmp.collectFromClusterV1(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Errorf("cannot obtain %s from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, err
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]v1.Ingress) []types.KubeObjectsDifference {
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
				diff := compareIngressesSpecs(ctx, name, &obj1, &obj2)

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

func ingressesV1Beta1ToV1(ctx context.Context, ing1, ing2 map[string]v1beta12.Ingress) (map[string]v1.Ingress, map[string]v1.Ingress) {
	var (
		log     = logging.FromContext(ctx)
		newMap1 = make(map[string]v1.Ingress, len(ing1))
		newMap2 = make(map[string]v1.Ingress, len(ing2))
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return nil, nil
	default:
		for key, value := range ing1 {
			newMap1[key] = convertV1Beta1toV1(ctx, value)
		}

		for key, value := range ing2 {
			newMap2[key] = convertV1Beta1toV1(ctx, value)
		}

		return newMap1, newMap2
	}

}

func convertV1Beta1toV1(ctx context.Context, V1Beta1 v1beta12.Ingress) v1.Ingress {
	var (
		log       = logging.FromContext(ctx)
		v1Ingress v1.Ingress
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return v1Ingress
	default:
		v1Ingress.TypeMeta = V1Beta1.TypeMeta
		v1Ingress.ObjectMeta = V1Beta1.ObjectMeta
		v1Ingress.Status = v1.IngressStatus(V1Beta1.Status)
		if V1Beta1.Spec.Backend != nil {
			v1Ingress.Spec.DefaultBackend.Resource = V1Beta1.Spec.Backend.Resource
			v1Ingress.Spec.DefaultBackend.Service.Port.Name = V1Beta1.Spec.Backend.ServicePort.StrVal
			v1Ingress.Spec.DefaultBackend.Service.Port.Number = V1Beta1.Spec.Backend.ServicePort.IntVal
			v1Ingress.Spec.DefaultBackend.Service.Name = V1Beta1.Spec.Backend.ServiceName
		}

		if V1Beta1.Spec.IngressClassName != nil {
			v1Ingress.Spec.IngressClassName = V1Beta1.Spec.IngressClassName
		}

		v1Ingress.Spec.Rules = make([]v1.IngressRule, len(V1Beta1.Spec.Rules))
		for i, rule := range V1Beta1.Spec.Rules {
			v1Ingress.Spec.Rules[i].Host = rule.Host
			if rule.IngressRuleValue.HTTP != nil {
				v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP = &v1.HTTPIngressRuleValue{Paths: make([]v1.HTTPIngressPath, len(rule.IngressRuleValue.HTTP.Paths))}
				for ii, path := range rule.IngressRuleValue.HTTP.Paths {
					v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Path = path.Path
					if path.PathType != nil {
						v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].PathType = (*v1.PathType)(path.PathType)
					}
					v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Backend.Service = &v1.IngressServiceBackend{
						Name: path.Backend.ServiceName,
						Port: v1.ServiceBackendPort{
							Name:   path.Backend.ServicePort.StrVal,
							Number: path.Backend.ServicePort.IntVal,
						},
					}
					//v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Backend.Service.Name = path.Backend.ServiceName
					//v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Backend.Service.Port.Name = path.Backend.ServicePort.StrVal
					//v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Backend.Service.Port.Number = path.Backend.ServicePort.IntVal
					if path.Backend.Resource != nil {
						v1Ingress.Spec.Rules[i].IngressRuleValue.HTTP.Paths[ii].Backend.Resource = path.Backend.Resource
					}

				}
			}

		}

		if V1Beta1.Spec.TLS != nil {
			v1Ingress.Spec.TLS = make([]v1.IngressTLS, len(V1Beta1.Spec.TLS))

			for i, tls := range V1Beta1.Spec.TLS {
				v1Ingress.Spec.TLS[i].SecretName = tls.SecretName
				if tls.Hosts != nil {
					v1Ingress.Spec.TLS[i].Hosts = make([]string, len(tls.Hosts))
					for ii, host := range tls.Hosts {
						v1Ingress.Spec.TLS[i].Hosts[ii] = host
					}
				}
			}
		}

		return v1Ingress
	}

}

//compareIngressesSpecs set information about services
func compareIngressesSpecs(ctx context.Context, name string, ing1, ing2 *v1.Ingress) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx)
	)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("ingress/%s compare started", name)
	defer func() {
		log.Debugf("ingress/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, ing1.ObjectMeta, ing2.ObjectMeta)
	err := compareSpecInIngresses(ctx, *ing1, *ing2)
	if err != nil {
		log.Warnf(err.Error())
	}

	return nil
}

// compareSpecInIngresses compare spec in the ingresses
func compareSpecInIngresses(ctx context.Context, ingress1, ingress2 v1.Ingress) error { //nolint
	var (
		log = logging.FromContext(ctx)
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return nil
	default:
		if ingress1.Spec.TLS != nil && ingress2.Spec.TLS != nil {

			if len(ingress1.Spec.TLS) != len(ingress2.Spec.TLS) {
				//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d TLS. In second ingress - %d TLS", ErrorTLSCountDifferent, ingress1.Name, len(ingress1.Spec.TLS), len(ingress2.Spec.TLS))
				log.Warnf("%s. %d vs %d", ErrorTLSCountDifferent.Error(), len(ingress1.Spec.TLS), len(ingress2.Spec.TLS))
				return nil
			}

			for index, value := range ingress1.Spec.TLS {

				if value.SecretName != ingress2.Spec.TLS[index].SecretName {
					//return fmt.Errorf("%w. Name ingress: '%s'. First ingress: '%s'. Second ingress: '%s'", ErrorSecretNameInTLSDifferent, ingress1.Name, value.SecretName, ingress2.Spec.TLS[index].SecretName)
					log.Warnf("%s. %s vs %s", ErrorSecretNameInTLSDifferent.Error(), value.SecretName, ingress2.Spec.TLS[index].SecretName)
					return nil
				}

				if value.Hosts != nil && ingress2.Spec.TLS[index].Hosts != nil {
					if len(value.Hosts) != len(ingress2.Spec.TLS[index].Hosts) {
						//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d hosts. In second ingress - %d hosts", ErrorHostsCountDifferent, ingress1.Name, len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
						log.Warnf("%s. %d vs %d", ErrorHostsCountDifferent.Error(), len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
						return nil
					}

					for i := 0; i < len(value.Hosts); i++ {
						if value.Hosts[i] != ingress2.Spec.TLS[index].Hosts[i] {
							//return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorNameHostDifferent, ingress1.Name, value.Hosts[i], ingress2.Spec.TLS[index].Hosts[i])
							log.Warnf("%s. %s vs %s", ErrorNameHostDifferent.Error(), value.Hosts[i], ingress2.Spec.TLS[index].Hosts[i])
							return nil
						}
					}

				} else if value.Hosts != nil || ingress2.Spec.TLS[index].Hosts != nil {
					//return fmt.Errorf("%w", ErrorHostsInIngressesDifferent)
					log.Warnf("%s", ErrorHostsInIngressesDifferent.Error())
					return nil
				}
			}
		} else if ingress1.Spec.TLS != nil || ingress2.Spec.TLS != nil {
			//return fmt.Errorf("%w", ErrorTLSInIngressesDifferent)
			log.Warnf("%s", ErrorTLSInIngressesDifferent.Error())
			return nil
		}

		if ingress1.Spec.DefaultBackend != nil && ingress2.Spec.DefaultBackend != nil {
			err := compareIngressesBackend(ctx, *ingress1.Spec.DefaultBackend, *ingress2.Spec.DefaultBackend, ingress1.Name)
			if err != nil {
				return err
			}
		} else if ingress1.Spec.DefaultBackend != nil || ingress2.Spec.DefaultBackend != nil {
			//return fmt.Errorf("%w", ErrorBackendInIngressesDifferent)
			log.Warnf("%s", ErrorBackendInIngressesDifferent.Error())
			return nil
		}

		if ingress1.Spec.Rules != nil && ingress2.Spec.Rules != nil {
			if len(ingress1.Spec.Rules) != len(ingress2.Spec.Rules) {
				//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' rules. In second ingress - '%d' rules", ErrorRulesCountDifferent, ingress1.Name, len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
				log.Warnf("%s. %d vs %d", ErrorRulesCountDifferent.Error(), len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
				return nil
			}

			for index, value := range ingress1.Spec.Rules {
				if value.Host != ingress2.Spec.Rules[index].Host {
					//return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorHostNameInRuleDifferent, ingress1.Name, value.Host, ingress2.Spec.Rules[index].Host)
					log.Warnf("%s. %s vs %s", ErrorHostNameInRuleDifferent.Error(), value.Host, ingress2.Spec.Rules[index].Host)
					return nil
				}

				if value.HTTP != nil && ingress2.Spec.Rules[index].HTTP != nil {
					err := compareIngressesHTTP(ctx, *value.HTTP, *ingress2.Spec.Rules[index].HTTP, ingress1.Name)
					if err != nil {
						return err
					}
				} else if value.HTTP != nil || ingress2.Spec.Rules[index].HTTP != nil {
					//return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
					log.Warnf("%s", ErrorHTTPInIngressesDifferent.Error())
					return nil
				}

				if value.IngressRuleValue.HTTP != nil && ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
					err := compareIngressesHTTP(ctx, *value.IngressRuleValue.HTTP, *ingress2.Spec.Rules[index].IngressRuleValue.HTTP, ingress1.Name)
					if err != nil {
						return err
					}
				} else if value.IngressRuleValue.HTTP != nil || ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
					//return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
					log.Warnf("%s", ErrorHTTPInIngressesDifferent.Error())
					return nil
				}

			}
		} else if ingress1.Spec.Rules != nil || ingress2.Spec.Rules != nil {
			//return fmt.Errorf("%w", ErrorRulesInIngressesDifferent)
			log.Warnf("%s", ErrorRulesInIngressesDifferent.Error())
			return nil
		}

		return nil
	}
}

// compareIngressesBackend compare backend in ingresses
func compareIngressesBackend(ctx context.Context, backend1, backend2 v1.IngressBackend, name string) error {
	var (
		log = logging.FromContext(ctx)
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return nil
	default:

		if backend1.Service != nil && backend2.Service != nil {
			if backend1.Service.Name != backend2.Service.Name {
				log.Warnf("%s. %s vs %s", ErrorServiceNameInBackendDifferent.Error(), backend1.Service.Name, backend2.Service.Name)
				return nil
			}
			if backend1.Service.Port != backend2.Service.Port {
				log.Warnf("%s", ErrorBackendServicePortDifferent.Error())
				return nil
			}
		} else if backend1.Service != nil || backend2.Service != nil {
			log.Warnf("%s", ErrorBackendServiceIsMissing.Error())
			return nil
		}

		if backend1.Resource != nil && backend2.Resource != nil {
			if backend1.Resource.APIGroup != nil && backend2.Resource.APIGroup != nil {
				if *backend1.Resource.APIGroup != *backend2.Resource.APIGroup {
					log.Warnf("%s. %s vs %s", ErrorBackendResourceApiGroup.Error(), *backend1.Resource.APIGroup, *backend2.Resource.APIGroup)
					return nil
				}
			} else if backend1.Resource.APIGroup != nil || backend2.Resource.APIGroup != nil {
				log.Warnf("%s", ErrorBackendResourceApiGroupIsMissing.Error())
				return nil
			}

			if backend1.Resource.Name != backend2.Resource.Name {
				log.Warnf("%s. %s vs %s", ErrorBackendResourceName.Error(), backend1.Resource.Name, backend2.Resource.Name)
				return nil
			}

			if backend1.Resource.Kind != backend2.Resource.Kind {
				log.Warnf("%s. %s vs %s", ErrorBackendResourceKind.Error(), backend1.Resource.Name, backend2.Resource.Name)
				return nil
			}

		}

		return nil
	}

}

// compareIngressesHTTP compare http in ingresses
func compareIngressesHTTP(ctx context.Context, http1, http2 v1.HTTPIngressRuleValue, name string) error {
	var (
		log = logging.FromContext(ctx)
	)
	if len(http1.Paths) != len(http2.Paths) {
		//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' paths. In second ingress - '%d' paths", ErrorPathsCountDifferent, name, len(http1.Paths), len(http2.Paths))
		log.Warnf("%s. %d vs %d", ErrorPathsCountDifferent.Error(), len(http1.Paths), len(http2.Paths))
		return nil
	}

	for i := 0; i < len(http1.Paths); i++ {
		if http1.Paths[i].Path != http2.Paths[i].Path {
			//return fmt.Errorf("%w. Name ingress: '%s'. Name path in first ingress - '%s'. Name path in second ingress - '%s'", ErrorPathValueDifferent, name, http1.Paths[i].Path, http2.Paths[i].Path)
			log.Warnf("%s. %s vs %s", ErrorPathValueDifferent.Error(), http1.Paths[i].Path, http2.Paths[i].Path)
			return nil
		}

		err := compareIngressesBackend(ctx, http1.Paths[i].Backend, http2.Paths[i].Backend, name)
		if err != nil {
			return err
		}
	}
	return nil
}
