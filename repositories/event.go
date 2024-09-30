package repositories

import (
	"context"
	"fmt"

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

func (rc *EventRepository) List(ctx context.Context, opts *model.EventFindOpts) ([]model.Event, error) {
	events := make([]model.Event, 0)
	filter := rc.fillFilter(opts)
	if filter == "" {
		err := rc.db.
			Model(&events).
			Select()
		if err != nil {
			return nil, err
		}
	} else {
		err := rc.db.
			Model(&events).
			Where(filter).
			Select()
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func (rc *EventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
	event := new(model.Event)

	fields := []string{}
	err := rc.db.
		Model(event).
		Where("id = ?", id).
		Column(fields...).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to find user [%s] id, error: %w", id, err)
	}

	return event, nil
}

func (rc *EventRepository) fillFilter(opts *model.EventFindOpts) string {
	var filter string
	if opts.Kind.IsSended {
		filter = addInFilter(filter, "kind", opts.Kind.Value)
	}
	if opts.EventKind.IsSended {
		filter = addInFilter(filter, "event_kind", opts.EventKind.Value)
	}
	if opts.CreationTime.IsSended {
		filter = addInFilter(filter, "creation_time", opts.CreationTime.Value)
	}
	if opts.OwnerID.IsSended {
		filter = addInFilter(filter, "owner.id", opts.OwnerID.Value)
	}
	if opts.OwnerUsername.IsSended {
		filter = addInFilter(filter, "owner.username", opts.OwnerUsername.Value)
	}

	return filter
}

func addInFilter(filter, key, value string) string {
	if filter == "" {
		filter = fmt.Sprintf("%s = %s", key, fmt.Sprintf(`'%s'`, value))
	} else {
		filter = fmt.Sprintf("%s AND %s = %s", filter, key, fmt.Sprintf(`'%s'`, value))
	}

	return filter
}
