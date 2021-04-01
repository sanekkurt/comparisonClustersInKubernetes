package healthcheck

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"

	"k8s-cluster-comparator/internal/utils"
	v1 "k8s.io/api/core/v1"
)

func compareExecProbes(ctx context.Context, probe1, probe2 v1.Probe) {
	var (
		//log = logging.FromContext(ctx)
		diffsBatch = diff.DiffBatchFromContext(ctx)
	)

	if probe1.Exec != nil && probe2.Exec != nil {
		if bDiff, dif := utils.AreStringListsEqual(ctx, probe1.Exec.Command, probe2.Exec.Command); !bDiff {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentExecCommand, dif)
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s", ErrorContainerHealthCheckDifferentExecCommand.Error(), dif)
		}
	} else if probe1.Exec != nil || probe2.Exec != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentExec, probe1.Exec, probe2.Exec)
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentExec, probe1.Exec, probe2.Exec)
	}

}
