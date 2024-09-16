package repositories

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (rc *UserRepository) Update(ctx context.Context, user model.User) (*model.User, error) {
	_, err := rc.db.Model(&user).WherePK().Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (rc *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user = new(model.User)

	err := rc.db.
		Model(user).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to find user [%s] id, error: %w", id, err)
	}

	return user, nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user = new(model.User)

	err := rc.db.
		Model(user).
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to find user by [%s], error: %w", usernameOrEmail, err)
	}

	return user, nil
}
