package uc

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	loadEnv()
}

func TestPodsUC_GetByName(t *testing.T) {
	type fields struct {
		podsRepo interfaces.PodInterfaces
	}
	type args struct {
		ctx       context.Context
		namespace string
		name      string
		opts      metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *corev1.Pod
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				podsRepo: repositories.NewPodRepository(initTestKubernetes()),
			},
			args: args{
				ctx:       context.Background(),
				namespace: "default",
				name:      "testpod",
				opts:      metav1.ListOptions{},
			},
			want:    &corev1.Pod{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &PodUC{
				podsRepo: tt.fields.podsRepo,
			}
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
