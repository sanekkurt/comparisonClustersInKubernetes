package healthcheck

import (
	"errors"
)

var (
	ErrorContainerHealthCheckLivenessProbeDifferent  = errors.New("one of the containers does not have a liveness probe")
	ErrorContainerHealthCheckReadinessProbeDifferent = errors.New("one of the containers does not have a readiness probe")

	ErrorContainerHealthCheckDifferentExec        = errors.New("the exec command missing in one probe")
	ErrorContainerHealthCheckDifferentExecCommand = errors.New("the exec command in probe not equal")

	ErrorContainerHealthCheckDifferentTCPSocket           = errors.New("the TCPSocket missing in one probe")
	ErrorContainerHealthCheckDifferentTCPSocketHost       = errors.New("the TCPSocket.Host in probe not equal")
	ErrorContainerHealthCheckDifferentTCPSocketPortIntVal = errors.New("the TCPSocket.Port.IntVal in probe not equal")
	ErrorContainerHealthCheckDifferentTCPSocketPortStrVal = errors.New("the TCPSocket.Port.StrVal in probe not equal")
	ErrorContainerHealthCheckDifferentTCPSocketPortType   = errors.New("the TCPSocket.Port.Type in probe not equal")
	//	ErrorContainerHealthCheckDifferentTCPSocketPort = errors.New("the TCPSocket.Port in probe not equal")

	ErrorContainerHealthCheckDifferentHTTPGet            = errors.New("the HTTPGet missing in one probe")
	ErrorContainerHealthCheckDifferentHTTPGetHost        = errors.New("the HTTPGet.Host in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders = errors.New("the HTTPGet.HTTPHeaders in probe not equal")

	ErrorContainerHealthCheckDifferentHTTPGetHeaderName  = errors.New("the name header in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetHeaderValue = errors.New("the value header in probe not equal")

	ErrorContainerHealthCheckHTTPGetMissingHeader       = errors.New("one of the containers is missing headers")
	ErrorContainerHealthCheckDifferentHTTPGetPortIntVal = errors.New("the HTTPGet.Port.IntVal in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetPortStrVal = errors.New("the HTTPGet.Port.StrVal in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetPortType   = errors.New("the HTTPGet.Port.Type in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetPath       = errors.New("the HTTPGet.Path in probe not equal")
	//	ErrorContainerHealthCheckDifferentHTTPGetPort   = errors.New("the HTTPGet.Port in probe not equal")
	ErrorContainerHealthCheckDifferentHTTPGetScheme = errors.New("the HTTPGet.Scheme in probe not equal")

	ErrorContainerHealthCheckDifferentFailureThreshold    = errors.New("the FailureThreshold in probe not equal")
	ErrorContainerHealthCheckDifferentInitialDelaySeconds = errors.New("the InitialDelaySeconds in probe not equal")
	ErrorContainerHealthCheckDifferentPeriodSeconds       = errors.New("the PeriodSeconds in probe not equal")
	ErrorContainerHealthCheckDifferentSuccessThreshold    = errors.New("the SuccessThreshold in probe not equal")
	ErrorContainerHealthCheckDifferentTimeoutSeconds      = errors.New("the TimeoutSeconds in probe not equal")
)
