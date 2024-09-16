package repositories

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/pkg"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func initTestKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}

func TestPodsRepository_Get(t *testing.T) {
	type fields struct {
		client *kubernetes.Clientset
	}
	type args struct {
		ctx       context.Context
		namespace string
		opts      metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *corev1.PodList
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				client: initTestKubernetes(),
			},
			args: args{
				ctx:       context.Background(),
				namespace: "default",
				opts: metav1.ListOptions{
					TypeMeta: metav1.TypeMeta{
						Kind:       "pod",
						APIVersion: "v1",
					},
					Limit: 2,
				},
			},
			want:    &corev1.PodList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &PodRepository{
				client: tt.fields.client,
			}
			got, err := rc.List(tt.args.ctx, tt.args.namespace, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodsRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodsRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
