package healthcheck

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s-cluster-comparator/internal/kubernetes/types"

	corev1 "k8s.io/api/core/v1"
)

func compareTCPSocketProbes(ctx context.Context, probe1, probe2 corev1.Probe) {
	var (
		//log = logging.FromContext(ctx)

		diffsBatch = ctx.Value("diffBatch").(*diff.DiffsBatch)
		meta       = ctx.Value("apcMeta").(types.AbstractObjectMetadata)
	)

	if probe1.TCPSocket != nil && probe2.TCPSocket != nil {
		if probe1.TCPSocket.Host != probe2.TCPSocket.Host {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketHost, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
			diffsBatch.Add(ctx, &meta.Type, &meta.Meta, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentTCPSocketHost.Error(), probe1.TCPSocket.Host, probe2.TCPSocket.Host)
		}

		if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal ||
			probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal ||
			probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
			//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocketPort, probe1.TCPSocket.Port, probe2.TCPSocket.Port)
			diffsBatch.Add(ctx, &meta.Type, &meta.Meta, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerHealthCheckDifferentTCPSocketPort.Error(), probe1.TCPSocket.Port, probe2.TCPSocket.Port)
		}
	} else if probe1.TCPSocket != nil || probe2.TCPSocket != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckDifferentTCPSocket, probe1.TCPSocket, probe2.TCPSocket)
		diffsBatch.Add(ctx, &meta.Type, &meta.Meta, false, zap.WarnLevel, "%s", ErrorContainerHealthCheckDifferentTCPSocket.Error())
	}

}
