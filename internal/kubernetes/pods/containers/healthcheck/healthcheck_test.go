package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func initCtx() context.Context {
	var (
		ctx = context.Background()
	)
	err := logging.ConfigureForTests()
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	return ctx
}
func newCtxWithCleanStorage(ctx context.Context) context.Context {
	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	return ctx
}

func initContainersForTest1() (corev1.Container, corev1.Container) {
	container1 := corev1.Container{
		LivenessProbe: &corev1.Probe{},
	}

	container2 := corev1.Container{}
	return container1, container2
}

func initContainersForTest2() (corev1.Container, corev1.Container) {
	container1 := corev1.Container{
		ReadinessProbe: &corev1.Probe{},
	}

	container2 := corev1.Container{}
	return container1, container2
}

func TestCompare(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	container1, container2 := initContainersForTest1()

	Compare(ctx, container1, container2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckLivenessProbeDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckLivenessProbeDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckLivenessProbeDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	container1, container2 = initContainersForTest2()

	Compare(ctx, container1, container2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckReadinessProbeDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckReadinessProbeDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckReadinessProbeDifferent.Error())
	}
}
