package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initSecretForTest1() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {

	secret1 := &v1.SecretVolumeSource{
		SecretName: "secretName",
	}

	secret2 := &v1.SecretVolumeSource{
		SecretName: "different secretName",
	}

	return secret1, secret2
}

func initSecretForTest2() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	flag := true
	secret1 := &v1.SecretVolumeSource{
		Optional: &flag,
	}

	secret2 := &v1.SecretVolumeSource{}

	return secret1, secret2
}

func initSecretForTest3() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	flagTrue := true
	flagFalse := false

	secret1 := &v1.SecretVolumeSource{
		Optional: &flagTrue,
	}

	secret2 := &v1.SecretVolumeSource{
		Optional: &flagFalse,
	}

	return secret1, secret2
}

func initSecretForTest4() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	var defMode1 int32
	var defMode2 int32
	defMode1 = 1
	defMode2 = 2
	secret1 := &v1.SecretVolumeSource{
		DefaultMode: &defMode1,
	}

	secret2 := &v1.SecretVolumeSource{
		DefaultMode: &defMode2,
	}

	return secret1, secret2
}

func initSecretForTest5() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	var defMode1 int32

	defMode1 = 1

	secret1 := &v1.SecretVolumeSource{
		DefaultMode: &defMode1,
	}

	secret2 := &v1.SecretVolumeSource{}

	return secret1, secret2
}

func initSecretForTest6() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {

	secret1 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "key1",
			},
			{
				Key: "key2",
			},
		},
	}

	secret2 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "key1",
			},
		},
	}

	return secret1, secret2
}

func initSecretForTest7() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {

	secret1 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "key",
			},
		},
	}

	secret2 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Key: "differentKey",
			},
		},
	}

	return secret1, secret2
}

func initSecretForTest8() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {

	secret1 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Path: "path",
			},
		},
	}

	secret2 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Path: "differentPath",
			},
		},
	}

	return secret1, secret2
}

func initSecretForTest9() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	var mode1 int32
	var mode2 int32
	mode1 = 1
	mode2 = 2

	secret1 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode1,
			},
		},
	}

	secret2 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode2,
			},
		},
	}

	return secret1, secret2
}

func initSecretForTest10() (*v1.SecretVolumeSource, *v1.SecretVolumeSource) {
	var mode1 int32
	mode1 = 1

	secret1 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{
				Mode: &mode1,
			},
		},
	}

	secret2 := &v1.SecretVolumeSource{
		Items: []v1.KeyToPath{
			{},
		},
	}

	return secret1, secret2
}

func TestCompareVolumeSecret(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 := initSecretForTest1()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretName) {
				t.Errorf("Error expected: '%s. 'secretName' vs 'different secretName''. But it was returned: %s", ErrorVolumeSecretName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'secretName' vs 'different secretName''. But the function found no errors", ErrorVolumeSecretName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest2()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeSecretOptional) {
				t.Errorf("Error expected: '%s. 'secretName' vs 'different secretName''. But it was returned: %s", ErrorMissingVolumeSecretOptional.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'secretName' vs 'different secretName''. But the function found no errors", ErrorMissingVolumeSecretOptional.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest3()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretOptional) {
				t.Errorf("Error expected: '%s. 'true' vs 'false''. But it was returned: %s", ErrorVolumeSecretOptional.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'true' vs 'false''. But the function found no errors", ErrorVolumeSecretOptional.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest4()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretDefaultMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorVolumeSecretDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorVolumeSecretDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest5()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeSecretDefaultMode) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingVolumeSecretDefaultMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingVolumeSecretDefaultMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest6()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretItemsLen) {
				t.Errorf("Error expected: '%s. '2' vs '1''. But it was returned: %s", ErrorVolumeSecretItemsLen.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorVolumeSecretItemsLen.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest7()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretItemsKey) {
				t.Errorf("Error expected: '%s. 'key' vs 'differentKey''. But it was returned: %s", ErrorVolumeSecretItemsKey.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'key' vs 'differentKey''. But the function found no errors", ErrorVolumeSecretItemsKey.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest8()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretItemsPath) {
				t.Errorf("Error expected: '%s. 'path' vs 'differentPath''. But it was returned: %s", ErrorVolumeSecretItemsPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'path' vs 'differentPath''. But the function found no errors", ErrorVolumeSecretItemsPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest9()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeSecretItemsMode) {
				t.Errorf("Error expected: '%s. '1' vs '2''. But it was returned: %s", ErrorVolumeSecretItemsMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorVolumeSecretItemsMode.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	secret1, secret2 = initSecretForTest10()

	CompareVolumeSecret(ctx, secret1, secret2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingVolumeSecretItemsMode) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingVolumeSecretItemsMode.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingVolumeSecretItemsMode.Error())
	}
}
