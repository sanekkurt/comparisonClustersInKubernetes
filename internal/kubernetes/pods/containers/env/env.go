package env

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"reflect"

	"k8s-cluster-comparator/internal/kubernetes/kv_maps/configmap"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps/secret"

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

func compareEnvVarValueFroms(ctx context.Context, env1, env2 v12.EnvVar) error {
	var (
		//diffsBatch = diff.BatchFromContext(ctx)
		diffsChannel = diff.ChanFromContext(ctx)
	)

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil ||
		env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {
		//log.Warnf("variable %s has different value sources: configMapKeyRef vs secretKeyRef", env1.Name)
		//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value sources: configMapKeyRef vs secretKeyRef", env1.Name)
		*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value sources: configMapKeyRef vs secretKeyRef", append(make([]interface{}, 0, 0), env1.Name)}
	}

	if env1.ValueFrom.ConfigMapKeyRef != nil && env2.ValueFrom.ConfigMapKeyRef != nil {
		if env1.ValueFrom.ConfigMapKeyRef.Name != env2.ValueFrom.ConfigMapKeyRef.Name {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value source ConfigMaps: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env2.ValueFrom.ConfigMapKeyRef.Name)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value source ConfigMaps: %s vs %s", append(make([]interface{}, 0, 0), env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env2.ValueFrom.ConfigMapKeyRef.Name)}

		}

		if env1.ValueFrom.ConfigMapKeyRef.Key != env2.ValueFrom.ConfigMapKeyRef.Key {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value source ConfigMap %s keys: %s vs %s", env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env1.ValueFrom.ConfigMapKeyRef.Key, env2.ValueFrom.ConfigMapKeyRef.Key)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value source ConfigMap %s keys: %s vs %s", append(make([]interface{}, 0, 0), env1.Name, env1.ValueFrom.ConfigMapKeyRef.Name, env1.ValueFrom.ConfigMapKeyRef.Key, env2.ValueFrom.ConfigMapKeyRef.Key)}

		}
	}

	if env1.ValueFrom.SecretKeyRef != nil && env2.ValueFrom.SecretKeyRef != nil {
		if env1.ValueFrom.SecretKeyRef.Name != env2.ValueFrom.SecretKeyRef.Name {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value source Secrets: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env2.ValueFrom.SecretKeyRef.Name)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value source Secrets: %s vs %s", append(make([]interface{}, 0, 0), env1.Name, env1.ValueFrom.SecretKeyRef.Name, env2.ValueFrom.SecretKeyRef.Name)}

		}

		if env1.ValueFrom.SecretKeyRef.Key != env2.ValueFrom.SecretKeyRef.Key {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value source Secret %s keys: %s vs %s", env1.Name, env1.ValueFrom.SecretKeyRef.Name, env1.ValueFrom.SecretKeyRef.Key, env2.ValueFrom.SecretKeyRef.Key)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value source Secret %s keys: %s vs %s", append(make([]interface{}, 0, 0), env1.Name, env1.ValueFrom.SecretKeyRef.Name, env1.ValueFrom.SecretKeyRef.Key, env2.ValueFrom.SecretKeyRef.Key)}

		}
	}

	if env1.ValueFrom.FieldRef != nil && env2.ValueFrom.FieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.FieldRef, *env2.ValueFrom.FieldRef); bDiff {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different fieldRef value sources", env1.Name)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different fieldRef value sources", append(make([]interface{}, 0, 0), env1.Name)}

		}
	}

	if env1.ValueFrom.ResourceFieldRef != nil && env2.ValueFrom.ResourceFieldRef != nil {
		if bDiff := reflect.DeepEqual(*env1.ValueFrom.ResourceFieldRef, *env2.ValueFrom.ResourceFieldRef); bDiff {
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different resourceFieldRef value sources", env1.Name)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different resourceFieldRef value sources", append(make([]interface{}, 0, 0), env1.Name)}

		}
	}

	return nil
}

func compareEnvVarValueSources(ctx context.Context, env1, env2 v12.EnvVar) error {
	var (
		//diffsBatch = diff.BatchFromContext(ctx)
		diffsChannel = diff.ChanFromContext(ctx)
	)

	if env1.ValueFrom == nil && env2.ValueFrom != nil || env1.ValueFrom != nil && env2.ValueFrom == nil {
		//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different value sources: raw value vs ValueFrom", env1.Name)
		*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different value sources: raw value vs ValueFrom", append(make([]interface{}, 0, 0), env1.Name)}

	}

	if env1.ValueFrom != nil && env2.ValueFrom != nil {
		err := compareEnvVarValueFroms(ctx, env1, env2)
		if err != nil {
			return err
		}

	}

	if env1.Value != env2.Value {
		//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable %s has different values: '%s' vs '%s'", env1.Name, env1.Value, env2.Value)
		*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "variable %s has different values: '%s' vs '%s'", append(make([]interface{}, 0, 0), env1.Name, env1.Value, env2.Value)}

	}

	return nil
}

