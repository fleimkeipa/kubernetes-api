package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type DeploymentUC struct {
	deploymentRepo interfaces.DeploymentInterfaces
	eventUC        *EventUC
}

func NewDeploymentUC(deploymentRepo interfaces.DeploymentInterfaces, eventUC *EventUC) *DeploymentUC {
	return &DeploymentUC{
		deploymentRepo: deploymentRepo,
		eventUC:        eventUC,
	}
}

func (rc *DeploymentUC) Create(ctx context.Context, request *model.DeploymentCreateRequest) (*model.Deployment, error) {
	request.Deployment.TypeMeta.Kind = "deployment"
	if request.Deployment.ObjectMeta.Namespace == "" {
		request.Deployment.ObjectMeta.Namespace = "default"
	}

	event := model.Event{
		Category: model.DeploymentCategory,
		Type:     model.CreateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	return rc.deploymentRepo.Create(ctx, &request.Deployment, request.Opts)
}

func (rc *DeploymentUC) Update(ctx context.Context, namespace, id string, request *model.DeploymentUpdateRequest) (*model.Deployment, error) {
	event := model.Event{
		Category: model.DeploymentCategory,
		Type:     model.UpdateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeDeployment := rc.fillDeployment(request)
	kubeDeployment.Namespace = namespace

	return rc.deploymentRepo.Update(ctx, namespace, id, kubeDeployment, request.Opts)
}

func (rc *DeploymentUC) List(ctx context.Context, namespace string, opts model.ListOptions) (*model.DeploymentList, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	return rc.deploymentRepo.List(ctx, namespace, opts)
}

func (rc *DeploymentUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Deployment, error) {
	return rc.deploymentRepo.GetByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *DeploymentUC) Delete(ctx context.Context, namespace, nameOrUID string, opts model.DeleteOptions) error {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	event := model.Event{
		Category: model.DeploymentCategory,
		Type:     model.DeleteEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return err
	}

	return rc.deploymentRepo.Delete(ctx, namespace, nameOrUID, opts)
}

func (rc *DeploymentUC) fillDeployment(request *model.DeploymentUpdateRequest) *model.Deployment {
	return &model.Deployment{
		Spec: model.DeploymentSpec{
			Strategy:                request.Deployment.Spec.Strategy,
			Replicas:                request.Deployment.Spec.Replicas,
			ProgressDeadlineSeconds: request.Deployment.Spec.ProgressDeadlineSeconds,
			Template:                request.Deployment.Spec.Template,
			MinReadySeconds:         request.Deployment.Spec.MinReadySeconds,
		},
		ObjectMeta: model.ObjectMeta{
			Labels:      request.Deployment.Labels,
			Annotations: request.Deployment.Annotations,
		},
	}
}
