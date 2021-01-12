package containers

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
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

func compareContainerHealthCheckParams(ctx context.Context, container1, container2 corev1.Container) (bool, error) {
	log := logging.FromContext(ctx)

	if container1.LivenessProbe != nil && container2.LivenessProbe != nil {

		log.Debugf("ComparePodSpecs: start checking LivenessProbe in container - %s", container1.Name)
		if err := compareContainerProbes(ctx, *container1.LivenessProbe, *container2.LivenessProbe, container1.Name, ErrorContainerHealthCheckLivenessProbeDifferent); err != nil {
			return true, err
		}

	} else if container1.LivenessProbe != nil || container2.LivenessProbe != nil {

		log.Debugf("ComparePodSpecs: start checking LivenessProbe in container - %s", container1.Name)
		return true, fmt.Errorf("%w. Name container: %s. One of the containers is missing Liveness probe", ErrorContainerHealthCheckLivenessProbeDifferent, container1.Name)

	} else {
		log.Debugf("ComparePodSpecs: start checking LivenessProbe in container - %s, but unfortunately they are equal to nil", container1.Name)
	}

	if container1.ReadinessProbe != nil && container2.ReadinessProbe != nil {

		log.Debugf("ComparePodSpecs: start checking ReadinessProbe in container - %s", container1.Name)
		if err := compareContainerProbes(ctx, *container1.ReadinessProbe, *container2.ReadinessProbe, container1.Name, ErrorContainerHealthCheckReadinessProbeDifferent); err != nil {
			return true, err
		}

	} else if container1.ReadinessProbe != nil || container2.ReadinessProbe != nil {

		log.Debugf("ComparePodSpecs: start checking ReadinessProbe in container - %s", container1.Name)
		return true, fmt.Errorf("%w. Name container: %s. One of the containers is missing Readiness probe", ErrorContainerHealthCheckReadinessProbeDifferent, container1.Name)

	} else {
		log.Debugf("ComparePodSpecs: start checking ReadinessProbe in container - %s, but unfortunately they are equal to nil", container1.Name)
	}

	return false, nil
}

