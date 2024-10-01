package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	corev1 "k8s.io/api/core/v1"
)

type PodInterfaces interface {
	Create(ctx context.Context, pod *corev1.Pod, opts model.CreateOptions) (*corev1.Pod, error)
	Update(ctx context.Context, namespace string, pod *corev1.Pod, opts model.UpdateOptions) (*corev1.Pod, error)
	List(ctx context.Context, namespace string, opts model.ListOptions) (*corev1.PodList, error)
	Delete(ctx context.Context, namespace string, podID string, opts model.DeleteOptions) error
}
