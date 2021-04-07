package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initClaimVolumeSourceForTest1() (*v1.PersistentVolumeClaimVolumeSource, *v1.PersistentVolumeClaimVolumeSource) {

	claim1 := &v1.PersistentVolumeClaimVolumeSource{
		ReadOnly: false,
	}
	claim2 := &v1.PersistentVolumeClaimVolumeSource{
		ReadOnly: true,
	}
	return claim1, claim2
}

func initClaimVolumeSourceForTest2() (*v1.PersistentVolumeClaimVolumeSource, *v1.PersistentVolumeClaimVolumeSource) {

	claim1 := &v1.PersistentVolumeClaimVolumeSource{
		ClaimName: "claimName",
	}
	claim2 := &v1.PersistentVolumeClaimVolumeSource{
		ClaimName: "diffClaimName",
	}
	return claim1, claim2
}

func TestCompareVolumePersistentVolumeClaim(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	claim1, claim2 := initClaimVolumeSourceForTest1()

	CompareVolumePersistentVolumeClaim(ctx, claim1, claim2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPersistentVolumeClaimReadOnly) {
				t.Errorf("Error expected: '%s. 'false' vs 'true''. But it was returned: %s", ErrorPersistentVolumeClaimReadOnly.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'false' vs 'true''. But the function found no errors", ErrorPersistentVolumeClaimReadOnly.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	claim1, claim2 = initClaimVolumeSourceForTest2()

	CompareVolumePersistentVolumeClaim(ctx, claim1, claim2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPersistentVolumeClaimName) {
				t.Errorf("Error expected: '%s. 'claimName' vs 'diffClaimName''. But it was returned: %s", ErrorPersistentVolumeClaimName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'claimName' vs 'diffClaimName''. But the function found no errors", ErrorPersistentVolumeClaimName.Error())
	}
}
