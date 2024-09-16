package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type EventInterfaces interface {
	Create(context.Context, *model.Event) (*model.Event, error)
	List(context.Context, *model.EventFindOpts) ([]model.Event, error)
	// GetByID(context.Context, string) ([]model.Event, error)
}
