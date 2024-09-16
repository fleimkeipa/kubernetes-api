package model

import "time"

type User struct {
	ID        int64     `json:"id" pg:",pk"`
	Username  string    `json:"username" binding:"required" pg:",unique"`
	Email     string    `json:"email" binding:"required" pg:",unique"`
	Password  string    `json:"password" binding:"required"`
	RoleID    uint      `json:"role_id"`
	DeletedAt time.Time `json:"deleted_at" pg:",soft_delete"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `bson:"role_id" json:"role_id"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}
