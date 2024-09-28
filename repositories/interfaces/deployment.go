package interfaces

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentInterfaces interface {
	Create(ctx context.Context, deployment *v1.Deployment, opts metav1.CreateOptions) (*v1.Deployment, error)
	Update(ctx context.Context, deployment *v1.Deployment, opts metav1.UpdateOptions) (*v1.Deployment, error)
	List(ctx context.Context, namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error)
	Delete(ctx context.Context, namespace string, deploymentID string, opts metav1.DeleteOptions) error
}
