package main

import (
	"errors"
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

type nameComponents map[string]struct{}

var (
	variableForNamespaces []string
	kubeconfig1YamlStruct KubeconfigYaml
	kubeconfig2YamlStruct KubeconfigYaml
	client1               *kubernetes.Clientset
	client2               *kubernetes.Clientset

	entities              map[string]nameComponents
	skipTypes             = [3]v12.SecretType{"kubernetes.io/service-account-token", "kubernetes.io/dockercfg", "helm.sh/release.v1"}
)

// Opts structure describing input information about clusters and namespaces for comparison
var Opts struct {
	KubeConfig1 string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
	KubeConfig2 string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
	NameSpaces  []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
	Miss        string   `long:"miss" env:"MISS" required:"false" description:"Skipping an entity"`
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

	if strings.Contains(Opts.NameSpaces[0], ",") {
		variableForNamespaces = strings.Split(Opts.NameSpaces[0], ",")
	}

	if variableForNamespaces == nil {
		variableForNamespaces = Opts.NameSpaces
	}

	kubeconfig1 := Opts.KubeConfig1
	kubeconfig2 := Opts.KubeConfig2

	if Opts.Miss != "" {
		err = GetMissEntities()
		if err != nil {
			log.Error(err)
		}
	}

	client1 = GetClientSet(kubeconfig1)
	client2 = GetClientSet(kubeconfig2)

	// parse yaml files in global environment
	YamlToStruct(kubeconfig1, &kubeconfig1YamlStruct)
	YamlToStruct(kubeconfig2, &kubeconfig2YamlStruct)

	ret := 0

	isClusterDiffer, err := CompareClusters(client1, client2, variableForNamespaces)
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

// YamlToStruct parse yaml file into structure
func YamlToStruct(nameYamlFile string, nameStruct *KubeconfigYaml) {
	data, err := ioutil.ReadFile(nameYamlFile) //nolint
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(data, nameStruct)
	if err != nil {
		panic(err.Error())
	}
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

// GetMissEntities gets information about entities to skip from a environment
func GetMissEntities() error {
	if strings.Contains(Opts.Miss, ";") {
		temp:= strings.Split(Opts.Miss, ";")
		var tempSlice []string
		entities = make(map[string]nameComponents)
		tempMap := nameComponents{}
		for _, value := range temp {
			if strings.Contains(value, ":") {
				tempSlice = strings.Split(value, ":")
				if strings.Contains(tempSlice[1], ",") {
					for _, val := range strings.Split(tempSlice[1], ","){
						tempMap[val] = struct{}{}
					}
					entities[tempSlice[0]] = make(map[string]struct{})
					for key, value := range tempMap {
						entities[tempSlice[0]][key] = value
						delete(tempMap, key)
					}
				} else {
					return errors.New("does not contain valid data in the 'miss' variable. The enumeration of the names of entities do through ',' please")
				}
			} else {
				return errors.New("does not contain valid data in the 'miss' variable. The enumeration of the names of entities start after ':' please or don't finish the line ';'")
			}
		}
		return nil
	}
	return errors.New("does not contain valid data in the 'miss' variable. Between entities put ';' please")


}