package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func initEmptyDirForTest1() (*v1.EmptyDirVolumeSource, *v1.EmptyDirVolumeSource) {

	emptyDir1 := &v1.EmptyDirVolumeSource{
		Medium: "medium",
	}

	emptyDir2 := &v1.EmptyDirVolumeSource{
		Medium: "badMedium",
	}

	return emptyDir1, emptyDir2
}

func initEmptyDirForTest2() (*v1.EmptyDirVolumeSource, *v1.EmptyDirVolumeSource) {

	emptyDir1 := &v1.EmptyDirVolumeSource{
		SizeLimit: &resource.Quantity{
			Format: "format",
		},
	}

	emptyDir2 := &v1.EmptyDirVolumeSource{}

	return emptyDir1, emptyDir2
}

func initEmptyDirForTest3() (*v1.EmptyDirVolumeSource, *v1.EmptyDirVolumeSource) {

	emptyDir1 := &v1.EmptyDirVolumeSource{
		SizeLimit: &resource.Quantity{
			Format: "format",
		},
	}

	emptyDir2 := &v1.EmptyDirVolumeSource{
		SizeLimit: &resource.Quantity{
			Format: "badFormat",
		},
	}

	return emptyDir1, emptyDir2
}

func TestCompareVolumeEmptyDir(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	emptyDir1, emptyDir2 := initEmptyDirForTest1()

	CompareVolumeEmptyDir(ctx, emptyDir1, emptyDir2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeEmptyDirMedium) {
				t.Errorf("Error expected: '%s. 'medium' vs 'badMedium''. But it was returned: %s", ErrorVolumeEmptyDirMedium.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'medium' vs 'badMedium''. But the function found no errors", ErrorVolumeEmptyDirMedium.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	emptyDir1, emptyDir2 = initEmptyDirForTest2()

	CompareVolumeEmptyDir(ctx, emptyDir1, emptyDir2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeEmptyDirSizeLimit) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorVolumeEmptyDirSizeLimit.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorVolumeEmptyDirSizeLimit.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	emptyDir1, emptyDir2 = initEmptyDirForTest3()

	CompareVolumeEmptyDir(ctx, emptyDir1, emptyDir2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeEmptyDirSizeLimitFormat) {
				t.Errorf("Error expected: '%s. 'format' vs 'badFormat''. But it was returned: %s", ErrorVolumeEmptyDirSizeLimitFormat.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'format' vs 'badFormat''. But the function found no errors", ErrorVolumeEmptyDirSizeLimitFormat.Error())
	}
}
