package uc

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

func initTestKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return client
}

func TestPodsUC_Update(t *testing.T) {
	type fields struct {
		podsRepo interfaces.PodsInterfaces
	}
	type args struct {
		ctx  context.Context
		pod  *corev1.Pod
		opts metav1.UpdateOptions
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
				ctx: context.TODO(),
				pod: &corev1.Pod{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       corev1.PodSpec{},
					Status:     corev1.PodStatus{},
				},
				opts: metav1.UpdateOptions{},
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
			got, err := rc.Update(tt.args.ctx, tt.args.pod, tt.args.opts)
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
