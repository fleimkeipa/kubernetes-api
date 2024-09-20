package uc

import (
	"context"
	"fmt"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type DeploymentUC struct {
	deploymentRepo interfaces.DeploymentInterfaces
	eventUC        *EventUC
}

func NewDeploymentUC(deploymentRepo interfaces.DeploymentInterfaces, eventUC *EventUC) *DeploymentUC {
	return &DeploymentUC{
		deploymentRepo: deploymentRepo,
	}
}

func (rc *DeploymentUC) Create(ctx context.Context, request *model.DeploymentCreateRequest) (*v1.Deployment, error) {
	var newDeployment = request.Deployment

	newDeployment.TypeMeta.Kind = "deployment"
	if newDeployment.ObjectMeta.Namespace == "" {
		newDeployment.ObjectMeta.Namespace = "default"
	}

	var event = model.Event{
		Kind:      model.DeploymentKind,
		EventKind: model.CreateEventKind,
		Owner:     model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	var kubeDeployment = rc.fillDeployment(&newDeployment)

	return rc.deploymentRepo.Create(ctx, kubeDeployment, request.Opts)
}

func (rc *DeploymentUC) Update(ctx context.Context, id, namespace string, request *model.DeploymentUpdateRequest) (*v1.Deployment, error) {
	existDeployment, err := rc.GetByNameOrUID(ctx, namespace, id, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var event = model.Event{
		Kind:      model.DeploymentKind,
		EventKind: model.UpdateEventKind,
		Owner:     model.User{},
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	var kubeDeployment = rc.overwriteDeployment(&request.Deployment, existDeployment)

	return rc.deploymentRepo.Update(ctx, kubeDeployment, request.Opts)
}

func (rc *DeploymentUC) List(ctx context.Context, namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	return rc.deploymentRepo.List(ctx, namespace, opts)
}

func (rc *DeploymentUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts metav1.ListOptions) (*v1.Deployment, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	deployments, err := rc.deploymentRepo.List(ctx, namespace, opts)
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

	var event = model.Event{
		Kind:      model.DeploymentKind,
		EventKind: model.DeleteEventKind,
		Owner:     model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	return rc.deploymentRepo.Delete(ctx, namespace, name, opts)
}

func (rc *DeploymentUC) fillDeployment(newDeployment *model.Deployment) *v1.Deployment {
	return &v1.Deployment{
		TypeMeta: metav1.TypeMeta(newDeployment.TypeMeta),
		ObjectMeta: metav1.ObjectMeta{
			Name:                       newDeployment.ObjectMeta.Name,
			GenerateName:               newDeployment.ObjectMeta.GenerateName,
			Namespace:                  newDeployment.ObjectMeta.Namespace,
			ResourceVersion:            newDeployment.ObjectMeta.ResourceVersion,
			Generation:                 newDeployment.ObjectMeta.Generation,
			DeletionGracePeriodSeconds: newDeployment.ObjectMeta.DeletionGracePeriodSeconds,
			Labels:                     newDeployment.ObjectMeta.Labels,
			Annotations:                newDeployment.ObjectMeta.Annotations,
			Finalizers:                 newDeployment.ObjectMeta.Finalizers,
		},
		Spec: v1.DeploymentSpec{
			Replicas:                new(int32),
			Selector:                &metav1.LabelSelector{},
			Strategy:                v1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    new(int32),
			Paused:                  false,
			ProgressDeadlineSeconds: new(int32),
		},
		Status: v1.DeploymentStatus{},
	}
}

func (rc *DeploymentUC) overwriteDeployment(newDeployment *model.Deployment, existDeployment *v1.Deployment) *v1.Deployment {
	existDeployment.Name = newDeployment.Name

	return existDeployment
}
