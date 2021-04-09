package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initCSIVolumeSourceForTest1() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	csi1 := &v1.CSIVolumeSource{
		Driver: "driver",
	}
	csi2 := &v1.CSIVolumeSource{
		Driver: "diffDriver",
	}
	return csi1, csi2
}

func initCSIVolumeSourceForTest2() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	var readOnly1 bool
	var readOnly2 bool
	readOnly2 = true

	csi1 := &v1.CSIVolumeSource{
		ReadOnly: &readOnly1,
	}
	csi2 := &v1.CSIVolumeSource{
		ReadOnly: &readOnly2,
	}
	return csi1, csi2
}

func initCSIVolumeSourceForTest3() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	var readOnly1 bool

	csi1 := &v1.CSIVolumeSource{
		ReadOnly: &readOnly1,
	}
	csi2 := &v1.CSIVolumeSource{}
	return csi1, csi2
}

func initCSIVolumeSourceForTest4() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	var fstype1 string
	var fstype2 string

	fstype1 = "fstype1"
	fstype2 = "fstype2"

	csi1 := &v1.CSIVolumeSource{
		FSType: &fstype1,
	}
	csi2 := &v1.CSIVolumeSource{
		FSType: &fstype2,
	}
	return csi1, csi2
}

func initCSIVolumeSourceForTest5() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	var fstype1 string

	csi1 := &v1.CSIVolumeSource{
		FSType: &fstype1,
	}
	csi2 := &v1.CSIVolumeSource{}
	return csi1, csi2
}

func initCSIVolumeSourceForTest6() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	csi1 := &v1.CSIVolumeSource{
		NodePublishSecretRef: &v1.LocalObjectReference{
			Name: "name",
		},
	}
	csi2 := &v1.CSIVolumeSource{
		NodePublishSecretRef: &v1.LocalObjectReference{
			Name: "diffName",
		},
	}

	return csi1, csi2

}

func initCSIVolumeSourceForTest7() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	csi1 := &v1.CSIVolumeSource{}
	csi2 := &v1.CSIVolumeSource{
		NodePublishSecretRef: &v1.LocalObjectReference{
			Name: "diffName",
		},
	}

	return csi1, csi2

}

func initCSIVolumeSourceForTest8() (*v1.CSIVolumeSource, *v1.CSIVolumeSource) {

	map1 := make(map[string]string)
	map2 := make(map[string]string)

	map1["1"] = "1"
	map2["1"] = "1"
	map2["2"] = "2"

	csi1 := &v1.CSIVolumeSource{
		VolumeAttributes: map1,
	}
	csi2 := &v1.CSIVolumeSource{
		VolumeAttributes: map2,
	}

	return csi1, csi2

}

func TestCompareVolumeCSI(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 := initCSIVolumeSourceForTest1()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeCSIDriver) {
				t.Errorf("Error expected: '%s. 'driver' vs 'diffDriver''. But it was returned: %s", ErrorVolumeCSIDriver.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'driver' vs 'diffDriver''. But the function found no errors", ErrorVolumeCSIDriver.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest2()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeCSIReadOnly) {
				t.Errorf("Error expected: '%s. 'false' vs 'true''. But it was returned: %s", ErrorVolumeCSIReadOnly.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'false' vs 'true''. But the function found no errors", ErrorVolumeCSIReadOnly.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest3()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingCSIReadOnly) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingCSIReadOnly.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingCSIReadOnly.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest4()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeCSIFSType) {
				t.Errorf("Error expected: '%s. 'fstype1' vs 'fstype2''. But it was returned: %s", ErrorVolumeCSIFSType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'fstype1' vs 'fstype2''. But the function found no errors", ErrorVolumeCSIFSType.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest5()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingCSIFSType) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingCSIFSType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingCSIFSType.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest6()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeCSIName) {
				t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But it was returned: %s", ErrorVolumeCSIName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But the function found no errors", ErrorVolumeCSIName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest7()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingCSINodePublishSecretRef) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingCSINodePublishSecretRef.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingCSINodePublishSecretRef.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	csi1, csi2 = initCSIVolumeSourceForTest8()

	CompareVolumeCSI(ctx, csi1, csi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeCSIVolumeAttributes) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorVolumeCSIVolumeAttributes.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorVolumeCSIVolumeAttributes.Error())
	}
}
