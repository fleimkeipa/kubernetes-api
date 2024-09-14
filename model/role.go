package model

const (
	AdminRole  = 7
	ViewerRole = 5
	EditorRole = 1
)

// Role model
type Role struct {
	ID          string `bson:"_id,omitempty"` // MongoDB ObjectID, primary key
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}
