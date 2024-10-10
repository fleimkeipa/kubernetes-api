package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

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

func (rc *DeploymentInterfaces) Create(ctx context.Context, deployment *v1.Deployment, opts model.CreateOptions) (*v1.Deployment, error) {
	metaOpts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metaOpts)
}

func (rc *DeploymentInterfaces) Update(ctx context.Context, deployment *v1.Deployment, opts model.UpdateOptions) (*v1.Deployment, error) {
	metaOpts := metav1.UpdateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metaOpts)
}

func (rc *DeploymentInterfaces) List(ctx context.Context, namespace string, opts model.ListOptions) (*v1.DeploymentList, error) {
	metaOpts := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		LabelSelector:        opts.LabelSelector,
		FieldSelector:        opts.FieldSelector,
		Watch:                opts.Watch,
		AllowWatchBookmarks:  opts.AllowWatchBookmarks,
		ResourceVersion:      opts.ResourceVersion,
		ResourceVersionMatch: metav1.ResourceVersionMatch(opts.ResourceVersionMatch),
		TimeoutSeconds:       opts.TimeoutSeconds,
		Limit:                opts.Limit,
		Continue:             opts.Continue,
		SendInitialEvents:    opts.SendInitialEvents,
	}
	return rc.client.AppsV1().Deployments(namespace).List(ctx, metaOpts)
}

func (rc *DeploymentInterfaces) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	gracePeriodSeconds := new(int64)
	if opts.GracePeriodSeconds != nil {
		gracePeriodSeconds = opts.GracePeriodSeconds
	}

	preconditions := new(model.Preconditions)
	if opts.Preconditions != nil {
		preconditions = opts.Preconditions
	}

	metaOpts := metav1.DeleteOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		GracePeriodSeconds: gracePeriodSeconds,
		Preconditions: &metav1.Preconditions{
			UID:             preconditions.UID,
			ResourceVersion: preconditions.ResourceVersion,
		},
		DryRun: opts.DryRun,
	}

	return rc.client.AppsV1().Deployments(namespace).Delete(ctx, name, metaOpts)
}
