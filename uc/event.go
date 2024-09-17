package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type EventUC struct {
	eventRepo interfaces.EventInterfaces
}

func NewEventUC(eventRepo interfaces.EventInterfaces) *EventUC {
	return &EventUC{
		eventRepo: eventRepo,
	}
}

func (rc *EventUC) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	return rc.eventRepo.Create(ctx, event)
}

func (rc *EventUC) List(ctx context.Context, opts *model.EventFindOpts) ([]model.Event, error) {
	return rc.eventRepo.List(ctx, opts)
}

func (rc *EventUC) GetByID(ctx context.Context, id string) (*model.Event, error) {
	return rc.eventRepo.GetByID(ctx, id)
}
