package ingress

import "errors"

var (
	ErrorTLSCountDifferent       = errors.New("the TLS count in the ingresses are different")
	ErrorTLSInIngressesDifferent = errors.New("the TLS in the ingresses are different")

	ErrorSecretNameInTLSDifferent      = errors.New("the secret name in the TLS are different in ingress specs")
	ErrorHostsCountDifferent           = errors.New("the hosts count in the TLS are different in ingress specs")
	ErrorHostsInIngressesDifferent     = errors.New("the hosts in the ingress specs are different")
	ErrorNameHostDifferent             = errors.New("the name host in the TLS are different in ingress specs")
	ErrorBackendInIngressesDifferent   = errors.New("the backend in the ingress specs are different")
	ErrorBackendServicePortDifferent   = errors.New("the service port in the backend are different in ingress specs")
	ErrorServiceNameInBackendDifferent = errors.New("the service name in the backend are different in ingress specs")
	ErrorRulesCountDifferent           = errors.New("the rules count in the ingress specs is different")
	ErrorRulesInIngressesDifferent     = errors.New("rules are different in the ingress specs")
	ErrorHostNameInRuleDifferent       = errors.New("ingress hosts names in specs are different")
	ErrorHTTPInIngressesDifferent      = errors.New("the HTTP in the ingress specs is different")
	ErrorPathsCountDifferent           = errors.New("the paths count in the ingresses is different in ingress specs")
	ErrorPathValueDifferent            = errors.New("the path value in the ingresses is different in ingress specs")
)
