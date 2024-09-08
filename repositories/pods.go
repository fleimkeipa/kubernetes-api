package repositories

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodsRepository struct {
	client *kubernetes.Clientset
}

func NewPodsRepository(client *kubernetes.Clientset) *PodsRepository {
	return &PodsRepository{client}
}

func (rc *PodsRepository) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, pod, opts)
}

func (rc *PodsRepository) Get(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	return rc.client.CoreV1().Pods(namespace).List(ctx, opts)
}

func (rc *PodsRepository) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return rc.client.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (rc *PodsRepository) Update(ctx context.Context, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Update(ctx, pod, opts)
}
