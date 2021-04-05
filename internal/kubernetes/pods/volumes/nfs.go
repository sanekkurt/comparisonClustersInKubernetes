package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeNFS(ctx context.Context, nfs1, nfs2 *v1.NFSVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if nfs1.ReadOnly != nfs2.ReadOnly {
		diffsBatch.Add(ctx, false, "%s. %t vs %t", ErrorVolumeNFSReadOnly.Error(), nfs1.ReadOnly, nfs2.ReadOnly)
	}

	if nfs1.Path != nfs2.Path {
		diffsBatch.Add(ctx, false, "%s. %s vs %s", ErrorVolumeNFSPath.Error(), nfs1.Path, nfs2.Path)
	}

	if nfs1.Server != nfs2.Server {
		diffsBatch.Add(ctx, false, "%s. %s vs %s", ErrorVolumeNFSServer.Error(), nfs1.Server, nfs2.Server)
	}

	return
}
