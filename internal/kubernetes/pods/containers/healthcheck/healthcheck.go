package healthcheck

import (
	"context"
	"reflect"
	"runtime"

	"k8s-cluster-comparator/internal/kubernetes/types"
	corev1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

func Compare(ctx context.Context, container1, container2 corev1.Container) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.ObjectsDiff, 0)
	)

	if container1.LivenessProbe != nil && container2.LivenessProbe != nil {
		log.Debugf("Compare liveness probes: started")

		diff, err := compareContainerProbes(ctx, *container1.LivenessProbe, *container2.LivenessProbe)
		if err != nil {
			return nil, err
		}

		diffs = append(diffs, diff...)
	} else if container1.LivenessProbe != nil || container2.LivenessProbe != nil {
		log.Debugf("ComparePodSpecs: start checking LivenessProbe in container - %s", container1.Name)
		logging.DiffLog(log, ErrorContainerHealthCheckLivenessProbeDifferent, "One of the containers has no liveness probe defined", *container1.LivenessProbe, *container2.LivenessProbe)
	}

	if container1.ReadinessProbe != nil && container2.ReadinessProbe != nil {
		log.Debugf("Compare readiness probes: started")

		diff, err := compareContainerProbes(ctx, *container1.ReadinessProbe, *container2.ReadinessProbe)
		if err != nil {
			return nil, err
		}

		diffs = append(diffs, diff...)
	} else if container1.ReadinessProbe != nil || container2.ReadinessProbe != nil {
		log.Debugf("ComparePodSpecs: start checking LivenessProbe in container - %s", container1.Name)
		logging.DiffLog(log, ErrorContainerHealthCheckLivenessProbeDifferent, "One of the containers has no readiness probe defined", *container1.ReadinessProbe, *container2.ReadinessProbe)
	}

	return diffs, nil
}

func compareContainerProbes(ctx context.Context, probe1, probe2 corev1.Probe) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)
		diffs = make([]types.ObjectsDiff, 0)
	)

	comparisons := []func(context.Context, corev1.Probe, corev1.Probe)([]types.ObjectsDiff, error){
		compareCommonProbeParams,
		compareTCPSocketProbes,
		compareHTTPGetProbes,
		compareExecProbes,
	}

	for _, cmp := range comparisons {
		log.Debugf("%s started", runtime.FuncForPC(reflect.ValueOf(cmp).Pointer()).Name())

		diff, err := cmp(ctx, probe1, probe2)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, diff...)
	}

	return diffs, nil
}
