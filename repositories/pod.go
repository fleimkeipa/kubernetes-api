package repositories

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodRepository struct {
	client *kubernetes.Clientset
}

func NewPodRepository(client *kubernetes.Clientset) *PodRepository {
	return &PodRepository{client}
}

func (rc *PodRepository) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, pod, opts)
}

func (rc *PodRepository) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	return rc.client.CoreV1().Pods(namespace).List(ctx, opts)
}

func (rc *PodRepository) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return rc.client.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (rc *PodRepository) Update(ctx context.Context, id string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Update(ctx, pod, opts)
}
