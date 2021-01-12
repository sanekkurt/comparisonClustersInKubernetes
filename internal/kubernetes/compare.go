package kubernetes

import (
	"context"
	"sync"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/jobs"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/kubernetes/networking"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

// CompareClusters main compare function, runs functions for comparing clusters by different parameters one at a time: Deployments, StatefulSets, DaemonSets, ConfigMaps
func CompareClusters(ctx context.Context, cfg *config.AppConfig) (bool, error) {
	type ResStr struct {
		IsClustersDiffer bool
		Err              error
	}

	var (
		wg = &sync.WaitGroup{}

		resCh = make(chan ResStr, len(cfg.Namespaces))
	)

	for _, namespace := range cfg.Namespaces {
		wg.Add(1)

		go func(wg *sync.WaitGroup, resCh chan ResStr, namespace string) {
			var (
				isClustersDifferFlag types.OnceSettableFlag

				kubeConns = &types.KubeConnections{
					C1:        cfg.Cluster1.Kubeconfig,
					C2:        cfg.Cluster2.Kubeconfig,
					Namespace: namespace,
				}
			)

			defer func() {
				wg.Done()
			}()

			ctx = config.With(ctx, kubeConns)

			isClustersDiffer, err := pod_controllers.CompareDeployments(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = pod_controllers.CompareStateFulSets(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = pod_controllers.CompareDaemonSets(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = kv_maps.CompareConfigMaps(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = kv_maps.CompareSecrets(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = jobs.CompareJobs(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = jobs.CompareCronJobs(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = networking.CompareServices(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = networking.CompareIngresses(ctx, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			resCh <- ResStr{
				Err:              nil,
				IsClustersDiffer: isClustersDifferFlag.GetFlag(),
			}
		}(wg, resCh, namespace)
	}

	wg.Wait()

	close(resCh)

	for res := range resCh {
		if res.Err != nil {
			return false, res.Err
		}
		if res.IsClustersDiffer {
			return res.IsClustersDiffer, nil
		}
	}

	return false, nil
}
