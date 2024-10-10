package model

import "time"

const (
	UserCategory       = "user"
	PodCategory        = "pod"
	DeploymentCategory = "deployment"
	NamespaceCategory  = "namespace"
)

const (
	CreateEventType = "create"
	UpdateEventType = "update"
	DeleteEventType = "delete"
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
