package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeNFS(ctx context.Context, nfs1, nfs2 *v1.NFSVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if nfs1.ReadOnly != nfs2.ReadOnly {
		log.Warnf("%s. %t vs %t", ErrorVolumeNFSReadOnly.Error(), nfs1.ReadOnly, nfs2.ReadOnly)
	}

	if nfs1.Path != nfs2.Path {
		log.Warnf("%s. %s vs %s", ErrorVolumeNFSPath.Error(), nfs1.Path, nfs2.Path)
	}

	if nfs1.Server != nfs2.Server {
		log.Warnf("%s. %s vs %s", ErrorVolumeNFSServer.Error(), nfs1.Server, nfs2.Server)
	}

	return nil, nil
}
