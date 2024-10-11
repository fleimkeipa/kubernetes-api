package model

import "time"

const ZeroCreds = "zeroCreds"

type User struct {
	DeletedAt time.Time `json:"deleted_at" pg:",soft_delete"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username" binding:"required" pg:",unique"`
	Email     string    `json:"email" binding:"required" pg:",unique"`
	Password  string    `json:"password" binding:"required"`
	ID        int64     `json:"id" pg:",pk"`
	RoleID    uint      `json:"role_id"`
}

type UserList struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	PaginationOpts
}

type Owner struct {
	Username string `json:"username" pg:",unique"`
	Email    string `json:"email" pg:",unique"`
	ID       int64  `json:"id" pg:",pk"`
	RoleID   uint   `json:"role_id"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GithubUser struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	ID            int    `json:"id"`
	VerifiedEmail bool   `json:"verified_email"`
}

type UserFindOpts struct {
	Username Filter
	Email    Filter
	RoleID   Filter
	FieldsOpts
	PaginationOpts
}
