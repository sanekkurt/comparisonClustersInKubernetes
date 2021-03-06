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
	var (
		debug bool
		ret   int
	)

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	ctx, doneFn := interrupt.Context()

	defer func() {
		doneFn()
		os.Exit(ret)
	}()

	err := logging.Configure(debug)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())

		ret = 1
		return
	}

	log := logging.FromContext(ctx)

	cfg, err := config.Parse(ctx)
	if err != nil {
		if err == config.ErrHelpShown {
			return
		}

		ret = 1
		return
	}

	log.Infow("Starting k8s-cluster-comparator")

	isClustersDiffer, err := kube.CompareClusters(ctx, cfg)
	if err != nil {
		log.Errorf("cannot compare clusters: %s", err.Error())

		ret = 2
		return
	}

	if isClustersDiffer {
		ret = 1
	}

	log.Infow("k8s-cluster-comparator completed")
}
