package ingress

import "errors"

var (
	ErrorTLSCountDifferent       = errors.New("the TLS count in the ingresses are different")
	ErrorTLSInIngressesDifferent = errors.New("the TLS in the ingresses are different")

	ErrorSecretNameInTLSDifferent         = errors.New("the secret name in the TLS are different in ingress specs")
	ErrorHostsCountDifferent              = errors.New("the hosts count in the TLS are different in ingress specs")
	ErrorHostsInIngressesDifferent        = errors.New("the hosts in the ingress specs are different")
	ErrorNameHostDifferent                = errors.New("the name host in the TLS are different in ingress specs")
	ErrorBackendInIngressesDifferent      = errors.New("the backend in the ingress specs are different")
	ErrorBackendServicePortDifferent      = errors.New("the service port in the backend are different in ingress specs")
	ErrorBackendServiceIsMissing          = errors.New("the service in the backend is missing in one of the ingress specs")
	ErrorBackendResourceApiGroup          = errors.New("the api group in the backend resource are different in ingress specs")
	ErrorBackendResourceApiGroupIsMissing = errors.New("the api group in the backend resource is missing in one of the ingress specs")
	ErrorBackendResourceName              = errors.New("the name in the backend resource are different in ingress specs")
	ErrorBackendResourceKind              = errors.New("the kind in the backend resource are different in ingress specs")
	ErrorServiceNameInBackendDifferent    = errors.New("the service name in the backend are different in ingress specs")
	ErrorRulesCountDifferent              = errors.New("the rules count in the ingress specs is different")
	ErrorRulesInIngressesDifferent        = errors.New("rules are different in the ingress specs")
	ErrorHostNameInRuleDifferent          = errors.New("ingress hosts names in specs are different")
	ErrorHTTPInIngressesDifferent         = errors.New("the HTTP in the ingress specs is different")
	ErrorPathsCountDifferent              = errors.New("the paths count in the ingresses is different in ingress specs")
	ErrorPathValueDifferent               = errors.New("the path value in the ingresses is different in ingress specs")
	ErrorApiV1Beta1NotSupported           = errors.New("the resource ingress api v1beta1 is not supported")
	ErrorApiV1NotSupported                = errors.New("the resource ingress api v1 is not supported")
)
