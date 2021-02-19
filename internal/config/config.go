package config

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/logging"
)

var (
	//// opts structure describing input information about clusters and namespaces for comparison
	//opts struct {
	//	KubeConfig1              string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
	//	KubeConfig2              string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
	//	NameSpaces               []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
	//	SkippedFullResourceNames string   `long:"skip" env:"SKIP" required:"false" description:"Skipping a resource by kind/name specification"`
	//	SkipAnyKindResourceNames string   `long:"skip-names" env:"SKIP_NAMES" required:"false" description:"Resource skipping by name"`
	//}

	ErrHelpShown = errors.New("help message shown")
)

// Parse performs configuration parsing from various sources and fills in the AppConfig struct
func Parse(ctx context.Context) (*AppConfig, error) {
	log := logging.FromContext(ctx)

	if len(os.Args) > 1 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			return nil, ErrHelpShown
		}
	}

	cfg, err := ParseConfig(ctx, "./config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	cfg.Connections.Cluster1.ClientSet = common.GetClientSet(cfg.Connections.Cluster1.KubeConfigFile)
	cfg.Connections.Cluster1.ConfigStruct = common.YamlToStruct(cfg.Connections.Cluster1.KubeConfigFile)

	cfg.Connections.Cluster2.ClientSet = common.GetClientSet(cfg.Connections.Cluster2.KubeConfigFile)
	cfg.Connections.Cluster2.ConfigStruct = common.YamlToStruct(cfg.Connections.Cluster2.KubeConfigFile)

	//fmt.Printf("%#v", cfg)

	if len(cfg.Connections.Namespaces) < 1 {
		return nil, fmt.Errorf("list of namespaces for comparison can not be empty")
	}

	log.Infof("Analyzing objects in %s namespace(s)", cfg.Connections.Namespaces)

	log.Debugw("Filling the skip list...")

	cfg.Skips, err = ParseSkipConfig(ctx, cfg.Common.Skips)
	if err != nil {
		log.Errorf("cannot parse skip entities list: %s", err.Error())
	}

	return cfg, nil
}
