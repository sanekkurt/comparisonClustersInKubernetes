package containers

import (
	"context"
	"errors"
	"reflect"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/logging"
)

var (
	ErrorContainerDifferentEnvVarsNumber = errors.New("different number of environment variables in container specs")
	ErrorContainerDifferentEnvVarNames   = errors.New("different environment variable names in container specs")

	ErrorContainerDifferentEnvVarValues       = errors.New("different values of environment variable in container specs")
	ErrorContainerDifferentEnvVarValueSources = errors.New("different environment variable value sources in container specs")

	ErrorContainerEnvValueFromComparisonNotImplemented = errors.New("environment variable ValueFrom type not implemented yet")
)

func getConfigMapKeyRefValue(ctx context.Context, clientSet kubernetes.Interface, namespace string, configMapName, configMapKeyRef string) (string, error) {
	log := logging.FromContext(ctx)

	configMap, err := configmap.GetConfigMapByName(ctx, clientSet, namespace, configMapName)
	if err != nil {
		return "", err
	}

	v, ok := configMap.Data[configMapKeyRef]
	if !ok {
		log.Warnf("%s does not exist in configmap/%s in %s", configMapKeyRef, configMapName, namespace)
		return "", kv_maps.ErrorKVMapNoSuchKey
	}

	return v, nil
}

func getSecretKeyRefValue(ctx context.Context, clientSet kubernetes.Interface, namespace string, secretName, secretKeyRef string) ([]byte, error) {
	log := logging.FromContext(ctx)

	secret, err := secret.GetSecretByName(ctx, clientSet, namespace, secretName)
	if err != nil {
		return nil, err
	}

	v, ok := secret.Data[secretKeyRef]
	if !ok {
		log.Warnf("%s does not exist in secret/%s in %s", secretKeyRef, secretName, namespace)
		return nil, kv_maps.ErrorKVMapNoSuchKey
	}

	return v, nil
}

func getEnvValue(ctx context.Context, clientSet kubernetes.Interface, namespace string, env v12.EnvVar) (interface{}, error) {
	log := logging.FromContext(ctx)

	if env.ValueFrom != nil {
		if env.ValueFrom.ConfigMapKeyRef != nil {
			return getConfigMapKeyRefValue(ctx, clientSet, namespace, env.ValueFrom.ConfigMapKeyRef.Name, env.ValueFrom.ConfigMapKeyRef.Key)
		} else if env.ValueFrom.SecretKeyRef != nil {
			return getSecretKeyRefValue(ctx, clientSet, namespace, env.ValueFrom.SecretKeyRef.Name, env.ValueFrom.SecretKeyRef.Key)
		} else {
			log.Warnf("unknown ValueFrom type: %#v", env.ValueFrom)
			return nil, ErrorContainerEnvValueFromComparisonNotImplemented
		}
	}

	return env.Value, nil
}

func compareEnvVarValueFroms(ctx context.Context, env1, env2 v12.EnvVar) (bool, error) {
	log := logging.FromContext(ctx)

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil ||
		env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {

		log.Warnf("variable %s has different value sources: configMapKeyRef vs secretKeyRef", env1.Name)
		return true, ErrorContainerDifferentEnvVarValueSources
	}

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {
		if env1.ValueFrom.ConfigMapKeyRef.Name != env2.ValueFrom.ConfigMapKeyRef.Name {
			log.Warnf("variable %s has different value source ConfigMaps: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env2.ValueFrom.ConfigMapKeyRef.Name)
			return true, ErrorContainerDifferentEnvVarValueSources
		}

		if env1.ValueFrom.ConfigMapKeyRef.Key != env2.ValueFrom.ConfigMapKeyRef.Key {
			log.Warnf("variable %s has different value source ConfigMap %s keys: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env1.ValueFrom.ConfigMapKeyRef.Key, env2.ValueFrom.ConfigMapKeyRef.Key)
			return true, ErrorContainerDifferentEnvVarValueSources
		}
	}

	if env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil {
		if env1.ValueFrom.SecretKeyRef.Name != env2.ValueFrom.SecretKeyRef.Name {
			log.Warnf("variable %s has different value source Secrets: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env2.ValueFrom.SecretKeyRef.Name)
			return true, ErrorContainerDifferentEnvVarValueSources
		}

		if env1.ValueFrom.SecretKeyRef.Key != env2.ValueFrom.SecretKeyRef.Key {
			log.Warnf("variable %s has different value source Secret %s keys: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env1.ValueFrom.SecretKeyRef.Key, env2.ValueFrom.SecretKeyRef.Key)
			return true, ErrorContainerDifferentEnvVarValueSources
		}
	}

	if env1.ValueFrom.FieldRef != nil && env2.ValueFrom.FieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.FieldRef, *env2.ValueFrom.FieldRef); bDiff {
			log.Warnf("variable %s has different fieldRef value sources", env1.Name)
			return true, ErrorContainerDifferentEnvVarValueSources
		}
	}

	if env1.ValueFrom.ResourceFieldRef != nil && env2.ValueFrom.ResourceFieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.ResourceFieldRef, *env2.ValueFrom.ResourceFieldRef); bDiff {
			log.Warnf("variable %s has different resourceFieldRef value sources", env1.Name)
			return true, ErrorContainerDifferentEnvVarValueSources
		}
	}

	return false, nil
}

