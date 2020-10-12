package common

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

var (
	SkippedKubeLabels = map[string]struct{}{
		"app.kubernetes.io/version": {},
	}
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
func GetPodsListOnMatchLabels(matchLabels map[string]string, namespace string, clientSet1, clientSet2 kubernetes.Interface) (*v12.PodList, *v12.PodList) { //nolint:gocritic,unused
	matchLabelsString := ConvertMatchLabelsToString(matchLabels)

	pods1, err := clientSet1.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}

	pods2, err := clientSet2.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: matchLabelsString})
	if err != nil {
		panic(err.Error())
	}

	return pods1, pods2
}

// ConvertMatchLabelsToString convert MatchLabels to string
func ConvertMatchLabelsToString(matchLabels map[string]string) string {
	var (
		keys   []string
		values []string
	)

	for key, _ := range matchLabels { //nolint
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for i := 0; i < len(keys); i++ {
		values = append(values, fmt.Sprintf("%s=%s", keys[i], matchLabels[keys[i]]))
	}
	return strings.Join(values, ",")
}
