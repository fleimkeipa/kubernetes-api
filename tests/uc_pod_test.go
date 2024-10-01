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

	corev1 "k8s.io/api/core/v1"
)

func TestPodsUC_GetByName(t *testing.T) {
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	type fields struct {
		podsRepo interfaces.PodInterfaces
		eventUC  *uc.EventUC
	}
	type args struct {
		ctx       context.Context
		namespace string
		name      string
		opts      model.ListOptions
	}
	tests := []struct {
		fields  fields
		want    *corev1.Pod
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				podsRepo: repositories.NewPodRepository(initTestKubernetes()),
				eventUC:  uc.NewEventUC(repositories.NewEventRepository(test_db)),
			},
			args: args{
				ctx:       context.Background(),
				namespace: "default",
				name:      "testpod",
				opts:      model.ListOptions{},
			},
			want:    &corev1.Pod{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := uc.NewPodUC(tt.fields.podsRepo, tt.fields.eventUC)
			got, err := rc.GetByNameOrUID(tt.args.ctx, tt.args.namespace, tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodsUC.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodsUC.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
