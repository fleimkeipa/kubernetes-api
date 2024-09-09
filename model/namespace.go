package model

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceRequest struct {
	Namespace corev1.Namespace
	Opts      metav1.CreateOptions
}
