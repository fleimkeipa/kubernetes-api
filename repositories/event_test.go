package repositories

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

func TestEventRepository_Create(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx   context.Context
		event *model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.Background(),
				event: &model.Event{
					Kind:         "pod",
					EventKind:    "create",
					CreationTime: time.Now(),
					Owner: model.User{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			want:    &model.Event{},
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.Background(),
				event: &model.Event{
					Kind:         "namespace",
					EventKind:    "create",
					CreationTime: time.Now(),
					Owner: model.User{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			want:    &model.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &EventRepository{
				db: tt.fields.db,
			}
			got, err := rc.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventRepository_List(t *testing.T) {
	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx    context.Context
		opts   *model.EventFindOpts
		events []model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Event
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.TODO(),
				opts: &model.EventFindOpts{
					PaginationOpts: model.PaginationOpts{},
					Kind: model.Filter{
						IsSended: true,
						Value:    "namespace",
					},
					EventKind: model.Filter{
						IsSended: true,
						Value:    "update",
					},
					CreationTime:  model.Filter{},
					OwnerID:       model.Filter{},
					OwnerUsername: model.Filter{},
				},
				events: []model.Event{
					{
						Kind:         "namespace",
						EventKind:    "create",
						CreationTime: time.Now(),
						Owner: model.User{
							ID:       1,
							Username: "test_username",
							Email:    "test@mail.com",
							RoleID:   1,
						},
					},
					{
						Kind:         "pod",
						EventKind:    "update",
						CreationTime: time.Now(),
						Owner: model.User{
							ID:       1,
							Username: "test_username",
							Email:    "test@mail.com",
							RoleID:   1,
						},
					},
				},
			},
			want:    []model.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.events {
				if err := addTempData(&v); (err != nil) != tt.wantErr {
					t.Errorf("EventRepository.List() addTempData error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			rc := &EventRepository{
				db: tt.fields.db,
			}
			got, err := rc.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
