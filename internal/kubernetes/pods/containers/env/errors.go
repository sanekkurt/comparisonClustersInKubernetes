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
	ErrorVarDifferentValues                            = errors.New("variable has different values")
	ErrorVarDifferentValSources                        = errors.New("variable has different value sources")
	ErrorVarDifferentValSourceConfigMaps               = errors.New("variable has different value source ConfigMaps")
	ErrorVarDifferentKeyConfigMaps                     = errors.New("variable has different value source ConfigMap keys")

	ErrorVarDifferentValSourceSecrets = errors.New("variable has different value source Secrets")
	ErrorVarDifferentKeySecrets       = errors.New("variable has different value source Secret keys")

	ErrorVarDifferentValSourceFieldRef         = errors.New("variable has different fieldRef value sources")
	ErrorVarDifferentValSourceResourceFieldRef = errors.New("variable has different resourceFieldRef value sources")

	ErrorVarDoesNotExistInOtherCluster = errors.New("env variable does not exist in other cluster")
)
