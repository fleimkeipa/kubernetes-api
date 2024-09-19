package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceCreateRequest struct {
	Namespace Namespace            `json:"namespace"`
	Opts      metav1.CreateOptions `json:"opts"`
}

type NamespaceUpdateRequest struct {
	Namespace Namespace            `json:"namespace"`
	Opts      metav1.UpdateOptions `json:"opts"`
}
