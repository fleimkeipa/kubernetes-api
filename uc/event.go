package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type EventUC struct {
	eventRepo interfaces.EventsInterfaces
}

func NewEventUC(eventRepo interfaces.EventsInterfaces) *EventUC {
	return &EventUC{
		eventRepo: eventRepo,
	}
}

func (rc *EventUC) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	return rc.eventRepo.Create(ctx, event)
}
