package repositories

import (
	"context"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestNamespaceRepository_Get(t *testing.T) {
	type fields struct {
		client *kubernetes.Clientset
	}
	type args struct {
		ctx  context.Context
		opts metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *corev1.NamespaceList
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				client: initTestKubernetes(),
			},
			args: args{
				ctx:  context.Background(),
				opts: metav1.ListOptions{},
			},
			want:    &corev1.NamespaceList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &NamespaceRepository{
				client: tt.fields.client,
			}
			got, err := rc.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
