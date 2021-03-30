package ingress

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1beta12 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingressKind = "ingress"
)

func ingressesRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Networking.Ingresses.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v1beta12.IngressList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch     *v1beta12.IngressList
		ingresses = &v1beta12.IngressList{
			Items: make([]v1beta12.Ingress, 0),
		}

		continueToken string

		err error
	)

	log.Debugf("fillInComparisonMap started")
	defer log.Debugf("fillInComparisonMap completed")

	for {
		batch, err = clientSet.NetworkingV1beta1().Ingresses(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

		ingresses.Items = append(ingresses.Items, batch.Items...)

		ingresses.TypeMeta = batch.TypeMeta
		ingresses.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	ingresses.Continue = ""

	return ingresses, err
}

type IngressesComparator struct {
}

func NewIngressesComparator(ctx context.Context, namespace string) IngressesComparator {
	return IngressesComparator{}
}

// Compare compares list of ingresses objects in two given k8s-clusters
func (cmd IngressesComparator) Compare(ctx context.Context, namespace string) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", ingressKind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Networking.Enabled ||
		!cfg.Networking.Ingresses.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", ingressKind)
		return nil, nil
	}

	ingresses1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, ingressesRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain ingresses list from 1st cluster: %w", err)
	}

	ingresses2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, ingressesRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain ingresses list from 2st cluster: %w", err)
	}

	mapIngresses1, mapIngresses2 := prepareIngressMaps(ctx, ingresses1, ingresses2)

	_ = setInformationAboutIngresses(ctx, mapIngresses1, mapIngresses2, ingresses1, ingresses2)

	return nil, nil
}

// prepareIngressMaps add value secrets in map
func prepareIngressMaps(ctx context.Context, ingresses1, ingresses2 *v1beta12.IngressList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		mapIngresses1 = make(map[string]types.IsAlreadyComparedFlag)
		mapIngresses2 = make(map[string]types.IsAlreadyComparedFlag)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range ingresses1.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(ingressKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("ingress/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		mapIngresses1[value.Name] = indexCheck

	}
	for index, value := range ingresses2.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(ingressKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("ingress/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		mapIngresses2[value.Name] = indexCheck

	}
	return mapIngresses1, mapIngresses2
}

func compareIngressSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, ing1, ing2 *v1beta12.Ingress) {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))

		flag bool
	)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking ingress/%s -----", name)

	if !metadata.IsMetadataDiffers(ctx, ing1.ObjectMeta, ing2.ObjectMeta) {
		channel <- true
		return
	}

	err := compareSpecInIngresses(ctx, *ing1, *ing2)
	if err != nil {
		log.Warnw(err.Error())
		flag = true
	}

	log.Debugf("----- End checking ingress/%s -----", name)
	channel <- flag
}

// setInformationAboutIngresses set information about ingresses
func setInformationAboutIngresses(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, ingresses1, ingresses2 *v1beta12.IngressList) bool {
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
		ctx := logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareIngressSpecInternals(ctx, wg, channel, name, &ingresses1.Items[index1.Index], &ingresses2.Items[index2.Index])
		} else {
			log.With(zap.String("objectName", name)).Warn("ingress does not exist in 2nd cluster")
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
			log.With(zap.String("objectName", name)).Warn("ingress does not exist in 1st cluster")
			flag = true
		}
	}

	return flag
}

// compareSpecInIngresses compare spec in the ingresses
func compareSpecInIngresses(ctx context.Context, ingress1, ingress2 v1beta12.Ingress) error { //nolint
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
		err := compareIngressesBackend(ctx, *ingress1.Spec.Backend, *ingress2.Spec.Backend, ingress1.Name)
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
				err := compareIngressesHTTP(ctx, *value.HTTP, *ingress2.Spec.Rules[index].HTTP, ingress1.Name)
				if err != nil {
					return err
				}
			} else if value.HTTP != nil || ingress2.Spec.Rules[index].HTTP != nil {
				return fmt.Errorf("%w", ErrorHTTPInIngressesDifferent)
			}

			if value.IngressRuleValue.HTTP != nil && ingress2.Spec.Rules[index].IngressRuleValue.HTTP != nil {
				err := compareIngressesHTTP(ctx, *value.IngressRuleValue.HTTP, *ingress2.Spec.Rules[index].IngressRuleValue.HTTP, ingress1.Name)
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
func compareIngressesBackend(ctx context.Context, backend1, backend2 v1beta12.IngressBackend, name string) error {
	if backend1.ServiceName != backend2.ServiceName {
		return fmt.Errorf("%w. Name ingress: '%s'. Service name in first ingress: '%s'. Service name in second ingress: '%s'", ErrorServiceNameInBackendDifferent, name, backend1.ServiceName, backend2.ServiceName)
	}

	if backend1.ServicePort.Type != backend2.ServicePort.Type || backend1.ServicePort.IntVal != backend2.ServicePort.IntVal || backend1.ServicePort.StrVal != backend2.ServicePort.StrVal {
		return fmt.Errorf("%w. Name ingress: '%s'", ErrorBackendServicePortDifferent, name)
	}
	return nil
}

// compareIngressesHTTP compare http in ingresses
func compareIngressesHTTP(ctx context.Context, http1, http2 v1beta12.HTTPIngressRuleValue, name string) error {
	if len(http1.Paths) != len(http2.Paths) {
		return fmt.Errorf("%w. Name ingress: '%s'. In first ingress - '%d' paths. In second ingress - '%d' paths", ErrorPathsCountDifferent, name, len(http1.Paths), len(http2.Paths))
	}

	for i := 0; i < len(http1.Paths); i++ {
		if http1.Paths[i].Path != http2.Paths[i].Path {
			return fmt.Errorf("%w. Name ingress: '%s'. Name path in first ingress - '%s'. Name path in second ingress - '%s'", ErrorPathValueDifferent, name, http1.Paths[i].Path, http2.Paths[i].Path)
		}

		err := compareIngressesBackend(ctx, http1.Paths[i].Backend, http2.Paths[i].Backend, name)
		if err != nil {
			return err
		}
	}
	return nil
}
