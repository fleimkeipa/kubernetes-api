package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceInterfaces interface {
	Create(context.Context, *corev1.Namespace, metav1.CreateOptions) (*corev1.Namespace, error)
	Update(context.Context, *corev1.Namespace, metav1.UpdateOptions) (*corev1.Namespace, error)
	List(context.Context, metav1.ListOptions) (*corev1.NamespaceList, error)
	Delete(context.Context, string, metav1.DeleteOptions) error
}
