package main

import (
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"os"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/interrupt"
	kube "k8s-cluster-comparator/internal/kubernetes"
	"k8s-cluster-comparator/internal/kubernetes/discovery"
	"k8s-cluster-comparator/internal/logging"
)

func main() {
	var (
		debug bool
	)

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	ctx, doneFn := interrupt.Context()

	defer func() {
		doneFn()
	}()

	err := logging.Configure(debug)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	log := logging.FromContext(ctx)

	cfg, err := config.Parse(ctx)
	if err != nil {
		if err == config.ErrHelpShown {
			return
		}

		log.Error(err)
		return
	}

	log.Infow("Starting k8s-cluster-comparator")
	defer func() {
		log.Infow("k8s-cluster-comparator completed")
	}()

	ctx = config.With(ctx, cfg)
	ctx = diff.With(ctx, &diff.DiffsStorage{
		Batches: make([]diff.DiffsBatch, 0, 0),
	})

	err = discovery.DetectKubeVersions(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	err = kube.CompareClusters(ctx)
	if err != nil {
		log.Errorf("cannot compare clusters: %s", err.Error())
		return
	}

	log.Debugw("k8s-cluster-comparator completed")
}
