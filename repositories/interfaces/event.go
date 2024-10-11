package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type EventInterfaces interface {
	Create(ctx context.Context, event *model.Event) (*model.Event, error)
	List(ctx context.Context, event *model.EventFindOpts) (*model.EventList, error)
	GetByID(ctx context.Context, eventID string) (*model.Event, error)
}
