package nodeSelectors

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
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

func initNodeSelectorsForTest1() (map[string]string, map[string]string) {

	nodeSelector1 := make(map[string]string)
	nodeSelector2 := make(map[string]string)

	nodeSelector1["1"] = "one"
	nodeSelector1["2"] = "two"

	nodeSelector2["1"] = "one"
	nodeSelector2["2"] = "bad"

	return nodeSelector1, nodeSelector2
}

func initNodeSelectorsForTest2() (map[string]string, map[string]string) {

	nodeSelector1 := make(map[string]string)
	nodeSelector2 := make(map[string]string)

	nodeSelector1["1"] = "one"
	nodeSelector1["2"] = "two"

	nodeSelector2["1"] = "one"
	nodeSelector2["2"] = "two"
	nodeSelector2["3"] = "three"

	return nodeSelector1, nodeSelector2
}

func initNodeSelectorsForTest3() (map[string]string, map[string]string) {

	nodeSelector1 := make(map[string]string)
	nodeSelector2 := make(map[string]string)

	nodeSelector1["1"] = "one"
	nodeSelector1["2"] = "two"
	nodeSelector1["3"] = "three"

	nodeSelector2["1"] = "one"
	nodeSelector2["2"] = "two"

	return nodeSelector1, nodeSelector2
}

func TestCompareNodeSelectors(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	nodeSelector1, nodeSelector2 := initNodeSelectorsForTest1()

	CompareNodeSelectors(ctx, nodeSelector1, nodeSelector2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorDiffersNodeSelectorsInTemplates) {
				t.Errorf("Error expected: '%s. '2-two' vs '2-bad''. But it was returned: %s", ErrorDiffersNodeSelectorsInTemplates.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorDiffersNodeSelectorsInTemplates.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	nodeSelector1, nodeSelector2 = initNodeSelectorsForTest2()

	CompareNodeSelectors(ctx, nodeSelector1, nodeSelector2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorNodeSelectorsDoesNotExist) {
				t.Errorf("Error expected: '%s. Cluster number: 1. 3-three'. But it was returned: %s", ErrorNodeSelectorsDoesNotExist.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. Cluster number: 1. 3-three. But the function found no errors", ErrorNodeSelectorsDoesNotExist.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	nodeSelector1, nodeSelector2 = initNodeSelectorsForTest3()

	CompareNodeSelectors(ctx, nodeSelector1, nodeSelector2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorNodeSelectorsDoesNotExist) {
				t.Errorf("Error expected: '%s. Cluster number: 2. 3-three'. But it was returned: %s", ErrorNodeSelectorsDoesNotExist.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. Cluster number: 2. 3-three. But the function found no errors", ErrorNodeSelectorsDoesNotExist.Error())
	}

}
