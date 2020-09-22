package kubernetes

import (
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type NameComponents map[string]struct{}

var (
	VariableForNamespaces []string
	Kubeconfig1YamlStruct KubeconfigYaml
	Kubeconfig2YamlStruct KubeconfigYaml
	Client1               *kubernetes.Clientset
	Client2               *kubernetes.Clientset

	Entities              map[string]NameComponents
	SkipTypes             = [3]v12.SecretType{"kubernetes.io/service-account-token", "kubernetes.io/dockercfg", "helm.sh/release.v1"}
)
