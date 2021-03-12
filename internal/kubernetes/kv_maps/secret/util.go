package secret

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetSecretByName(ctx context.Context, clientSet kubernetes.Interface, namespace, configMapName string) (*corev1.Secret, error) {
	configMap, err := clientSet.CoreV1().Secrets(namespace).Get(configMapName, metav1.GetOptions{})
	return configMap, err
}
