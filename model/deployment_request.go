package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentCreateRequest struct {
	Opts       metav1.CreateOptions
	Deployment Deployment
}

type DeploymentUpdateRequest struct {
	Opts       metav1.UpdateOptions
	Deployment Deployment
}
