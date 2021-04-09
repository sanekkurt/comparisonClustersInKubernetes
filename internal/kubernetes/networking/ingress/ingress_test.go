package ingress

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/logging"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
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

func initIngressForTest1() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{}, {},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest2() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					SecretName: "secretName",
				},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					SecretName: "otherSecretName",
				},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest3() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts: []string{
						"",
					},
				},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts: []string{
						"", "",
					},
				},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest4() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts: []string{
						"hostname1",
					},
				},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts: []string{
						"hostname2",
					},
				},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest5() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts: []string{
						"hostname2",
					},
				},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest6() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest7() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			DefaultBackend: &v1.IngressBackend{},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{},
	}

	return ingress1, ingress2
}

func initIngressForTest8() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{}, {},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest9() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: "hostname1",
				},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: "hostname2",
				},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest10() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{},
					},
				},
			},
		},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{},
			},
		},
	}

	return ingress1, ingress2
}

func initIngressForTest11() (v1.Ingress, v1.Ingress) {

	ingress1 := v1.Ingress{
		Spec: v1.IngressSpec{},
	}

	ingress2 := v1.Ingress{
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{},
			},
		},
	}

	return ingress1, ingress2
}

func TestCompareSpecInIngresses(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 := initIngressForTest1()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorTLSCountDifferent) {
				t.Errorf("Error expected. '%s: '2' vs '1''. But it was returned: %s", ErrorTLSCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorTLSCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest2()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorSecretNameInTLSDifferent) {
				t.Errorf("Error expected. '%s: 'secretName' vs 'otherSecretName''. But it was returned: %s", ErrorSecretNameInTLSDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'secretName' vs 'otherSecretName''. But the function found no errors", ErrorSecretNameInTLSDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest3()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHostsCountDifferent) {
				t.Errorf("Error expected. '%s: '1' vs '2''. But it was returned: %s", ErrorHostsCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorHostsCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest4()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorNameHostDifferent) {
				t.Errorf("Error expected. '%s: 'hostname1' vs 'hostname2''. But it was returned: %s", ErrorNameHostDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'hostname1' vs 'hostname2''. But the function found no errors", ErrorNameHostDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest5()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHostsInIngressesDifferent) {
				t.Errorf("Error expected. '%s''. But it was returned: %s", ErrorHostsInIngressesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorHostsInIngressesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest6()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorTLSInIngressesDifferent) {
				t.Errorf("Error expected. '%s''. But it was returned: %s", ErrorTLSInIngressesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorTLSInIngressesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest7()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendInIngressesDifferent) {
				t.Errorf("Error expected. '%s''. But it was returned: %s", ErrorBackendInIngressesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorBackendInIngressesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest8()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorRulesCountDifferent) {
				t.Errorf("Error expected. '%s. '1' vs '2''. But it was returned: %s", ErrorRulesCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '1' vs '2''. But the function found no errors", ErrorRulesCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest9()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHostNameInRuleDifferent) {
				t.Errorf("Error expected. '%s. 'hostname1' vs 'hostname2''. But it was returned: %s", ErrorHostNameInRuleDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'hostname1' vs 'hostname2''. But the function found no errors", ErrorHostNameInRuleDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest10()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorHTTPInIngressesDifferent) {
				t.Errorf("Error expected. '%s'. But it was returned: %s", ErrorHTTPInIngressesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorHTTPInIngressesDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingress1, ingress2 = initIngressForTest11()

	compareSpecInIngresses(ctx, ingress1, ingress2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorRulesInIngressesDifferent) {
				t.Errorf("Error expected. '%s'. But it was returned: %s", ErrorRulesInIngressesDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorRulesInIngressesDifferent.Error())
	}
}

func initIngressBackendForTest1() (v1.IngressBackend, v1.IngressBackend) {
	ingressBackend1 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Name: "serviceBackendName1",
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Name: "serviceBackendName2",
		},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest2() (v1.IngressBackend, v1.IngressBackend) {
	ingressBackend1 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Port: v1.ServiceBackendPort{
				Name:   "http",
				Number: 80,
			},
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Port: v1.ServiceBackendPort{
				Name:   "https",
				Number: 80,
			},
		},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest3() (v1.IngressBackend, v1.IngressBackend) {
	ingressBackend1 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Port: v1.ServiceBackendPort{
				Name:   "https",
				Number: 8080,
			},
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{
			Port: v1.ServiceBackendPort{
				Name:   "https",
				Number: 8088,
			},
		},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest4() (v1.IngressBackend, v1.IngressBackend) {
	ingressBackend1 := v1.IngressBackend{
		Service: &v1.IngressServiceBackend{},
	}

	ingressBackend2 := v1.IngressBackend{}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest5() (v1.IngressBackend, v1.IngressBackend) {

	apiGr1 := "apiGr1"
	apiGr2 := "apiGr2"

	ingressBackend1 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			APIGroup: &apiGr1,
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			APIGroup: &apiGr2,
		},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest6() (v1.IngressBackend, v1.IngressBackend) {

	apiGr1 := "apiGr1"

	ingressBackend1 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			APIGroup: &apiGr1,
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest7() (v1.IngressBackend, v1.IngressBackend) {

	ingressBackend1 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			Name: "resName1",
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			Name: "resName2",
		},
	}

	return ingressBackend1, ingressBackend2
}

func initIngressBackendForTest8() (v1.IngressBackend, v1.IngressBackend) {

	ingressBackend1 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			Kind: "testKind1",
		},
	}

	ingressBackend2 := v1.IngressBackend{
		Resource: &v12.TypedLocalObjectReference{
			Kind: "testKind2",
		},
	}

	return ingressBackend1, ingressBackend2
}

func TestCompareIngressesBackend(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 := initIngressBackendForTest1()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorServiceNameInBackendDifferent) {
				t.Errorf("Error expected. '%s. 'serviceBackendName1' vs 'serviceBackendName2''. But it was returned: %s", ErrorServiceNameInBackendDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'serviceBackendName1' vs 'serviceBackendName2''. But the function found no errors", ErrorServiceNameInBackendDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest2()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendServicePortDifferent) {
				t.Errorf("Error expected. '%s. 'http-80' vs 'https-80''. But it was returned: %s", ErrorBackendServicePortDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'http-80' vs 'https-80''. But the function found no errors", ErrorBackendServicePortDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest3()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendServicePortDifferent) {
				t.Errorf("Error expected. '%s. 'https-8080' vs 'https-8088''. But it was returned: %s", ErrorBackendServicePortDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'https-8080' vs 'https-8088''. But the function found no errors", ErrorBackendServicePortDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest4()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendServiceIsMissing) {
				t.Errorf("Error expected. '%s'. But it was returned: %s", ErrorBackendServiceIsMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorBackendServiceIsMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest5()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendResourceApiGroup) {
				t.Errorf("Error expected. '%s. 'apiGr1' vs 'apiGr2''. But it was returned: %s", ErrorBackendResourceApiGroup.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'apiGr1' vs 'apiGr2''. But the function found no errors", ErrorBackendResourceApiGroup.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest6()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendResourceApiGroupIsMissing) {
				t.Errorf("Error expected. '%s'. But it was returned: %s", ErrorBackendResourceApiGroupIsMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorBackendResourceApiGroupIsMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest7()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendResourceName) {
				t.Errorf("Error expected. '%s. 'resName1' vs 'resName2''. But it was returned: %s", ErrorBackendResourceName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'resName1' vs 'resName2''. But the function found no errors", ErrorBackendResourceName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressBackend1, ingressBackend2 = initIngressBackendForTest8()

	compareIngressesBackend(ctx, ingressBackend1, ingressBackend2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorBackendResourceKind) {
				t.Errorf("Error expected. '%s. 'testKind1' vs 'testKind2''. But it was returned: %s", ErrorBackendResourceKind.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'testKind1' vs 'testKind2''. But the function found no errors", ErrorBackendResourceKind.Error())
	}
}

func initIngressesHTTPForTest1() (v1.HTTPIngressRuleValue, v1.HTTPIngressRuleValue) {
	ingressHTTP1 := v1.HTTPIngressRuleValue{
		Paths: []v1.HTTPIngressPath{
			{}, {},
		},
	}

	ingressHTTP2 := v1.HTTPIngressRuleValue{
		Paths: []v1.HTTPIngressPath{
			{},
		},
	}

	return ingressHTTP1, ingressHTTP2
}

func initIngressesHTTPForTest2() (v1.HTTPIngressRuleValue, v1.HTTPIngressRuleValue) {
	ingressHTTP1 := v1.HTTPIngressRuleValue{
		Paths: []v1.HTTPIngressPath{
			{
				Path: "path1",
			},
		},
	}

	ingressHTTP2 := v1.HTTPIngressRuleValue{
		Paths: []v1.HTTPIngressPath{
			{
				Path: "path2",
			},
		},
	}

	return ingressHTTP1, ingressHTTP2
}

func TestCompareIngressesHTTP(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	ingressHTTP1, ingressHTTP2 := initIngressesHTTPForTest1()

	compareIngressesHTTP(ctx, ingressHTTP1, ingressHTTP2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPathsCountDifferent) {
				t.Errorf("Error expected. '%s. '2' vs '1''. But it was returned: %s", ErrorPathsCountDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. '2' vs '1''. But the function found no errors", ErrorPathsCountDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	ingressHTTP1, ingressHTTP2 = initIngressesHTTPForTest2()

	compareIngressesHTTP(ctx, ingressHTTP1, ingressHTTP2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPathValueDifferent) {
				t.Errorf("Error expected. '%s. 'path1' vs 'path2''. But it was returned: %s", ErrorPathValueDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'path1' vs 'path2''. But the function found no errors", ErrorPathValueDifferent.Error())
	}
}
