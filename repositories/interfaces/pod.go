package interfaces

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInterfaces interface {
	Create(context.Context, *corev1.Pod, metav1.CreateOptions) (*corev1.Pod, error)
	Update(context.Context, string, *corev1.Pod, metav1.UpdateOptions) (*corev1.Pod, error)
	List(context.Context, string, metav1.ListOptions) (*corev1.PodList, error)
	Delete(context.Context, string, string, metav1.DeleteOptions) error
}
