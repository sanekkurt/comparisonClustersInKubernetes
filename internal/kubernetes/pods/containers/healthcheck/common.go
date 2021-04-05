package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func compareCommonProbeParams(ctx context.Context, probe1, probe2 v1.Probe) {
	var (
		//log = logging.FromContext(ctx)

		diffsBatch = diff.BatchFromContext(ctx)

		//diffs = make([]types.objectsDiff, 0)

	)

	if probe1.FailureThreshold != probe2.FailureThreshold {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentFailureThreshold, probe1.FailureThreshold, probe2.FailureThreshold)
		//diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentFailureThreshold.Error(), probe1.FailureThreshold, probe2.FailureThreshold)
		diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentFailureThreshold.Error(), probe1.FailureThreshold, probe2.FailureThreshold)
	}

	if probe1.InitialDelaySeconds != probe2.InitialDelaySeconds {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentInitialDelaySeconds, probe1.InitialDelaySeconds, probe2.InitialDelaySeconds)
		//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentInitialDelaySeconds.Error(), probe1.InitialDelaySeconds, probe2.InitialDelaySeconds)
		diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentInitialDelaySeconds.Error(), probe1.InitialDelaySeconds, probe2.InitialDelaySeconds)
	}

	if probe1.PeriodSeconds != probe2.PeriodSeconds {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentPeriodSeconds, probe1.PeriodSeconds, probe2.PeriodSeconds)
		//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentPeriodSeconds.Error(), probe1.PeriodSeconds, probe2.PeriodSeconds)
		diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentPeriodSeconds.Error(), probe1.PeriodSeconds, probe2.PeriodSeconds)
	}

	if probe1.SuccessThreshold != probe2.SuccessThreshold {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentSuccessThreshold, probe1.SuccessThreshold, probe2.SuccessThreshold)
		//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentSuccessThreshold.Error(), probe1.SuccessThreshold, probe2.SuccessThreshold)
		diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentSuccessThreshold.Error(), probe1.SuccessThreshold, probe2.SuccessThreshold)
	}

	if probe1.TimeoutSeconds != probe2.TimeoutSeconds {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTimeoutSeconds, probe1.TimeoutSeconds, probe2.TimeoutSeconds)
		//diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentTimeoutSeconds.Error(), probe1.TimeoutSeconds, probe2.TimeoutSeconds)
		diffsBatch.Add(ctx, false, "%s: %s vs %s", ErrorContainerHealthCheckDifferentTimeoutSeconds.Error(), probe1.TimeoutSeconds, probe2.TimeoutSeconds)
	}
}
