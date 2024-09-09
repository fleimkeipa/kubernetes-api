package model

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentRequest struct {
	Deployment v1.Deployment
	Opts       metav1.CreateOptions
}
