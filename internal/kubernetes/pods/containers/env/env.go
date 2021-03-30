package env

import (
	"context"
	"reflect"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"
	"k8s-cluster-comparator/internal/kubernetes/types"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/logging"
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

func compareEnvVarValueFroms(ctx context.Context, env1, env2 v12.EnvVar) ([]types.ObjectsDiff, error) {
	log := logging.FromContext(ctx)

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil ||
		env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {

		log.Warnf("variable %s has different value sources: configMapKeyRef vs secretKeyRef", env1.Name)
	}

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {
		if env1.ValueFrom.ConfigMapKeyRef.Name != env2.ValueFrom.ConfigMapKeyRef.Name {
			log.Warnf("variable %s has different value source ConfigMaps: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env2.ValueFrom.ConfigMapKeyRef.Name)
		}

		if env1.ValueFrom.ConfigMapKeyRef.Key != env2.ValueFrom.ConfigMapKeyRef.Key {
			log.Warnf("variable %s has different value source ConfigMap %s keys: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env1.ValueFrom.ConfigMapKeyRef.Key, env2.ValueFrom.ConfigMapKeyRef.Key)
		}
	}

	if env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil {
		if env1.ValueFrom.SecretKeyRef.Name != env2.ValueFrom.SecretKeyRef.Name {
			log.Warnf("variable %s has different value source Secrets: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env2.ValueFrom.SecretKeyRef.Name)
		}

		if env1.ValueFrom.SecretKeyRef.Key != env2.ValueFrom.SecretKeyRef.Key {
			log.Warnf("variable %s has different value source Secret %s keys: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env1.ValueFrom.SecretKeyRef.Key, env2.ValueFrom.SecretKeyRef.Key)
		}
	}

	if env1.ValueFrom.FieldRef != nil && env2.ValueFrom.FieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.FieldRef, *env2.ValueFrom.FieldRef); bDiff {
			log.Warnf("variable %s has different fieldRef value sources", env1.Name)
		}
	}

	if env1.ValueFrom.ResourceFieldRef != nil && env2.ValueFrom.ResourceFieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.ResourceFieldRef, *env2.ValueFrom.ResourceFieldRef); bDiff {
			log.Warnf("variable %s has different resourceFieldRef value sources", env1.Name)
		}
	}

	return nil, nil
}

func compareEnvVarValueSources(ctx context.Context, env1, env2 v12.EnvVar) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)
		diffs = make([]types.ObjectsDiff, 0)
	)

	if env1.ValueFrom == nil && env2.ValueFrom != nil || env1.ValueFrom != nil && env2.ValueFrom == nil {
		log.Warnf("variable %s has different value sources: raw value vs ValueFrom", env1.Name)
	}

	if env1.ValueFrom != nil && env2.ValueFrom != nil {
		diff, err := compareEnvVarValueFroms(ctx, env1, env2)
		if err != nil  {
			return nil, err
		}

		diffs = append(diffs, diff...)
	}

	if env1.Value != env2.Value {
		log.Warnf("variable %s has different values: '%s' vs '%s'", env1.Name, env1.Value, env2.Value)
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

	return diffs, nil
}

func compareEnvVars(ctx context.Context, envIdx int, env1, env2 v12.EnvVar) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)

		//clientSet1, clientSet2, namespace = config.FromContext(ctx)
	)

	if env1.Name != env2.Name {
		log.Warnf("variable #%d: %s: %s vs %s", envIdx+1, ErrorContainerDifferentEnvVarNames.Error(), env1.Name, env2.Name)
	}

	diff, err := compareEnvVarValueSources(ctx, env1, env2)
	if err != nil {
		return nil, err
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

	return diff, nil
}

// Compare compare environment variables in container specs
func Compare(ctx context.Context, env1, env2 []v12.EnvVar) ([]types.ObjectsDiff, error) {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.ObjectsDiff, 0)
	)

	log.Debugf("CompareEnvVars: started")
	defer log.Debugf("CompareEnvVars: completed")

	if len(env1) != len(env2) {
		log.Warnf("%s: %d vs %d", ErrorContainerDifferentEnvVarsNumber.Error(), len(env1), len(env2))
	}

	for pod1EnvIdx := range env1 {
		if pod1EnvIdx > len(env2) - 1 {
			log.Warnf("CompareEnvVars: there are only %d envVars in 2nd cluster", len(env2))
			break
		}

		err := compareEnvVars(ctx, pod1EnvIdx, env1[pod1EnvIdx], env2[pod1EnvIdx])
		if err != nil {
			return nil, err
		}

		if batch.IsFinal() {
			break
		}

		diffs = append(diffs, diff...)
	}

	if len(env2) > len(env1) {
		for idx := 1 + (len(env2) - len(env1)); idx < len(env2); idx++ {
			log.Warnf("env variable #%d '%s' does not exist in 1st cluster", idx + 1, env2[idx].Name)
		}
	}

	return diffs, nil
}
