package kv_maps

import (
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"

	"fmt"
	"sync"
)

var (
	skipSecretTypes = [3]v12.SecretType{
		"kubernetes.io/service-account-token",
		"kubernetes.io/dockercfg",
		"helm.sh/release.v1",
	}
)

// CompareSecrets compares list of secret objects in two given k8s-clusters
func CompareSecrets(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)
	secrets1, err := clientSet1.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain secrets list from 1st cluster: %w", err)
	}

	secrets2, err := clientSet2.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain secrets list from 2nd cluster: %w", err)
	}

	mapSecrets1, mapSecrets2 := prepareSecretMaps(secrets1, secrets2, skipEntityList.GetByKind("secrets"))

	isClustersDiffer = compareSecretsSpecs(mapSecrets1, mapSecrets2, secrets1, secrets2)

	return isClustersDiffer, nil
}

// prepareSecretMaps add value secrets in map
func prepareSecretMaps(secrets1, secrets2 *v12.SecretList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapSecrets1 := make(map[string]types.IsAlreadyComparedFlag)
	mapSecrets2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range secrets1.Items {
		if checkContinueTypes(value.Type) {
			log.Debugf("secret %s is skipped from comparison due to its '%s' type", value.Name, value.Type)
			continue
		}
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("secret %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapSecrets1[value.Name] = indexCheck

	}
	for index, value := range secrets2.Items {
		if checkContinueTypes(value.Type) {
			log.Debugf("secret %s is skipped from comparison due to its '%s' type", value.Name, value.Type)
			continue
		}
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("secret %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapSecrets2[value.Name] = indexCheck

	}
	return mapSecrets1, mapSecrets2
}

func compareSecretSpecInternals(wg *sync.WaitGroup, channel chan bool, name string, secret1, secret2 *v12.Secret) {
	var (
		flag bool
	)

	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking secret: '%s' -----", name)

	if len(secret1.Data) != len(secret2.Data) {
		log.Infof("secret '%s' in 1st cluster has '%d' keys but the 2nd - '%d'", name, len(secret1.Data), len(secret2.Data))
		flag = true
	} else {
		for key, value := range secret1.Data {
			v1 := string(value)
			v2 := string(secret2.Data[key])

			if v1 != v2 {
				log.Infof("secret '%s', values by key '%s' do not match: '%s' and %s", name, key, v1, v2)
				flag = true
			}
		}
	}

	log.Debugf("----- End checking secret: '%s' -----", name)

	channel <- flag
}

// compareSecretsSpecs set information about secrets
func compareSecretsSpecs(map1, map2 map[string]types.IsAlreadyComparedFlag, secrets1, secrets2 *v12.SecretList) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		log.Infof("secret counts are different")
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareSecretSpecInternals(wg, channel, name, &secrets1.Items[index1.Index], &secrets2.Items[index2.Index])
		} else {
			log.Infof("secret '%s' does not exist in 2nd cluster", name)
			flag = true
		}
	}

	wg.Wait()

	close(channel)

	for ch := range channel {
		if ch {
			flag = true
		}
	}
	for name, index := range map2 {
		if !index.Check {

			log.Infof("secret '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}
	return flag
}

func checkContinueTypes(secretType v12.SecretType) bool {
	var skip bool
	for _, skipType := range skipSecretTypes {
		if secretType == skipType {
			skip = true
		}
	}
	return skip
}
