package healthcheck

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initProbesForTestExec() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"comm1", "comm2"},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{},
	}

	return probe1, probe2
}

func TestCompareExecProbes(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 := initProbesForTestExec()

	compareExecProbes(ctx, probe1, probe2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentExec) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckDifferentExec.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckDifferentExec.Error())
	}
}
