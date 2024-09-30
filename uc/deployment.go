package uc

import (
	"context"

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
		eventUC:        eventUC,
	}
}

func (rc *DeploymentUC) Create(ctx context.Context, request *model.DeploymentCreateRequest) (*v1.Deployment, error) {
	newDeployment := request.Deployment

	newDeployment.TypeMeta.Kind = "deployment"
	if newDeployment.ObjectMeta.Namespace == "" {
		newDeployment.ObjectMeta.Namespace = "default"
	}

	event := model.Event{
		Category: model.DeploymentKind,
		Type:     model.CreateEventKind,
		Owner:    model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeDeployment := rc.fillDeployment(&newDeployment)

	return rc.deploymentRepo.Create(ctx, kubeDeployment, request.Opts)
}

func (rc *DeploymentUC) Update(ctx context.Context, id, namespace string, request *model.DeploymentUpdateRequest) (*v1.Deployment, error) {
	existDeployment, err := rc.GetByNameOrUID(ctx, namespace, id, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	event := model.Event{
		Category: model.DeploymentKind,
		Type:     model.UpdateEventKind,
		Owner:    model.User{},
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeDeployment := rc.overwriteDeployment(&request.Deployment, existDeployment)

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

	event := model.Event{
		Category: model.DeploymentKind,
		Type:     model.DeleteEventKind,
		Owner:    model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return err
	}

	return rc.deploymentRepo.Delete(ctx, namespace, name, opts)
}

func (rc *DeploymentUC) fillDeployment(newDeployment *model.Deployment) *v1.Deployment {
	matchExpressions := make([]metav1.LabelSelectorRequirement, 0)
	for _, v := range newDeployment.Spec.Selector.MatchExpressions {
		matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
			Key:      v.Key,
			Operator: metav1.LabelSelectorOperator(v.Operator),
			Values:   v.Values,
		})
	}
	selector := metav1.LabelSelector{
		MatchLabels:      newDeployment.Spec.Selector.MatchLabels,
		MatchExpressions: matchExpressions,
	}

	conditions := make([]v1.DeploymentCondition, 0)
	for _, v := range newDeployment.Status.Conditions {
		conditions = append(conditions, v1.DeploymentCondition{
			Type:               v1.DeploymentConditionType(v.Type),
			Status:             v.Status,
			LastUpdateTime:     metav1.Time{Time: v.LastUpdateTime},
			LastTransitionTime: metav1.Time{Time: v.LastTransitionTime},
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}

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
			Replicas: newDeployment.Spec.Replicas,
			Selector: &selector,
			Strategy: v1.DeploymentStrategy{
				Type:          v1.DeploymentStrategyType(newDeployment.Spec.Strategy.Type),
				RollingUpdate: (*v1.RollingUpdateDeployment)(newDeployment.Spec.Strategy.RollingUpdate),
			},
			MinReadySeconds:         newDeployment.Spec.MinReadySeconds,
			RevisionHistoryLimit:    newDeployment.Spec.RevisionHistoryLimit,
			Paused:                  newDeployment.Spec.Paused,
			ProgressDeadlineSeconds: newDeployment.Spec.ProgressDeadlineSeconds,
		},
		Status: v1.DeploymentStatus{
			ObservedGeneration:  newDeployment.Status.ObservedGeneration,
			Replicas:            newDeployment.Status.Replicas,
			UpdatedReplicas:     newDeployment.Status.UpdatedReplicas,
			ReadyReplicas:       newDeployment.Status.ReadyReplicas,
			AvailableReplicas:   newDeployment.Status.AvailableReplicas,
			UnavailableReplicas: newDeployment.Status.UnavailableReplicas,
			Conditions:          conditions,
			CollisionCount:      newDeployment.Status.CollisionCount,
		},
	}
}

func (rc *DeploymentUC) overwriteDeployment(newDeployment *model.Deployment, existDeployment *v1.Deployment) *v1.Deployment {
	existDeployment.Name = newDeployment.Name

	return existDeployment
}
