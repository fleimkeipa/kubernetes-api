package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type UserUC struct {
	userRepo interfaces.UserInterfaces
}

func NewUserUC(repo interfaces.UserInterfaces) *UserUC {
	return &UserUC{
		userRepo: repo,
	}
}

func (rc *UserUC) Create(ctx context.Context, user model.User) (*model.User, error) {
	hashedPassword, err := model.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return rc.userRepo.Create(ctx, user)
}

func (rc *UserUC) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return rc.userRepo.GetUserByUsername(ctx, username)
}

func (rc *UserUC) deleteCreds(users []model.User) []model.User {
	for i := range users {
		users[i].Password = ""
	}
	return users
}