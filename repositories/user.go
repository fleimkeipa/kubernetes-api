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

func (rc *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user = new(model.User)

	err := rc.db.
		Model(user).
		Where("username = ?", username).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}
