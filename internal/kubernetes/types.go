package kubernetes

import (
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsAlreadyComparedFlag to indicate whether the information of this entity was compared
type IsAlreadyComparedFlag struct {
	Index int
	Check bool
}

// Container to describe the main information in the comparison container
type Container struct {
	Name    string
	Image   string
	ImageID string
}

// InformationAboutObject for generalizing the comparison function, which allows you to pass information to it from both deployment and statefulset
type InformationAboutObject struct {
	Template v12.PodTemplateSpec
	Selector *v1.LabelSelector
}

// KubeconfigYaml structure for describing the kubeconfig Yaml file
type KubeconfigYaml struct {
	APIVersion string `json:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data"`
			Server                   string `json:"server"`
		} `json:"cluster"`
		Name string `json:"name"`
	} `json:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
		Name string `json:"name"`
	} `json:"contexts"`
	CurrentContext string `json:"current-context"`
	Kind           string `json:"kind"`
	Preferences    struct {
	} `json:"preferences"`
	Users []struct {
		Name string `json:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data"`
		} `json:"user"`
	} `json:"users"`
}
