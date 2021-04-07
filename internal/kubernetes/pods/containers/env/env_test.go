package env

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	v12 "k8s.io/api/core/v1"
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

func initEnvsForTest1() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ConfigMapKeyRef: &v12.ConfigMapKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "confMap",
				},
				Key: "test confMap",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			SecretKeyRef: &v12.SecretKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "secret",
				},
				Key: "test secret",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest2() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ConfigMapKeyRef: &v12.ConfigMapKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "confMapName1",
				},
				Key: "test confMap",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ConfigMapKeyRef: &v12.ConfigMapKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "confMapName2",
				},
				Key: "test confMap",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest3() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ConfigMapKeyRef: &v12.ConfigMapKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "confMapName",
				},
				Key: "test confMap1",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ConfigMapKeyRef: &v12.ConfigMapKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "confMapName",
				},
				Key: "test confMap2",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest4() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			SecretKeyRef: &v12.SecretKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "secretName1",
				},
				Key: "test secret",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			SecretKeyRef: &v12.SecretKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "secretName2",
				},
				Key: "test secret",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest5() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			SecretKeyRef: &v12.SecretKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "secretName1",
				},
				Key: "test secret1",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			SecretKeyRef: &v12.SecretKeySelector{
				LocalObjectReference: v12.LocalObjectReference{
					Name: "secretName1",
				},
				Key: "test secret2",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest6() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			FieldRef: &v12.ObjectFieldSelector{
				APIVersion: "v1",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			FieldRef: &v12.ObjectFieldSelector{
				APIVersion: "v2",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest7() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ResourceFieldRef: &v12.ResourceFieldSelector{
				ContainerName: "containerName1",
				Resource:      "resource",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ResourceFieldRef: &v12.ResourceFieldSelector{
				ContainerName: "containerName2",
				Resource:      "resource",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest8a() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ResourceFieldRef: &v12.ResourceFieldSelector{
				ContainerName: "containerName1",
				Resource:      "resource",
			},
		},
	}

	env2 := v12.EnvVar{
		Name: "env1",
	}

	return env1, env2
}

func initEnvsForTest8b() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
	}

	env2 := v12.EnvVar{
		Name: "env1",
		ValueFrom: &v12.EnvVarSource{
			ResourceFieldRef: &v12.ResourceFieldSelector{
				ContainerName: "containerName1",
				Resource:      "resource",
			},
		},
	}

	return env1, env2
}

func initEnvsForTest9() (v12.EnvVar, v12.EnvVar) {
	env1 := v12.EnvVar{
		Name:  "env1",
		Value: "hello",
	}

	env2 := v12.EnvVar{
		Name:  "env1",
		Value: "world",
	}

	return env1, env2
}

func initEnvsForTest10() ([]v12.EnvVar, []v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
	}

	env2 := v12.EnvVar{
		Name: "env2",
	}

	env3 := v12.EnvVar{
		Name: "env3",
	}

	envs1 := []v12.EnvVar{env1, env2, env3}
	envs2 := []v12.EnvVar{env1, env2}

	return envs1, envs2
}

func initEnvsForTest11() ([]v12.EnvVar, []v12.EnvVar) {
	env1 := v12.EnvVar{
		Name: "env1",
	}

	env2 := v12.EnvVar{
		Name: "env2",
	}

	env3 := v12.EnvVar{
		Name: "env3",
	}

	envs1 := []v12.EnvVar{env1, env2}
	envs2 := []v12.EnvVar{env1, env2, env3}

	return envs1, envs2
}

