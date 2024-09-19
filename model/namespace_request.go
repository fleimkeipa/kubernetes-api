package model

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceCreateRequest struct {
	Namespace corev1.Namespace     `json:"namespace"`
	Opts      metav1.CreateOptions `json:"opts"`
}

type NamespaceUpdateRequest struct {
	Namespace corev1.Namespace     `json:"namespace"`
	Opts      metav1.UpdateOptions `json:"opts"`
}
