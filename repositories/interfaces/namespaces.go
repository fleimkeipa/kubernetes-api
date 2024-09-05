package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceInterfaces interface {
	Get(context.Context, metav1.ListOptions) (*corev1.NamespaceList, error)
}
