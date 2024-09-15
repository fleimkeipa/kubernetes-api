package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsInterfaces interface {
	Create(context.Context, *model.Pod, metav1.CreateOptions) (*corev1.Pod, error)
	Get(context.Context, string, metav1.ListOptions) (*corev1.PodList, error)
	Delete(context.Context, string, string, metav1.DeleteOptions) error
	Update(context.Context, *model.Pod, metav1.UpdateOptions) (*corev1.Pod, error)
}
