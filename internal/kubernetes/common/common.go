package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

var (
	SkippedKubeLabels = map[string]struct{}{
		"app.kubernetes.io/version": {},
	}
)

const (
	matchLabelsStringSep = ","
)

// YamlToStruct parse yaml file into structure
func YamlToStruct(yamlFileName string) *types.KubeconfigYaml {
	kubeconfigYaml := &types.KubeconfigYaml{}

	data, err := ioutil.ReadFile(yamlFileName) //nolint
	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(data, kubeconfigYaml)
	if err != nil {
		panic(err.Error())
	}

	return kubeconfigYaml
}

// GetClientSet reads the configuration from the yaml file using the passed path
func GetClientSet(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

// GetPodsListOnMatchLabels get pods list
func GetPodsListOnMatchLabels(ctx context.Context, matchLabels map[string]string, namespace string, clientSet kubernetes.Interface) (*v12.PodList, error) { //nolint:gocritic,unused
	matchLabelsString := ConvertMatchLabelsToString(ctx, matchLabels)

	pods, err := clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		return nil, err
	}

	return pods, nil
}

// ConvertMatchLabelsToString convert MatchLabels to string
func ConvertMatchLabelsToString(ctx context.Context, matchLabels map[string]string) string {
	var (
		keys   = make([]string, 0, len(matchLabels))
		values = make([]string, 0, len(matchLabels))
	)

	for key, _ := range matchLabels { //nolint
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, k := range keys {
		values = append(values, fmt.Sprintf("%s=%s", k, matchLabels[k]))
	}

	return strings.Join(values, matchLabelsStringSep)
}

type KubeObject struct {
	Type     metav1.TypeMeta
	Metadata metav1.ObjectMeta
	Item     interface{}
}

func CreateKubeObjectsMap(ctx context.Context, items []KubeObject, skipEntities skipper.SkipComponentNames) map[string]types.IsAlreadyComparedFlag {
	var (
		log = logging.FromContext(ctx)

		objMap = make(map[string]types.IsAlreadyComparedFlag)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, itm := range items {
		//if checkContinueTypes(value.Type) {
		//	log.Debugf("secret/%s is skipped from comparison due to its '%s' type", value.Name, value.Type)
		//	continue
		//}

		if skipEntities.IsSkippedEntity(itm.Metadata.Name) {
			log.Debugf("%s/%s is skipped from comparison due to its name", strings.ToLower(itm.Type.Kind), itm.Metadata.Name)
			continue
		}

		indexCheck.Index = index
		objMap[itm.Metadata.Name] = indexCheck
	}

	return objMap
}