func compareContainerProbes(ctx context.Context, probe1, probe2 corev1.Probe, nameContainer string, er error) error {
	if probe1.Exec != nil && probe2.Exec != nil {

		err := CompareMassStringsInContainers(ctx, probe1.Exec.Command, probe2.Exec.Command)
		if err != nil {
			return fmt.Errorf("%s. Containers name: %s. %w: %s", er, nameContainer, ErrorContainerHealthCheckDifferentExecCommand, err)
		}

	} else if probe1.Exec != nil || probe2.Exec != nil {
		return fmt.Errorf("%s. Containers name: %s. %w", er, nameContainer, ErrorContainerHealthCheckDifferentExec)
	}

	if probe1.TCPSocket != nil && probe2.TCPSocket != nil {

		if probe1.TCPSocket.Host != probe2.TCPSocket.Host {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 host - %s, container 2 host - %s", er, nameContainer, ErrorContainerHealthCheckDifferentTCPSocketHost, probe1.TCPSocket.Host, probe2.TCPSocket.Host)
		}

		if probe1.TCPSocket.Port.IntVal != probe2.TCPSocket.Port.IntVal || probe1.TCPSocket.Port.StrVal != probe2.TCPSocket.Port.StrVal || probe1.TCPSocket.Port.Type != probe2.TCPSocket.Port.Type {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 port - %s, container 2 port - %s", er, nameContainer, ErrorContainerHealthCheckDifferentTCPSocketPort, fmt.Sprintln(probe1.TCPSocket.Port), fmt.Sprintln(probe2.TCPSocket.Port))
		}
	} else if probe1.TCPSocket != nil || probe2.TCPSocket != nil {
		return fmt.Errorf("%s. Containers name: %s. %w", er, nameContainer, ErrorContainerHealthCheckDifferentTCPSocket)
	}

	if probe1.HTTPGet != nil && probe2.HTTPGet != nil {

		if probe1.HTTPGet.Host != probe2.HTTPGet.Host {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 host - %s, container 2 host - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetHost, probe1.HTTPGet.Host, probe2.HTTPGet.Host)
		}

		if probe1.HTTPGet.HTTPHeaders != nil && probe2.HTTPGet.HTTPHeaders != nil {
			if len(probe1.HTTPGet.HTTPHeaders) != len(probe2.HTTPGet.HTTPHeaders) {
				return fmt.Errorf("%s. Containers name: %s. %w: container 1 count - %d, container 2 count - %d", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetHTTPHeaders, len(probe1.HTTPGet.HTTPHeaders), len(probe2.HTTPGet.HTTPHeaders))
			}

			for index, value := range probe1.HTTPGet.HTTPHeaders {
				if value.Name != probe2.HTTPGet.HTTPHeaders[index].Name {
					return fmt.Errorf("%s. Containers name: %s. %w: container 1 header name - %s, container 2 header name - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetHeaderName, value.Name, probe2.HTTPGet.HTTPHeaders[index].Name)
				}

				if value.Value != probe2.HTTPGet.HTTPHeaders[index].Value {
					return fmt.Errorf("%s. Containers name: %s. %w. Name header - %s. Container 1 header value - %s, container 2 header value - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetHeaderValue, value.Name, value.Value, probe2.HTTPGet.HTTPHeaders[index].Value)
				}
			}

		} else if probe1.HTTPGet.HTTPHeaders != nil || probe2.HTTPGet.HTTPHeaders != nil {
			return fmt.Errorf("%s. Containers name: %s. %w", er, nameContainer, ErrorContainerHealthCheckHTTPGetMissingHeader)
		}

		if probe1.HTTPGet.Path != probe2.HTTPGet.Path {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 path - %s, container 2 path - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetPath, probe1.HTTPGet.Path, probe2.HTTPGet.Path)
		}

		if probe1.HTTPGet.Port.IntVal != probe2.HTTPGet.Port.IntVal || probe1.HTTPGet.Port.StrVal != probe2.HTTPGet.Port.StrVal || probe1.HTTPGet.Port.Type != probe2.HTTPGet.Port.Type {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 port - %s, container 2 port - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetPort, fmt.Sprintln(probe1.HTTPGet.Port), fmt.Sprintln(probe2.HTTPGet.Port))
		}

		if probe1.HTTPGet.Scheme != probe2.HTTPGet.Scheme {
			return fmt.Errorf("%s. Containers name: %s. %w: container 1 port - %s, container 2 port - %s", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGetScheme, fmt.Sprintln(probe1.HTTPGet.Port), fmt.Sprintln(probe2.HTTPGet.Port))
		}

	} else if probe1.HTTPGet != nil || probe2.HTTPGet != nil {
		return fmt.Errorf("%s. Containers name: %s. %w", er, nameContainer, ErrorContainerHealthCheckDifferentHTTPGet)
	}

	if probe1.FailureThreshold != probe2.FailureThreshold {
		return fmt.Errorf("%s. Containers name: %s. %w: container 1 - %d, container 2 - %d", er, nameContainer, ErrorContainerHealthCheckDifferentFailureThreshold, probe1.FailureThreshold, probe2.FailureThreshold)
	}

	if probe1.InitialDelaySeconds != probe2.InitialDelaySeconds {
		return fmt.Errorf("%s. Containers name: %s. %w: container 1 - %d, container 2 - %d", er, nameContainer, ErrorContainerHealthCheckDifferentInitialDelaySeconds, probe1.InitialDelaySeconds, probe2.InitialDelaySeconds)
	}

	if probe1.PeriodSeconds != probe2.PeriodSeconds {
		return fmt.Errorf("%s. Containers name: %s. %w: container 1 - %d, container 2 - %d", er, nameContainer, ErrorContainerHealthCheckDifferentPeriodSeconds, probe1.PeriodSeconds, probe2.PeriodSeconds)
	}

	if probe1.SuccessThreshold != probe2.SuccessThreshold {
		return fmt.Errorf("%s. Containers name: %s. %w: container 1 - %d, container 2 - %d", er, nameContainer, ErrorContainerHealthCheckDifferentSuccessThreshold, probe1.SuccessThreshold, probe2.SuccessThreshold)
	}

	if probe1.TimeoutSeconds != probe2.TimeoutSeconds {
		return fmt.Errorf("%s. Containers name: %s.  %w: container 1 - %d, container 2 - %d", er, nameContainer, ErrorContainerHealthCheckDifferentTimeoutSeconds, probe1.TimeoutSeconds, probe2.TimeoutSeconds)
	}

	return nil
}
