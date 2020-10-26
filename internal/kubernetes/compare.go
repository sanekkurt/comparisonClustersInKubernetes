package kubernetes

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/jobs"
	"sync"

	"k8s-cluster-comparator/internal/config"
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

		clientSet1 = cfg.Cluster1.Kubeconfig
		clientSet2 = cfg.Cluster2.Kubeconfig
	)

	if err := pod_controllers.Init(ctx); err != nil {
		return false, fmt.Errorf("cannot init pod_controllers package: %w", err)
	}
	if err := kv_maps.Init(ctx); err != nil {
		return false, fmt.Errorf("cannot init kv_maps package: %w", err)
	}
	if err := networking.Init(ctx); err != nil {
		return false, fmt.Errorf("cannot init networking package: %w", err)
	}
	if err := jobs.Init(ctx); err != nil {
		return false, fmt.Errorf("cannot init jobs package: %w", err)
	}

	for _, namespace := range cfg.Namespaces {
		wg.Add(1)

		go func(wg *sync.WaitGroup, resCh chan ResStr, namespace string) {
			var (
				isClustersDifferFlag types.OnceSettableFlag
			)

			defer func() {
				wg.Done()
			}()

			isClustersDiffer, err := pod_controllers.CompareDeployments(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = pod_controllers.CompareStateFulSets(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = pod_controllers.CompareDaemonSets(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = kv_maps.CompareConfigMaps(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = kv_maps.CompareSecrets(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = networking.CompareServices(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = networking.CompareIngresses(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = jobs.CompareJobs(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
			if err != nil {
				resCh <- ResStr{
					IsClustersDiffer: isClustersDiffer,
					Err:              err,
				}
				return
			}
			isClustersDifferFlag.SetFlag(isClustersDiffer)

			isClustersDiffer, err = jobs.CompareCronJobs(clientSet1, clientSet2, namespace, cfg.SkipEntitiesList)
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
