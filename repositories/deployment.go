package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type DeploymentRepository struct {
	client *kubernetes.Clientset
}

func NewDeploymentInterfaces(client *kubernetes.Clientset) *DeploymentRepository {
	return &DeploymentRepository{client}
}

func (rc *DeploymentRepository) Create(ctx context.Context, deployment *model.Deployment, opts model.CreateOptions) (*model.Deployment, error) {
	metaOpts := convertCreateOptsToKube(opts)

	kubeDeployment := rc.fillRequestDeployment(deployment)

	createdDeployment, err := rc.client.AppsV1().Deployments(deployment.Namespace).Create(ctx, kubeDeployment, metaOpts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseDeployment(createdDeployment), nil
}

func (rc *DeploymentRepository) Update(ctx context.Context, namespace, deploymentID string, deployment *model.Deployment, opts model.UpdateOptions) (*model.Deployment, error) {
	metaOpts := convertUpdateOptsToKube(opts)

	existDeployment, err := rc.getByNameOrUID(ctx, namespace, deploymentID, model.ListOptions{})
	if err != nil {
		return nil, err
	}

	kubeDeployment := rc.overwriteOnKubeDeployment(deployment, existDeployment)

	updatedDeployment, err := rc.client.AppsV1().Deployments(namespace).Update(ctx, kubeDeployment, metaOpts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseDeployment(updatedDeployment), nil
}

func (rc *DeploymentRepository) List(ctx context.Context, namespace string, opts model.ListOptions) (*model.DeploymentList, error) {
	kubeDeployments, err := rc.list(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}

	deploymentList := model.DeploymentList{}
	for _, deployment := range kubeDeployments.Items {
		deploymentList.Items = append(deploymentList.Items, *rc.fillResponseDeployment(&deployment))
	}

	deploymentList.ListMeta = model.ListMeta{
		RemainingItemCount: kubeDeployments.ListMeta.RemainingItemCount,
		ResourceVersion:    kubeDeployments.ListMeta.ResourceVersion,
		Continue:           kubeDeployments.ListMeta.Continue,
	}
	deploymentList.TypeMeta = model.TypeMeta(kubeDeployments.TypeMeta)

	return &deploymentList, nil
}

func (rc *DeploymentRepository) Delete(ctx context.Context, namespace, nameOrUID string, opts model.DeleteOptions) error {
	metaOpts := convertDeleteOptsToKube(opts)

	if metaOpts.Preconditions != nil && metaOpts.Preconditions.UID != nil {
		return rc.client.AppsV1().Deployments(namespace).DeleteCollection(ctx, metaOpts, metav1.ListOptions{})
	}

	return rc.client.AppsV1().Deployments(namespace).Delete(ctx, nameOrUID, metaOpts)
}

func (rc *DeploymentRepository) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Deployment, error) {
	deployment, err := rc.getByNameOrUID(ctx, namespace, nameOrUID, opts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseDeployment(deployment), nil
}

func (rc *DeploymentRepository) getByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*v1.Deployment, error) {
	opts.TypeMeta.Kind = "deployment"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	deployments, err := rc.list(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range deployments.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if deployments.ListMeta.Continue == "" {
		return nil, fmt.Errorf("deployment %s not found", nameOrUID)
	}

	opts.Continue = deployments.ListMeta.Continue
	return rc.getByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *DeploymentRepository) list(ctx context.Context, namespace string, opts model.ListOptions) (*v1.DeploymentList, error) {
	metaOpts := convertListOptsToKube(opts)

	return rc.client.AppsV1().Deployments(namespace).List(ctx, metaOpts)
}

func (rc *DeploymentRepository) fillRequestDeployment(deployment *model.Deployment) *v1.Deployment {
	selector := fillSelector(deployment)

	conditions := fillConditions(deployment)

	template := fillTemplate(deployment)

	return &v1.Deployment{
		TypeMeta: metav1.TypeMeta(deployment.TypeMeta),
		ObjectMeta: metav1.ObjectMeta{
			Name:                       deployment.ObjectMeta.Name,
			GenerateName:               deployment.ObjectMeta.GenerateName,
			Namespace:                  deployment.ObjectMeta.Namespace,
			ResourceVersion:            deployment.ObjectMeta.ResourceVersion,
			Generation:                 deployment.ObjectMeta.Generation,
			DeletionGracePeriodSeconds: deployment.ObjectMeta.DeletionGracePeriodSeconds,
			Labels:                     deployment.ObjectMeta.Labels,
			Annotations:                deployment.ObjectMeta.Annotations,
			Finalizers:                 deployment.ObjectMeta.Finalizers,
		},
		Spec: v1.DeploymentSpec{
			Replicas: deployment.Spec.Replicas,
			Selector: &selector,
			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyType(deployment.Spec.Strategy.Type),
				// RollingUpdate: (*v1.RollingUpdateDeployment)(deployment.Spec.Strategy.RollingUpdate),
			},
			Template:                template,
			MinReadySeconds:         deployment.Spec.MinReadySeconds,
			RevisionHistoryLimit:    deployment.Spec.RevisionHistoryLimit,
			Paused:                  deployment.Spec.Paused,
			ProgressDeadlineSeconds: deployment.Spec.ProgressDeadlineSeconds,
		},
		Status: v1.DeploymentStatus{
			ObservedGeneration:  deployment.Status.ObservedGeneration,
			Replicas:            deployment.Status.Replicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			ReadyReplicas:       deployment.Status.ReadyReplicas,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
			Conditions:          conditions,
			CollisionCount:      deployment.Status.CollisionCount,
		},
	}
}

func fillTemplate(deployment *model.Deployment) corev1.PodTemplateSpec {
	containers := convertContainersToKube(deployment.Spec.Template.Spec.Containers)
	initContainers := convertContainersToKube(deployment.Spec.Template.Spec.InitContainers)

	tolerations := convertTolerationsToKube(deployment.Spec.Template.Spec.Tolerations)

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deployment.Spec.Template.ObjectMeta.Name,
			Labels:      deployment.Spec.Template.ObjectMeta.Labels,
			Annotations: deployment.Spec.Template.ObjectMeta.Annotations,
		},
		Spec: corev1.PodSpec{
			Containers:            containers,
			InitContainers:        initContainers,
			Tolerations:           tolerations,
			ActiveDeadlineSeconds: deployment.Spec.Template.Spec.ActiveDeadlineSeconds,
		},
	}
	return template
}

