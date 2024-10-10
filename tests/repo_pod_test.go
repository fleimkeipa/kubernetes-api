package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func TestPodRepository_Create(t *testing.T) {
	client := initTestKubernetes()
	defer deleteTestNamespace(client)

	type fields struct {
		client *kubernetes.Clientset
	}
	type args struct {
		ctx  context.Context
		pod  *model.Pod
		opts model.CreateOptions
	}
	tests := []struct {
		fields  fields
		want    *model.Pod
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: client,
			},
			args: args{
				ctx: context.TODO(),
				pod: &model.Pod{
					TypeMeta: model.TypeMeta{
						Kind:       "pod",
						APIVersion: "v1",
					},
					ObjectMeta: model.ObjectMeta{
						Name:      "pod1",
						Namespace: "test",
					},
					Spec: model.PodSpec{
						Containers: []model.Container{
							{
								Name:  "nginx",
								Image: "nginx:latest",
							},
						},
					},
					Status: model.PodStatus{},
				},
				opts: model.CreateOptions{},
			},
			want: &model.Pod{
				TypeMeta: model.TypeMeta{},
				ObjectMeta: model.ObjectMeta{
					Name: "pod1",
				},
				Spec:   model.PodSpec{},
				Status: model.PodStatus{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := repositories.NewPodRepository(tt.fields.client)
			got, err := rc.Create(tt.args.ctx, tt.args.pod, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ObjectMeta.Name != tt.args.pod.ObjectMeta.Name {
				t.Errorf("PodRepository.Create() = %v, want %v", got.ObjectMeta.Name, tt.want.ObjectMeta.Name)
			}
		})
	}
}

func TestPodsRepository_Get(t *testing.T) {
	type fields struct {
		client *kubernetes.Clientset
	}
	type args struct {
		ctx       context.Context
		namespace string
		opts      model.ListOptions
	}
	tests := []struct {
		fields  fields
		want    *corev1.PodList
		name    string
		args    args
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
				opts: model.ListOptions{
					TypeMeta: model.TypeMeta{
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
			rc := repositories.NewPodRepository(tt.fields.client)
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
