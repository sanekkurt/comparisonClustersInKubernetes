package service

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

func initServicesForTests1() (v12.Service, v12.Service) {
	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{}, {},
			},
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{},
			},
		},
	}

	return service1, service2
}

func initServicesForTests2() (v12.Service, v12.Service) {
	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "name",
					Port:     80,
					Protocol: "http",
				},
			},
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "diffName",
					Port:     80,
					Protocol: "http",
				},
			},
		},
	}

	return service1, service2
}

func initServicesForTests3() (v12.Service, v12.Service) {
	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "name",
					Port:     80,
					Protocol: "http",
				},
			},
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "name",
					Port:     81,
					Protocol: "http",
				},
			},
		},
	}

	return service1, service2
}

func initServicesForTests4() (v12.Service, v12.Service) {
	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "name",
					Port:     80,
					Protocol: "http",
				},
			},
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:     "name",
					Port:     80,
					Protocol: "https",
				},
			},
		},
	}

	return service1, service2
}

func initServicesForTests5() (v12.Service, v12.Service) {
	sel1 := make(map[string]string)
	sel2 := make(map[string]string)

	sel1["1"] = "selector-key1"
	sel1["2"] = "selector-key2"
	sel2["1"] = "selector-key1"

	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Selector: sel1,
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Selector: sel2,
		},
	}

	return service1, service2
}

func initServicesForTests6() (v12.Service, v12.Service) {
	sel1 := make(map[string]string)
	sel2 := make(map[string]string)

	sel1["key1"] = "value1"
	sel1["key2"] = "value2"
	sel2["key1"] = "diffValue"
	sel2["key2"] = "value2"

	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Selector: sel1,
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Selector: sel2,
		},
	}

	return service1, service2
}

func initServicesForTests7() (v12.Service, v12.Service) {

	service1 := v12.Service{
		Spec: v12.ServiceSpec{
			Type: "type",
		},
	}

	service2 := v12.Service{
		Spec: v12.ServiceSpec{
			Type: "diffType",
		},
	}

	return service1, service2
}

func TestCompareSpecInServices(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	service1, service2 := initServicesForTests1()

	compareSpecInServices(ctx, service1, service2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPortsCountDifferent) {
				t.Errorf("Error expected: '%s. '2' vs '1''. But it was returned: %s", ErrorPortsCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorPortsCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests2()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPortInServicesDifferent) {
				t.Errorf("Error expected: '%s. 'name-80-http' vs 'diffName-80-http''. But it was returned: %s", ErrorPortInServicesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name-80-http' vs 'diffName-80-http''. But the function found no errors", ErrorPortInServicesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests3()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPortInServicesDifferent) {
				t.Errorf("Error expected: '%s. 'name-80-http' vs 'name-81-http''. But it was returned: %s", ErrorPortInServicesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name-80-http' vs 'name-81-http''. But the function found no errors", ErrorPortInServicesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests4()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPortInServicesDifferent) {
				t.Errorf("Error expected: '%s. 'name-80-http' vs 'name-80-https''. But it was returned: %s", ErrorPortInServicesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name-80-http' vs 'name-80-https''. But the function found no errors", ErrorPortInServicesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests5()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorSelectorsCountDifferent) {
				t.Errorf("Error expected: '%s. '2' vs '1''. But it was returned: %s", ErrorSelectorsCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorSelectorsCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests6()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorSelectorInServicesDifferent) {
				t.Errorf("Error expected: '%s. 'key1-value1' vs 'key1-diffValue''. But it was returned: %s", ErrorSelectorInServicesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'key1-value1' vs 'key1-diffValue''. But the function found no errors", ErrorSelectorInServicesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	service1, service2 = initServicesForTests7()

	compareSpecInServices(ctx, service1, service2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorTypeInServicesDifferent) {
				t.Errorf("Error expected: '%s. 'type' vs 'diffType''. But it was returned: %s", ErrorTypeInServicesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'type' vs 'diffType''. But the function found no errors", ErrorTypeInServicesDifferent.Error())
	}
}
