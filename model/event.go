package model

import "time"

const (
	PodKind        = "pod"
	DeploymentKind = "deployment"
	NamespaceKind  = "namespace"
)

const (
	CreateEventKind = "create"
	UpdateEventKind = "update"
	DeleteEventKind = "delete"
)

type Event struct {
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at" pg:",soft_delete"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Owner     Owner     `json:"owner" pg:"rel:has-one"`
}

type EventFindOpts struct {
	Kind          Filter
	EventKind     Filter
	CreationTime  Filter
	OwnerID       Filter
	OwnerUsername Filter
	FieldsOpts
	PaginationOpts
}
