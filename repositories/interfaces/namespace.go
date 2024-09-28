package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceInterfaces interface {
	Create(ctx context.Context, namespace *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error)
	Update(ctx context.Context, namespace *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error)
	List(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error)
	Delete(ctx context.Context, namespaceID string, opts metav1.DeleteOptions) error
}
