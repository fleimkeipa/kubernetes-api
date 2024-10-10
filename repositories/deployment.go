package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentInterfaces struct {
	client *kubernetes.Clientset
}

func NewDeploymentInterfaces(client *kubernetes.Clientset) *DeploymentInterfaces {
	return &DeploymentInterfaces{client}
}

func (rc *DeploymentInterfaces) Create(ctx context.Context, deployment *v1.Deployment, opts model.CreateOptions) (*v1.Deployment, error) {
	metaOpts := convertCreateOptsToKube(opts)

	return rc.client.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metaOpts)
}

func (rc *DeploymentInterfaces) Update(ctx context.Context, deployment *v1.Deployment, opts model.UpdateOptions) (*v1.Deployment, error) {
	metaOpts := convertUpdateOptsToKube(opts)

	return rc.client.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metaOpts)
}

func (rc *DeploymentInterfaces) List(ctx context.Context, namespace string, opts model.ListOptions) (*v1.DeploymentList, error) {
	metaOpts := convertListOptsToKube(opts)

	return rc.client.AppsV1().Deployments(namespace).List(ctx, metaOpts)
}

func (rc *DeploymentInterfaces) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	metaOpts := convertDeleteOptsToKube(opts)

	return rc.client.AppsV1().Deployments(namespace).Delete(ctx, name, metaOpts)
}
