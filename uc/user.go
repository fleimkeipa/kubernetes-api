package uc

import (
	"context"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type UserUC struct {
	userRepo interfaces.UserInterfaces
	eventUC  *EventUC
}

func NewUserUC(repo interfaces.UserInterfaces, eventUC *EventUC) *UserUC {
	return &UserUC{
		userRepo: repo,
		eventUC:  eventUC,
	}
}

func (rc *UserUC) Create(ctx context.Context, user model.User) (*model.User, error) {
	hashedPassword, err := model.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	event := model.Event{
		Category: model.UserCategory,
		Type:     model.CreateEventType,
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = time.Now()

	return rc.userRepo.Create(ctx, user)
}

func (rc *UserUC) Update(ctx context.Context, id string, user model.User) (*model.User, error) {
	// user exist control
	existUser, err := rc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.ID = existUser.ID

	event := model.Event{
		Category: model.UserCategory,
		Type:     model.UpdateEventType,
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	return rc.userRepo.Update(ctx, user)
}

func (rc *UserUC) List(ctx context.Context, opts *model.UserFindOpts) ([]model.User, error) {
	return rc.userRepo.List(ctx, opts)
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := rc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	user, err := rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rc *UserUC) Delete(ctx context.Context, id string) error {
	event := model.Event{
		Category: model.UserCategory,
		Type:     model.DeleteEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return err
	}

	return rc.userRepo.Delete(ctx, id)
}
