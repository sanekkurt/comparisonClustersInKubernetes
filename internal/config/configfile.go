package config

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

type PodsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type DeploymentsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type StatefulSetsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type DaemonSetsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type PodControllersComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Deployments  DeploymentsComparisonConfiguration  `yaml:"deployments"`
	StatefulSets StatefulSetsComparisonConfiguration `yaml:"statefulSets"`
	DaemonSets   DaemonSetsComparisonConfiguration   `yaml:"daemonSets"`
}

type WorkloadsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Pods           PodsComparisonConfiguration           `yaml:"pods"`
	PodControllers PodControllersComparisonConfiguration `yaml:"podControllers"`
}

type JobsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type CronJobsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type TasksComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type ServicesComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type IngressesComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type NetworkingComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Services  ServicesComparisonConfiguration  `yaml:"services"`
	Ingresses IngressesComparisonConfiguration `yaml:"ingresses"`
}

// ClusterConnection represents a k8s-cluster config
type ClusterConnection struct {
	KubeConfigFile string `yaml:"kubeConfig"`

	ClientSet    kubernetes.Interface  `yaml:"-"`
	ConfigStruct *types.KubeconfigYaml `yaml:"-"`
}

type ConnectionsConfigurations struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Cluster1 ClusterConnection `yaml:"cluster1"`
	Cluster2 ClusterConnection `yaml:"cluster2"`

	Namespaces []string `yaml:"namespaces,omitempty"`
}

type CommonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	WorkMode WorkMode `yaml:"mode,omitempty"`

	Skips SkipConfiguration `yaml:"skips" json:"skips"`
}

// AppConfig is the main application configuration storage
type AppConfig struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Version string `yaml:"version" yaml:"version,omitempty"`

	Connections ConnectionsConfigurations `yaml:"connections" json:"connections"`

	Common CommonConfiguration `yaml:"common" yaml:"common"`

	Workloads WorkloadsComparisonConfiguration `yaml:"workloads"`
	Tasks     TasksComparisonConfiguration     `yaml:"tasks"`

	Networking NetworkingComparisonConfiguration `yaml:"networking"`

	Skips skipper.SkipEntitiesList `yaml:"-"`
}

func ParseConfig(ctx context.Context, cfgPath string) (*AppConfig, error) {
	log := logging.FromContext(ctx)

	f, err := os.Open(cfgPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("cannot open config '%s' for reading: %w", cfgPath, err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Warnf("cannot close config '%s': %s", cfgPath, err)
		}
	}()

	dec := yaml.NewDecoder(f)
	dec.SetStrict(true)

	cfg := &AppConfig{}

	err = dec.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config from '%s': %w", cfgPath, err)
	}

	return cfg, nil
}
