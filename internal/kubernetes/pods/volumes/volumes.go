package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"reflect"

	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumes(ctx context.Context, volume1, volume2 v1.Volume) error {
	var (
		log = logging.FromContext(ctx)

		diffsBatch = diff.BatchFromContext(ctx)
	)

	if volume1.Name != volume2.Name {
		//log.Warnf("%s: %s vs %s", ErrorVolumeDifferentNames.Error(), volume1.Name, volume2.Name)
		diffsBatch.Add(ctx, true, "%s: %s vs %s", ErrorVolumeDifferentNames.Error(), volume1.Name, volume2.Name)
		return nil
	}

	log = log.With(zap.String("volumeName", volume1.Name))

	if volume1.HostPath != nil && volume2.HostPath != nil {

		CompareVolumeHostPath(ctx, volume1.HostPath, volume2.HostPath)

	} else if volume1.HostPath != nil || volume2.HostPath != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorMissingHostPathVolumeSource.Error())
	}

	if volume1.VolumeSource.EmptyDir != nil && volume2.VolumeSource.EmptyDir != nil {

		CompareVolumeEmptyDir(ctx, volume1.VolumeSource.EmptyDir, volume2.VolumeSource.EmptyDir)

	} else if volume1.VolumeSource.EmptyDir != nil || volume2.VolumeSource.EmptyDir != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesEmptyDir.Error())
	}

	if volume1.Secret != nil && volume2.Secret != nil {

		CompareVolumeSecret(ctx, volume1.Secret, volume2.Secret)

	} else if volume1.Secret != nil || volume2.Secret != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesSecrets.Error())
	}

	if volume1.NFS != nil && volume2.NFS != nil {

		CompareVolumeNFS(ctx, volume1.NFS, volume2.NFS)

	} else if volume1.NFS != nil || volume2.NFS != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesNFS.Error())
	}

	if volume1.PersistentVolumeClaim != nil && volume2.PersistentVolumeClaim != nil {

		CompareVolumePersistentVolumeClaim(ctx, volume1.PersistentVolumeClaim, volume2.PersistentVolumeClaim)

	} else if volume1.PersistentVolumeClaim != nil || volume2.PersistentVolumeClaim != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingPersistentVolumeClaim.Error())
	}

	if volume1.DownwardAPI != nil && volume2.DownwardAPI != nil {

		CompareVolumeDownwardAPI(ctx, volume1.DownwardAPI, volume2.DownwardAPI)

	} else if volume1.DownwardAPI != nil || volume2.DownwardAPI != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesDownwardAPI.Error())
	}

	if volume1.VolumeSource.ConfigMap != nil && volume2.VolumeSource.ConfigMap != nil {

		CompareVolumeConfigMap(ctx, volume1.VolumeSource.ConfigMap, volume2.VolumeSource.ConfigMap)

	} else if volume1.VolumeSource.ConfigMap != nil || volume2.VolumeSource.ConfigMap != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesConfigMap.Error())
	}

	if volume1.CSI != nil && volume2.CSI != nil {

		CompareVolumeCSI(ctx, volume1.CSI, volume2.CSI)

	} else if volume1.CSI != nil || volume2.CSI != nil {

		diffsBatch.Add(ctx, false, "%s", ErrorPodMissingVolumesCSI.Error())

	}

	if volume1.ISCSI != nil && volume2.ISCSI != nil {
		if !reflect.DeepEqual(*volume1.ISCSI, *volume2.ISCSI) {
			diffsBatch.Add(ctx, false, "%s", ErrorISCSIDifferent.Error())
		}
	} else if volume1.ISCSI != nil || volume2.ISCSI != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorISCSIMissing.Error())
	}

	if volume1.CephFS != nil && volume2.CephFS != nil {
		if !reflect.DeepEqual(*volume1.CephFS, *volume2.CephFS) {
			diffsBatch.Add(ctx, false, "%s", ErrorCephFSDifferent.Error())
		}
	} else if volume1.CephFS != nil || volume2.CephFS != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorCephFSMissing.Error())
	}

	if volume1.GCEPersistentDisk != nil && volume2.GCEPersistentDisk != nil {
		if !reflect.DeepEqual(*volume1.GCEPersistentDisk, *volume2.GCEPersistentDisk) {
			diffsBatch.Add(ctx, false, "%s", ErrorGCEPersistentDiskDifferent.Error())
		}
	} else if volume1.GCEPersistentDisk != nil || volume2.GCEPersistentDisk != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorGCEPersistentDiskMissing.Error())
	}

	if volume1.AWSElasticBlockStore != nil && volume2.AWSElasticBlockStore != nil {
		if !reflect.DeepEqual(*volume1.AWSElasticBlockStore, *volume2.AWSElasticBlockStore) {
			diffsBatch.Add(ctx, false, "%s", ErrorAWSElasticBlockStoreDifferent.Error())
		}
	} else if volume1.AWSElasticBlockStore != nil || volume2.AWSElasticBlockStore != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorAWSElasticBlockStoreMissing.Error())
	}

	if volume1.Glusterfs != nil && volume2.Glusterfs != nil {
		if !reflect.DeepEqual(*volume1.Glusterfs, *volume2.Glusterfs) {
			diffsBatch.Add(ctx, false, "%s", ErrorGlusterfsDifferent.Error())
		}
	} else if volume1.Glusterfs != nil || volume2.Glusterfs != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorGlusterfsMissing.Error())
	}

	if volume1.RBD != nil && volume2.RBD != nil {
		if !reflect.DeepEqual(*volume1.RBD, *volume2.RBD) {
			diffsBatch.Add(ctx, false, "%s", ErrorRBDDifferent.Error())
		}
	} else if volume1.RBD != nil || volume2.RBD != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorRBDMissing.Error())
	}

	if volume1.FlexVolume != nil && volume2.FlexVolume != nil {
		if !reflect.DeepEqual(*volume1.FlexVolume, *volume2.FlexVolume) {
			diffsBatch.Add(ctx, false, "%s", ErrorFlexVolumeDifferent.Error())
		}
	} else if volume1.FlexVolume != nil || volume2.FlexVolume != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorFlexVolumeMissing.Error())
	}

	if volume1.Cinder != nil && volume2.Cinder != nil {
		if !reflect.DeepEqual(*volume1.Cinder, *volume2.Cinder) {
			diffsBatch.Add(ctx, false, "%s", ErrorCinderDifferent.Error())
		}

	} else if volume1.Cinder != nil || volume2.Cinder != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorCinderMissing.Error())
	}

	if volume1.Flocker != nil && volume2.Flocker != nil {
		if !reflect.DeepEqual(*volume1.Flocker, *volume2.Flocker) {
			diffsBatch.Add(ctx, false, "%s", ErrorFlockerDifferent.Error())
		}

	} else if volume1.Flocker != nil || volume2.Flocker != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorFlockerMissing.Error())
	}

	if volume1.FC != nil && volume2.FC != nil {
		if !reflect.DeepEqual(*volume1.FC, *volume2.FC) {
			diffsBatch.Add(ctx, false, "%s", ErrorFCDifferent.Error())
		}

	} else if volume1.FC != nil || volume2.FC != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorFCMissing.Error())
	}

	if volume1.AzureFile != nil && volume2.AzureFile != nil {
		if !reflect.DeepEqual(*volume1.AzureFile, *volume2.AzureFile) {
			diffsBatch.Add(ctx, false, "%s", ErrorAzureFileDifferent.Error())
		}

	} else if volume1.AzureFile != nil || volume2.AzureFile != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorAzureFileMissing.Error())
	}

	if volume1.VsphereVolume != nil && volume2.VsphereVolume != nil {
		if !reflect.DeepEqual(*volume1.VsphereVolume, *volume2.VsphereVolume) {
			diffsBatch.Add(ctx, false, "%s", ErrorVsphereVolumeDifferent.Error())
		}

	} else if volume1.VsphereVolume != nil || volume2.VsphereVolume != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorVsphereVolumeMissing.Error())
	}

	if volume1.Quobyte != nil && volume2.Quobyte != nil {
		if !reflect.DeepEqual(*volume1.Quobyte, *volume2.Quobyte) {
			diffsBatch.Add(ctx, false, "%s", ErrorQuobyteDifferent.Error())
		}

	} else if volume1.Quobyte != nil || volume2.Quobyte != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorQuobyteMissing.Error())
	}

	if volume1.AzureDisk != nil && volume2.AzureDisk != nil {
		if !reflect.DeepEqual(*volume1.AzureDisk, *volume2.AzureDisk) {
			diffsBatch.Add(ctx, false, "%s", ErrorAzureDiskDifferent.Error())
		}

	} else if volume1.AzureDisk != nil || volume2.AzureDisk != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorAzureDiskMissing.Error())
	}

	if volume1.PhotonPersistentDisk != nil && volume2.PhotonPersistentDisk != nil {
		if !reflect.DeepEqual(*volume1.PhotonPersistentDisk, *volume2.PhotonPersistentDisk) {
			diffsBatch.Add(ctx, false, "%s", ErrorPhotonPersistentDiskDifferent.Error())
		}

	} else if volume1.PhotonPersistentDisk != nil || volume2.PhotonPersistentDisk != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPhotonPersistentDiskMissing.Error())
	}

	if volume1.Projected != nil && volume2.Projected != nil {
		if !reflect.DeepEqual(*volume1.Projected, *volume2.Projected) {
			diffsBatch.Add(ctx, false, "%s", ErrorProjectedDifferent.Error())
		}

	} else if volume1.Projected != nil || volume2.Projected != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorProjectedMissing.Error())
	}

	if volume1.PortworxVolume != nil && volume2.PortworxVolume != nil {
		if !reflect.DeepEqual(*volume1.PortworxVolume, *volume2.PortworxVolume) {
			diffsBatch.Add(ctx, false, "%s", ErrorPortworxVolumeDifferent.Error())
		}

	} else if volume1.PortworxVolume != nil || volume2.PortworxVolume != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorPortworxVolumeMissing.Error())
	}

	if volume1.ScaleIO != nil && volume2.ScaleIO != nil {
		if !reflect.DeepEqual(*volume1.ScaleIO, *volume2.ScaleIO) {
			diffsBatch.Add(ctx, false, "%s", ErrorScaleIODifferent.Error())
		}

	} else if volume1.ScaleIO != nil || volume2.ScaleIO != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorScaleIOMissing.Error())
	}

	if volume1.StorageOS != nil && volume2.StorageOS != nil {
		if !reflect.DeepEqual(*volume1.StorageOS, *volume2.StorageOS) {
			diffsBatch.Add(ctx, false, "%s", ErrorStorageOSDifferent.Error())
		}

	} else if volume1.StorageOS != nil || volume2.StorageOS != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorStorageOSMissing.Error())
	}

	if volume1.Ephemeral != nil && volume2.Ephemeral != nil {
		if !reflect.DeepEqual(*volume1.Ephemeral, *volume2.Ephemeral) {
			diffsBatch.Add(ctx, false, "%s", ErrorEphemeralDifferent.Error())
		}

	} else if volume1.Ephemeral != nil || volume2.Ephemeral != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorEphemeralMissing.Error())
	}

	return nil
}
