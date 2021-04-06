package kubernetes

import (
	"context"
	"sync"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"
	"k8s-cluster-comparator/internal/kubernetes/networking/ingress"
	"k8s-cluster-comparator/internal/kubernetes/networking/service"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/daemonset"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/statefulset"
	"k8s-cluster-comparator/internal/kubernetes/tasks/cronjob"
	"k8s-cluster-comparator/internal/kubernetes/tasks/job"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/deployment"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	"go.uber.org/zap"
)

type ResStr struct {
	IsClustersDiffer bool
	Err              error
}

func runComparatorsAsynchronouslyAndWaitForCompletion(ctx context.Context, comparators []types.KubeResourceComparator) {
	var (
		log = logging.FromContext(ctx)
		wg  = &sync.WaitGroup{}
	)

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

// CompareClusters main compare function, runs functions for comparing clusters by different parameters one at a time: Deployments, StatefulSets, DaemonSets, ConfigMaps
func CompareClusters(ctx context.Context) error {
	var (
		log   = logging.FromContext(ctx)
		cfg   = config.FromContext(ctx)
		diffs = diff.NewDiffsStorage(ctx)
	)

	ctx = diff.WithDiffStorage(ctx, diffs)

	for _, namespace := range cfg.Connections.Namespaces {
		ctx := logging.WithLogger(ctx, log.With(zap.String("namespace", namespace))) //nolint:govet

		comparators := []types.KubeResourceComparator{
			deployment.NewComparator(ctx, namespace),
			statefulset.NewComparator(ctx, namespace),
			daemonset.NewComparator(ctx, namespace),
		}

		runComparatorsAsynchronouslyAndWaitForCompletion(ctx, comparators)

		comparators = []types.KubeResourceComparator{
			job.NewComparator(ctx, namespace),
			cronjob.NewComparator(ctx, namespace),

			configmap.NewComparator(ctx, namespace),
			secret.NewComparator(ctx, namespace),

			service.NewComparator(ctx, namespace),
			ingress.NewComparator(ctx, namespace),
		}
		runComparatorsAsynchronouslyAndWaitForCompletion(ctx, comparators)
	}

	diffs.Finalize(ctx)

	return nil
}
