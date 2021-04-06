package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"

	"k8s-cluster-comparator/internal/consts"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/logging"
)

var (
	opts struct {
		//	KubeConfig1              string   `long:"kube-config1" env:"KUBECONFIG1" required:"true" description:"Path to Kubernetes client1 config file"`
		//	KubeConfig2              string   `long:"kube-config2" env:"KUBECONFIG2" required:"true" description:"Path to Kubernetes client2 config file"`
		//	NameSpaces               []string `long:"ns" env:"NAMESPACES" required:"true" description:"Configmaps massive"`
		//	SkippedFullResourceNames string   `long:"skip" env:"SKIP" required:"false" description:"Skipping a resource by kind/name specification"`
		//	SkipAnyKindResourceNames string   `long:"skip-names" env:"SKIP_NAMES" required:"false" description:"Resource skipping by name"`

		ConfigPath string `long:"config" short:"c" env:"CONFIG_PATH" description:"Path to config.yaml file" required:"true"`
	}

	ErrHelpShown = errors.New("help message shown")
)

// Parse performs configuration parsing from various sources and fills in the AppConfig struct
func Parse(ctx context.Context, args []string) (*AppConfig, error) {
	log := logging.FromContext(ctx)

	_, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).ParseArgs(args[1:])
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				return nil, ErrHelpShown
			}

			return nil, fmt.Errorf("cannot parse arguments: %w", flagsErr)
		}

		return nil, fmt.Errorf("cannot parse arguments: %w", err)
	}

	cfg, err := ParseConfig(ctx, opts.ConfigPath)
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

	log.Infof("Analyzing objects in [%s] namespace(s)", strings.Join(cfg.Connections.Namespaces, ", "))

	log.Debugw("Filling the skip list...")

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		cfg.ExcludesIncludes, err = ParseSkipConfig(ctx, cfg.Common.Excludes)
		if err != nil {
			log.Errorf("cannot parse exclude entities list: %s", err.Error())
		}
	} else {
		cfg.ExcludesIncludes, err = ParseSkipConfig(ctx, cfg.Common.Includes)
		if err != nil {
			log.Errorf("cannot parse include entities list: %s", err.Error())
		}
	}
	cfg.ExcludesIncludes.WorkMode = cfg.Common.WorkMode

	log.Infof("%s work mode", cfg.Common.WorkMode)

	return cfg, nil
}

//func ParseForTests(ctx context.Context, args []string) (*AppConfig, error) {
//
//}
