package configmap

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetConfigMapByName(ctx context.Context, clientSet kubernetes.Interface, namespace, configMapName string) (*corev1.ConfigMap, error) {
	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	return configMap, err
}
