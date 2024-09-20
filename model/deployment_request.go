package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentCreateRequest struct {
	Deployment Deployment
	Opts       metav1.CreateOptions
}

type DeploymentUpdateRequest struct {
	Deployment Deployment
	Opts       metav1.UpdateOptions
}
