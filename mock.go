package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	kube "k8s-cluster-comparator/internal/kubernetes"
	"k8s-cluster-comparator/internal/logging"
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

	logging.Log.Infow("Starting k8s-cluster-comparator mocker")

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

	clusters := map[string]*kubernetes.Clientset{
		"cluster_1": GetClientSet(Opts.KubeConfig1),
		"cluster_2": GetClientSet(Opts.KubeConfig2),
	}

	for _, ns := range kube.VariableForNamespaces {
		for cName, clientSet := range clusters {
			logging.Log.Info("Writing deployments...")
			deployments, err := clientSet.AppsV1().Deployments(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "deployments", cName), deployments)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing statefulsets...")
			statefulset, err := clientSet.AppsV1().StatefulSets(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "statefulset", cName), statefulset)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing daemonsets...")
			daemonsets, err := clientSet.AppsV1().DaemonSets(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "daemonsets", cName), daemonsets)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing configmaps...")
			configmaps, err := clientSet.CoreV1().ConfigMaps(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "configmaps", cName), configmaps)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing secrets...")
			secrets, err := clientSet.CoreV1().Secrets(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "secrets", cName), secrets)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing services...")
			services, err := clientSet.CoreV1().Services(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "services", cName), services)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			logging.Log.Info("Writing ingresses...")
			ingresses, err := clientSet.NetworkingV1beta1().Ingresses(ns).List(metav1.ListOptions{})
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}

			err = writeGob(fmt.Sprintf("gobs/%s_%s.gob", "ingresses", cName), ingresses)
			if err != nil {
				logging.Log.Fatalw(err.Error())
			}
		}
	}

	//obj := &v1.DeploymentList{}
	//err = readGob(f, obj)
	//if err != nil {
	//	logging.Log.Fatalw(err.Error())
	//}

	//fmt.Printf("%#v", obj)
}
func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err := encoder.Encode(object)
		if err != nil {
			return fmt.Errorf("cannot serialize struct: %w", err)
		}
	}

	if err = file.Close(); err != nil {
		logging.Log.Infof("cannot close file: %s", err.Error())
	}

	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	if err = file.Close(); err != nil {
		logging.Log.Infof("cannot close file: %s", err.Error())
	}

	return err
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
