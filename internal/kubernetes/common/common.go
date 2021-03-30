package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s-cluster-comparator/internal/kubernetes/types"
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

//// ConvertMatchLabelsToString convert MatchLabels to string
//func ConvertMatchLabelsToString(ctx context.Context, matchLabels map[string]string) string {
//	var (
//		keys   = make([]string, 0, len(matchLabels))
//		values = make([]string, 0, len(matchLabels))
//	)
//
//	for key, _ := range matchLabels { //nolint
//		keys = append(keys, key)
//	}
//
//	sort.Strings(keys)
//
//	for _, k := range keys {
//		values = append(values, fmt.Sprintf("%s=%s", k, matchLabels[k]))
//	}
//
//	return strings.Join(values, matchLabelsStringSep)
//}
