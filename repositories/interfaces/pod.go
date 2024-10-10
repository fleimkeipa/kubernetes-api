package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type PodInterfaces interface {
	Create(ctx context.Context, pod *model.Pod, opts model.CreateOptions) (*model.Pod, error)
	Update(ctx context.Context, id string, pod *model.Pod, opts model.UpdateOptions) (*model.Pod, error)
	List(ctx context.Context, namespace string, opts model.ListOptions) (*model.PodList, error)
	Delete(ctx context.Context, namespace string, podID string, opts model.DeleteOptions) error
	GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Pod, error)
}
