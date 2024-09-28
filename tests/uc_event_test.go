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
	type args struct {
		ctx  context.Context
		opts *model.EventFindOpts
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
				eventRepo: repositories.NewEventRepository(test_db),
			},
			args: args{
				ctx:  context.TODO(),
				opts: &model.EventFindOpts{},
			},
			want:    []model.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
