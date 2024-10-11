package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type NamespaceUC struct {
	namespaceRepo interfaces.NamespaceInterfaces
	eventUC       *EventUC
}

func NewNamespaceUC(namespaceRepo interfaces.NamespaceInterfaces, eventUC *EventUC) *NamespaceUC {
	return &NamespaceUC{
		namespaceRepo: namespaceRepo,
		eventUC:       eventUC,
	}
}

func (rc *NamespaceUC) Create(ctx context.Context, request model.NamespaceCreateRequest) (*model.Namespace, error) {
	request.Namespace.TypeMeta.Kind = "namespace"

	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.CreateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	return rc.namespaceRepo.Create(ctx, &request.Namespace, request.Opts)
}

func (rc *NamespaceUC) Update(ctx context.Context, nameOrUID string, request *model.NamespaceUpdateRequest) (*model.Namespace, error) {
	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.UpdateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeNamespace := rc.fillNamespace(request)

	return rc.namespaceRepo.Update(ctx, nameOrUID, kubeNamespace, request.Opts)
}

func (rc *NamespaceUC) List(ctx context.Context, opts model.ListOptions) (*model.NamespaceList, error) {
	opts.TypeMeta.Kind = "namespace"

	return rc.namespaceRepo.List(ctx, opts)
}

func (rc *NamespaceUC) GetByNameOrUID(ctx context.Context, nameOrUID string, opts model.ListOptions) (*model.Namespace, error) {
	return rc.namespaceRepo.GetByNameOrUID(ctx, nameOrUID, opts)
}

func (rc *NamespaceUC) Delete(ctx context.Context, name string, opts model.DeleteOptions) error {
	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.DeleteEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return err
	}

	return rc.namespaceRepo.Delete(ctx, name, opts)
}

func (rc *NamespaceUC) fillNamespace(request *model.NamespaceUpdateRequest) *model.Namespace {
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Labels:      request.Namespace.Labels,
			Annotations: request.Namespace.Annotations,
		},
		Spec: request.Namespace.Spec,
	}
}