func fillConditions(newDeployment *model.Deployment) []v1.DeploymentCondition {
	conditions := make([]v1.DeploymentCondition, 0)

	for _, v := range newDeployment.Status.Conditions {
		conditions = append(conditions, v1.DeploymentCondition{
			Type:               v1.DeploymentConditionType(v.Type),
			Status:             corev1.ConditionStatus(v.Status),
			LastUpdateTime:     metav1.Time{Time: v.LastUpdateTime},
			LastTransitionTime: metav1.Time{Time: v.LastTransitionTime},
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}
	return conditions
}

func fillSelector(newDeployment *model.Deployment) metav1.LabelSelector {
	matchExpressions := make([]metav1.LabelSelectorRequirement, 0)
	if newDeployment.Spec.Selector == nil {
		return metav1.LabelSelector{
			MatchLabels:      nil,
			MatchExpressions: matchExpressions,
		}
	}

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

	return selector
}

func (rc *DeploymentRepository) fillResponseDeployment(deployment *v1.Deployment) *model.Deployment {
	matchExpressions := make([]model.LabelSelectorRequirement, 0)
	for _, v := range deployment.Spec.Selector.MatchExpressions {
		matchExpressions = append(matchExpressions, model.LabelSelectorRequirement{
			Key:      v.Key,
			Operator: model.LabelSelectorOperator(v.Operator),
			Values:   v.Values,
		})
	}
	selector := model.LabelSelector{
		MatchLabels:      deployment.Spec.Selector.MatchLabels,
		MatchExpressions: matchExpressions,
	}

	conditions := make([]model.DeploymentCondition, 0)
	for _, v := range deployment.Status.Conditions {
		conditions = append(conditions, model.DeploymentCondition{
			Type:               model.DeploymentConditionType(v.Type),
			Status:             model.ConditionStatus(v.Status),
			LastUpdateTime:     v.LastUpdateTime.Time,
			LastTransitionTime: v.LastTransitionTime.Time,
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}

	template := convertTemplateToModel(deployment)

	ownerReferences := make([]model.OwnerReference, 0, len(deployment.OwnerReferences))
	for _, v := range deployment.OwnerReferences {
		ownerReferences = append(ownerReferences, model.OwnerReference{
			Controller:         v.Controller,
			BlockOwnerDeletion: v.BlockOwnerDeletion,
			APIVersion:         v.APIVersion,
			Kind:               v.Kind,
			Name:               v.Name,
		})
	}

	deletionTimestamp := new(time.Time)
	if deletionTime := deployment.DeletionTimestamp; deletionTime != nil {
		deletionTimestamp = &deletionTime.Time
	}

	return &model.Deployment{
		TypeMeta: model.TypeMeta(deployment.TypeMeta),
		ObjectMeta: model.ObjectMeta{
			UID:                        string(deployment.UID),
			CreationTimestamp:          deployment.CreationTimestamp.Time,
			DeletionTimestamp:          deletionTimestamp,
			DeletionGracePeriodSeconds: deployment.ObjectMeta.DeletionGracePeriodSeconds,
			Labels:                     deployment.ObjectMeta.Labels,
			Annotations:                deployment.ObjectMeta.Annotations,
			Name:                       deployment.ObjectMeta.Name,
			GenerateName:               deployment.ObjectMeta.GenerateName,
			Namespace:                  deployment.ObjectMeta.Namespace,
			ResourceVersion:            deployment.ObjectMeta.ResourceVersion,
			OwnerReferences:            ownerReferences,
			Finalizers:                 deployment.ObjectMeta.Finalizers,
			Generation:                 deployment.ObjectMeta.Generation,
		},
		Spec: model.DeploymentSpec{
			Replicas: deployment.Spec.Replicas,
			Selector: &selector,
			Strategy: model.DeploymentStrategy{
				Type: model.DeploymentStrategyType(deployment.Spec.Strategy.Type),
				// RollingUpdate: (*model.RollingUpdateDeployment)(newDeployment.Spec.Strategy.RollingUpdate),
			},
			Template:                template,
			MinReadySeconds:         deployment.Spec.MinReadySeconds,
			RevisionHistoryLimit:    deployment.Spec.RevisionHistoryLimit,
			Paused:                  deployment.Spec.Paused,
			ProgressDeadlineSeconds: deployment.Spec.ProgressDeadlineSeconds,
		},
		Status: model.DeploymentStatus{
			ObservedGeneration:  deployment.Status.ObservedGeneration,
			Replicas:            deployment.Status.Replicas,
			UpdatedReplicas:     deployment.Status.UpdatedReplicas,
			ReadyReplicas:       deployment.Status.ReadyReplicas,
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
			Conditions:          conditions,
			CollisionCount:      deployment.Status.CollisionCount,
		},
	}
}

func (rc *DeploymentRepository) overwriteOnKubeDeployment(newDeployment *model.Deployment, existDeployment *v1.Deployment) *v1.Deployment {
	existDeployment.ObjectMeta.Labels = newDeployment.Labels
	existDeployment.ObjectMeta.Annotations = newDeployment.Annotations

	existDeployment.Spec.Replicas = newDeployment.Spec.Replicas
	existDeployment.Spec.Strategy = v1.DeploymentStrategy{
		Type: v1.DeploymentStrategyType(newDeployment.Spec.Strategy.Type),
		// RollingUpdate: (*v1.RollingUpdateDeployment)(newDeployment.Spec.Strategy.RollingUpdate),
	}
	existDeployment.Spec.Template = convertTemplateToKube(&newDeployment.Spec.Template)
	existDeployment.Spec.MinReadySeconds = newDeployment.Spec.MinReadySeconds
	existDeployment.Spec.ProgressDeadlineSeconds = newDeployment.Spec.ProgressDeadlineSeconds

	return existDeployment
}
