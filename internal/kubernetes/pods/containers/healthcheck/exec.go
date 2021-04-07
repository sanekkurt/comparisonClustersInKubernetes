package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/diff"

	"k8s-cluster-comparator/internal/utils"
	v1 "k8s.io/api/core/v1"
)

func compareExecProbes(ctx context.Context, probe1, probe2 v1.Probe) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if probe1.Exec != nil && probe2.Exec != nil {
		err := utils.AreStringListsEqual(ctx, probe1.Exec.Command, probe2.Exec.Command)
		if err != nil {
			diffsBatch.Add(ctx, false, "%w: '%s'", ErrorContainerHealthCheckDifferentExecCommand, err)
		}
		//if bDiff, diff := ; !bDiff {
		//	//logging.DiffLog(log, ErrorContainerHealthCheckDifferentExecCommand, dif)
		//	//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s", ErrorContainerHealthCheckDifferentExecCommand.Error(), dif)
		//	diffsBatch.Add(ctx, false, "%s: %s", ErrorContainerHealthCheckDifferentExecCommand.Error(), diff)
		//}
	} else if probe1.Exec != nil || probe2.Exec != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentExec, probe1.Exec, probe2.Exec)
		//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentExec, probe1.Exec, probe2.Exec)
		diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckDifferentExec)
	}
}
