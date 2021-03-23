package ingress

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/consts"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1beta12 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
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

func NewIngressesComparator(ctx context.Context, namespace string) *Comparator {
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

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]v1beta12.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]v1beta12.Ingress)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).Get(string(name), metav1.GetOptions{})
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
		obj, err := clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).Get(string(name), metav1.GetOptions{})
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

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]v1beta12.Ingress, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *v1beta12.IngressList
		objects = make(map[string]v1beta12.Ingress)

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
			batch, err = clientSet.NetworkingV1beta1().Ingresses(cmp.Namespace).List(metav1.ListOptions{
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

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]v1beta12.Ingress, error) {
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
}

func (cmp *Comparator) collect(ctx context.Context) ([]map[string]v1beta12.Ingress, error) {
	var (
		log = logging.FromContext(ctx)
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

			objects[idx], err = cmp.collectFromCluster(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Fatalf("cannot obtain %ss from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]v1beta12.Ingress) []types.KubeObjectsDifference {
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

//compareIngressesSpecs set information about services
func compareIngressesSpecs(ctx context.Context, name string, ing1, ing2 *v1beta12.Ingress) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))
	)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("service/%s compare started", name)
	defer func() {
		log.Debugf("service/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, ing1.ObjectMeta, ing2.ObjectMeta)
	compareSpecInIngresses(ctx, *ing1, *ing2)

	return nil
}

// compareSpecInIngresses compare spec in the ingresses
func compareSpecInIngresses(ctx context.Context, ingress1, ingress2 v1beta12.Ingress) error { //nolint
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
				log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %d vs %d", ErrorTLSCountDifferent.Error(), len(ingress1.Spec.TLS), len(ingress2.Spec.TLS))
				return nil
			}

			for index, value := range ingress1.Spec.TLS {

				if value.SecretName != ingress2.Spec.TLS[index].SecretName {
					//return fmt.Errorf("%w. Name ingress: '%s'. First ingress: '%s'. Second ingress: '%s'", ErrorSecretNameInTLSDifferent, ingress1.Name, value.SecretName, ingress2.Spec.TLS[index].SecretName)
					log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %s vs %s", ErrorSecretNameInTLSDifferent.Error(), value.SecretName, ingress2.Spec.TLS[index].SecretName)
					return nil
				}

				if value.Hosts != nil && ingress2.Spec.TLS[index].Hosts != nil {
					if len(value.Hosts) != len(ingress2.Spec.TLS[index].Hosts) {
						//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - %d hosts. In second ingress - %d hosts", ErrorHostsCountDifferent, ingress1.Name, len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
						log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %d vs %d", ErrorHostsCountDifferent.Error(), len(value.Hosts), len(ingress2.Spec.TLS[index].Hosts))
						return nil
					}

					for i := 0; i < len(value.Hosts); i++ {
						if value.Hosts[i] != ingress2.Spec.TLS[index].Hosts[i] {
							//return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorNameHostDifferent, ingress1.Name, value.Hosts[i], ingress2.Spec.TLS[index].Hosts[i])
							log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %s vs %s", ErrorNameHostDifferent.Error(), value.Hosts[i], ingress2.Spec.TLS[index].Hosts[i])
							return nil
						}
					}

				} else if value.Hosts != nil || ingress2.Spec.TLS[index].Hosts != nil {
					//return fmt.Errorf("%w", ErrorHostsInIngressesDifferent)
					log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorHostsInIngressesDifferent.Error())
					return nil
				}
			}
		} else if ingress1.Spec.TLS != nil || ingress2.Spec.TLS != nil {
			//return fmt.Errorf("%w", ErrorTLSInIngressesDifferent)
			log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorTLSInIngressesDifferent.Error())
			return nil
		}

		if ingress1.Spec.Backend != nil && ingress2.Spec.Backend != nil {
			err := compareIngressesBackend(ctx, *ingress1.Spec.Backend, *ingress2.Spec.Backend, ingress1.Name)
			if err != nil {
				return err
			}
		} else if ingress1.Spec.Backend != nil || ingress2.Spec.Backend != nil {
			//return fmt.Errorf("%w", ErrorBackendInIngressesDifferent)
			log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorBackendInIngressesDifferent.Error())
			return nil
		}

		if ingress1.Spec.Rules != nil && ingress2.Spec.Rules != nil {
			if len(ingress1.Spec.Rules) != len(ingress2.Spec.Rules) {
				//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' rules. In second ingress - '%d' rules", ErrorRulesCountDifferent, ingress1.Name, len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
				log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %d vs %d", ErrorRulesCountDifferent.Error(), len(ingress1.Spec.Rules), len(ingress2.Spec.Rules))
				return nil
			}

			for index, value := range ingress1.Spec.Rules {
				if value.Host != ingress2.Spec.Rules[index].Host {
					//return fmt.Errorf("%w. Name ingress: '%s'. Name host in first ingress - '%s'. Name host in second ingress - '%s'", ErrorHostNameInRuleDifferent, ingress1.Name, value.Host, ingress2.Spec.Rules[index].Host)
					log.With(zap.String("objectName", ingress1.Name)).Warnf("%s. %s vs %s", ErrorHostNameInRuleDifferent.Error(), value.Host, ingress2.Spec.Rules[index].Host)
					return nil
				}

				if value.HTTP != nil && ingress2.Spec.Rules[index].HTTP != nil {
					err := compareIngressesHTTP(ctx, *value.HTTP, *ingress2.Spec.Rules[index].HTTP, ingress1.Name)
					if err != nil {
						return err
					}
				} else if value.HTTP != nil || ingress2.Spec.Rules[index].HTTP != nil {
					//return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
					log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorHTTPInIngressesDifferent.Error())
					return nil
				}

				if value.IngressRuleValue.HTTP != nil && ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
					err := compareIngressesHTTP(ctx, *value.IngressRuleValue.HTTP, *ingress2.Spec.Rules[index].IngressRuleValue.HTTP, ingress1.Name)
					if err != nil {
						return err
					}
				} else if value.IngressRuleValue.HTTP != nil || ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
					//return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
					log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorHTTPInIngressesDifferent.Error())
					return nil
				}

			}
		} else if ingress1.Spec.Rules != nil || ingress2.Spec.Rules != nil {
			//return fmt.Errorf("%w", ErrorRulesInIngressesDifferent)
			log.With(zap.String("objectName", ingress1.Name)).Warnf("%s", ErrorRulesInIngressesDifferent.Error())
			return nil
		}

		return nil
	}
}

