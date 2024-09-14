package model

type User struct {
	ID       int64
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `bson:"role_id" json:"role_id"`
}
