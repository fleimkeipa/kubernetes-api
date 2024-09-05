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

func (rc *NamespaceRepository) Get(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error) {
	return rc.client.CoreV1().Namespaces().List(ctx, opts)
}
