package kubernetes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/networking/ingress"
	"k8s-cluster-comparator/internal/kubernetes/networking/service"
	"sync"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/daemonset"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/deployment"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/statefulset"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

type ResStr struct {
	IsClustersDiffer bool
	Err              error
}

func compareKubeNamespaces(ctx context.Context, ns string) (*types.KubeObjectsDifference, error) {
	log := logging.FromContext(ctx).With(zap.String("namespace", ns))

	log.Debugf("Processing namespace/%s", ns)
	defer func() {
		log.Debugf("End of namespace/%s processing", ns)
	}()

	return nil, nil
}

// CompareClusters main compare function, runs functions for comparing clusters by different parameters one at a time: Deployments, StatefulSets, DaemonSets, ConfigMaps
func CompareClusters(ctx context.Context) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		diffs = make([]types.KubeObjectsDifference, 0)
	)

	for _, namespace := range cfg.Connections.Namespaces {
		log := log.With(zap.String("namespace", namespace))

		comparators := []types.KubeResourceComparator{
			deployment.NewDeploymentsComparator(ctx, namespace),
			statefulset.NewStatefulSetsComparator(ctx, namespace),
			daemonset.NewDaemonSetsComparator(ctx, namespace),

			//job.NewJobsComparator(ctx, namespace),
			//cronjob.NewCronJobsComparator(ctx, namespace),

			configmap.NewConfigMapsComparator(ctx, namespace),
			secret.NewSecretsComparator(ctx, namespace),

			service.NewServicesComparator(ctx, namespace),
			ingress.NewIngressesComparator(ctx, namespace),
		}

		wg := &sync.WaitGroup{}
		diffsCh := make(chan []types.KubeObjectsDifference, len(comparators))

		for _, cmp := range comparators {
			select {
			case <-ctx.Done():
				break
			default:
				cmp := cmp

				wg.Add(1)

				go func(wg *sync.WaitGroup) {
					defer wg.Done()

					log.Debugf("%T started", cmp)
					defer func() {
						log.Debugf("%T completed", cmp)
					}()

					diffs, err := cmp.Compare(ctx)
					if err != nil {
						log.Errorf("cannot call %T: %s", cmp, err.Error())
						return
					}

					diffsCh <- diffs
				}(wg)
			}
		}

		wg.Wait()
		close(diffsCh)

		for diff := range diffsCh {
			diffs = append(diffs, diff...)
		}
	}

	return diffs, nil
}
