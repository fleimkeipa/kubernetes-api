package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type NamespaceInterfaces interface {
	Create(ctx context.Context, namespace *model.Namespace, opts model.CreateOptions) (*model.Namespace, error)
	Update(ctx context.Context, namespace *model.Namespace, opts model.UpdateOptions) (*model.Namespace, error)
	List(ctx context.Context, opts model.ListOptions) (*model.NamespaceList, error)
	Delete(ctx context.Context, namespaceID string, opts model.DeleteOptions) error
	GetByNameOrUID(ctx context.Context, nameOrUID string, opts model.ListOptions) (*model.Namespace, error)
}
