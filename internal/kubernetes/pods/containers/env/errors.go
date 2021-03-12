package env

import (
	"errors"
)

var (
	ErrorContainerDifferentEnvVarsNumber = errors.New("different number of environment variables in container specs")
	ErrorContainerDifferentEnvVarNames   = errors.New("different environment variable names in container specs")

	ErrorContainerDifferentEnvVarValues       = errors.New("different values of environment variable in container specs")
	ErrorContainerDifferentEnvVarValueSources = errors.New("different environment variable value sources in container specs")

	ErrorContainerEnvValueFromComparisonNotImplemented = errors.New("environment variable ValueFrom type not implemented yet")
)

