package interfaces

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentInterfaces interface {
	Get(context.Context, string, metav1.ListOptions) (*v1.DeploymentList, error)
	Create(context.Context, *v1.Deployment, metav1.CreateOptions) (*v1.Deployment, error)
	Delete(context.Context, string, string, metav1.DeleteOptions) error
	Update(context.Context, *v1.Deployment, metav1.UpdateOptions) (*v1.Deployment, error)
}
