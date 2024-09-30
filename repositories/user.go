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

func (rc *UserRepository) Create(ctx context.Context, newUser model.User) (*model.User, error) {
	_, err := rc.db.Model(&newUser).Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (rc *UserRepository) Update(ctx context.Context, updatedUser model.User) (*model.User, error) {
	result, err := rc.db.Model(&updatedUser).WherePK().Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no user updated")
	}

	return &updatedUser, nil
}

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) ([]model.User, error) {
	var users []model.User

	filter := rc.fillFilter(opts)
	fields := rc.fillFields(opts)

	q := rc.db.Model(&users).Column(fields...)

	if filter != "" {
		q = q.Where(filter)
	}

	q = q.Limit(opts.Limit).Offset(opts.Skip)

	err := q.Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (rc *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := rc.db.
		Model(&user).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id [%s]: %w", id, err)
	}

	return &user, nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user model.User

	err := rc.db.
		Model(&user).
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get user by [%s]: %w", usernameOrEmail, err)
	}

	return &user, nil
}

func (rc *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := rc.db.Model(&model.User{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user deleted")
	}

	return nil
}

func (rc *UserRepository) fillFields(opts *model.UserFindOpts) []string {
	fields := opts.Fields

	if len(fields) == 0 {
		return nil
	}

	if len(fields) == 1 && fields[0] == model.ZeroCreds {
		return []string{
			"id",
			"username",
			"email",
			"role_id",
			"deleted_at",
		}
	}

	return fields
}

func (rc *UserRepository) fillFilter(opts *model.UserFindOpts) string {
	filter := ""

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
