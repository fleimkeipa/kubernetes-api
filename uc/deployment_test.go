package uc

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deploymentTestRepo interfaces.DeploymentInterfaces

func init() {
	deploymentTestRepo = repositories.NewDeploymentInterfaces(initTestKubernetes())
}

func TestDeploymentUC_List(t *testing.T) {
	type fields struct {
		deploymentRepo interfaces.DeploymentInterfaces
		eventUC        *EventUC
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
		want    *v1.DeploymentList
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				deploymentRepo: deploymentTestRepo,
				eventUC:        &EventUC{},
			},
			args: args{
				ctx:       context.TODO(),
				namespace: "",
				opts:      metav1.ListOptions{},
			},
			want:    &v1.DeploymentList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &DeploymentUC{
				deploymentRepo: tt.fields.deploymentRepo,
				eventUC:        tt.fields.eventUC,
			}
			got, err := rc.List(tt.args.ctx, tt.args.namespace, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentUC.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentUC.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentUC_GetByNameOrUID(t *testing.T) {
	type fields struct {
		deploymentRepo interfaces.DeploymentInterfaces
		eventUC        *EventUC
	}
	type args struct {
		ctx       context.Context
		namespace string
		nameOrUID string
		opts      metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.Deployment
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				deploymentRepo: deploymentTestRepo,
				eventUC:        &EventUC{},
			},
			args: args{
				ctx:       context.TODO(),
				namespace: "",
				nameOrUID: "",
				opts:      metav1.ListOptions{},
			},
			want:    &v1.Deployment{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &DeploymentUC{
				deploymentRepo: tt.fields.deploymentRepo,
				eventUC:        tt.fields.eventUC,
			}
			got, err := rc.GetByNameOrUID(tt.args.ctx, tt.args.namespace, tt.args.nameOrUID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentUC.GetByNameOrUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentUC.GetByNameOrUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
