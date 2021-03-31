package volumes

import (
	"context"
	"reflect"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumes(ctx context.Context, volume1, volume2 v1.Volume) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	//return nil, nil

	if volume1.Name != volume2.Name {
		log.Warnf("%s: %s vs %s", ErrorVolumeDifferentNames.Error(), volume1.Name, volume2.Name)
		return nil, nil
	}

	log = log.With(zap.String("volumeName", volume1.Name))

	if volume1.HostPath != nil && volume2.HostPath != nil {

		_, _ = CompareVolumeHostPath(ctx, volume1.HostPath, volume2.HostPath)

	} else if volume1.HostPath != nil || volume2.HostPath != nil {
		log.Warnf("%s", ErrorMissingHostPathVolumeSource.Error())
	}

	if volume1.VolumeSource.EmptyDir != nil && volume2.VolumeSource.EmptyDir != nil {

		_, _ = CompareVolumeEmptyDir(ctx, volume1.VolumeSource.EmptyDir, volume2.VolumeSource.EmptyDir)

	} else if volume1.VolumeSource.EmptyDir != nil || volume2.VolumeSource.EmptyDir != nil {
		log.Warnf("%s", ErrorPodMissingVolumesEmptyDir.Error())
	}

	if volume1.Secret != nil && volume2.Secret != nil {

		_, _ = CompareVolumeSecret(ctx, volume1.Secret, volume2.Secret)

	} else if volume1.Secret != nil || volume2.Secret != nil {
		log.Warnf("%s", ErrorPodMissingVolumesSecrets.Error())
	}

	if volume1.NFS != nil && volume2.NFS != nil {

		_, _ = CompareVolumeNFS(ctx, volume1.NFS, volume2.NFS)

	} else if volume1.NFS != nil || volume2.NFS != nil {
		log.Warnf("%s", ErrorPodMissingVolumesNFS.Error())
	}

	if volume1.PersistentVolumeClaim != nil && volume2.PersistentVolumeClaim != nil {

		_, _ = CompareVolumePersistentVolumeClaim(ctx, volume1.PersistentVolumeClaim, volume2.PersistentVolumeClaim)

	} else if volume1.PersistentVolumeClaim != nil || volume2.PersistentVolumeClaim != nil {
		log.Warnf("%s", ErrorPodMissingPersistentVolumeClaim.Error())
	}

	if volume1.DownwardAPI != nil && volume2.DownwardAPI != nil {

		_, _ = CompareVolumeDownwardAPI(ctx, volume1.DownwardAPI, volume2.DownwardAPI)

	} else if volume1.DownwardAPI != nil || volume2.DownwardAPI != nil {
		log.Warnf("%s", ErrorPodMissingVolumesDownwardAPI.Error())
	}

	if volume1.VolumeSource.ConfigMap != nil && volume2.VolumeSource.ConfigMap != nil {

		_, _ = CompareVolumeConfigMap(ctx, volume1.VolumeSource.ConfigMap, volume2.VolumeSource.ConfigMap)

	} else if volume1.VolumeSource.ConfigMap != nil || volume2.VolumeSource.ConfigMap != nil {
		log.Warnf("%s", ErrorPodMissingVolumesConfigMap.Error())
	}

	if volume1.CSI != nil && volume2.CSI != nil {

		_, _ = CompareVolumeCSI(ctx, volume1.CSI, volume2.CSI)

	} else if volume1.CSI != nil || volume2.CSI != nil {

		log.Warnf("%s", ErrorPodMissingVolumesCSI.Error())

	}

	if volume1.ISCSI != nil && volume2.ISCSI != nil {
		if !reflect.DeepEqual(*volume1.ISCSI, *volume2.ISCSI) {
			log.Warnf("%s", ErrorISCSIDifferent.Error())
		}
	} else if volume1.ISCSI != nil || volume2.ISCSI != nil {
		log.Warnf("%s", ErrorISCSIMissing.Error())
	}

	if volume1.CephFS != nil && volume2.CephFS != nil {
		if !reflect.DeepEqual(*volume1.CephFS, *volume2.CephFS) {
			log.Warnf("%s", ErrorCephFSDifferent.Error())
		}
	} else if volume1.CephFS != nil || volume2.CephFS != nil {
		log.Warnf("%s", ErrorCephFSMissing.Error())
	}

	if volume1.GCEPersistentDisk != nil && volume2.GCEPersistentDisk != nil {
		if !reflect.DeepEqual(*volume1.GCEPersistentDisk, *volume2.GCEPersistentDisk) {
			log.Warnf("%s", ErrorGCEPersistentDiskDifferent.Error())
		}
	} else if volume1.GCEPersistentDisk != nil || volume2.GCEPersistentDisk != nil {
		log.Warnf("%s", ErrorGCEPersistentDiskMissing.Error())
	}

	if volume1.AWSElasticBlockStore != nil && volume2.AWSElasticBlockStore != nil {
		if !reflect.DeepEqual(*volume1.AWSElasticBlockStore, *volume2.AWSElasticBlockStore) {
			log.Warnf("%s", ErrorAWSElasticBlockStoreDifferent.Error())
		}
	} else if volume1.AWSElasticBlockStore != nil || volume2.AWSElasticBlockStore != nil {
		log.Warnf("%s", ErrorAWSElasticBlockStoreMissing.Error())
	}

	if volume1.Glusterfs != nil && volume2.Glusterfs != nil {
		if !reflect.DeepEqual(*volume1.Glusterfs, *volume2.Glusterfs) {
			log.Warnf("%s", ErrorGlusterfsDifferent.Error())
		}
	} else if volume1.Glusterfs != nil || volume2.Glusterfs != nil {
		log.Warnf("%s", ErrorGlusterfsMissing.Error())
	}

	if volume1.RBD != nil && volume2.RBD != nil {
		if !reflect.DeepEqual(*volume1.RBD, *volume2.RBD) {
			log.Warnf("%s", ErrorRBDDifferent.Error())
		}
	} else if volume1.RBD != nil || volume2.RBD != nil {
		log.Warnf("%s", ErrorRBDMissing.Error())
	}

	if volume1.FlexVolume != nil && volume2.FlexVolume != nil {
		if !reflect.DeepEqual(*volume1.FlexVolume, *volume2.FlexVolume) {
			log.Warnf("%s", ErrorFlexVolumeDifferent.Error())
		}
	} else if volume1.FlexVolume != nil || volume2.FlexVolume != nil {
		log.Warnf("%s", ErrorFlexVolumeMissing.Error())
	}

	if volume1.Cinder != nil && volume2.Cinder != nil {
		if !reflect.DeepEqual(*volume1.Cinder, *volume2.Cinder) {
			log.Warnf("%s", ErrorCinderDifferent.Error())
		}

	} else if volume1.Cinder != nil || volume2.Cinder != nil {
		log.Warnf("%s", ErrorCinderMissing.Error())
	}

	if volume1.Flocker != nil && volume2.Flocker != nil {
		if !reflect.DeepEqual(*volume1.Flocker, *volume2.Flocker) {
			log.Warnf("%s", ErrorFlockerDifferent.Error())
		}

	} else if volume1.Flocker != nil || volume2.Flocker != nil {
		log.Warnf("%s", ErrorFlockerMissing.Error())
	}

	if volume1.FC != nil && volume2.FC != nil {
		if !reflect.DeepEqual(*volume1.FC, *volume2.FC) {
			log.Warnf("%s", ErrorFCDifferent.Error())
		}

	} else if volume1.FC != nil || volume2.FC != nil {
		log.Warnf("%s", ErrorFCMissing.Error())
	}

	if volume1.AzureFile != nil && volume2.AzureFile != nil {
		if !reflect.DeepEqual(*volume1.AzureFile, *volume2.AzureFile) {
			log.Warnf("%s", ErrorAzureFileDifferent.Error())
		}

	} else if volume1.AzureFile != nil || volume2.AzureFile != nil {
		log.Warnf("%s", ErrorAzureFileMissing.Error())
	}

	if volume1.VsphereVolume != nil && volume2.VsphereVolume != nil {
		if !reflect.DeepEqual(*volume1.VsphereVolume, *volume2.VsphereVolume) {
			log.Warnf("%s", ErrorVsphereVolumeDifferent.Error())
		}

	} else if volume1.VsphereVolume != nil || volume2.VsphereVolume != nil {
		log.Warnf("%s", ErrorVsphereVolumeMissing.Error())
	}

	if volume1.Quobyte != nil && volume2.Quobyte != nil {
		if !reflect.DeepEqual(*volume1.Quobyte, *volume2.Quobyte) {
			log.Warnf("%s", ErrorQuobyteDifferent.Error())
		}

	} else if volume1.Quobyte != nil || volume2.Quobyte != nil {
		log.Warnf("%s", ErrorQuobyteMissing.Error())
	}

	if volume1.AzureDisk != nil && volume2.AzureDisk != nil {
		if !reflect.DeepEqual(*volume1.AzureDisk, *volume2.AzureDisk) {
			log.Warnf("%s", ErrorAzureDiskDifferent.Error())
		}

	} else if volume1.AzureDisk != nil || volume2.AzureDisk != nil {
		log.Warnf("%s", ErrorAzureDiskMissing.Error())
	}

	if volume1.PhotonPersistentDisk != nil && volume2.PhotonPersistentDisk != nil {
		if !reflect.DeepEqual(*volume1.PhotonPersistentDisk, *volume2.PhotonPersistentDisk) {
			log.Warnf("%s", ErrorPhotonPersistentDiskDifferent.Error())
		}

	} else if volume1.PhotonPersistentDisk != nil || volume2.PhotonPersistentDisk != nil {
		log.Warnf("%s", ErrorPhotonPersistentDiskMissing.Error())
	}

	if volume1.Projected != nil && volume2.Projected != nil {
		if !reflect.DeepEqual(*volume1.Projected, *volume2.Projected) {
			log.Warnf("%s", ErrorProjectedDifferent.Error())
		}

	} else if volume1.Projected != nil || volume2.Projected != nil {
		log.Warnf("%s", ErrorProjectedMissing.Error())
	}

	if volume1.PortworxVolume != nil && volume2.PortworxVolume != nil {
		if !reflect.DeepEqual(*volume1.PortworxVolume, *volume2.PortworxVolume) {
			log.Warnf("%s", ErrorPortworxVolumeDifferent.Error())
		}

	} else if volume1.PortworxVolume != nil || volume2.PortworxVolume != nil {
		log.Warnf("%s", ErrorPortworxVolumeMissing.Error())
	}

	if volume1.ScaleIO != nil && volume2.ScaleIO != nil {
		if !reflect.DeepEqual(*volume1.ScaleIO, *volume2.ScaleIO) {
			log.Warnf("%s", ErrorScaleIODifferent.Error())
		}

	} else if volume1.ScaleIO != nil || volume2.ScaleIO != nil {
		log.Warnf("%s", ErrorScaleIOMissing.Error())
	}

	if volume1.StorageOS != nil && volume2.StorageOS != nil {
		if !reflect.DeepEqual(*volume1.StorageOS, *volume2.StorageOS) {
			log.Warnf("%s", ErrorStorageOSDifferent.Error())
		}

	} else if volume1.StorageOS != nil || volume2.StorageOS != nil {
		log.Warnf("%s", ErrorStorageOSMissing.Error())
	}

	if volume1.Ephemeral != nil && volume2.Ephemeral != nil {
		if !reflect.DeepEqual(*volume1.Ephemeral, *volume2.Ephemeral) {
			log.Warnf("%s", ErrorEphemeralDifferent.Error())
		}

	} else if volume1.Ephemeral != nil || volume2.Ephemeral != nil {
		log.Warnf("%s", ErrorEphemeralMissing.Error())
	}

	return nil, nil
}