// compareIngressesBackend compare backend in ingresses
func compareIngressesBackend(ctx context.Context, backend1, backend2 v1beta12.IngressBackend, name string) error {
	var (
		log = logging.FromContext(ctx)
	)

	select {
	case <-ctx.Done():
		log.Warnw(context.Canceled.Error())
		return nil
	default:
		if backend1.ServiceName != backend2.ServiceName {
			//return fmt.Errorf("%w. Name ingress: '%s'. Service name in first ingress: '%s'. Service name in second ingress: '%s'", ErrorServiceNameInBackendDifferent, name, backend1.ServiceName, backend2.ServiceName)
			log.With(zap.String("objectName", name)).Warnf("%s. %s vs %s", ErrorServiceNameInBackendDifferent.Error(), backend1.ServiceName, backend2.ServiceName)
			return nil
		}

		if backend1.ServicePort.Type != backend2.ServicePort.Type || backend1.ServicePort.IntVal != backend2.ServicePort.IntVal || backend1.ServicePort.StrVal != backend2.ServicePort.StrVal {
			//return fmt.Errorf("%w. Name ingress: '%s'", ErrorBackendServicePortDifferent, name)
			log.With(zap.String("objectName", name)).Warnf("%s", ErrorBackendServicePortDifferent.Error())
			return nil
		}
		return nil
	}

}

// compareIngressesHTTP compare http in ingresses
func compareIngressesHTTP(ctx context.Context, http1, http2 v1beta12.HTTPIngressRuleValue, name string) error {
	var (
		log = logging.FromContext(ctx)
	)
	if len(http1.Paths) != len(http2.Paths) {
		//return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' paths. In second ingress - '%d' paths", ErrorPathsCountDifferent, name, len(http1.Paths), len(http2.Paths))
		log.With(zap.String("objectName", name)).Warnf("%s. %d vs %d", ErrorPathsCountDifferent.Error(), len(http1.Paths), len(http2.Paths))
		return nil
	}

	for i := 0; i < len(http1.Paths); i++ {
		if http1.Paths[i].Path != http2.Paths[i].Path {
			//return fmt.Errorf("%w. Name ingress: '%s'. Name path in first ingress - '%s'. Name path in second ingress - '%s'", ErrorPathValueDifferent, name, http1.Paths[i].Path, http2.Paths[i].Path)
			log.With(zap.String("objectName", name)).Warnf("%s. %s vs %s", ErrorPathValueDifferent.Error(), http1.Paths[i].Path, http2.Paths[i].Path)
			return nil
		}

		err := compareIngressesBackend(ctx, http1.Paths[i].Backend, http2.Paths[i].Backend, name)
		if err != nil {
			return err
		}
	}
	return nil
}
