package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
	"github.com/fleimkeipa/kubernetes-api/uc"
)

func TestEventUC_List(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	type fields struct {
		eventRepo interfaces.EventInterfaces
	}
	type tempData struct {
		ticket *model.Event
	}
	type args struct {
		ctx  context.Context
		opts *model.EventFindOpts
	}
	tests := []struct {
		name      string
		fields    fields
		tempDatas []tempData
		args      args
		want      []model.Event
		wantErr   bool
	}{
		{
			name: "success",
			fields: fields{
				eventRepo: repositories.NewEventRepository(test_db),
			},
			tempDatas: []tempData{
				{
					ticket: &model.Event{
						Kind:      "pod",
						EventKind: "create",
						Owner: model.User{
							ID:       1,
							Username: "test_username",
							Email:    "test@mail.com",
							RoleID:   1,
						},
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				opts: &model.EventFindOpts{},
			},
			want: []model.Event{
				{
					Kind:      "pod",
					EventKind: "create",
					Owner: model.User{
						ID:       1,
						Username: "test_username",
						Email:    "test@mail.com",
						RoleID:   1,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.tempDatas {
				if err := addTempData(v.ticket); (err != nil) != tt.wantErr {
					t.Errorf("EventUC.List() addTempData error = %v", err)
					return
				}
			}
			rc := uc.NewEventUC(tt.fields.eventRepo)
			got, err := rc.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventUC.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventUC.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
