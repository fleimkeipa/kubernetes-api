package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"

	"github.com/go-pg/pg"
)

func TestUserRepository_GetUserByUsername(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx      context.Context
		user     *model.User
		username string
	}
	tests := []struct {
		args    args
		fields  fields
		want    *model.User
		name    string
		wantErr bool
	}{
		{
			name: "success - username search",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin",
				user: &model.User{
					ID:       1,
					Username: "admin",
					Email:    "admin@admin.com",
					Password: "password",
					RoleID:   7,
				},
			},
			want: &model.User{
				ID:       1,
				Username: "admin",
				Email:    "admin@admin.com",
				Password: "password",
				RoleID:   7,
			},
			wantErr: false,
		},
		{
			name: "success - email search",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin2@admin.com",
				user: &model.User{
					ID:       1,
					Username: "admin2",
					Email:    "admin2@admin.com",
					Password: "password",
					RoleID:   7,
				},
			},
			want: &model.User{
				ID:       1,
				Username: "admin2",
				Email:    "admin2@admin.com",
				Password: "password",
				RoleID:   7,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addTempData(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetUserByUsername() addTempData error = %v", err)
				return
			}
			rc := repositories.NewUserRepository(tt.fields.db)
			got, err := rc.GetByUsernameOrEmail(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetUserByUsername() = %v, want %v", got, tt.want)
			}
			if err := clearTable("users"); err != nil {
				t.Errorf("UserRepository.GetUserByUsername() clearTable error = %v", err)
				return
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx  context.Context
		user *model.User
		id   string
	}
	tests := []struct {
		args    args
		fields  fields
		want    *model.User
		name    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				user: &model.User{
					ID:       1,
					Username: "admin3",
					Email:    "admin3@admin.com",
					Password: "password",
					RoleID:   7,
				},
			},
			want: &model.User{
				ID:       1,
				Username: "admin3",
				Email:    "admin3@admin.com",
				Password: "password",
				RoleID:   7,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addTempData(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetByID() addTempData error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rc := repositories.NewUserRepository(tt.fields.db)
			got, err := rc.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx  context.Context
		user *model.User
		id   string
	}
	tests := []struct {
		args    args
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.TODO(),
				id:  "1",
				user: &model.User{
					ID:       0,
					Username: "admin3",
					Email:    "admin3@admin.com",
					Password: "password",
					RoleID:   7,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := addTempData(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Delete() addTempData error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rc := repositories.NewUserRepository(tt.fields.db)
			if err := rc.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