func TestCompareEnvVarValueFroms(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)

	env1, env2 := initEnvsForTest1()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSources) {
				t.Errorf("Error expected: '%s. env1: configMapKeyRef vs secretKeyRef. But it was returned: %s", ErrorVarDifferentValSources.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1: configMapKeyRef vs secretKeyRef. But the function found no errors", ErrorVarDifferentValSources.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest2()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSourceConfigMaps) {
				t.Errorf("Error expected: '%s. env1: 'confMapName1' vs 'confMapName2'. But it was returned: %s", ErrorVarDifferentValSourceConfigMaps.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1: 'confMapName1' vs 'confMapName2'. But the function found no errors", ErrorVarDifferentValSourceConfigMaps.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest3()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentKeyConfigMaps) {
				t.Errorf("Error expected: '%s. var-configMap: env1 - confMapName. Keys: 'test confMap1' vs 'test confMap2'. But it was returned: %s", ErrorVarDifferentKeyConfigMaps.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. var-configMap: env1 - confMapName. Keys: 'test confMap1' vs 'test confMap2'. But the function found no errors", ErrorVarDifferentKeyConfigMaps.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest4()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSourceSecrets) {
				t.Errorf("Error expected: '%s. env1: 'secretName1' vs 'secretName2'. But it was returned: %s", ErrorVarDifferentValSourceSecrets.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1: 'secretName1' vs 'secretName2'. But the function found no errors", ErrorVarDifferentValSourceSecrets.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest5()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentKeySecrets) {
				t.Errorf("Error expected: '%s. var-secret: env1 - secretName1. Keys: 'test secret1' vs 'test secret2'. But it was returned: %s", ErrorVarDifferentKeySecrets.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. var-secret: env1 - secretName1. Keys: 'test secret1' vs 'test secret2'. But the function found no errors", ErrorVarDifferentKeySecrets.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest6()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSourceFieldRef) {
				t.Errorf("Error expected: '%s. env1. But it was returned: %s", ErrorVarDifferentValSourceFieldRef.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1. But the function found no errors", ErrorVarDifferentValSourceFieldRef.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)

	env1, env2 = initEnvsForTest7()

	compareEnvVarValueFroms(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSourceResourceFieldRef) {
				t.Errorf("Error expected: '%s. env1. But it was returned: %s", ErrorVarDifferentValSourceResourceFieldRef.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1. But the function found no errors", ErrorVarDifferentValSourceResourceFieldRef.Error())
	}

}

func TestCompareEnvVarValueSources(t *testing.T) {
	cleanCtx := initCtx()
	ctx := newCtxWithCleanStorage(cleanCtx)

	env1, env2 := initEnvsForTest8a()

	compareEnvVarValueSources(ctx, env1, env2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSources) {
				t.Errorf("Error expected: '%s. env1: raw value vs ValueFrom. But it was returned: %s", ErrorVarDifferentValSources.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1. But the function found no errors", ErrorVarDifferentValSources.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	env1, env2 = initEnvsForTest8b()

	compareEnvVarValueSources(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValSources) {
				t.Errorf("Error expected: '%s. env1: raw value vs ValueFrom. But it was returned: %s", ErrorVarDifferentValSources.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1. But the function found no errors", ErrorVarDifferentValSources.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	env1, env2 = initEnvsForTest9()

	compareEnvVarValueSources(ctx, env1, env2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVarDifferentValues) {
				t.Errorf("Error expected: '%s. env1: 'hello' vs 'world'. But it was returned: %s", ErrorVarDifferentValues.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. env1: 'hello' vs 'world'. But the function found no errors", ErrorVarDifferentValues.Error())
	}
}

func TestCompare(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	envs1, envs2 := initEnvsForTest10()

	Compare(ctx, envs1, envs2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 2 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerDifferentEnvVarsNumber) {
				t.Errorf("Error expected: '%s: '3' vs '2'. But it was returned: %s", ErrorContainerDifferentEnvVarsNumber.Error(), diffs[0].Msg)
			}
			if !errors.Is(diffs[1].Msg.(error), ErrorVarDoesNotExistInOtherCluster) {
				t.Errorf("Error expected: '%s. Cluster number: '2'. varName: 'env3'. But it was returned: %s", ErrorVarDoesNotExistInOtherCluster.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Expected two errors. But the function found no errors")
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	envs1, envs2 = initEnvsForTest11()

	Compare(ctx, envs1, envs2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 2 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerDifferentEnvVarsNumber) {
				t.Errorf("Error expected: '%s: '2' vs '3'. But it was returned: %s", ErrorContainerDifferentEnvVarsNumber.Error(), diffs[0].Msg)
			}
			if !errors.Is(diffs[1].Msg.(error), ErrorVarDoesNotExistInOtherCluster) {
				t.Errorf("Error expected: '%s. Cluster number: '1'. varName: 'env3'. But it was returned: %s", ErrorVarDoesNotExistInOtherCluster.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Expected two errors. But the function found no errors")
	}
}
