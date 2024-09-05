package uc

import (
	"context"

	"kub/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
