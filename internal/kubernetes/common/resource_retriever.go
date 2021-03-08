package common

//import (
//	"context"
//
//	"go.uber.org/zap"
//	"k8s-cluster-comparator/internal/config"
//	"k8s-cluster-comparator/internal/logging"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//)
//
//const (
//	secretKind             = "secret"
//)
//
//var (
//	clientSet = struct{}{}
//)
//
//type KubeObjectsDifference struct {
//	ObjectType metav1.TypeMeta
//	ObjectMeta metav1.ObjectMeta
//
//	Message string
//
//	Critical bool
//}
//
//type KubeKindComparator interface {
//	Compare(context.Context, string) ([]KubeObjectsDifference, error)
//	collect()
//
//	FieldSelectorProvider (context.Context) string
//	LabelSelectorProvider (context.Context) string
//
//}
//
//type SecretsComparator struct {
//	Namespace string
//	BatchSize int64
//}
//
//func NewSecretsComparator(ctx context.Context, namespace string) SecretsComparator {
//	return SecretsComparator{
//		Namespace: namespace,
//		BatchSize: 25,
//	}
//}
//
//func (cmp *SecretsComparator) FieldSelectorProvider(ctx context.Context) string {
//	var (
//		cfg = config.FromContext(ctx)
//		fieldSelector = ""
//	)
//
//	for k := range cfg.Configs.Secrets.SkipTypesMap {
//		fieldSelector += "type!=" + k
//	}
//
//	return fieldSelector
//}
//
//func (cmp *SecretsComparator) LabelSelectorProvider(ctx context.Context) string {
//	return ""
//}
//
