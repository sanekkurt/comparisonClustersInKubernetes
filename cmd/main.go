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

	versionApiServerCluster1, _ := cfg.Connections.Cluster1.ClientSet.Discovery().ServerVersion()
	versionApiServerCluster2, _ := cfg.Connections.Cluster2.ClientSet.Discovery().ServerVersion()
	if *versionApiServerCluster1 == *versionApiServerCluster2 {
		log.Infof("discovered kube-apiserver version(s): %s", *versionApiServerCluster1)
	} else {
		log.Infof("discovered kube-apiserver version(s): %s vs %s", *versionApiServerCluster1, *versionApiServerCluster2)
	}

	_, err = kube.CompareClusters(ctx)
	if err != nil {
		log.Errorf("cannot compare clusters: %s", err.Error())
		return
	}

	log.Debugw("k8s-cluster-comparator completed")
}
