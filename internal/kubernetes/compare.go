package kubernetes

import (
	"context"
	"sync"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"
	"k8s-cluster-comparator/internal/kubernetes/networking/ingress"
	"k8s-cluster-comparator/internal/kubernetes/networking/service"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/daemonset"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/deployment"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/statefulset"
	"k8s-cluster-comparator/internal/kubernetes/tasks/cronjob"
	"k8s-cluster-comparator/internal/kubernetes/tasks/job"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	"go.uber.org/zap"
)

type ResStr struct {
	IsClustersDiffer bool
	Err              error
}

//func compareKubeNamespaces(ctx context.Context, ns string) (*types.ObjectsDiff, error) {
//	log := logging.FromContext(ctx).With(zap.String("namespace", ns))
//
//	log.Debugf("Processing namespace/%s", ns)
//	defer func() {
//		log.Debugf("End of namespace/%s processing", ns)
//	}()
//
//	return nil, nil
//}

// CompareClusters main compare function, runs functions for comparing clusters by different parameters one at a time: Deployments, StatefulSets, DaemonSets, ConfigMaps
func CompareClusters(ctx context.Context) (*diff.DiffsStorage, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		diffs = diff.FromContext(ctx)
	)

	for _, namespace := range cfg.Connections.Namespaces {
		log := log.With(zap.String("namespace", namespace))

		comparators := []types.KubeResourceComparator{
			deployment.NewComparator(ctx, namespace),
			statefulset.NewComparator(ctx, namespace),
			daemonset.NewComparator(ctx, namespace),

			job.NewComparator(ctx, namespace),
			cronjob.NewComparator(ctx, namespace),

			configmap.NewComparator(ctx, namespace),
			secret.NewComparator(ctx, namespace),

			service.NewComparator(ctx, namespace),
			ingress.NewComparator(ctx, namespace),
		}

		wg := &sync.WaitGroup{}

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

					_, err := cmp.Compare(ctx)
					if err != nil {
						log.Errorf("cannot call %T: %s", cmp, err.Error())
						return
					}
				}(wg)
			}
		}

		wg.Wait()
	}

	return diffs, nil
}
