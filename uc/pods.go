package uc

import (
	"context"

	"kub/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
