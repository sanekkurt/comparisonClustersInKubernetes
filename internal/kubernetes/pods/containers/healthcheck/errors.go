package healthcheck

import (
	"errors"
)

var (
	ErrorContainerHealthCheckLivenessProbeDifferent  = errors.New("livenessProbe in containers are different")
	ErrorContainerHealthCheckReadinessProbeDifferent = errors.New("readinessProbe in containers are different")

	ErrorContainerHealthCheckDifferentExec        = errors.New("the exec command missing in one probe")
	ErrorContainerHealthCheckDifferentExecCommand = errors.New("the exec command in probe not equal")

	ErrorContainerHealthCheckDifferentTCPSocket     = errors.New("the TCPSocket missing in one probe")
	ErrorContainerHealthCheckDifferentTCPSocketHost = errors.New("the TCPSocket.Host in probe not equal")
	ErrorContainerHealthCheckDifferentTCPSocketPort = errors.New("the TCPSocket.Port in probe not equal")

	ErrorContainerHealthCheckDifferentHTTPGet            = errors.New("the HTTPGet missing in one probe")
	ErrorContainerHealthCheckDifferentHTTPGetHost        = errors.New("the HTTPGet.Host in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders = errors.New("the HTTPGet.HTTPHeaders in probe not equal")

	ErrorContainerHealthCheckDifferentHTTPGetHeaderName  = errors.New("the name header in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetHeaderValue = errors.New("the value header in probe not equal")

	ErrorContainerHealthCheckHTTPGetMissingHeader = errors.New("one of the containers is missing headers")

	ErrorContainerHealthCheckDifferentHTTPGetPath   = errors.New("the HTTPGet.Path in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetPort   = errors.New("the HTTPGet.Port in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetScheme = errors.New("the HTTPGet.Scheme in probe not equal")

	ErrorContainerHealthCheckDifferentFailureThreshold    = errors.New("the FailureThreshold in probe not equal")
	ErrorContainerHealthCheckDifferentInitialDelaySeconds = errors.New("the InitialDelaySeconds in probe not equal")
	ErrorContainerHealthCheckDifferentPeriodSeconds       = errors.New("the PeriodSeconds in probe not equal")
	ErrorContainerHealthCheckDifferentSuccessThreshold    = errors.New("the SuccessThreshold in probe not equal")
	ErrorContainerHealthCheckDifferentTimeoutSeconds      = errors.New("the TimeoutSeconds in probe not equal")
)