func compareEnvVarValueSources(ctx context.Context, env1, env2 v12.EnvVar) (bool, error) {
	log := logging.FromContext(ctx)

	if env1.ValueFrom == nil && env2.ValueFrom != nil || env1.ValueFrom != nil && env2.ValueFrom == nil {
		log.Warnf("variable %s has different value sources: raw value vs ValueFrom", env1.Name)
		return true, ErrorContainerDifferentEnvVarValueSources
	}

	if env1.ValueFrom != nil && env2.ValueFrom != nil {
		bDiff, err := compareEnvVarValueFroms(ctx, env1, env2)

		if err != nil || bDiff {
			return bDiff, err
		}
	}

	if env1.Value != env2.Value {
		log.Warnf("variable %s has different values: '%s' vs '%s'", env1.Name, env1.Value, env2.Value)
		return true, ErrorContainerDifferentEnvVarValues
	}

	//	if env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef != nil && env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef != nil {
	//
	//		if env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key != env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key || env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name != env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name {
	//			return fmt.Errorf("%w. Different ValueFrom: ValueFrom ConfigMapKeyRef in container 1 - %s:%s. ValueFrom ConfigMapKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key)
	//		}
	//
	//		// logic check on configMapKey
	//		log.Debugf("compare environments in container %s: get configMap1", containerName)
	//		configMap1, err := clientSet1.CoreV1().ConfigMaps(namespace).Get(env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//
	//		log.Debugf("compare environments in container %s: get configMap2", containerName)
	//		configMap2, err := clientSet2.CoreV1().ConfigMaps(namespace).Get(env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Name, metav1.GetOptions{})
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//
	//		log.Debugf("compare environments in container %s: check env in config map %s", containerName, configMap1.Name)
	//		if configMap1.Data[env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key] != configMap2.Data[env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key] {
	//			return fmt.Errorf("%w. Environment in container 1: ConfigMapKeyRef.Key = %s, value = %s. Environment in container 2: ConfigMapKeyRef.Key = %s, value = %s", ErrorDifferentValueConfigMapKey, env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, configMap1.Data[env1[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key], env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key, configMap2.Data[env2[pod1EnvIdx].ValueFrom.ConfigMapKeyRef.Key])
	//		}
	//
	//	} else if env1[pod1EnvIdx].ValueFrom.SecretKeyRef != nil && env2[pod1EnvIdx].ValueFrom.SecretKeyRef != nil {
	//
	//		if env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key != env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key || env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name != env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name {
	//			return fmt.Errorf("%w. Different ValueFrom: ValueFrom SecretKeyRef in container 1 - %s:%s. ValueFrom SecretKeyRef in container 2 - %s:%s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key)
	//		}
	//
	//		// logic check on secretKey
	//		log.Debugf("compare environments in container %s: get secrets1", containerName)
	//		secret1, err := clientSet1.CoreV1().Secrets(namespace).Get(env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//
	//		log.Debugf("compare environments in container %s: get secrets2", containerName)
	//		secret2, err := clientSet2.CoreV1().Secrets(namespace).Get(env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//
	//		log.Debugf("compare environments in container %s: check env in secret %s", containerName, secret1.Name)
	//		if string(secret1.Data[env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]) != string(secret2.Data[env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]) {
	//			return fmt.Errorf("%w. Environment in container 1: SecretKeyRef.Key = %s, value = %s. Environment in container 2: SecretKeyRef.Key = %s, value = %s", ErrorDifferentValueSecretKey, env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, string(secret1.Data[env1[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]), env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key, string(secret2.Data[env2[pod1EnvIdx].ValueFrom.SecretKeyRef.Key]))
	//		}
	//
	//	}
	//} else if env1[pod1EnvIdx].ValueFrom != nil || env2[pod1EnvIdx].ValueFrom != nil {
	//return fmt.Errorf("%w. Different ValueFrom: ValueFrom in container 1 - %s. ValueFrom in container 2 - %s", ErrorEnvironmentNotEqual, env1[pod1EnvIdx].ValueFrom, env2[pod1EnvIdx].ValueFrom)
	//}

	return false, nil
}

func compareEnvVars(ctx context.Context, envIdx int, env1, env2 v12.EnvVar) (bool, error) {
	var (
		log = logging.FromContext(ctx)

		//clientSet1, clientSet2, namespace = config.FromContext(ctx)
	)

	if env1.Name != env2.Name {
		log.Warnf("variable #%d has different names: %s vs %s", envIdx, env1.Name, env2.Name)
		return true, ErrorContainerDifferentEnvVarNames
	}

	bDiff, err := compareEnvVarValueSources(ctx, env1, env2)
	if err != nil {
		return false, err
	}
	if bDiff {
		return bDiff, err
	}

	//envValue1, err := getEnvValue(ctx, clientSet1, namespace, env1)
	//if err != nil {
	//	return false, err
	//}
	//
	//envValue2, err := getEnvValue(ctx, clientSet2, namespace, env2)
	//if err != nil {
	//	return false, err
	//}

	return false, nil
}

// compareContainerEnvVars compare environment in containers
func compareContainerEnvVars(ctx context.Context, env1, env2 []v12.EnvVar) (bool, error) {
	var (
		log = logging.FromContext(ctx)
	)

	log.Debugf("Start compare container environments")

	if len(env1) != len(env2) {
		log.Warnf("%s: %d vs %d", ErrorContainerDifferentEnvVarsNumber.Error(), len(env1), len(env2))
		//return fmt.Errorf()
	}

	for pod1EnvIdx := range env1 {
		bDiff, err := compareEnvVars(ctx, pod1EnvIdx, env1[pod1EnvIdx], env2[pod1EnvIdx])
		if err != nil {
			return false, err
		}

		if bDiff {
			return true, err
		}
	}
	return false, nil
}
