package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type UserInterfaces interface {
	Create(context.Context, model.User) (*model.User, error)
	Update(context.Context, model.User) (*model.User, error)
	List(context.Context, *model.UserFindOpts) ([]model.User, error)
	GetByUsernameOrEmail(context.Context, string) (*model.User, error)
	GetByID(context.Context, string) (*model.User, error)
}
