package healthcheck

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func initProbesForTestHTTPGet1() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path: "path",
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet2() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Host: "Host1",
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Host: "Host2",
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet3() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeader1",
						Value: "valueForHTTPHeader1",
					},
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet4() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeader1",
						Value: "valueForHTTPHeader1",
					},
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeader1",
						Value: "valueForHTTPHeader1",
					},
					{
						Name:  "HTTPHeader2",
						Value: "valueForHTTPHeader2",
					},
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet5() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeaderName",
						Value: "valueForHTTPHeader",
					},
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeaderBadName",
						Value: "valueForHTTPHeader",
					},
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet6() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeaderName",
						Value: "valueForHTTPHeader",
					},
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				HTTPHeaders: []v1.HTTPHeader{
					{
						Name:  "HTTPHeaderName",
						Value: "badValueForHTTPHeader",
					},
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet7() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path: "path",
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Path: "badPath",
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet8a() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					IntVal: 1,
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					IntVal: 2,
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet8b() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					StrVal: "value1",
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					StrVal: "value2",
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet8c() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					Type: 1,
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.IntOrString{
					Type: 2,
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestHTTPGet9() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Scheme: "scheme",
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Scheme: "badScheme",
			},
		},
	}

	return probe1, probe2
}

func TestCompareHTTPGetProbes(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 := initProbesForTestHTTPGet1()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGet) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGet.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGet.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet2()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetHost) {
				t.Errorf("Error expected: '%s: 'Host1' vs 'Host2''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetHost.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'Host1' vs 'Host2''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetHost.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet3()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckHTTPGetMissingHeader) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckHTTPGetMissingHeader.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckHTTPGetMissingHeader.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet4()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet5()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetHeaderName) {
				t.Errorf("Error expected: '%s, 1: 'HTTPHeaderName' vs 'HTTPHeaderBadName''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetHeaderName.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s, 1: 'HTTPHeaderName' vs 'HTTPHeaderBadName''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetHeaderName.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet6()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetHeaderValue) {
				t.Errorf("Error expected: '%s, HTTPHeaderName: 'valueForHTTPHeader' vs 'badValueForHTTPHeader''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetHeaderValue.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s, HTTPHeaderName: 'valueForHTTPHeader' vs 'badValueForHTTPHeader''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetHeaderValue.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet7()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetPath) {
				t.Errorf("Error expected: '%s: 'path' vs 'badPath''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'path' vs 'badPath''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet8a()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetPortIntVal) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetPortIntVal.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetPortIntVal.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet8b()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetPortStrVal) {
				t.Errorf("Error expected: '%s: 'value1' vs 'value2''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetPortStrVal.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'value1' vs 'value2''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetPortStrVal.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet8c()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetPortType) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetPortType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetPortType.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestHTTPGet9()

	compareHTTPGetProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentHTTPGetScheme) {
				t.Errorf("Error expected: '%s: 'scheme' vs 'badScheme''. But it was returned: %s", ErrorContainerHealthCheckDifferentHTTPGetScheme.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'scheme' vs 'badScheme''. But the function found no errors", ErrorContainerHealthCheckDifferentHTTPGetScheme.Error())
	}
}
