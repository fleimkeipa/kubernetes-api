package uc

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	"github.com/go-pg/pg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	loadEnv()
}

func initTestKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}

func initTestDB() *pg.DB {
	return pkg.NewPSQLClient()
}

func TestPodsUC_GetByName(t *testing.T) {
	type fields struct {
		podsRepo interfaces.PodsInterfaces
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
				podsRepo: repositories.NewPodsRepository(initTestKubernetes()),
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
			rc := &PodsUC{
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

func TestPodsUC_Update(t *testing.T) {
	var second int64 = 128
	type fields struct {
		podsRepo  interfaces.PodsInterfaces
		eventRepo interfaces.EventsInterfaces
	}
	type args struct {
		ctx  context.Context
		pod  *model.Pod
		opts metav1.UpdateOptions
		id   string
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
				podsRepo:  repositories.NewPodsRepository(initTestKubernetes()),
				eventRepo: repositories.NewEventRepository(initTestDB()),
			},
			args: args{
				ctx: context.TODO(),
				pod: &model.Pod{
					ObjectMeta: model.ObjectMeta{
						Name:      "pod1",
						Namespace: "default",
					},
					Spec: model.PodSpec{
						Containers: []model.Container{
							{
								Name:       "nginx",
								Image:      "nginx:1.16",
								WorkingDir: "dir1",
							},
						},
						ActiveDeadlineSeconds: &second,
					},
				},
				opts: metav1.UpdateOptions{},
				id:   "pod1",
			},
			want:    &corev1.Pod{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &PodsUC{
				podsRepo:  tt.fields.podsRepo,
				eventRepo: tt.fields.eventRepo,
			}
			got, err := rc.Update(tt.args.ctx, tt.args.id, tt.args.pod, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodsUC.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodsUC.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
