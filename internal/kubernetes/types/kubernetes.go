package types

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ObjectKind string
type ObjectName string

// ObjectKindWrapper is a wrapper function that transforms object kind name to a canonical form
func ObjectKindWrapper(kind string) string {
	return strings.ToLower(kind)
}

// Container to describe the main information in the comparison container
type Container struct {
	Name    string
	Image   string
	ImageID string
}

// InformationAboutObject for generalizing the comparison function, which allows you to pass information to it from both deployment and statefulset
type InformationAboutObject struct {
	Template corev1.PodTemplateSpec
	Selector *metav1.LabelSelector
}

type AbstractObjectMetadata struct {
	Type metav1.TypeMeta
	Meta metav1.ObjectMeta
}

type KubeConnections struct {
	C1 kubernetes.Interface
	C2 kubernetes.Interface

	Namespace string
}
