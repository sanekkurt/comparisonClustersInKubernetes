package main

import (
	"os"
	"testing"
)

func TestEntireComparison(t *testing.T) {
	if err := os.Setenv("DEBUG", "true"); err != nil {
		t.Error(err)
	}

	Run([]string{os.Args[0], "-c", "../../config.yaml"})
}
