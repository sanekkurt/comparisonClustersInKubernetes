package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initConfigMapSourceForTest1() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	configMap1 := &v1.ConfigMapVolumeSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "name",
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "diffName",
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest2() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	configMap1 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{},
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{}, {},
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest3() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	configMap1 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Path: "path",
			},
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Path: "diffPath",
			},
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest4() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	configMap1 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "key",
			},
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "diffKey",
			},
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest5() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var mode1 int32
	var mode2 int32
	mode1 = 1
	mode2 = 2

	configMap1 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode1,
			},
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode2,
			},
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest6() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var mode1 int32
	mode1 = 1

	configMap1 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode1,
			},
		},
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Items: []v1.KeyToPath{
			{},
		},
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest7() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var defMode1 int32
	var defMode2 int32
	defMode1 = 1
	defMode2 = 2

	configMap1 := &v1.ConfigMapVolumeSource{
		DefaultMode: &defMode1,
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		DefaultMode: &defMode2,
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest8() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var defMode1 int32

	defMode1 = 1

	configMap1 := &v1.ConfigMapVolumeSource{
		DefaultMode: &defMode1,
	}
	configMap2 := &v1.ConfigMapVolumeSource{}
	return configMap1, configMap2
}

func initConfigMapSourceForTest9() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var optional1 bool
	var optional2 bool

	optional2 = true

	configMap1 := &v1.ConfigMapVolumeSource{
		Optional: &optional1,
	}
	configMap2 := &v1.ConfigMapVolumeSource{
		Optional: &optional2,
	}
	return configMap1, configMap2
}

func initConfigMapSourceForTest10() (*v1.ConfigMapVolumeSource, *v1.ConfigMapVolumeSource) {

	var optional1 bool

	configMap1 := &v1.ConfigMapVolumeSource{
		Optional: &optional1,
	}
	configMap2 := &v1.ConfigMapVolumeSource{}
	return configMap1, configMap2
}

func TestCompareVolumeConfigMap(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 := initConfigMapSourceForTest1()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapName) {
				t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But it was returned: %s", ErrorVolumeConfigMapName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But the function found no errors", ErrorVolumeConfigMapName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest2()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapItemsLen) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorVolumeConfigMapItemsLen.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorVolumeConfigMapItemsLen.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest3()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapPath) {
				t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But it was returned: %s", ErrorVolumeConfigMapPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But the function found no errors", ErrorVolumeConfigMapPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest4()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapKey) {
				t.Errorf("Error expected: '%s. 'key' vs 'diffKey''. But it was returned: %s", ErrorVolumeConfigMapKey.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'key' vs 'diffKey''. But the function found no errors", ErrorVolumeConfigMapKey.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest5()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorVolumeConfigMapMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorVolumeConfigMapMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest6()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeConfigMapMode) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingVolumeConfigMapMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingVolumeConfigMapMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest7()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapDefaultMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorVolumeConfigMapDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorVolumeConfigMapDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest8()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeConfigMapDefaultMode) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingVolumeConfigMapDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingVolumeConfigMapDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest9()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeConfigMapOptional) {
				t.Errorf("Error expected: '%s. 'false' vs 'true''. But it was returned: %s", ErrorVolumeConfigMapOptional.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'false' vs 'true''. But the function found no errors", ErrorVolumeConfigMapOptional.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	configMap1, configMap2 = initConfigMapSourceForTest10()

	CompareVolumeConfigMap(ctx, configMap1, configMap2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeConfigMapOptional) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingVolumeConfigMapOptional.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingVolumeConfigMapOptional.Error())
	}
}
