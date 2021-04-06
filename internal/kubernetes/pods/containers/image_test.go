package containers

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

func initLoggingAndConfig() context.Context {
	var (
		debug bool

		ctx = context.Background()
	)

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	err := logging.Configure(debug)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	log := logging.FromContext(ctx)

	args := []string{os.Args[0], "-c", "C:\\Users\\Александр\\go\\src\\comparisonClustersInKubernetes\\config.yaml" /*"../../config.yaml"*/}

	cfg, err := config.Parse(ctx, args)
	if err != nil {
		if err == config.ErrHelpShown {
			return nil
		}

		log.Error(err)
		return nil
	}

	ctx = config.With(ctx, cfg)

	return ctx
}

func initEnvForImageTagTest(ctx context.Context) (context.Context, v1.Container, v1.Container) {

	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	container1 := v1.Container{
		Image:           "k8s.gcr.io/pause:3.1",
		ImagePullPolicy: "always",
	}

	container2 := v1.Container{
		Image:           "k8s.gcr.io/pause:3.4",
		ImagePullPolicy: "always",
	}

	return ctx, container1, container2
}

func initEnvForImageLabelsTest(ctx context.Context) (context.Context, v1.Container, v1.Container) {

	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	container1 := v1.Container{
		Image:           "k8s.gcr.io/pause:3.1",
		ImagePullPolicy: "always",
	}

	container2 := v1.Container{
		Image:           "pause:3.1",
		ImagePullPolicy: "always",
	}

	return ctx, container1, container2
}

func initEnvForImagePullPolicyTest(ctx context.Context) (context.Context, v1.Container, v1.Container) {

	diffs := diff.NewDiffsStorage(ctx)
	ctx = diff.WithDiffStorage(ctx, diffs)

	batch := diffs.NewLazyBatch(metav1.TypeMeta{Kind: "", APIVersion: ""}, metav1.ObjectMeta{})
	ctx = diff.WithDiffBatch(ctx, batch)

	container1 := v1.Container{
		Image:           "k8s.gcr.io/pause:3.1",
		ImagePullPolicy: "IfNotPresent",
	}

	container2 := v1.Container{
		Image:           "k8s.gcr.io/pause:3.1",
		ImagePullPolicy: "always",
	}

	return ctx, container1, container2
}

func TestCompareContainerSpecImages(t *testing.T) {
	ctx := initLoggingAndConfig()
	ctxTagTest, container1, container2 := initEnvForImageTagTest(ctx)

	compareContainerSpecImages(ctxTagTest, container1, container2)

	diffStorage := diff.StorageFromContext(ctxTagTest)
	diffStorage.Finalize(ctxTagTest)
	batch := diff.BatchFromContext(ctxTagTest)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerDifferentImageTags) {
				t.Errorf("Error expected: '%s: 3.1 vs 3.4'. But it was returned: %s", ErrorContainerDifferentImageTags.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 3.1 vs 3.4'. But the function found no errors", ErrorContainerDifferentImageTags.Error())
	}

	ctxLabelsTest, container1, container2 := initEnvForImageLabelsTest(ctx)

	compareContainerSpecImages(ctxLabelsTest, container1, container2)

	diffStorage = diff.StorageFromContext(ctxLabelsTest)
	diffStorage.Finalize(ctxLabelsTest)
	batch = diff.BatchFromContext(ctxLabelsTest)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerDifferentImageLabels) {
				t.Errorf("Error expected: '%s: k8s.gcr.io/pause vs pause'. But it was returned: %s", ErrorContainerDifferentImageLabels.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: k8s.gcr.io/pause vs pause'. But the function found no errors", ErrorContainerDifferentImageLabels.Error())
	}

	ctxImagePullPolicy, container1, container2 := initEnvForImagePullPolicyTest(ctx)

	compareContainerSpecImages(ctxImagePullPolicy, container1, container2)

	diffStorage = diff.StorageFromContext(ctxImagePullPolicy)
	diffStorage.Finalize(ctxImagePullPolicy)
	batch = diff.BatchFromContext(ctxImagePullPolicy)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerDifferentImagePolicies) {
				t.Errorf("Error expected: '%s: IfNotPresent vs always'. But it was returned: %s", ErrorContainerDifferentImagePolicies.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: IfNotPresent vs always'. But the function found no errors", ErrorContainerDifferentImagePolicies.Error())
	}

}

//func TestCompareContainerSpecImages(t *testing.T) {
//	ctx, container1, container2 := initEnv()
//
//	compareContainerSpecImages(ctx, container1, container2)
//
//	diffStorage := diff.StorageFromContext(ctx)
//	diffStorage.Finalize(ctx)
//
//	batch := diff.BatchFromContext(ctx)
//
//	if batch.Diffs != nil {
//		if len(batch.Diffs) == 1 {
//			if batch.Diffs[0].Msg != "different container image tags in Pod specs: 3.1 vs 3.4" {
//				t.Errorf("Error expected: 'different container image tags in Pod specs: 3.1 vs 3.4'. But it was returned: %s", batch.Diffs[0].Msg)
//			}
//		} else {
//			t.Errorf("1 error was expected. But it was returned: %d", len(batch.Diffs))
//		}
//	} else {
//		t.Errorf("Error expected: 'different container image tags in Pod specs: 3.1 vs 3.4'. But the function found no errors")
//	}
//
//}
