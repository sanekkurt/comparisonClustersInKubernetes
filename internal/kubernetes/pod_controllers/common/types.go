package common

import (
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

// AbstractPodController is an generalized abstraction above deployment/statefulset/daemonset/etc kubernetes pod controllers
type AbstractPodController struct {
	Name string

	Metadata types.AbstractObjectMetadata

	Labels      map[string]string
	Annotations map[string]string

	Replicas *int32

	PodLabelSelector *v12.LabelSelector
	PodTemplateSpec  v1.PodTemplateSpec
}
