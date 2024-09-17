package repositories

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/go-pg/pg"
)

func TestUserRepository_GetUserByUsername(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx      context.Context
		username string
		user     *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "username search",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin",
				user: &model.User{
					ID:       0,
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
			name: "email search",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin2@admin.com",
				user: &model.User{
					ID:       0,
					Username: "admin2",
					Email:    "admin2@admin.com",
					Password: "password",
					RoleID:   7,
				},
			},
			want: &model.User{
				ID:       2,
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
				t.Errorf("UserRepository.GetUserByUsername() addTempData error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rc := &UserRepository{
				db: tt.fields.db,
			}
			got, err := rc.GetByUsernameOrEmail(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx  context.Context
		id   string
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "",
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
			rc := &UserRepository{
				db: tt.fields.db,
			}
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
		id   string
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
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
			rc := &UserRepository{
				db: tt.fields.db,
			}
			if err := rc.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
