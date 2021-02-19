package kubernetes

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/jobs"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/kubernetes/networking"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers"
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
		//log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		diffs = make([]types.KubeObjectsDifference, 0)
	)

	for _, namespace := range cfg.Connections.Namespaces {
		comparators := []types.KubeResourceComparator{
			pod_controllers.NewDeploymentsComparator(ctx, namespace),
			pod_controllers.NewStatefulSetsComparator(ctx, namespace),
			pod_controllers.NewDaemonSetsComparator(ctx, namespace),

			jobs.NewJobsComparator(ctx, namespace),
			jobs.NewCronJobsComparator(ctx, namespace),

			kv_maps.NewConfigMapsComparator(ctx, namespace),
			kv_maps.NewSecretsComparator(ctx, namespace),

			networking.NewServicesComparator(ctx, namespace),
			networking.NewIngressesComparator(ctx, namespace),
		}

		diffs := make([]types.KubeObjectsDifference, len(comparators))
		for _, cmp := range comparators {
			diff, err := cmp.Compare(ctx, namespace)
			if err != nil {
				return nil, fmt.Errorf("cannot call %t: %w", cmp, err)
			}

			diffs = append(diffs, diff...)
		}
	}

	return diffs, nil
}
