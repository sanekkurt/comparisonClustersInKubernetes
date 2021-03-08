package ingress

import "errors"

var (
	ErrorTLSCountDifferent       = errors.New("the TLS count in the ingresses are different")
	ErrorTLSInIngressesDifferent = errors.New("the TLS in the ingresses are different")

	ErrorSecretNameInTLSDifferent      = errors.New("the secret name in the TLS are different")
	ErrorHostsCountDifferent           = errors.New("the hosts count in the TLS are different")
	ErrorHostsInIngressesDifferent     = errors.New("the hosts in the ingresses are different")
	ErrorNameHostDifferent             = errors.New("the name host in the TLS are different")
	ErrorBackendInIngressesDifferent   = errors.New("the backend in the ingresses are different")
	ErrorBackendServicePortDifferent   = errors.New("the service port in the backend are different")
	ErrorServiceNameInBackendDifferent = errors.New("the service name in the backend are different")
	ErrorRulesCountDifferent           = errors.New("the rules count in the ingresses is different")
	ErrorRulesInIngressesDifferent     = errors.New("ingress rules are different")
	ErrorHostNameInRuleDifferent       = errors.New("ingress hosts names are different")
	ErrorHTTPInIngressesDifferent      = errors.New("the HTTP in the ingresses is different")
	ErrorPathsCountDifferent           = errors.New("the paths count in the ingresses is different")
	ErrorPathValueDifferent            = errors.New("the path value in the ingresses is different")
)
