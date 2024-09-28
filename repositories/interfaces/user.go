package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type UserInterfaces interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, user model.User) (*model.User, error)
	List(ctx context.Context, opts *model.UserFindOpts) ([]model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	Delete(ctx context.Context, userID string) error
}
