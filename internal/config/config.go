package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

// SkipComponentNames is a map of blank structs of k8s objects that must be skipped from the comparison
type SkipComponentNames map[types.ObjectName]struct{}

type SkipEntitiesList map[types.ObjectKind]SkipComponentNames

const (
	NamespacesListSep = ","
)

var (
	// Opts structure describing input information about clusters and namespaces for comparison
	opts struct {
		KubeConfig1 string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
		KubeConfig2 string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
		NameSpaces  []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
		Skip        string   `long:"skip" env:"SKIP" required:"false" description:"Skipping an entity"`
	}
)

type ClusterConfig struct {
	Kubeconfig   kubernetes.Interface
	ConfigStruct *types.KubeconfigYaml
}

type AppConfig struct {
	Cluster1 ClusterConfig
	Cluster2 ClusterConfig

	Namespaces []string

	SkipEntitiesList SkipEntitiesList
}

func Parse(ctx context.Context) (*AppConfig, error) {
	log := logging.FromContext(ctx)

	_, err := flags.Parse(&opts)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	appConfig := &AppConfig{
		Cluster1: ClusterConfig{
			Kubeconfig:   common.GetClientSet(opts.KubeConfig1),
			ConfigStruct: common.YamlToStruct(opts.KubeConfig1),
		},
		Cluster2: ClusterConfig{
			Kubeconfig:   common.GetClientSet(opts.KubeConfig2),
			ConfigStruct: common.YamlToStruct(opts.KubeConfig2),
		},
	}

	if strings.Contains(opts.NameSpaces[0], ",") {
		appConfig.Namespaces = strings.Split(opts.NameSpaces[0], NamespacesListSep)
	} else {
		appConfig.Namespaces = opts.NameSpaces
	}

	if opts.Skip != "" {
		skipEntities, err := parseSkipEntities()
		if err != nil {
			log.Errorf("cannot parse skip entities list: %s", err.Error())
		}

		appConfig.SkipEntitiesList = skipEntities
	}

	return appConfig, nil
}

// parseSkipEntities parses information about entities to skip from a environment
func parseSkipEntities() (SkipEntitiesList, error) {
	if !strings.Contains(opts.Skip, ";") {
		return nil, errors.New("does not contain valid data in the 'skip' variable. Between entities put ';' please")
	}

	var (
		tempSlice []string

		temp         = strings.Split(opts.Skip, ";")
		skipEntities = make(map[types.ObjectKind]SkipComponentNames)

		tempMap = SkipComponentNames{}
	)

	for _, value := range temp {
		if !strings.Contains(value, ":") {
			return nil, errors.New("does not contain valid data in the 'skip' variable. The enumeration of the names of entities start after ':' please or don't finish the line ';'")
		}

		tempSlice = strings.Split(value, ":")

		if strings.Contains(tempSlice[1], ",") {
			for _, val := range strings.Split(tempSlice[1], ",") {
				tempMap[types.ObjectName(val)] = struct{}{}
			}

			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})

			for key, value := range tempMap {
				skipEntities[types.ObjectKind(tempSlice[0])][key] = value
				delete(tempMap, key)
			}
		} else {
			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})
			skipEntities[types.ObjectKind(tempSlice[0])][types.ObjectName(tempSlice[0])] = struct{}{}
		}
	}

	return skipEntities, nil
}
