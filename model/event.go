package model

import "time"

const (
	CreateEventKind = "create"
	UpdateEventKind = "update"
	DeleteEventKind = "delete"
)

type Event struct {
	Kind         string    `json:"kind"`
	EventKind    string    `json:"event_kind"`
	CreationTime time.Time `json:"creation_time"`
	Owner        User      `json:"owner" pg:"rel:has-one"`
	DeletedAt    time.Time `json:"deleted_at" pg:",soft_delete"`
}

type EventFindOpts struct {
	PaginationOpts
	Kind          Filter
	EventKind     Filter
	CreationTime  Filter
	OwnerID       Filter
	OwnerUsername Filter
}
