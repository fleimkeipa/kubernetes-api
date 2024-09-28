package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInterfaces interface {
	Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error)
	Update(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
	List(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error)
	Delete(ctx context.Context, namespace string, podID string, opts metav1.DeleteOptions) error
}
