package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func initDownwardAPIForTest1() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {
	var defMode1 int32
	var defMode2 int32

	defMode1 = 1
	defMode2 = 2

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		DefaultMode: &defMode1,
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		DefaultMode: &defMode2,
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest2() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {
	var defMode1 int32

	defMode1 = 1

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		DefaultMode: &defMode1,
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest3() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{}, {},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest4() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				Path: "path",
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				Path: "diffPath",
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest5() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	var mode1 int32
	var mode2 int32

	mode1 = 1
	mode2 = 2

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				Mode: &mode1,
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				Mode: &mode2,
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest6() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	var mode1 int32
	mode1 = 1

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				Mode: &mode1,
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest7() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					Resource: "resource",
				},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					Resource: "diffResource",
				},
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest8() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					ContainerName: "name",
				},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					ContainerName: "diffName",
				},
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest9() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					Divisor: resource.Quantity{
						Format: "format",
					},
				},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					Divisor: resource.Quantity{
						Format: "diffFormat",
					},
				},
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest10() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				ResourceFieldRef: &v1.ResourceFieldSelector{},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest11() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				FieldRef: &v1.ObjectFieldSelector{
					APIVersion: "v1",
				},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				FieldRef: &v1.ObjectFieldSelector{
					APIVersion: "v2",
				},
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest12() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "fieldPath",
				},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "diffFieldPath",
				},
			},
		},
	}
	return downwardApi1, downwardApi2
}

func initDownwardAPIForTest13() (*v1.DownwardAPIVolumeSource, *v1.DownwardAPIVolumeSource) {

	downwardApi1 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{
				FieldRef: &v1.ObjectFieldSelector{},
			},
		},
	}
	downwardApi2 := &v1.DownwardAPIVolumeSource{
		Items: []v1.DownwardAPIVolumeFile{
			{},
		},
	}
	return downwardApi1, downwardApi2
}

func TestCompareVolumeDownwardAPI(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 := initDownwardAPIForTest1()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIDefaultMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorDownwardAPIDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorDownwardAPIDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest2()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingDownwardAPIDefaultMode) {
				t.Errorf("Error expected: '%s.'. But it was returned: %s", ErrorMissingDownwardAPIDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s.'. But the function found no errors", ErrorMissingDownwardAPIDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest3()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeDownwardAPIItemsLen) {
				t.Errorf("Error expected: '%s. '2' vs '1''. But it was returned: %s", ErrorVolumeDownwardAPIItemsLen.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorVolumeDownwardAPIItemsLen.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest4()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsPath) {
				t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But it was returned: %s", ErrorDownwardAPIItemsPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But the function found no errors", ErrorDownwardAPIItemsPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest5()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorDownwardAPIItemsMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorDownwardAPIItemsMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest6()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingDownwardAPIItemsMode) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingDownwardAPIItemsMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingDownwardAPIItemsMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest7()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsResFieldRefResource) {
				t.Errorf("Error expected: '%s. 'resource' vs 'diffResource''. But it was returned: %s", ErrorDownwardAPIItemsResFieldRefResource.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'resource' vs 'diffResource''. But the function found no errors", ErrorDownwardAPIItemsResFieldRefResource.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest8()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsResFieldRefContainerName) {
				t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But it was returned: %s", ErrorDownwardAPIItemsResFieldRefContainerName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But the function found no errors", ErrorDownwardAPIItemsResFieldRefContainerName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest9()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsResFieldRefFormat) {
				t.Errorf("Error expected: '%s. 'format' vs 'diffFormat''. But it was returned: %s", ErrorDownwardAPIItemsResFieldRefFormat.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'format' vs 'diffFormat''. But the function found no errors", ErrorDownwardAPIItemsResFieldRefFormat.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest10()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingDownwardAPIItemsResourceFieldRef) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingDownwardAPIItemsResourceFieldRef.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingDownwardAPIItemsResourceFieldRef.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest11()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsFieldRefAPIVersion) {
				t.Errorf("Error expected: '%s. 'v1' vs 'v2''. But it was returned: %s", ErrorDownwardAPIItemsFieldRefAPIVersion.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'v1' vs 'v2''. But the function found no errors", ErrorDownwardAPIItemsFieldRefAPIVersion.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest12()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDownwardAPIItemsFieldRefFieldPath) {
				t.Errorf("Error expected: '%s. 'fieldPath' vs 'diffFieldPath''. But it was returned: %s", ErrorDownwardAPIItemsFieldRefFieldPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'fieldPath' vs 'diffFieldPath''. But the function found no errors", ErrorDownwardAPIItemsFieldRefFieldPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	downwardApi1, downwardApi2 = initDownwardAPIForTest13()

	CompareVolumeDownwardAPI(ctx, downwardApi1, downwardApi2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingDownwardAPIItemsFieldRef) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingDownwardAPIItemsFieldRef.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingDownwardAPIItemsFieldRef.Error())
	}
}
