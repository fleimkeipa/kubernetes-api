package repositories

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentInterfaces struct {
	client *kubernetes.Clientset
}

func NewDeploymentInterfaces(client *kubernetes.Clientset) *DeploymentInterfaces {
	return &DeploymentInterfaces{client}
}

func (rc *DeploymentInterfaces) Create(ctx context.Context, deployment *v1.Deployment, opts metav1.CreateOptions) (*v1.Deployment, error) {
	return rc.client.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, opts)
}

func (rc *DeploymentInterfaces) Update(ctx context.Context, deployment *v1.Deployment, opts metav1.UpdateOptions) (*v1.Deployment, error) {
	return rc.client.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, opts)
}

func (rc *DeploymentInterfaces) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	return rc.client.AppsV1().Deployments(namespace).List(ctx, opts)
}

func (rc *DeploymentInterfaces) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return rc.client.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}
