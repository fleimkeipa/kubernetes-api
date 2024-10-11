package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type DeploymentInterfaces interface {
	Create(ctx context.Context, deployment *model.Deployment, opts model.CreateOptions) (*model.Deployment, error)
	Update(ctx context.Context, namespace, id string, deployment *model.Deployment, opts model.UpdateOptions) (*model.Deployment, error)
	List(ctx context.Context, namespace string, opts model.ListOptions) (*model.DeploymentList, error)
	Delete(ctx context.Context, namespace string, deploymentID string, opts model.DeleteOptions) error
	GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Deployment, error)
}
