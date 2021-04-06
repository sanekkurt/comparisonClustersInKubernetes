package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
)

func main() {
	cfg := config.AppConfig{
		Version:          "v1",
		Connections:      config.ConnectionsConfigurations{},
		Common:           config.CommonConfiguration{},
		Configs:          config.ConfigsConfiguration{},
		Workloads:        config.WorkloadsComparisonConfiguration{},
		Tasks:            config.TasksComparisonConfiguration{},
		Networking:       config.NetworkingComparisonConfiguration{},
		ExcludesIncludes: skipper.SkipEntitiesList{},
	}

	enc := yaml.NewEncoder(os.Stdout)
	defer func() {
		err := enc.Close()
		if err != nil {
			fmt.Println("[ERROR] ", err)
		}
	}()

	err := enc.Encode(cfg)
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}
}
