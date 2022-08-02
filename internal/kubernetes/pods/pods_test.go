package pods

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"k8s-cluster-comparator/internal/config"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func initCtx() context.Context {
	var (
		ctx = context.Background()
	)
	err := logging.ConfigureForTests()
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return ctx
	}

	args := []string{os.Args[0], "-c", strings.Join([]string{cwd, "..", "..", "..", "config.yaml"}, string(os.PathSeparator))}

	log := logging.FromContext(ctx)
	cfg, err := config.Parse(ctx, args)
	if err != nil {
		if err == config.ErrHelpShown {
			return nil
		}
		log.Error(err)
		return nil
	}

	ctx = config.With(ctx, cfg)
	ctx = kubectx.WithNamespace(ctx, "namespace")
	return ctx
}
func newCtxWithCleanStorage(ctx context.Context) context.Context {
	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	return ctx
}

func initPodSpecsForTest1() (types.InformationAboutObject, types.InformationAboutObject) {

	spec1 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{}, {},
				},
			},
		},
	}

	spec2 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{},
				},
			},
		},
	}

	return spec1, spec2
}

func initPodSpecsForTest2() (types.InformationAboutObject, types.InformationAboutObject) {

	map1 := make(map[string]string)
	map2 := make(map[string]string)

	map1["1"] = "1"
	map2["1"] = "1"
	map2["2"] = "2"

	spec1 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				NodeSelector: map1,
			},
		},
	}

	spec2 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				NodeSelector: map2,
			},
		},
	}

	return spec1, spec2
}

func initPodSpecsForTest3() (types.InformationAboutObject, types.InformationAboutObject) {
	map2 := make(map[string]string)

	map2["1"] = "1"
	map2["2"] = "2"

	spec1 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{},
		},
	}

	spec2 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				NodeSelector: map2,
			},
		},
	}

	return spec1, spec2
}

func initPodSpecsForTest4() (types.InformationAboutObject, types.InformationAboutObject) {

	spec1 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{}, {},
				},
			},
		},
	}

	spec2 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{},
				},
			},
		},
	}

	return spec1, spec2
}

func initPodSpecsForTest5() (types.InformationAboutObject, types.InformationAboutObject) {

	spec1 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{}, {},
				},
			},
		},
	}

	spec2 := types.InformationAboutObject{
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{},
		},
	}

	return spec1, spec2
}

func TestComparePodSpecs(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	spec1, spec2 := initPodSpecsForTest1()

	if err := ComparePodSpecs(ctx, spec1, spec2); err != nil {
		t.Fatalf("cannot complete: %v", err)
	}

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDiffersContainersNumberInTemplates) {
				t.Errorf("Error expected: '%s: '2' vs '1''. But it was returned: %s", ErrorDiffersContainersNumberInTemplates.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '2' vs '1''. But the function found no errors", ErrorDiffersContainersNumberInTemplates.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	spec1, spec2 = initPodSpecsForTest2()

	if err := ComparePodSpecs(ctx, spec1, spec2); err != nil {
		t.Fatalf("cannot complete: %v", err)
	}

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDiffersNodeSelectorsNumberInTemplates) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorDiffersNodeSelectorsNumberInTemplates.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorDiffersNodeSelectorsNumberInTemplates.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	spec1, spec2 = initPodSpecsForTest3()

	if err := ComparePodSpecs(ctx, spec1, spec2); err != nil {
		t.Fatalf("cannot complete: %v", err)
	}

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingNodeSelectors) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingNodeSelectors.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingNodeSelectors.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	spec1, spec2 = initPodSpecsForTest4()

	if err := ComparePodSpecs(ctx, spec1, spec2); err != nil {
		t.Fatalf("cannot complete: %v", err)
	}

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDiffersVolumesNumberInTemplates) {
				t.Errorf("Error expected: '%s: '2' vs '1''. But it was returned: %s", ErrorDiffersVolumesNumberInTemplates.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '2' vs '1''. But the function found no errors", ErrorDiffersVolumesNumberInTemplates.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	spec1, spec2 = initPodSpecsForTest5()

	if err := ComparePodSpecs(ctx, spec1, spec2); err != nil {
		t.Fatalf("cannot complete: %v", err)
	}

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumes) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumes.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumes.Error())
	}
}
