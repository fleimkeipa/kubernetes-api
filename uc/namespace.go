package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type NamespaceUC struct {
	namespaceRepo interfaces.NamespaceInterfaces
}

func NewNamespaceUC(namespaceRepo interfaces.NamespaceInterfaces) *NamespaceUC {
	return &NamespaceUC{
		namespaceRepo: namespaceRepo,
	}
}

func (rc *NamespaceUC) Get(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error) {
	return rc.namespaceRepo.Get(ctx, opts)
}

func (rc *NamespaceUC) Create(ctx context.Context, namespace *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return rc.namespaceRepo.Create(ctx, namespace, opts)
}

func (rc *NamespaceUC) GetByNameOrUID(ctx context.Context, nameOrUID string, opts metav1.ListOptions) (*corev1.Namespace, error) {
	opts.TypeMeta.Kind = "namespace"

	opts.Limit = 100
	namespaces, err := rc.namespaceRepo.Get(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range namespaces.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if namespaces.ListMeta.Continue == "" {
		return &corev1.Namespace{}, nil
	}

	opts.Continue = namespaces.ListMeta.Continue
	return rc.GetByNameOrUID(ctx, nameOrUID, opts)
}

func (rc *NamespaceUC) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	opts.TypeMeta.Kind = "namespace"

	return rc.namespaceRepo.Delete(ctx, name, opts)
}

func (rc *NamespaceUC) Update(ctx context.Context, namespace *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	namespace.TypeMeta.Kind = "namespace"

	return rc.namespaceRepo.Update(ctx, namespace, opts)
}
