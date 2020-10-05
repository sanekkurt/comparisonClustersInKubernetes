package main

import (
	"fmt"
	"os"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/interrupt"
	kube "k8s-cluster-comparator/internal/kubernetes"
	"k8s-cluster-comparator/internal/logging"
)

func main() {
	var debug bool
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	ctx, doneFn := interrupt.Context()
	defer doneFn()

	err := logging.Configure(debug)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		os.Exit(1)
	}

	log := logging.FromContext(ctx)

	log.Infow("Starting k8s-cluster-comparator")

	cfg, err := config.Parse(ctx)
	if err != nil {
		log.Fatalf("configuration error: %s", err.Error())
	}

	ret := 0

	isClustersDiffer, err := kube.CompareClusters(cfg)
	if err != nil {
		log.Errorf("cannot compare clusters: %s", err.Error())
		os.Exit(2)
	}

	if isClustersDiffer {
		ret = 1
	}

	log.Infow("k8s-cluster-comparator completed")

	os.Exit(ret)
}
