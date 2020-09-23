package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	kube "k8s-cluster-comparator/internal/kubernetes"
	"k8s-cluster-comparator/internal/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

// Opts structure describing input information about clusters and namespaces for comparison
var Opts struct {
	KubeConfig1 string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
	KubeConfig2 string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
	NameSpaces  []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
	Skip        string   `long:"skip" env:"SKIP" required:"false" description:"Skipping an entity"`
}

func main() {
	if err := logging.SetupLogging(); err != nil {
		fmt.Println("[ERROR] ", err.Error())
		os.Exit(1)
	}

	logging.Log.Infow("Starting k8s-cluster-comparator")

	_, err := flags.Parse(&Opts)
	if err != nil {
		panic(err.Error())
	}

	if strings.Contains(Opts.NameSpaces[0], ",") {
		kube.VariableForNamespaces = strings.Split(Opts.NameSpaces[0], ",")
	}

	if kube.VariableForNamespaces == nil {
		kube.VariableForNamespaces = Opts.NameSpaces
	}

	kubeconfig1 := Opts.KubeConfig1
	kubeconfig2 := Opts.KubeConfig2

	if Opts.Skip != "" {
		err = GetMissEntities()
		if err != nil {
			logging.Log.Error(err)
		}
	}

	kube.Client1 = GetClientSet(kubeconfig1)
	kube.Client2 = GetClientSet(kubeconfig2)

	// parse yaml files in global environment
	YamlToStruct(kubeconfig1, &kube.Kubeconfig1YamlStruct)
	YamlToStruct(kubeconfig2, &kube.Kubeconfig2YamlStruct)

	ret := 0

	isClusterDiffer, err := kube.CompareClusters(kube.Client1, kube.Client2, kube.VariableForNamespaces)
	if err != nil {
		logging.Log.Errorf("cannot compare clusters: %s", err.Error())
		os.Exit(2)
	}

	if isClusterDiffer {
		ret = 1
	}

	logging.Log.Infow("k8s-cluster-comparator completed")

	os.Exit(ret)
}

// YamlToStruct parse yaml file into structure
func YamlToStruct(nameYamlFile string, nameStruct *kube.KubeconfigYaml) {
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
	if strings.Contains(Opts.Skip, ";") {
		temp := strings.Split(Opts.Skip, ";")
		var tempSlice []string
		kube.ToSkipEntities = make(map[string]kube.ToSkipComponentNames)
		tempMap := kube.ToSkipComponentNames{}
		for _, value := range temp {
			if strings.Contains(value, ":") {
				tempSlice = strings.Split(value, ":")
				if strings.Contains(tempSlice[1], ",") {
					for _, val := range strings.Split(tempSlice[1], ",") {
						tempMap[val] = struct{}{}
					}
					kube.ToSkipEntities[tempSlice[0]] = make(map[string]struct{})
					for key, value := range tempMap {
						kube.ToSkipEntities[tempSlice[0]][key] = value
						delete(tempMap, key)
					}
				} else {
					kube.ToSkipEntities[tempSlice[0]] = make(map[string]struct{})
					kube.ToSkipEntities[tempSlice[0]][tempSlice[1]] = struct{}{}
				}
			} else {
				return errors.New("does not contain valid data in the 'skip' variable. The enumeration of the names of entities start after ':' please or don't finish the line ';'")
			}
		}
		return nil
	}
	return errors.New("does not contain valid data in the 'skip' variable. Between entities put ';' please")

}
