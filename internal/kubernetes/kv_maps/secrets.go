package kv_maps

import (
	"context"

	"go.uber.org/zap"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

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

const (
	secretObjectBatchLimit = 25
)

func GetSecretMapByName(ctx context.Context, clientSet kubernetes.Interface, namespace, configMapName string) (*v12.Secret, error) {
	secret, err := clientSet.CoreV1().Secrets(namespace).Get(configMapName, metav1.GetOptions{})
	return secret, err
}

func addItemsToSecretList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v12.SecretList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToSecretList started")
	defer log.Debugf("addItemsToSecretList completed")

	var (
		batch   *v12.SecretList
		secrets = &v12.SecretList{
			Items: make([]v12.Secret, 0),
		}

		continueToken string

		err error
	)

forLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.CoreV1().Secrets(namespace).List(metav1.ListOptions{
				Limit:    limit,
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("addItemsToSecretList: %d objects received", len(batch.Items))

			secrets.Items = append(secrets.Items, batch.Items...)

			secrets.TypeMeta = batch.TypeMeta
			secrets.ListMeta = batch.ListMeta

			if batch.Continue == "" {
				break forLoop
			}

			continueToken = batch.Continue
		}
	}

	secrets.Continue = ""

	return secrets, err
}

// CompareSecrets compares list of secret objects in two given k8s-clusters
func CompareSecrets(ctx context.Context, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", "secret"))

		clientSet1, clientSet2, namespace = config.FromContext(ctx)

		isClustersDiffer bool
	)
	ctx = logging.WithLogger(ctx, log)

	secrets1, err := addItemsToSecretList(ctx, clientSet1, namespace, secretObjectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain secrets list from 1st cluster: %w", err)
	}

	secrets2, err := addItemsToSecretList(ctx, clientSet2, namespace, secretObjectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain secrets list from 2st cluster: %w", err)
	}

	mapSecrets1, mapSecrets2 := prepareSecretMaps(ctx, secrets1, secrets2, skipEntityList.GetByKind("secrets"))

	isClustersDiffer = compareSecretsSpecs(ctx, mapSecrets1, mapSecrets2, secrets1, secrets2)

	return isClustersDiffer, nil
}

// prepareSecretMaps add value secrets in map
func prepareSecretMaps(ctx context.Context, secrets1, secrets2 *v12.SecretList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	log := logging.FromContext(ctx)

	mapSecrets1 := make(map[string]types.IsAlreadyComparedFlag)
	mapSecrets2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range secrets1.Items {
		if checkContinueTypes(value.Type) {
			log.Debugf("secret/%s is skipped from comparison due to its '%s' type", value.Name, value.Type)
			continue
		}

		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("secret/%s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		mapSecrets1[value.Name] = indexCheck

	}
	for index, value := range secrets2.Items {
		if checkContinueTypes(value.Type) {
			log.Debugf("secret/%s is skipped from comparison due to its '%s' type", value.Name, value.Type)
			continue
		}

		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("secret/%s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		mapSecrets2[value.Name] = indexCheck

	}

	return mapSecrets1, mapSecrets2
}

func compareSecretSpecInternals(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name string, secret1, secret2 *v12.Secret) {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))

		flag bool
	)
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking secret/%s -----", name)

	if !metadata.IsMetadataDiffers(ctx, secret1.ObjectMeta, secret1.ObjectMeta) {
		channel <- true
		return
	}

	if !common.AreKVBytesMapsEqual(ctx, secret1.Data, secret2.Data, nil) {
		flag = true
	}
	log.Debugf("----- End checking secret/%s -----", name)

	channel <- flag
}

// compareSecretsSpecs set information about secrets
func compareSecretsSpecs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, secrets1, secrets2 *v12.SecretList) bool {
	var (
		log = logging.FromContext(ctx)

		flag bool
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return true
		default:
			if index2, ok := map2[name]; ok {
				wg.Add(1)

				index1.Check = true
				map1[name] = index1
				index2.Check = true
				map2[name] = index2

				compareSecretSpecInternals(ctx, wg, channel, name, &secrets1.Items[index1.Index], &secrets2.Items[index2.Index])
			} else {
				log.Infof("secret/%s does not exist in 2nd cluster", name)
				flag = true
			}
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
			log.Warnf("secret/%s does not exist in 1st cluster", name)
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
