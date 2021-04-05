package comparatormain

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/interrupt"
	kube "k8s-cluster-comparator/internal/kubernetes"
	"k8s-cluster-comparator/internal/kubernetes/discovery"
	"k8s-cluster-comparator/internal/logging"
	"os"
)

func Main(args []string) {
	var (
		debug bool
		ctx   = context.Background()
	)

	os.Args = args

	ctx, doneFn := interrupt.Context(ctx)
	defer func() {
		doneFn()
	}()

	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

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

	ctx = config.With(ctx, cfg)

	log.Infow("Starting k8s-cluster-comparator")
	defer func() {
		log.Infow("k8s-cluster-comparator completed")
	}()

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
