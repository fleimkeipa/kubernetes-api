package repositories

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceRepository struct {
	client *kubernetes.Clientset
}

func NewNamespaceRepository(client *kubernetes.Clientset) *NamespaceRepository {
	return &NamespaceRepository{
		client: client,
	}
}

func (rc *NamespaceRepository) Create(ctx context.Context, namespace *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return rc.client.CoreV1().Namespaces().Create(ctx, namespace, opts)
}

func (rc *NamespaceRepository) Update(ctx context.Context, namespace *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	return rc.client.CoreV1().Namespaces().Update(ctx, namespace, opts)
}

func (rc *NamespaceRepository) List(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error) {
	return rc.client.CoreV1().Namespaces().List(ctx, opts)
}

func (rc *NamespaceRepository) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return rc.client.CoreV1().Namespaces().Delete(ctx, name, opts)
}
