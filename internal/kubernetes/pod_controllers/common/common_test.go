package common

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/config"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	"os"
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

	args := []string{os.Args[0], "-c", "C:\\Users\\Александр\\go\\src\\comparisonClustersInKubernetes\\config.yaml" /*"../../config.yaml"*/}

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
	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	return ctx
}

func initPodControllerSpecsForTest1() (*AbstractPodController, *AbstractPodController) {
	var repl1 int32
	var repl2 int32
	repl1 = 5
	repl2 = 9

	apc1 := &AbstractPodController{
		Replicas: &repl1,
	}

	apc2 := &AbstractPodController{
		Replicas: &repl2,
	}

	return apc1, apc2
}

func initPodControllerSpecsForTest2() (*AbstractPodController, *AbstractPodController) {

	var repl2 int32

	repl2 = 9

	apc1 := &AbstractPodController{}

	apc2 := &AbstractPodController{
		Replicas: &repl2,
	}

	return apc1, apc2
}

func TestComparePodControllerSpecs(t *testing.T) {
	ctx := initCtx()

	apc1, apc2 := initPodControllerSpecsForTest1()

	ComparePodControllerSpecs(ctx, "objName", apc1, apc2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batches := diffStorage.GetBatches()
	diffs := batches[0].GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDifferentNumberReplicas) {
				t.Errorf("Error expected: '%s: '5' vs '9''. But it was returned: %s", ErrorDifferentNumberReplicas.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '5' vs '9''. But the function found no errors", ErrorDifferentNumberReplicas.Error())
	}

	ctx = initCtx()

	apc1, apc2 = initPodControllerSpecsForTest2()

	ComparePodControllerSpecs(ctx, "objName", apc1, apc2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batches = diffStorage.GetBatches()
	diffs = batches[0].GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingReplicas) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingReplicas.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingReplicas.Error())
	}

}
