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

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) ([]model.User, error) {
	var users = make([]model.User, 0)
	var filter = rc.fillFilter(opts)
	var fields = rc.fillFields(opts)
	if filter == "" {
		err := rc.db.
			Model(&users).
			Column(fields...).
			Limit(opts.Limit).
			Offset(opts.Skip).
			Select()
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	err := rc.db.
		Model(&users).
		Column(fields...).
		Where(filter).
		Limit(opts.Limit).
		Offset(opts.Skip).
		Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (rc *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user = new(model.User)

	var fields = []string{}
	err := rc.db.
		Model(user).
		Where("id = ?", id).
		Column(fields...).
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

func (rc *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := rc.db.Model(&model.User{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (rc *UserRepository) fillFields(opts *model.UserFindOpts) []string {
	if len(opts.Fields) == 0 {
		return []string{}
	}

	if len(opts.Fields) == 1 {
		if opts.Fields[0] == model.ZeroCreds {
			return []string{
				"id", "username", "email", "role_id", "deleted_at",
			}
		}
	}

	return opts.Fields
}

func (rc *UserRepository) fillFilter(opts *model.UserFindOpts) string {
	var filter string
	if opts.Username.IsSended {
		filter = addInFilter(filter, "username", opts.Username.Value)
	}
	if opts.Email.IsSended {
		filter = addInFilter(filter, "email", opts.Email.Value)
	}
	if opts.RoleID.IsSended {
		filter = addInFilter(filter, "role_id", opts.RoleID.Value)
	}

	return filter
}
