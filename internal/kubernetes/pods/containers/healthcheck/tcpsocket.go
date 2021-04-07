package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/diff"

	corev1 "k8s.io/api/core/v1"
)

func compareTCPSocketProbes(ctx context.Context, probe1, probe2 corev1.Probe) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if probe1.TCPSocket != nil && probe2.TCPSocket != nil {
		if probe1.TCPSocket.Host != probe2.TCPSocket.Host {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketHost, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentTCPSocketHost, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
		}

		if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal {
			diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorContainerHealthCheckDifferentTCPSocketPortIntVal, probe1.TCPSocket.Port.IntVal, probe2.TCPSocket.Port.IntVal)
		}

		if probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal {
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorContainerHealthCheckDifferentTCPSocketPortStrVal, probe1.TCPSocket.Port.StrVal, probe2.TCPSocket.Port.StrVal)
		}

		if probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
			diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorContainerHealthCheckDifferentTCPSocketPortType, probe1.TCPSocket.Port.Type, probe2.TCPSocket.Port.Type)
		}

		//if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal ||
		//	probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal ||
		//	probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
		//	//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketPort, probe1.TCPSocket.Port, probe2.TCPSocket.Port)
		//	diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckDifferentTCPSocketPort)
		//}
	} else if probe1.TCPSocket != nil || probe2.TCPSocket != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocket, probe1.TCPSocket, probe2.TCPSocket)
		diffsBatch.Add(ctx, false, "%w", ErrorContainerHealthCheckDifferentTCPSocket)
	}
}
