package uc

import (
	"context"
	"fmt"

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

func (rc *UserUC) Update(ctx context.Context, id string, user model.User) (*model.User, error) {
	// user exist control
	existUser, err := rc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.ID = existUser.ID

	return rc.userRepo.Update(ctx, user)
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	user, err := rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by [%s], error: %w", usernameOrEmail, err)
	}

	return user, nil
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := rc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user [%s] id, error: %w", id, err)
	}

	return user, nil
}
