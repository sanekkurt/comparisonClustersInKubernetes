package deployment

import (
	"context"
	"fmt"
	"sort"
	"time"

	"k8s-cluster-comparator/internal/config"
	pccommon "k8s-cluster-comparator/internal/kubernetes/pod_controllers/common"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/replicaset"
	"k8s-cluster-comparator/internal/logging"
	appsv1 "k8s.io/api/apps/v1"
)

func (cmp *Comparator) getFilteredDeploymentsListByUpdateTime(ctx context.Context, apcs []map[string]*pccommon.AbstractPodController) ([]map[string]*pccommon.AbstractPodController, error) {
	var (
		cfg = config.FromContext(ctx)
		log = logging.FromContext(ctx)

		rsComparator = replicaset.NewComparator(ctx, cmp.Namespace)

		limitTime = time.Now().Add(-1 * cfg.Workloads.PodControllers.Deployments.MinimumUpdateAgeMinutes)
	)

	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
		for apcName, apc := range apcs[0] {
			rsMaps, err := rsComparator.WithLabelSelector(apc.PodLabelSelector).Collect(ctx)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve ReplicaSets by selector: %s", err.Error())
			}

			for idx := 0; idx < len(apcs)-1; idx++ {
				if len(rsMaps[idx]) == 0 {
					continue
				}

				rsList := make([]appsv1.ReplicaSet, 0, len(rsMaps[idx]))

				for _, rs := range rsMaps[idx] {
					rsList = append(rsList, rs)
				}

				sort.SliceStable(rsList, func(i, j int) bool {
					return rsList[i].CreationTimestamp.Time.After(rsList[j].CreationTimestamp.Time)
				})

				if rsList[0].CreationTimestamp.Time.After(limitTime) {
					log.With("objectName", apcName).Infof("deployment/%s skipped from comparison: ReplicaSet '%s' created at %s, less than %d minutes ago", apcName, rsList[0].Name, rsList[0].CreationTimestamp.Time, cfg.Workloads.PodControllers.Deployments.MinimumUpdateAgeMinutes/time.Minute)
					delete(apcs[idx], apcName)
				}
			}
		}
	}

	return apcs, nil
}
