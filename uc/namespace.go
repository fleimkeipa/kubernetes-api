package uc

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type NamespaceUC struct {
	namespaceRepo interfaces.NamespaceInterfaces
	eventUC       *EventUC
}

func NewNamespaceUC(namespaceRepo interfaces.NamespaceInterfaces, eventUC *EventUC) *NamespaceUC {
	return &NamespaceUC{
		namespaceRepo: namespaceRepo,
		eventUC:       eventUC,
	}
}

func (rc *NamespaceUC) Create(ctx context.Context, request model.NamespaceCreateRequest) (*corev1.Namespace, error) {
	request.Opts.TypeMeta.Kind = "namespace"

	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.CreateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeNamespace := rc.fillNamespace(&request.Namespace)

	return rc.namespaceRepo.Create(ctx, kubeNamespace, request.Opts)
}

func (rc *NamespaceUC) Update(ctx context.Context, id string, request *model.NamespaceUpdateRequest) (*corev1.Namespace, error) {
	existNamespace, err := rc.GetByNameOrUID(ctx, id, model.ListOptions{})
	if err != nil {
		return nil, err
	}

	request.Namespace.TypeMeta.Kind = "namespace"

	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.UpdateEventType,
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubeNamespace := rc.overwriteNamespace(request, existNamespace)

	return rc.namespaceRepo.Update(ctx, kubeNamespace, request.Opts)
}

func (rc *NamespaceUC) List(ctx context.Context, opts model.ListOptions) (*corev1.NamespaceList, error) {
	opts.TypeMeta.Kind = "namespace"

	return rc.namespaceRepo.List(ctx, opts)
}

func (rc *NamespaceUC) GetByNameOrUID(ctx context.Context, nameOrUID string, opts model.ListOptions) (*corev1.Namespace, error) {
	opts.TypeMeta.Kind = "namespace"

	opts.Limit = 100
	namespaces, err := rc.namespaceRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range namespaces.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if namespaces.ListMeta.Continue == "" {
		return &corev1.Namespace{}, nil
	}

	opts.Continue = namespaces.ListMeta.Continue
	return rc.GetByNameOrUID(ctx, nameOrUID, opts)
}

func (rc *NamespaceUC) Delete(ctx context.Context, name string, opts model.DeleteOptions) error {
	opts.TypeMeta.Kind = "namespace"

	event := model.Event{
		Category: model.NamespaceCategory,
		Type:     model.DeleteEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return err
	}

	return rc.namespaceRepo.Delete(ctx, name, opts)
}

func (rc *NamespaceUC) fillNamespace(namespace *model.Namespace) *corev1.Namespace {
	conditions := make([]corev1.NamespaceCondition, 0)
	for _, v := range namespace.Status.Conditions {
		conditions = append(conditions, corev1.NamespaceCondition{
			Type:               corev1.NamespaceConditionType(v.Type),
			Status:             corev1.ConditionStatus(v.Status),
			LastTransitionTime: metav1.Time{Time: v.LastTransitionTime},
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}

	finalizers := make([]corev1.FinalizerName, 0)
	for _, v := range namespace.Finalizers {
		finalizers = append(finalizers, corev1.FinalizerName(v))
	}

	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta(namespace.TypeMeta),
		ObjectMeta: metav1.ObjectMeta{
			Name:                       namespace.Name,
			GenerateName:               namespace.GenerateName,
			Namespace:                  namespace.Namespace,
			ResourceVersion:            namespace.ResourceVersion,
			Generation:                 namespace.Generation,
			DeletionGracePeriodSeconds: namespace.DeletionGracePeriodSeconds,
			Labels:                     namespace.Labels,
			Annotations:                namespace.Annotations,
			Finalizers:                 namespace.Finalizers,
		},
		Spec: corev1.NamespaceSpec{
			Finalizers: finalizers,
		},
		Status: corev1.NamespaceStatus{
			Phase:      corev1.NamespacePhase(namespace.Status.Phase),
			Conditions: conditions,
		},
	}
}

func (rc *NamespaceUC) overwriteNamespace(newNamespace *model.NamespaceUpdateRequest, existNamespace *corev1.Namespace) *corev1.Namespace {
	existNamespace.Name = newNamespace.Namespace.Name

	existNamespace.Kind = "namespace"
	existNamespace.APIVersion = "v1"

	return existNamespace
}
