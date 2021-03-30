package main

import (
	"os"
	"testing"
)

func TestEntireComparison(t *testing.T) {
	err := os.Setenv("DEBUG", "true")
	if err != nil {
		t.Error(err)
	}

	main()
}
