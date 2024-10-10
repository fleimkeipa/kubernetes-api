package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
)

type NamespaceInterfaces interface {
	Create(ctx context.Context, namespace *corev1.Namespace, opts model.CreateOptions) (*corev1.Namespace, error)
	Update(ctx context.Context, namespace *corev1.Namespace, opts model.UpdateOptions) (*corev1.Namespace, error)
	List(ctx context.Context, opts model.ListOptions) (*corev1.NamespaceList, error)
	Delete(ctx context.Context, namespaceID string, opts model.DeleteOptions) error
}
