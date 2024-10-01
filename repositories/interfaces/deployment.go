package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	v1 "k8s.io/api/apps/v1"
)

type DeploymentInterfaces interface {
	Create(ctx context.Context, deployment *v1.Deployment, opts model.CreateOptions) (*v1.Deployment, error)
	Update(ctx context.Context, deployment *v1.Deployment, opts model.UpdateOptions) (*v1.Deployment, error)
	List(ctx context.Context, namespace string, opts model.ListOptions) (*v1.DeploymentList, error)
	Delete(ctx context.Context, namespace string, deploymentID string, opts model.DeleteOptions) error
}
