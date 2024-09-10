package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type DeploymentUC struct {
	deploymentRepo interfaces.DeploymentInterfaces
}

func NewDeploymentUC(deploymentRepo interfaces.DeploymentInterfaces) *DeploymentUC {
	return &DeploymentUC{
		deploymentRepo: deploymentRepo,
	}
}

func (rc *DeploymentUC) Get(ctx context.Context, namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	return rc.deploymentRepo.Get(ctx, namespace, opts)
}

func (rc *DeploymentUC) Create(ctx context.Context, deployment *v1.Deployment, opts metav1.CreateOptions) (*v1.Deployment, error) {
	deployment.TypeMeta.Kind = "deployment"
	if deployment.ObjectMeta.Namespace == "" {
		deployment.ObjectMeta.Namespace = "default"
	}

	return rc.deploymentRepo.Create(ctx, deployment, opts)
}

func (rc *DeploymentUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts metav1.ListOptions) (*v1.Deployment, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	deployments, err := rc.deploymentRepo.Get(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range deployments.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if deployments.ListMeta.Continue == "" {
		return &v1.Deployment{}, nil
	}

	opts.Continue = deployments.ListMeta.Continue
	return rc.GetByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *DeploymentUC) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	return rc.deploymentRepo.Delete(ctx, namespace, name, opts)
}

func (rc *DeploymentUC) Update(ctx context.Context, deployment *v1.Deployment, opts metav1.UpdateOptions) (*v1.Deployment, error) {
	deployment.TypeMeta.Kind = "deployment"
	if deployment.ObjectMeta.Namespace == "" {
		deployment.ObjectMeta.Namespace = "default"
	}

	return rc.deploymentRepo.Update(ctx, deployment, opts)
}
