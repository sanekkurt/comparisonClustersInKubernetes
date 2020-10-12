package config

import (
	"context"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	NamespacesListSep = ","
)

var (
	// opts structure describing input information about clusters and namespaces for comparison
	opts struct {
		KubeConfig1 string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
		KubeConfig2 string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
		NameSpaces  []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
		Skip        string   `long:"skip" env:"SKIP" required:"false" description:"Skipping an entity"`
	}

	ErrHelpShown = errors.New("help message shown")
)

// ClusterConfig represents a k8s-cluster config
type ClusterConfig struct {
	Kubeconfig   kubernetes.Interface
	ConfigStruct *types.KubeconfigYaml
}

// AppConfig is the main application configuration storage
type AppConfig struct {
	Cluster1 ClusterConfig
	Cluster2 ClusterConfig

	Namespaces []string

	SkipEntitiesList skipper.SkipEntitiesList
}

// Parse performs configuration parsing from various sources and fills in the AppConfig struct
func Parse(ctx context.Context) (*AppConfig, error) {
	log := logging.FromContext(ctx)

	_, err := flags.Parse(&opts)
	if err != nil {
		if len(os.Args) > 1 {
			if os.Args[1] == "--help" || os.Args[1] == "-h" {
				return nil, ErrHelpShown
			}
		}
		log.Debugf("configuration invalid: %w", err)
		return nil, err
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
		skipEntities, err := skipper.ParseSkipConfig(opts.Skip)
		if err != nil {
			log.Errorf("cannot parse skip entities list: %s", err.Error())
		}

		appConfig.SkipEntitiesList = skipEntities
	}

	return appConfig, nil
}
