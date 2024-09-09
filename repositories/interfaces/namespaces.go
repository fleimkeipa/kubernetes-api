package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceInterfaces interface {
	Get(context.Context, metav1.ListOptions) (*corev1.NamespaceList, error)
	Create(ctx context.Context, namespace *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error)
	Delete(context.Context, string, metav1.DeleteOptions) error
}
