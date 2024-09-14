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
			name: "",
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
			want:    &model.User{},
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
			got, err := rc.GetUserByUsername(tt.args.ctx, tt.args.username)
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
