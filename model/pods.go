package model

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsRequest struct {
	Pod  corev1.Pod           `json:"pod"`
	Opts metav1.CreateOptions `json:"opts"`
}
