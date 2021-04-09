package volumes

import (
	"errors"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func initVolumesForTest1() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		Name: "name",
	}
	volume2 := v1.Volume{
		Name: "diffName",
	}
	return volume1, volume2
}

func initVolumesForTest2() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest3() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest4() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest5() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			NFS: &v1.NFSVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest6() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest7() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			DownwardAPI: &v1.DownwardAPIVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest8() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest9() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			CSI: &v1.CSIVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest10() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ISCSI: &v1.ISCSIVolumeSource{
				TargetPortal: "portal",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ISCSI: &v1.ISCSIVolumeSource{
				TargetPortal: "portal2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest11() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ISCSI: &v1.ISCSIVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest12() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			CephFS: &v1.CephFSVolumeSource{
				Path: "path1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			CephFS: &v1.CephFSVolumeSource{
				Path: "path2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest13() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			CephFS: &v1.CephFSVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest14() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			GCEPersistentDisk: &v1.GCEPersistentDiskVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			GCEPersistentDisk: &v1.GCEPersistentDiskVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest15() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			GCEPersistentDisk: &v1.GCEPersistentDiskVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest16() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest17() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest18() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Glusterfs: &v1.GlusterfsVolumeSource{
				Path: "path1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Glusterfs: &v1.GlusterfsVolumeSource{
				Path: "path2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest19() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Glusterfs: &v1.GlusterfsVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest20() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			RBD: &v1.RBDVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			RBD: &v1.RBDVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest21() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			RBD: &v1.RBDVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest22() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FlexVolume: &v1.FlexVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FlexVolume: &v1.FlexVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest23() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FlexVolume: &v1.FlexVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest24() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Cinder: &v1.CinderVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Cinder: &v1.CinderVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest25() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Cinder: &v1.CinderVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest26() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Flocker: &v1.FlockerVolumeSource{
				DatasetName: "name1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Flocker: &v1.FlockerVolumeSource{
				DatasetName: "name2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest27() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Flocker: &v1.FlockerVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest28() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FC: &v1.FCVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FC: &v1.FCVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest29() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			FC: &v1.FCVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest30() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureFile: &v1.AzureFileVolumeSource{
				ReadOnly: true,
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureFile: &v1.AzureFileVolumeSource{
				ReadOnly: false,
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest31() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureFile: &v1.AzureFileVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest32() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			VsphereVolume: &v1.VsphereVirtualDiskVolumeSource{
				FSType: "type1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			VsphereVolume: &v1.VsphereVirtualDiskVolumeSource{
				FSType: "type2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest33() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			VsphereVolume: &v1.VsphereVirtualDiskVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest34() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Quobyte: &v1.QuobyteVolumeSource{
				ReadOnly: false,
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Quobyte: &v1.QuobyteVolumeSource{
				ReadOnly: true,
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest35() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Quobyte: &v1.QuobyteVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest36() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureDisk: &v1.AzureDiskVolumeSource{
				DiskName: "disk1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureDisk: &v1.AzureDiskVolumeSource{
				DiskName: "disk2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest37() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			AzureDisk: &v1.AzureDiskVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest38() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			PhotonPersistentDisk: &v1.PhotonPersistentDiskVolumeSource{
				FSType: "FSType1",
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			PhotonPersistentDisk: &v1.PhotonPersistentDiskVolumeSource{
				FSType: "FSType2",
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest39() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			PhotonPersistentDisk: &v1.PhotonPersistentDiskVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest40() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Projected: &v1.ProjectedVolumeSource{
				Sources: []v1.VolumeProjection{
					{},
				},
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Projected: &v1.ProjectedVolumeSource{
				Sources: []v1.VolumeProjection{
					{}, {},
				},
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest41() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Projected: &v1.ProjectedVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest42() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ScaleIO: &v1.ScaleIOVolumeSource{
				ReadOnly: true,
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ScaleIO: &v1.ScaleIOVolumeSource{
				ReadOnly: false,
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest43() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			ScaleIO: &v1.ScaleIOVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest44() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			StorageOS: &v1.StorageOSVolumeSource{
				ReadOnly: true,
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			StorageOS: &v1.StorageOSVolumeSource{
				ReadOnly: false,
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest45() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			StorageOS: &v1.StorageOSVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func initVolumesForTest46() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Ephemeral: &v1.EphemeralVolumeSource{
				ReadOnly: true,
			},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Ephemeral: &v1.EphemeralVolumeSource{
				ReadOnly: false,
			},
		},
	}
	return volume1, volume2
}

func initVolumesForTest47() (v1.Volume, v1.Volume) {

	volume1 := v1.Volume{
		VolumeSource: v1.VolumeSource{
			Ephemeral: &v1.EphemeralVolumeSource{},
		},
	}
	volume2 := v1.Volume{
		VolumeSource: v1.VolumeSource{},
	}
	return volume1, volume2
}

func TestCompareVolumes(t *testing.T) {
	cleanCtx := initCtx()

	ctx := newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 := initVolumesForTest1()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage := diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch := diff.BatchFromContext(ctx)
	diffs := batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVolumeDifferentNames) {
				t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But it was returned: %s", ErrorVolumeDifferentNames.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s. 'name' vs 'diffName''. But the function found no errors", ErrorVolumeDifferentNames.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest2()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorMissingHostPathVolumeSource) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorMissingHostPathVolumeSource.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorMissingHostPathVolumeSource.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest3()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesEmptyDir) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesEmptyDir.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesEmptyDir.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest4()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesSecrets) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesSecrets.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesSecrets.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest5()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesNFS) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesNFS.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesNFS.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest6()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingPersistentVolumeClaim) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingPersistentVolumeClaim.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingPersistentVolumeClaim.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest7()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesDownwardAPI) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesDownwardAPI.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesDownwardAPI.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest8()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesConfigMap) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesConfigMap.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesConfigMap.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest9()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPodMissingVolumesCSI) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPodMissingVolumesCSI.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPodMissingVolumesCSI.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest10()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorISCSIDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorISCSIDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorISCSIDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest11()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorISCSIMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorISCSIMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorISCSIMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest12()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorCephFSDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorCephFSDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorCephFSDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest13()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorCephFSMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorCephFSMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorCephFSMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest14()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorGCEPersistentDiskDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorGCEPersistentDiskDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorGCEPersistentDiskDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest15()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorGCEPersistentDiskMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorGCEPersistentDiskMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorGCEPersistentDiskMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest16()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAWSElasticBlockStoreDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAWSElasticBlockStoreDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAWSElasticBlockStoreDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest17()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAWSElasticBlockStoreMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAWSElasticBlockStoreMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAWSElasticBlockStoreMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest18()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorGlusterfsDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorGlusterfsDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorGlusterfsDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest19()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorGlusterfsMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorGlusterfsMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorGlusterfsMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest20()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorRBDDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorRBDDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorRBDDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest21()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorRBDMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorRBDMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorRBDMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest22()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFlexVolumeDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFlexVolumeDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFlexVolumeDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest23()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFlexVolumeMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFlexVolumeMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFlexVolumeMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest24()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorCinderDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorCinderDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorCinderDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest25()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorCinderMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorCinderMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorCinderMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest26()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFlockerDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFlockerDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFlockerDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest27()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFlockerMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFlockerMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFlockerMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest28()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFCDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFCDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFCDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest29()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorFCMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorFCMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorFCMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest30()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAzureFileDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAzureFileDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAzureFileDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest31()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAzureFileMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAzureFileMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAzureFileMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest32()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVsphereVolumeDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorVsphereVolumeDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorVsphereVolumeDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest33()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorVsphereVolumeMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorVsphereVolumeMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorVsphereVolumeMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest34()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorQuobyteDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorQuobyteDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorQuobyteDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest35()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorQuobyteMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorQuobyteMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorQuobyteMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest36()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAzureDiskDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAzureDiskDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAzureDiskDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest37()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorAzureDiskMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorAzureDiskMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorAzureDiskMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest38()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPhotonPersistentDiskDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPhotonPersistentDiskDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPhotonPersistentDiskDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest39()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorPhotonPersistentDiskMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorPhotonPersistentDiskMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorPhotonPersistentDiskMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest40()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorProjectedDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorProjectedDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorProjectedDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest41()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorProjectedMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorProjectedMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorProjectedMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest42()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorScaleIODifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorScaleIODifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorScaleIODifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest43()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorScaleIOMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorScaleIOMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorScaleIOMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest44()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorStorageOSDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorStorageOSDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorStorageOSDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest45()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorStorageOSMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorStorageOSMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorStorageOSMissing.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest46()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorEphemeralDifferent) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorEphemeralDifferent.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorEphemeralDifferent.Error())
	}

	ctx = newCtxWithCleanStorage(cleanCtx)
	volume1, volume2 = initVolumesForTest47()

	CompareVolumes(ctx, volume1, volume2)

	diffStorage = diff.StorageFromContext(ctx)
	diffStorage.Finalize(ctx)
	batch = diff.BatchFromContext(ctx)
	diffs = batch.GetDiffs()

	if diffs != nil {
		if len(diffs) == 1 {
			if !errors.Is(diffs[0].Msg.(error), ErrorEphemeralMissing) {
				t.Errorf("Error expected: '%s'. But it was returned: %s", ErrorEphemeralMissing.Error(), diffs[0].Msg)
			}
		} else {
			t.Errorf("1 error was expected. But it was returned: %d", len(diffs))
		}
	} else {
		t.Errorf("Error expected: '%s'. But the function found no errors", ErrorEphemeralMissing.Error())
	}
}
