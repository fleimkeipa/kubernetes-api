package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
	"github.com/fleimkeipa/kubernetes-api/uc"

	v1 "k8s.io/api/apps/v1"
)

func TestDeploymentUC_List(t *testing.T) {
	deploymentTestRepo := repositories.NewDeploymentInterfaces(initTestKubernetes())

	type fields struct {
		deploymentRepo interfaces.DeploymentInterfaces
		eventUC        *uc.EventUC
	}
	type args struct {
		ctx       context.Context
		namespace string
		opts      model.ListOptions
	}
	tests := []struct {
		fields  fields
		want    *v1.DeploymentList
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				deploymentRepo: deploymentTestRepo,
				eventUC:        &uc.EventUC{},
			},
			args: args{
				ctx:       context.TODO(),
				namespace: "",
				opts:      model.ListOptions{},
			},
			want:    &v1.DeploymentList{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := uc.NewDeploymentUC(tt.fields.deploymentRepo, tt.fields.eventUC)
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
	deploymentTestRepo := repositories.NewDeploymentInterfaces(initTestKubernetes())
	type fields struct {
		deploymentRepo interfaces.DeploymentInterfaces
		eventUC        *uc.EventUC
	}
	type args struct {
		ctx       context.Context
		namespace string
		nameOrUID string
		opts      model.ListOptions
	}
	tests := []struct {
		fields  fields
		want    *v1.Deployment
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				deploymentRepo: deploymentTestRepo,
				eventUC:        &uc.EventUC{},
			},
			args: args{
				ctx:       context.TODO(),
				namespace: "",
				nameOrUID: "",
				opts:      model.ListOptions{},
			},
			want:    &v1.Deployment{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := uc.NewDeploymentUC(tt.fields.deploymentRepo, tt.fields.eventUC)
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
