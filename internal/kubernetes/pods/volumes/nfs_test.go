package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initNFSForTest1() (*v1.NFSVolumeSource, *v1.NFSVolumeSource) {
	nfs1 := &v1.NFSVolumeSource{
		ReadOnly: true,
	}

	nfs2 := &v1.NFSVolumeSource{
		ReadOnly: false,
	}

	return nfs1, nfs2
}

func initNFSForTest2() (*v1.NFSVolumeSource, *v1.NFSVolumeSource) {
	nfs1 := &v1.NFSVolumeSource{
		Path: "path",
	}

	nfs2 := &v1.NFSVolumeSource{
		Path: "diffPath",
	}

	return nfs1, nfs2
}

func initNFSForTest3() (*v1.NFSVolumeSource, *v1.NFSVolumeSource) {
	nfs1 := &v1.NFSVolumeSource{
		Server: "server",
	}

	nfs2 := &v1.NFSVolumeSource{
		Server: "diffServer",
	}

	return nfs1, nfs2
}

func TestCompareVolumeNFS(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	nfs1, nfs2 := initNFSForTest1()

	CompareVolumeNFS(ctx, nfs1, nfs2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeNFSReadOnly) {
				t.Errorf("Error expected: '%s. 'true' vs 'false''. But it was returned: %s", ErrorVolumeNFSReadOnly.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'true' vs 'false''. But the function found no errors", ErrorVolumeNFSReadOnly.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	nfs1, nfs2 = initNFSForTest2()

	CompareVolumeNFS(ctx, nfs1, nfs2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeNFSPath) {
				t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But it was returned: %s", ErrorVolumeNFSPath.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'path' vs 'diffPath''. But the function found no errors", ErrorVolumeNFSPath.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	nfs1, nfs2 = initNFSForTest3()

	CompareVolumeNFS(ctx, nfs1, nfs2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeNFSServer) {
				t.Errorf("Error expected: '%s. 'server' vs 'diffServer''. But it was returned: %s", ErrorVolumeNFSServer.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'server' vs 'diffServer''. But the function found no errors", ErrorVolumeNFSServer.Error())
	}
}
