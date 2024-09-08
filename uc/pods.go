package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PodsUC struct {
	podsRepo interfaces.PodsInterfaces
}

func NewPodsUC(podsRepo interfaces.PodsInterfaces) *PodsUC {
	return &PodsUC{
		podsRepo: podsRepo,
	}
}

func (rc *PodsUC) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	pod.TypeMeta.Kind = "pod"
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = "default"
	}

	return rc.podsRepo.Create(ctx, pod, opts)
}

func (rc *PodsUC) Update(ctx context.Context, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	pod.TypeMeta.Kind = "pod"
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = "default"
	}

	return rc.podsRepo.Update(ctx, pod, opts)
}

func (rc *PodsUC) Get(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	return rc.podsRepo.Get(ctx, namespace, opts)
}

func (rc *PodsUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts metav1.ListOptions) (*corev1.Pod, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	pods, err := rc.podsRepo.Get(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range pods.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if pods.ListMeta.Continue == "" {
		return &corev1.Pod{}, nil
	}

	opts.Continue = pods.ListMeta.Continue
	return rc.GetByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *PodsUC) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	return rc.podsRepo.Delete(ctx, namespace, name, opts)
}
