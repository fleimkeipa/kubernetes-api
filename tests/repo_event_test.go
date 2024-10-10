package tests

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

func TestEventRepository_Create(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	now := time.Now()

	type fields struct {
		db *pg.DB
	}
	type args struct {
		ctx   context.Context
		event *model.Event
	}
	tests := []struct {
		args    args
		fields  fields
		want    *model.Event
		name    string
		wantErr bool
	}{
		{
			name: "success - pod create",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.Background(),
				event: &model.Event{
					Category:  "pod",
					Type:      "create",
					CreatedAt: now,
					Owner: model.Owner{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			want: &model.Event{
				Category:  "pod",
				Type:      "create",
				CreatedAt: now,
				Owner: model.Owner{
					ID:       1,
					Username: "test_username",
					Email:    "test@mail.com",
					RoleID:   1,
				},
				DeletedAt: time.Time{},
			},
			wantErr: false,
		},
		{
			name: "success - namespace create",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.Background(),
				event: &model.Event{
					Category:  "namespace",
					Type:      "create",
					CreatedAt: now,
					Owner: model.Owner{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			want: &model.Event{
				Category:  "namespace",
				Type:      "create",
				CreatedAt: now,
				Owner: model.Owner{
					ID:       1,
					Username: "test_username",
					Email:    "test@mail.com",
					RoleID:   1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := repositories.NewEventRepository(tt.fields.db)
			got, err := rc.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.Create() = %v, want %v", got, tt.want)
			}
			if err := clearTable("events"); err != nil {
				t.Errorf("EventRepository.Create() clearTable error = %v", err)
				return
			}
		})
	}
}

func TestEventRepository_List(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

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
			name: "success",
			fields: fields{
				db: test_db,
			},
			args: args{
				ctx: context.TODO(),
				opts: &model.EventFindOpts{
					PaginationOpts: model.PaginationOpts{},
					Category: model.Filter{
						IsSended: true,
						Value:    "namespace",
					},
					Type: model.Filter{
						IsSended: true,
						Value:    "update",
					},
					CreatedAt:     model.Filter{},
					OwnerID:       model.Filter{},
					OwnerUsername: model.Filter{},
				},
				events: []model.Event{
					{
						Category:  "namespace",
						Type:      "create",
						CreatedAt: time.Now(),
						Owner: model.Owner{
							ID:       1,
							Username: "test_username",
							Email:    "test@mail.com",
							RoleID:   1,
						},
					},
					{
						Category:  "pod",
						Type:      "update",
						CreatedAt: time.Now(),
						Owner: model.Owner{
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
			rc := repositories.NewEventRepository(tt.fields.db)
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
