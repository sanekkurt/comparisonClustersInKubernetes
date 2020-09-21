package main

import (
	"flag"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

var (
	kubeconfig            = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	variableForNamespaces []string
	kubeconfig1YamlStruct KubeconfigYaml
	kubeconfig2YamlStruct KubeconfigYaml
	client1               *kubernetes.Clientset
	client2               *kubernetes.Clientset

	skipTypes = []v12.SecretType{"kubernetes.io/service-account-token", "kubernetes.io/dockercfg", "helm.sh/release.v1"}
)

var Opts struct {
	KubeConfig1 string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
	KubeConfig2 string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
	NameSpaces  []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
}

func main() {
	if err := SetupLogging(); err != nil {
		fmt.Println("[ERROR] ", err.Error())
		os.Exit(1)
	}

	log.Infow("Starting k8s-cluster-comparator")

	_, err := flags.Parse(&Opts)
	if err != nil {
		panic(err.Error())
	}
	//home, err := os.Getwd()
	//if err != nil {
	//	panic(err.Error())
	//}
	//kubeconfig1 := flag.String("kubeconfig1", filepath.Join(home, "kubeconfig1.yaml"), "(optional) absolute path to the kubeconfig file")
	//kubeconfig2 := flag.String("kubeconfig2", filepath.Join(home, "kubeconfig2.yaml"), "(optional) absolute path to the kubeconfig file")

	if strings.Contains(Opts.NameSpaces[0], ",") {
		variableForNamespaces = strings.Split(Opts.NameSpaces[0], ",")
	}

	//for _, value := range Opts.NameSpaces[0] {
	//	if value == 44 {
	//		variableForNamespaces = strings.Split(Opts.NameSpaces[0], ",")
	//		break
	//	}
	//}
	if variableForNamespaces == nil {
		variableForNamespaces = Opts.NameSpaces
	}

	kubeconfig1 := &Opts.KubeConfig1
	kubeconfig2 := &Opts.KubeConfig2

	client1 = GetClientSet(kubeconfig1)
	client2 = GetClientSet(kubeconfig2)

	//распарсинг yaml файлов в глобальные переменные, чтобы в будущем получить из них URL
	YamlToStruct(*kubeconfig1, &kubeconfig1YamlStruct)
	YamlToStruct(*kubeconfig2, &kubeconfig2YamlStruct)

	ret := 0

	isClusterDiffer, err := Compare(client1, client2 /*"default"*/, variableForNamespaces)
	if err != nil {
		log.Errorf("cannot compare clusters: %s", err.Error())
		os.Exit(2)
	}

	if isClusterDiffer {
		ret = 1
	}

	log.Infow("k8s-cluster-comparator completed")

	os.Exit(ret)
}

//переводит yaml в структуру
func YamlToStruct(nameYamlFile string, nameStruct *KubeconfigYaml) {
	data, err := ioutil.ReadFile(nameYamlFile)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(data, nameStruct)
	if err != nil {
		panic(err.Error())
	}
}

//читает конфигурацию из yaml файла по переданному пути
func GetClientSet(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