func compareEnvVars(ctx context.Context, env1, env2 v12.EnvVar) error {

	//if env1.Name != env2.Name {
	//	diffsBatch.Add(ctx, false, zapcore.WarnLevel, "variable #%d: %s: %s vs %s", envIdx+1, ErrorContainerDifferentEnvVarNames.Error(), env1.Name, env2.Name)
	//	//log.Warnf("variable #%d: %s: %s vs %s", envIdx+1, ErrorContainerDifferentEnvVarNames.Error(), env1.Name, env2.Name)
	//}

	err := compareEnvVarValueSources(ctx, env1, env2)
	if err != nil {
		return err
	}

	return nil
}

// Compare compare environment variables in container specs
func Compare(ctx context.Context, envs1, envs2 []v12.EnvVar) error {
	var (
		log = logging.FromContext(ctx)

		//diffsBatch = diff.BatchFromContext(ctx)
		diffsChannel = diff.ChanFromContext(ctx)
	)

	log.Debugf("CompareEnvVars: started")
	defer log.Debugf("CompareEnvVars: completed")

	if len(envs1) != len(envs2) {
		//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "%s: %d vs %d", ErrorContainerDifferentEnvVarsNumber.Error(), len(envs1), len(envs2))
		*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "%s: %d vs %d", append(make([]interface{}, 0, 0), ErrorContainerDifferentEnvVarsNumber.Error(), len(envs1), len(envs2))}
	}

	//for pod1EnvIdx := range env1 {
	//	if pod1EnvIdx > len(env2)-1 {
	//		//log.Warnf("CompareEnvVars: there are only %d envVars in 2nd cluster", len(env2))
	//		diffsBatch.Add(ctx, false, zapcore.WarnLevel, "CompareEnvVars: there are only %d envVars in 2nd cluster", len(env2))
	//		break
	//	}
	//	err := compareEnvVars(ctx, pod1EnvIdx, env1[pod1EnvIdx], env2[pod1EnvIdx])
	//	if err != nil {
	//		return err
	//	}
	//
	//}
	//
	//if len(envs2) > len(envs1) {
	//	for idx := 1 + (len(envs2) - len(envs1)); idx < len(envs2); idx++ {
	//		log.Warnf("env variable #%d '%s' does not exist in 1st cluster", idx+1, envs2[idx].Name)
	//	}
	//}

	mapEnv1 := makeEnvMap(envs1)
	mapEnv2 := makeEnvMap(envs2)

	//if _, ok := mapEnv1["AB_TEST_REDIS_TIMEOUT"]; ok {
	//	fmt.Println("yes")
	//}

	for key, value := range mapEnv1 {
		if _, ok := mapEnv2[key]; ok {
			err := compareEnvVars(ctx, value, mapEnv2[key])
			if err != nil {
				return err
			}

			delete(mapEnv1, key)
			delete(mapEnv2, key)
		}
	}

	if len(mapEnv1) > 0 {
		for key, _ := range mapEnv1 {
			//log.Warnf("env variable '%s' does not exist in 2st cluster", key)
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "env variable '%s' does not exist in 2st cluster", key)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "env variable '%s' does not exist in 2st cluster", append(make([]interface{}, 0, 0), key)}

		}
	}

	if len(mapEnv2) > 0 {
		for key, _ := range mapEnv2 {
			//log.Warnf("env variable '%s' does not exist in 1st cluster", key)
			//diffsBatch.Add(ctx, false, zapcore.WarnLevel, "env variable '%s' does not exist in 1st cluster", key)
			*diffsChannel <- diff.Diff{ctx, false, zap.WarnLevel, "env variable '%s' does not exist in 1st cluster", append(make([]interface{}, 0, 0), key)}
		}
	}

	return nil
}

func makeEnvMap(envs []v12.EnvVar) map[string]v12.EnvVar {
	newEnvMap := make(map[string]v12.EnvVar, len(envs))

	for _, value := range envs {
		newEnvMap[value.Name] = value
	}

	return newEnvMap
}
