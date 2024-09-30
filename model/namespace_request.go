package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceCreateRequest struct {
	Opts      metav1.CreateOptions `json:"opts"`
	Namespace Namespace            `json:"namespace"`
}

type NamespaceUpdateRequest struct {
	Opts      metav1.UpdateOptions `json:"opts"`
	Namespace Namespace            `json:"namespace"`
}
