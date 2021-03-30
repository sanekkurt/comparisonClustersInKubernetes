package pods

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodsListOnMatchLabels retrieves list of PODs according to a given selector
func GetPodsListOnMatchLabels(ctx context.Context, clientSet kubernetes.Interface, namespace string, ls *metav1.LabelSelector) ([]corev1.Pod, error) {
	matchLabels, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		return nil, err
	}

	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: matchLabels.String(),
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
