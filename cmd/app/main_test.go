package main

import (
	"os"
	"testing"
)

func TestEntireComparison(t *testing.T) {
	if err := os.Setenv("DEBUG", "true"); err != nil {
		t.Error(err)
	}

	run([]string{os.Args[0], "-c", "C:\\Users\\Александр\\go\\src\\comparisonClustersInKubernetes\\config.yaml"})
	//run([]string{os.Args[0], "-c", "../../config.yaml"})
}
