package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

func CompareVolumes(ctx context.Context, volume1, volume2 v1.Volume) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if volume1.Name != volume2.Name {
		log.Warnf("%s: %s vs %s", ErrorVolumeDifferentNames.Error(), volume1.Name, volume2.Name)
		return nil, nil
	}

	log = log.With(zap.String("volumeName", volume1.Name))

	if volume1.HostPath != nil && volume2.HostPath != nil {

		_, _ = CompareVolumeHostPath(ctx, volume1.HostPath, volume2.HostPath)
		//if volume1.HostPath.Path != volume2.HostPath.Path {
		//	log.Warnf("%s: %s vs %s", ErrorHostPath.Error(), volume1.HostPath.Path, volume2.HostPath.Path)
		//}
		//
		//if volume1.HostPath.Type != nil && volume2.HostPath.Type != nil {
		//	if *volume1.HostPath.Type != *volume2.HostPath.Type {
		//		log.Warnf("%s: %s vs %s", ErrorHostPathType.Error(), *volume1.HostPath.Type, *volume2.HostPath.Type)
		//	}
		//} else if volume1.HostPath.Type != nil || volume2.HostPath.Type != nil {
		//	log.Warnf("%s", ErrorMissingHostPathType.Error())
		//}

	} else if volume1.HostPath != nil || volume2.HostPath != nil {
		log.Warnf("%s", ErrorMissingHostPathVolumeSource.Error())
	}

	if volume1.VolumeSource.EmptyDir != nil && volume2.VolumeSource.EmptyDir != nil {

		_, _ = CompareVolumeEmptyDir(ctx, volume1.VolumeSource.EmptyDir, volume2.VolumeSource.EmptyDir)
		//if volume1.VolumeSource.EmptyDir.Medium != volume2.VolumeSource.EmptyDir.Medium {
		//	log.Warnf("%s. %s vs %s", ErrorVolumeEmptyDirMedium.Error(), volume1.VolumeSource.EmptyDir.Medium, volume2.VolumeSource.EmptyDir.Medium)
		//}
		//
		//if volume1.VolumeSource.EmptyDir.SizeLimit != nil && volume2.VolumeSource.EmptyDir.SizeLimit != nil {
		//
		//	if volume1.VolumeSource.EmptyDir.SizeLimit.Format != volume2.VolumeSource.EmptyDir.SizeLimit.Format {
		//		log.Warnf("%s. %s vs %s", ErrorVolumeEmptyDirSizeLimitFormat.Error(), volume1.VolumeSource.EmptyDir.SizeLimit.Format, volume2.VolumeSource.EmptyDir.SizeLimit.Format)
		//	}
		//
		//} else if volume1.VolumeSource.EmptyDir.SizeLimit != nil || volume2.VolumeSource.EmptyDir.SizeLimit != nil {
		//	log.Warnf("%s", ErrorVolumeEmptyDirSizeLimit.Error())
		//}

	} else if volume1.VolumeSource.EmptyDir != nil || volume2.VolumeSource.EmptyDir != nil {
		log.Warnf("%s", ErrorPodMissingVolumesEmptyDir.Error())
	}

	if volume1.Secret != nil && volume2.Secret != nil {

		_, _ = CompareVolumeSecret(ctx, volume1.Secret, volume2.Secret)
		//if volume1.Secret.SecretName != volume2.Secret.SecretName {
		//	log.Warnf("%s. %s vs %s", ErrorVolumeSecretName.Error(), volume1.Secret.SecretName, volume2.Secret.SecretName)
		//}
		//
		//if volume1.Secret.Optional != nil && volume2.Secret.Optional != nil {
		//
		//	if *volume1.Secret.Optional != *volume2.Secret.Optional {
		//		log.Warnf("%s. %t vs %t", ErrorVolumeSecretOptional.Error(), *volume1.Secret.Optional, *volume2.Secret.Optional)
		//	}
		//
		//} else if volume1.Secret.Optional != nil || volume2.Secret.Optional != nil {
		//	log.Warnf("%s", ErrorMissingVolumeSecretOptional.Error())
		//}
		//
		//if volume1.Secret.DefaultMode != nil && volume2.Secret.DefaultMode != nil {
		//
		//	if *volume1.Secret.DefaultMode != *volume2.Secret.DefaultMode {
		//		log.Warnf("%s. %d vs %d", ErrorVolumeSecretDefaultMode.Error(), *volume1.Secret.DefaultMode, *volume2.Secret.DefaultMode)
		//	}
		//
		//} else if volume1.Secret.DefaultMode != nil || volume2.Secret.DefaultMode != nil {
		//	log.Warnf("%s", ErrorMissingVolumeSecretDefaultMode.Error())
		//}
		//
		//if len(volume1.Secret.Items) != len(volume2.Secret.Items) {
		//	log.Warnf("%s. %d vs %d", ErrorVolumeSecretItemsLen.Error(), len(volume1.Secret.Items), len(volume2.Secret.Items))
		//} else {
		//
		//	for index, item := range volume1.Secret.Items {
		//
		//		if item.Path != volume2.Secret.Items[index].Path {
		//			log.Warnf("%s. %s vs %s", ErrorVolumeSecretItemsPath.Error(), item.Path, volume2.Secret.Items[index].Path)
		//		}
		//
		//		if item.Key != volume2.Secret.Items[index].Key {
		//			log.Warnf("%s. %s vs %s", ErrorVolumeSecretItemsKey.Error(), item.Key, volume2.Secret.Items[index].Key)
		//		}
		//
		//		if item.Mode != nil && volume2.Secret.Items[index].Mode != nil {
		//			if *item.Mode != *volume2.Secret.Items[index].Mode {
		//				log.Warnf("%s. %d vs %d", ErrorVolumeSecretItemsMode.Error(), item.Mode, volume2.Secret.Items[index].Mode)
		//			}
		//		} else if item.Mode != nil || volume2.Secret.Items[index].Mode != nil {
		//			log.Warnf("%s", ErrorMissingVolumeSecretItemsMode.Error())
		//		}
		//
		//	}
		//
		//}

	} else if volume1.Secret != nil || volume2.Secret != nil {
		log.Warnf("%s", ErrorPodMissingVolumesSecrets.Error())
	}

	if volume1.NFS != nil && volume2.NFS != nil {

		_, _ = CompareVolumeNFS(ctx, volume1.NFS, volume2.NFS)

		//if volume1.NFS.ReadOnly != volume2.NFS.ReadOnly {
		//	log.Warnf("%s. %t vs %t", ErrorVolumeNFSReadOnly.Error(), volume1.NFS.ReadOnly, volume2.NFS.ReadOnly)
		//}
		//
		//if volume1.NFS.Path != volume2.NFS.Path {
		//	log.Warnf("%s. %s vs %s", ErrorVolumeNFSPath.Error(), volume1.NFS.Path, volume2.NFS.Path)
		//}
		//
		//if volume1.NFS.Server != volume2.NFS.Server {
		//	log.Warnf("%s. %s vs %s", ErrorVolumeNFSServer.Error(), volume1.NFS.Server, volume2.NFS.Server)
		//}

	} else if volume1.NFS != nil || volume2.NFS != nil {
		log.Warnf("%s", ErrorPodMissingVolumesNFS.Error())
	}

	if volume1.PersistentVolumeClaim != nil && volume2.PersistentVolumeClaim != nil {

		_, _ = CompareVolumePersistentVolumeClaim(ctx, volume1.PersistentVolumeClaim, volume2.PersistentVolumeClaim)
		//if volume1.PersistentVolumeClaim.ReadOnly != volume2.PersistentVolumeClaim.ReadOnly {
		//	log.Warnf("%s. %t vs %t", ErrorPersistentVolumeClaimReadOnly.Error(), volume1.PersistentVolumeClaim.ReadOnly, volume2.PersistentVolumeClaim.ReadOnly)
		//}
		//
		//if volume1.PersistentVolumeClaim.ClaimName != volume2.PersistentVolumeClaim.ClaimName {
		//	log.Warnf("%s. %s vs %s", ErrorPersistentVolumeClaimName.Error(), volume1.PersistentVolumeClaim.ClaimName, volume2.PersistentVolumeClaim.ClaimName)
		//}

	} else if volume1.PersistentVolumeClaim != nil || volume2.PersistentVolumeClaim != nil {
		log.Warnf("%s", ErrorPodMissingPersistentVolumeClaim.Error())
	}

	if volume1.DownwardAPI != nil && volume2.DownwardAPI != nil {

		_, _ = CompareVolumeDownwardAPI(ctx, volume1.DownwardAPI, volume2.DownwardAPI)
		//if volume1.DownwardAPI.DefaultMode != nil && volume2.DownwardAPI.DefaultMode != nil {
		//
		//	if *volume1.DownwardAPI.DefaultMode != *volume2.DownwardAPI.DefaultMode {
		//		log.Warnf("%s. %d vs %d", ErrorDownwardAPIDefaultMode.Error(), *volume1.DownwardAPI.DefaultMode , *volume2.DownwardAPI.DefaultMode )
		//	}
		//
		//} else if volume1.DownwardAPI.DefaultMode != nil || volume2.DownwardAPI.DefaultMode != nil {
		//	log.Warnf("%s", ErrorMissingDownwardAPIDefaultMode.Error())
		//}
		//
		//if len(volume1.DownwardAPI.Items) != len(volume2.DownwardAPI.Items) {
		//	log.Warnf("%s. %d vs %d", ErrorVolumeDownwardAPIItemsLen.Error(), len(volume1.DownwardAPI.Items), len(volume2.DownwardAPI.Items))
		//} else {
		//
		//	for index, item := range volume1.DownwardAPI.Items {
		//
		//		if item.Path != volume2.DownwardAPI.Items[index].Path {
		//			log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsPath.Error(), item.Path, volume2.DownwardAPI.Items[index].Path)
		//		}
		//
		//		if item.Mode != nil && volume2.DownwardAPI.Items[index].Mode != nil {
		//
		//			if *item.Mode != *volume2.DownwardAPI.Items[index].Mode {
		//				log.Warnf("%s. %d vs %d", ErrorDownwardAPIItemsMode.Error(), *item.Mode, *volume2.DownwardAPI.Items[index].Mode)
		//			}
		//
		//		} else if item.Mode != nil || volume2.DownwardAPI.Items[index].Mode != nil {
		//			log.Warnf("%s", ErrorMissingDownwardAPIItemsMode.Error())
		//		}
		//
		//		if item.ResourceFieldRef != nil && volume2.DownwardAPI.Items[index].ResourceFieldRef != nil {
		//
		//			if item.ResourceFieldRef.Resource != volume2.DownwardAPI.Items[index].ResourceFieldRef.Resource {
		//				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefResource.Error(), item.ResourceFieldRef.Resource, volume2.DownwardAPI.Items[index].ResourceFieldRef.Resource)
		//			}
		//
		//			if item.ResourceFieldRef.ContainerName != volume2.DownwardAPI.Items[index].ResourceFieldRef.ContainerName {
		//				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefContainerName.Error(), item.ResourceFieldRef.ContainerName, volume2.DownwardAPI.Items[index].ResourceFieldRef.ContainerName)
		//			}
		//
		//			if item.ResourceFieldRef.Divisor.Format != volume2.DownwardAPI.Items[index].ResourceFieldRef.Divisor.Format {
		//				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefFormat.Error(), item.ResourceFieldRef.Divisor.Format, volume2.DownwardAPI.Items[index].ResourceFieldRef.Divisor.Format)
		//			}
		//
		//		} else if item.ResourceFieldRef != nil || volume2.DownwardAPI.Items[index].ResourceFieldRef != nil {
		//			log.Warnf("%s", ErrorMissingDownwardAPIItemsResourceFieldRef.Error())
		//		}
		//
		//		if item.FieldRef != nil && volume2.DownwardAPI.Items[index].FieldRef != nil {
		//
		//			if item.FieldRef.APIVersion != volume2.DownwardAPI.Items[index].FieldRef.APIVersion {
		//				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsFieldRefAPIVersion.Error(), item.FieldRef.APIVersion, volume2.DownwardAPI.Items[index].FieldRef.APIVersion)
		//			}
		//
		//			if item.FieldRef.FieldPath != volume2.DownwardAPI.Items[index].FieldRef.FieldPath {
		//				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsFieldRefFieldPath.Error(), item.FieldRef.FieldPath, volume2.DownwardAPI.Items[index].FieldRef.FieldPath)
		//			}
		//
		//		} else if item.FieldRef != nil || volume2.DownwardAPI.Items[index].FieldRef != nil {
		//			log.Warnf("%s", ErrorMissingDownwardAPIItemsFieldRef.Error())
		//		}
		//
		//	}
		//}

	} else if volume1.DownwardAPI != nil || volume2.DownwardAPI != nil {
		log.Warnf("%s", ErrorPodMissingVolumesDownwardAPI.Error())
	}

	if volume1.VolumeSource.ConfigMap != nil && volume2.VolumeSource.ConfigMap != nil {

		_, _ = CompareVolumeConfigMap(ctx, volume1.VolumeSource.ConfigMap, volume2.VolumeSource.ConfigMap)
		//if volume1.VolumeSource.ConfigMap.Name != volume2.VolumeSource.ConfigMap.Name {
		//	log.Warnf("%s: %s vs %s", ErrorVolumeConfigMapName.Error(), volume1.VolumeSource.ConfigMap.Name, volume2.VolumeSource.ConfigMap.Name)
		//}
		//
		//if len(volume1.VolumeSource.ConfigMap.Items) != len(volume2.VolumeSource.ConfigMap.Items) {
		//	log.Warnf("%s", ErrorVolumeConfigMapItemsLen.Error())
		//} else {
		//	for i, item := range volume1.VolumeSource.ConfigMap.Items {
		//		if item.Path != volume2.VolumeSource.ConfigMap.Items[i].Path {
		//			log.Warnf("%s. %s vs %s", ErrorVolumeConfigMapPath.Error(), item.Path, volume2.VolumeSource.ConfigMap.Items[i].Path)
		//		}
		//
		//		if item.Key != volume2.VolumeSource.ConfigMap.Items[i].Key {
		//			log.Warnf("%s. %s vs %s", ErrorVolumeConfigMapKey.Error(), item.Key, volume2.VolumeSource.ConfigMap.Items[i].Key)
		//		}
		//
		//		if item.Mode != nil && volume2.VolumeSource.ConfigMap.Items[i].Mode != nil {
		//			if *item.Mode != *volume2.VolumeSource.ConfigMap.Items[i].Mode {
		//				log.Warnf("%s. %d vs %d", ErrorVolumeConfigMapMode.Error(),  *item.Mode, *volume2.VolumeSource.ConfigMap.Items[i].Mode)
		//			}
		//		} else if item.Mode != nil || volume2.VolumeSource.ConfigMap.Items[i].Mode != nil {
		//			log.Warnf("%s", ErrorVolumeConfigMapMode.Error())
		//		}
		//
		//	}
		//}
		//
		//if volume1.VolumeSource.ConfigMap.DefaultMode != nil && volume2.VolumeSource.ConfigMap.DefaultMode != nil {
		//	if *volume1.VolumeSource.ConfigMap.DefaultMode != *volume2.VolumeSource.ConfigMap.DefaultMode {
		//		log.Warnf("%s: %d vs %d", ErrorVolumeConfigMapDefaultMode.Error(), *volume1.VolumeSource.ConfigMap.DefaultMode, *volume2.VolumeSource.ConfigMap.DefaultMode)
		//	}
		//
		//} else if volume1.VolumeSource.ConfigMap.DefaultMode != nil || volume2.VolumeSource.ConfigMap.DefaultMode != nil{
		//	log.Warnf("%s", ErrorVolumeConfigMapDefaultMode.Error())
		//}
		//
		//if volume1.VolumeSource.ConfigMap.Optional != nil && volume2.VolumeSource.ConfigMap.Optional != nil {
		//	if *volume1.VolumeSource.ConfigMap.Optional != *volume2.VolumeSource.ConfigMap.Optional {
		//		log.Warnf("%s: %t vs %t", ErrorVolumeConfigMapOptional.Error(), *volume1.VolumeSource.ConfigMap.Optional, *volume2.VolumeSource.ConfigMap.Optional)
		//	}
		//
		//} else if volume1.VolumeSource.ConfigMap.DefaultMode != nil || volume2.VolumeSource.ConfigMap.DefaultMode != nil{
		//	log.Warnf("%s", ErrorVolumeConfigMapOptional.Error())
		//}

	} else if volume1.VolumeSource.ConfigMap != nil || volume2.VolumeSource.ConfigMap != nil {
		log.Warnf("%s", ErrorPodMissingVolumesConfigMap.Error())
	}

	if volume1.CSI != nil && volume2.CSI != nil {

		_, _ = CompareVolumeCSI(ctx, volume1.CSI, volume2.CSI)

		//if volume1.CSI.Driver != volume2.CSI.Driver {
		//	log.Warnf("%s. %s vs %s", ErrorVolumeCSIDriver.Error(), volume1.CSI.Driver, volume2.CSI.Driver)
		//}
		//
		//if volume1.CSI.ReadOnly != nil && volume2.CSI.ReadOnly != nil {
		//
		//	if *volume1.CSI.ReadOnly != *volume2.CSI.ReadOnly {
		//		log.Warnf("%s. %t vs %t", ErrorVolumeCSIReadOnly.Error(), *volume1.CSI.ReadOnly, *volume2.CSI.ReadOnly)
		//	}
		//
		//} else if volume1.CSI.ReadOnly != nil || volume2.CSI.ReadOnly != nil {
		//	log.Warnf("%s", ErrorMissingCSIReadOnly.Error())
		//}
		//
		//if volume1.CSI.FSType != nil && volume2.CSI.FSType != nil {
		//
		//	if *volume1.CSI.FSType != *volume2.CSI.FSType {
		//		log.Warnf("%s. %s vs %s", ErrorVolumeCSIFSType.Error(), *volume1.CSI.FSType, *volume2.CSI.FSType)
		//	}
		//
		//} else if volume1.CSI.FSType != nil || volume2.CSI.FSType != nil {
		//	log.Warnf("%s", ErrorMissingCSIFSType.Error())
		//}
		//
		//if volume1.CSI.NodePublishSecretRef != nil && volume2.CSI.NodePublishSecretRef != nil {
		//
		//	if volume1.CSI.NodePublishSecretRef.Name != volume2.CSI.NodePublishSecretRef.Name {
		//		log.Warnf("%s. %s vs %s", ErrorVolumeCSIName.Error(), volume1.CSI.NodePublishSecretRef.Name, volume2.CSI.NodePublishSecretRef.Name)
		//	}
		//
		//} else if volume1.CSI.NodePublishSecretRef != nil || volume2.CSI.NodePublishSecretRef != nil {
		//	log.Warnf("%s", ErrorMissingCSINodePublishSecretRef.Error())
		//}
		//
		//if !reflect.DeepEqual(volume1.CSI.VolumeAttributes, volume2.CSI.VolumeAttributes) {
		//	log.Warnf("%s", ErrorVolumeCSIVolumeAttributes.Error())
		//}

	} else if volume1.CSI != nil || volume2.CSI != nil {

		log.Warnf("%s", ErrorPodMissingVolumesCSI.Error())

	}

	if !reflect.DeepEqual(*volume1.ISCSI, *volume2.ISCSI) {
		log.Warnf("%s", ErrorISCSIDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.CephFS, *volume2.CephFS) {
		log.Warnf("%s", ErrorCephFSDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.GCEPersistentDisk, *volume2.GCEPersistentDisk) {
		log.Warnf("%s", ErrorGCEPersistentDiskDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.AWSElasticBlockStore, *volume2.AWSElasticBlockStore) {
		log.Warnf("%s", ErrorAWSElasticBlockStoreDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Glusterfs, *volume2.Glusterfs) {
		log.Warnf("%s", ErrorGlusterfsDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.RBD, *volume2.RBD) {
		log.Warnf("%s", ErrorRBDDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.FlexVolume, *volume2.FlexVolume) {
		log.Warnf("%s", ErrorFlexVolumeDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Cinder, *volume2.Cinder) {
		log.Warnf("%s", ErrorCinderDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Flocker, *volume2.Flocker) {
		log.Warnf("%s", ErrorFlockerDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.FC, *volume2.FC) {
		log.Warnf("%s", ErrorFCDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.AzureFile, *volume2.AzureFile) {
		log.Warnf("%s", ErrorAzureFileDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.VsphereVolume, *volume2.VsphereVolume) {
		log.Warnf("%s", ErrorVsphereVolumeDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Quobyte, *volume2.Quobyte) {
		log.Warnf("%s", ErrorQuobyteDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.AzureDisk, *volume2.AzureDisk) {
		log.Warnf("%s", ErrorAzureDiskDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.PhotonPersistentDisk, *volume2.PhotonPersistentDisk) {
		log.Warnf("%s", ErrorPhotonPersistentDiskDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Projected, *volume2.Projected) {
		log.Warnf("%s", ErrorProjectedDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.PortworxVolume, *volume2.PortworxVolume) {
		log.Warnf("%s", ErrorPortworxVolumeDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.ScaleIO, *volume2.ScaleIO) {
		log.Warnf("%s", ErrorScaleIODifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.StorageOS, *volume2.StorageOS) {
		log.Warnf("%s", ErrorStorageOSDifferent.Error())
	}

	if !reflect.DeepEqual(*volume1.Ephemeral, *volume2.Ephemeral) {
		log.Warnf("%s", ErrorEphemeralDifferent.Error())
	}

	return nil, nil
}
