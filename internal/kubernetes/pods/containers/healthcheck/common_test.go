package healthcheck

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initProbesForTestCommon1() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		FailureThreshold: 1,
	}

	probe2 := v1.Probe{
		FailureThreshold: 2,
	}

	return probe1, probe2
}

func initProbesForTestCommon2() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		InitialDelaySeconds: 1,
	}

	probe2 := v1.Probe{
		InitialDelaySeconds: 2,
	}

	return probe1, probe2
}

func initProbesForTestCommon3() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		PeriodSeconds: 1,
	}

	probe2 := v1.Probe{
		PeriodSeconds: 2,
	}

	return probe1, probe2
}

func initProbesForTestCommon4() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		SuccessThreshold: 1,
	}

	probe2 := v1.Probe{
		SuccessThreshold: 2,
	}

	return probe1, probe2
}

func initProbesForTestCommon5() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		TimeoutSeconds: 1,
	}

	probe2 := v1.Probe{
		TimeoutSeconds: 2,
	}

	return probe1, probe2
}

func TestCompareCommonProbeParams(t *testing.T) {

	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 := initProbesForTestCommon1()

	compareCommonProbeParams(ctx, probe1, probe2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentFailureThreshold) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentFailureThreshold.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2'. But the function found no errors", ErrorContainerHealthCheckDifferentFailureThreshold.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestCommon2()

	compareCommonProbeParams(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentInitialDelaySeconds) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentInitialDelaySeconds.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2'. But the function found no errors", ErrorContainerHealthCheckDifferentInitialDelaySeconds.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestCommon3()

	compareCommonProbeParams(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentPeriodSeconds) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentPeriodSeconds.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2'. But the function found no errors", ErrorContainerHealthCheckDifferentPeriodSeconds.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestCommon4()

	compareCommonProbeParams(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentSuccessThreshold) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentSuccessThreshold.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2'. But the function found no errors", ErrorContainerHealthCheckDifferentSuccessThreshold.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestCommon5()

	compareCommonProbeParams(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTimeoutSeconds) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentTimeoutSeconds.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2'. But the function found no errors", ErrorContainerHealthCheckDifferentTimeoutSeconds.Error())
	}
}
