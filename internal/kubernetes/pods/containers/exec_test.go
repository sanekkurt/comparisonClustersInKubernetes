package containers

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
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

func TestCompareContainerExecParams(t *testing.T) {
	ctx := initCtx()
	ctx1 := newCtxWithCleanStorage(ctx)

	container1 := v1.Container{Name: "one", Command: []string{"comm1", "comm3"}}
	container2 := v1.Container{Name: "two", Command: []string{"comm1", "comm2"}}

	compareContainerExecParams(ctx1, container1, container2)

	diffStorage := diff.StorageFromContext(ctx1)
	diffStorage.Finalize(ctx1)
	batch := diff.BatchFromContext(ctx1)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerCommandsDifferent) {
				t.Errorf("Error expected: '%s. container 'one': different values in string lists at position #2: comm3 vs comm2'. But it was returned: %s", ErrorContainerCommandsDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. container 'one': different values in string lists at position #2: comm3 vs comm2'. But the function found no errors", ErrorContainerCommandsDifferent.Error())
	}

	ctx2 := newCtxWithCleanStorage(ctx)
	container1 = v1.Container{Name: "one", Args: []string{"arg1", "arg2"}}
	container2 = v1.Container{Name: "two", Args: []string{"arg8", "arg2"}}

	compareContainerExecParams(ctx2, container1, container2)

	diffStorage = diff.StorageFromContext(ctx2)
	diffStorage.Finalize(ctx2)
	batch = diff.BatchFromContext(ctx2)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerArgumentsDifferent) {
				t.Errorf("Error expected: '%s. container 'one': different values in string lists at position #1: arg1 vs arg8'. But it was returned: %s", ErrorContainerArgumentsDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. container 'one': different values in string lists at position #1: arg1 vs arg8'. But the function found no errors", ErrorContainerArgumentsDifferent.Error())
	}

}
