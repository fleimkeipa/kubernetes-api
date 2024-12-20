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
	q := rc.db.Model(&newUser)

	_, err := q.Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (rc *UserRepository) Update(ctx context.Context, updatedUser model.User) (*model.User, error) {
	q := rc.db.Model(&updatedUser).WherePK()

	result, err := q.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("no user updated")
	}

	return &updatedUser, nil
}

func (rc *UserRepository) List(ctx context.Context, opts *model.UserFindOpts) (*model.UserList, error) {
	var users []model.User

	filter := rc.fillFilter(opts)
	fields := rc.fillFields(opts)

	q := rc.db.Model(&users).Column(fields...)

	if filter != "" {
		q = q.Where(filter)
	}

	q = q.Limit(opts.Limit).Offset(opts.Skip)

	count, err := q.SelectAndCount()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return &model.UserList{
		Users: users,
		Total: count,
		PaginationOpts: model.PaginationOpts{
			Skip:  opts.Skip,
			Limit: opts.Limit,
		},
	}, nil
}

func (rc *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	q := rc.db.Model(&user)

	if id == "0" || id == "" {
		return nil, fmt.Errorf("invalid user id")
	}

	q = q.Where("id = ?", id)

	if err := q.Select(); err != nil {
		return nil, fmt.Errorf("failed to find user by id [%s]: %w", id, err)
	}

	return &user, nil
}

func (rc *UserRepository) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user model.User

	q := rc.db.Model(&user)

	if usernameOrEmail == "" {
		return nil, fmt.Errorf("invalid username or email")
	}

	q = q.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	if err := q.Select(); err != nil {
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
		filter = addFilterClause(filter, "username", opts.Username.Value)
	}

	if opts.Email.IsSended {
		filter = addFilterClause(filter, "email", opts.Email.Value)
	}

	if opts.RoleID.IsSended {
		filter = addFilterClause(filter, "role_id", opts.RoleID.Value)
	}

	return filter
}
