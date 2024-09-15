package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (rc *UserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	_, err := rc.db.Model(&user).Insert()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (rc *UserRepository) GetUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user = new(model.User)

	err := rc.db.
		Model(user).
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}
