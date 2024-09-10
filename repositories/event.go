package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/go-pg/pg"
)

type EventRepository struct {
	db *pg.DB
}

func NewEventRepository(db *pg.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (rc *EventRepository) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	_, err := rc.db.Model(event).Insert()
	if err != nil {
		return nil, err
	}

	return event, nil
}
