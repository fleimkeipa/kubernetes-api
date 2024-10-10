package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
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

func (rc *NamespaceRepository) Create(ctx context.Context, namespace *corev1.Namespace, opts model.CreateOptions) (*corev1.Namespace, error) {
	metaOpts := convertCreateOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().Create(ctx, namespace, metaOpts)
}

func (rc *NamespaceRepository) Update(ctx context.Context, namespace *corev1.Namespace, opts model.UpdateOptions) (*corev1.Namespace, error) {
	metaOpts := convertUpdateOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().Update(ctx, namespace, metaOpts)
}

func (rc *NamespaceRepository) List(ctx context.Context, opts model.ListOptions) (*corev1.NamespaceList, error) {
	metaOpts := convertListOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().List(ctx, metaOpts)
}

func (rc *NamespaceRepository) Delete(ctx context.Context, name string, opts model.DeleteOptions) error {
	metaOpts := convertDeleteOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().Delete(ctx, name, metaOpts)
}
