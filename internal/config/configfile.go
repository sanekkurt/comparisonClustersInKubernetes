package config

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/consts"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	"k8s-cluster-comparator/internal/utils"
)

type ConfigurationSection interface {
	Parse(ctx context.Context) error
}

type CompareContainersConfiguration struct {
	RollingTags struct {
		TagsList    []string            `yaml:"tagsList"`
		TagsListMap map[string]struct{} `yaml:"-"`

		WarnOnRollingTag bool `yaml:"warnOnRollingTag"`
	} `yaml:"rollingTags"`

	Env struct {
		EnvFrom struct {
			DeepCompareAlways       bool `yaml:"deepCompareAlways"`
			DeepCompareOnRollingTag bool `yaml:"deepCompareOnRollingTag"`
		} `yaml:"envFrom"`
	} `yaml:"env"`

	Image struct {
		Mirrors []struct {
			From string `yaml:"from"`
			To   string `yaml:"to"`
		} `yaml:"mirrors"`
	} `yaml:"image""`
}

func (cfg *CompareContainersConfiguration) Parse(ctx context.Context) error {
	var err error

	cfg.RollingTags.TagsListMap, err = utils.StringsListToMap(ctx, cfg.RollingTags.TagsList, false)
	if err != nil {
		return fmt.Errorf("cannot parse rolling tags list: %w", err)
	}

	return nil
}

type PodsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`
}

type DeploymentsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`

	DiscardDeploymentsUpdatedLaterTime int64 `yaml:"discardDeploymentsUpdatedLaterTime"`
}

type StatefulSetsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type DaemonSetsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type PodControllersComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	CompareImageDigestsAlways       bool `yaml:"compareImageDigestsAlways"`
	CompareImageDigestsOnRollingTag bool `yaml:"compareImageDigestsOnRollingTag"`

	Deployments  DeploymentsComparisonConfiguration  `yaml:"deployments"`
	StatefulSets StatefulSetsComparisonConfiguration `yaml:"statefulSets"`
	DaemonSets   DaemonSetsComparisonConfiguration   `yaml:"daemonSets"`
}

type WorkloadsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	Containers CompareContainersConfiguration `yaml:"containers"`

	Pods           PodsComparisonConfiguration           `yaml:"pods"`
	PodControllers PodControllersComparisonConfiguration `yaml:"podControllers"`
}

type JobsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type CronJobsComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type TasksComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	Jobs     JobsComparisonConfiguration     `yaml:"jobs"`
	CronJobs CronJobsComparisonConfiguration `yaml:"cronJobs"`
}

type ServicesComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type IngressesComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	BatchSize int64 `yaml:"batchSize"`
}

type NetworkingComparisonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	Services  ServicesComparisonConfiguration  `yaml:"services"`
	Ingresses IngressesComparisonConfiguration `yaml:"ingresses"`
}

// ClusterConnection represents a k8s-cluster config
type ClusterConnection struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	Name string `yaml:"name"`

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

type MetadataCompareConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	DumpDifferentValues bool `yaml:"dumpDifferentValues"`

	SkipLabels    []string            `yaml:"skipLabelsList"`
	SkipLabelsMap map[string]struct{} `yaml:"-"`

	SkipAnnotations    []string            `yaml:"skipAnnotationsList"`
	SkipAnnotationsMap map[string]struct{} `yaml:"-"`
}

func (cfg *MetadataCompareConfiguration) Parse(ctx context.Context) error {
	var err error

	cfg.SkipLabelsMap, err = utils.StringsListToMap(ctx, cfg.SkipLabels, false)
	if err != nil {
		return fmt.Errorf("cannot parse skip labels list: %w", err)
	}

	cfg.SkipAnnotationsMap, err = utils.StringsListToMap(ctx, cfg.SkipAnnotations, false)
	if err != nil {
		return fmt.Errorf("cannot parse skip annotations list: %w", err)
	}

	return nil
}

type CommonConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	WorkMode string `yaml:"mode""`

	Excludes ExcludeIncludeSpec `yaml:"excludes"`
	Includes ExcludeIncludeSpec `yaml:"includes"`

	MetadataCompareConfiguration MetadataCompareConfiguration `yaml:"metadata"`

	DefaultBatchSize int64 `yaml:"defaultBatchSize"`

	CheckingCreationTimestampDeploymentsLimit bool `yaml:"checkingCreationTimestampDeploymentsLimit"`
}

func (cfg *CommonConfiguration) Parse(ctx context.Context) error {
	allowedModes := map[string]struct{}{
		consts.EverythingButNotExcludesWorkMode: {},
		consts.NothingButIncludesWorkMode:       {},
	}

	if _, ok := allowedModes[cfg.WorkMode]; !ok {
		return fmt.Errorf("unknown work mode '%s'", cfg.WorkMode)
	}

	return nil
}

type ConfigMapsConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	DumpDifferentValues bool `yaml:"dumpDifferentValues"`

	BatchSize int64 `yaml:"batchSize"`
}

type SecretsConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	DumpDifferentValues bool `yaml:"dumpDifferentValues"`

	BatchSize int64 `yaml:"batchSize"`

	SkipTypes    []string            `yaml:"skipTypesList"`
	SkipTypesMap map[string]struct{} `yaml:"-"`
}

func (cfg *SecretsConfiguration) Parse(ctx context.Context) error {
	var err error

	cfg.SkipTypesMap, err = utils.StringsListToMap(ctx, cfg.SkipTypes, false)
	if err != nil {
		return fmt.Errorf("cannot parse skip secret types list: %w", err)
	}

	return nil
}

type ConfigsConfiguration struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty" json:"_,omitempty"`

	Enabled bool `yaml:"enabled"`

	ConfigMaps ConfigMapsConfiguration `yaml:"configMaps"`
	Secrets    SecretsConfiguration    `yaml:"secrets"`
}

// AppConfig is the main application configuration storage
type AppConfig struct {
	// Prevent comparisons
	_ [0][]byte `yaml:"_,omitempty"`

	Version string `yaml:"version" yaml:"version,omitempty"`

	Connections ConnectionsConfigurations `yaml:"connections" json:"connections"`

	Common CommonConfiguration `yaml:"common" yaml:"common"`

	Configs ConfigsConfiguration `yaml:"configs"`

	Workloads WorkloadsComparisonConfiguration `yaml:"workloads"`
	Tasks     TasksComparisonConfiguration     `yaml:"tasks"`

	Networking NetworkingComparisonConfiguration `yaml:"networking"`

	ExcludesIncludes skipper.SkipEntitiesList `yaml:"-"`
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

	var extraCfgSections = []ConfigurationSection{
		&cfg.Common.MetadataCompareConfiguration,
		&cfg.Configs.Secrets,
		&cfg.Workloads.Containers,
	}

	for _, section := range extraCfgSections {
		err := section.Parse(ctx)
		if err != nil {
			return nil, fmt.Errorf("cannot parse '%T' configuration section: %w", section, err)
		}
	}

	return cfg, nil
}
