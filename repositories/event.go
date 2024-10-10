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

func (rc *EventRepository) Create(ctx context.Context, newEvent *model.Event) (*model.Event, error) {
	q := rc.db.Model(newEvent)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create event (kind: [%s], event kind: [%s]): %v", newEvent.Category, newEvent.Type, err)
	}

	return newEvent, nil
}

func (rc *EventRepository) List(ctx context.Context, opts *model.EventFindOpts) ([]model.Event, error) {
	var events []model.Event

	filter := rc.fillFilter(opts)
	fields := rc.fillFields(opts)

	q := rc.db.Model(&events).Column(fields...)

	if filter != "" {
		q = q.Where(filter)
	}

	q = q.Limit(opts.Limit).Offset(opts.Skip)

	err := q.Select()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}

func (rc *EventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event

	q := rc.db.Model(event)

	if id == "0" || id == "" {
		return nil, fmt.Errorf("invalid event id")
	}

	q = q.Where("id = ?", id)

	if err := q.Select(); err != nil {
		return nil, fmt.Errorf("failed to find event [%s] id: %w", id, err)
	}

	return &event, nil
}

func (rc *EventRepository) fillFields(opts *model.EventFindOpts) []string {
	fields := opts.Fields

	if len(fields) == 0 {
		return nil
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return []string{
			"kind",
			"event_kind",
			"creation_time",
			"owner.id",
			"deleted_at",
		}
	}

	return fields
}

func (rc *EventRepository) fillFilter(opts *model.EventFindOpts) string {
	filter := ""

	if opts.Category.IsSended {
		filter = addFilterClause(filter, "category", opts.Category.Value)
	}

	if opts.Type.IsSended {
		filter = addFilterClause(filter, "type", opts.Type.Value)
	}

	if opts.CreatedAt.IsSended {
		filter = addFilterClause(filter, "created_at", opts.CreatedAt.Value)
	}

	if opts.OwnerID.IsSended {
		filter = addFilterClause(filter, "owner.id", opts.OwnerID.Value)
	}

	if opts.OwnerUsername.IsSended {
		filter = addFilterClause(filter, "owner.username", opts.OwnerUsername.Value)
	}

	return filter
}

func addFilterClause(filter string, key string, value string) string {
	if filter == "" {
		return fmt.Sprintf("%s = '%s'", key, value)
	}

	return fmt.Sprintf("%s AND %s = '%s'", filter, key, value)
}
