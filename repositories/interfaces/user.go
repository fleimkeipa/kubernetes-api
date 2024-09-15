package interfaces

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
)

type UserInterfaces interface {
	Create(context.Context, model.User) (*model.User, error)
	GetUserByUsernameOrEmail(context.Context, string) (*model.User, error)
}
