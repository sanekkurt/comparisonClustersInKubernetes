package volumes

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

func initHostPathForTest1() (*v1.HostPathVolumeSource, *v1.HostPathVolumeSource) {
	hostPath1 := &v1.HostPathVolumeSource{
		Path: "path",
	}
	hostPath2 := &v1.HostPathVolumeSource{
		Path: "badPath",
	}

	return hostPath1, hostPath2
}

func initHostPathForTest2() (*v1.HostPathVolumeSource, *v1.HostPathVolumeSource) {
	var type1 v1.HostPathType
	type1 = "type1"

	var type2 v1.HostPathType
	type2 = "badType"

	hostPath1 := &v1.HostPathVolumeSource{
		Type: &type1,
	}
	hostPath2 := &v1.HostPathVolumeSource{
		Type: &type2,
	}

	return hostPath1, hostPath2
}

func initHostPathForTest3() (*v1.HostPathVolumeSource, *v1.HostPathVolumeSource) {

	var type2 v1.HostPathType
	type2 = "badType"

	hostPath1 := &v1.HostPathVolumeSource{}
	hostPath2 := &v1.HostPathVolumeSource{
		Type: &type2,
	}

	return hostPath1, hostPath2
}

func TestCompareVolumeHostPath(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2 := initHostPathForTest1()

	CompareVolumeHostPath(ctx, hostPath1, hostPath2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHostPath) {
				t.Errorf("Error expected: '%s: 'path' vs 'badPath''. But it was returned: %s", ErrorHostPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'path' vs 'badPath''. But the function found no errors", ErrorHostPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2 = initHostPathForTest2()

	CompareVolumeHostPath(ctx, hostPath1, hostPath2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHostPathType) {
				t.Errorf("Error expected: '%s: 'type1' vs 'badType''. But it was returned: %s", ErrorHostPathType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'type1' vs 'badType''. But the function found no errors", ErrorHostPathType.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	hostPath1, hostPath2 = initHostPathForTest3()

	CompareVolumeHostPath(ctx, hostPath1, hostPath2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingHostPathType) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingHostPathType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingHostPathType.Error())
	}
}
