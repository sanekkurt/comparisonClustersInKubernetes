package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	"k8s-cluster-comparator/internal/utils"
	v1 "k8s.io/api/core/v1"
)

func compareExecProbes(ctx context.Context, probe1, probe2 v1.Probe) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)
	)

	if probe1.Exec != nil && probe2.Exec != nil {
		if bDiff, diff := utils.AreStringListsEqual(ctx, probe1.Exec.Command, probe2.Exec.Command); !bDiff {
			logging.DiffLog(log, ErrorContainerHealthCheckDifferentExecCommand, diff)
		}
	} else if probe1.Exec != nil || probe2.Exec != nil {
		logging.DiffLog(log, ErrorContainerHealthCheckDifferentExec, probe1.Exec, probe2.Exec)
	}

	return nil, nil
}