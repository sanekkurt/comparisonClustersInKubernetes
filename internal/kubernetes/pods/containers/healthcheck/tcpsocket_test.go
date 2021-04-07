package healthcheck

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func initProbesForTestTcpSocket1() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Host: "host1",
			},
		},
	}

	probe2 := v1.Probe{}

	return probe1, probe2
}

func initProbesForTestTcpSocket2() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Host: "host1",
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Host: "host2",
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestTcpSocket3a() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 1,
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 2,
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestTcpSocket3b() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					StrVal: "value1",
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					StrVal: "value2",
				},
			},
		},
	}

	return probe1, probe2
}

func initProbesForTestTcpSocket3c() (v1.Probe, v1.Probe) {
	probe1 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type: 1,
				},
			},
		},
	}

	probe2 := v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type: 2,
				},
			},
		},
	}

	return probe1, probe2
}

func TestCompareTCPSocketProbes(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 := initProbesForTestTcpSocket1()

	compareTCPSocketProbes(ctx, probe1, probe2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTCPSocket) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorContainerHealthCheckDifferentTCPSocket.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorContainerHealthCheckDifferentTCPSocket.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestTcpSocket2()

	compareTCPSocketProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTCPSocketHost) {
				t.Errorf("Error expected: '%s: 'host1' vs 'host2''. But it was returned: %s", ErrorContainerHealthCheckDifferentTCPSocketHost.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'host1' vs 'host2''. But the function found no errors", ErrorContainerHealthCheckDifferentTCPSocketHost.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestTcpSocket3a()

	compareTCPSocketProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTCPSocketPortIntVal) {
				t.Errorf("Error expected: '%s: '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentTCPSocketPortIntVal.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: '1' vs '2''. But the function found no errors", ErrorContainerHealthCheckDifferentTCPSocketPortIntVal.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestTcpSocket3b()

	compareTCPSocketProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTCPSocketPortStrVal) {
				t.Errorf("Error expected: '%s: 'value1' vs 'value2''. But it was returned: %s", ErrorContainerHealthCheckDifferentTCPSocketPortStrVal.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s: 'value1' vs 'value2''. But the function found no errors", ErrorContainerHealthCheckDifferentTCPSocketPortStrVal.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	probe1, probe2 = initProbesForTestTcpSocket3c()

	compareTCPSocketProbes(ctx, probe1, probe2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorContainerHealthCheckDifferentTCPSocketPortType) {
				t.Errorf("Error expected: '%s:  '1' vs '2''. But it was returned: %s", ErrorContainerHealthCheckDifferentTCPSocketPortType.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s:  '1' vs '2''. But the function found no errors", ErrorContainerHealthCheckDifferentTCPSocketPortType.Error())
	}
}
