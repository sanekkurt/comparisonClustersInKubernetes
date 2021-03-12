package healthcheck

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	corev1 "k8s.io/api/core/v1"
)

func compareTCPSocketProbes(ctx context.Context, probe1, probe2 corev1.Probe) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)
	)

	if probe1.TCPSocket != nil && probe2.TCPSocket != nil {
		if probe1.TCPSocket.Host != probe2.TCPSocket.Host {
			logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketHost, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
		}

		if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal ||
			probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal ||
			probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
			logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketPort, probe1.TCPSocket.Port, probe2.TCPSocket.Port)
		}
	} else if probe1.TCPSocket != nil || probe2.TCPSocket != nil {
		logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocket, probe1.TCPSocket, probe2.TCPSocket)
	}

	return nil, nil
}