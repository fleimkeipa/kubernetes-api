package uc

import (
	"context"
	"fmt"

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

func (rc *NamespaceUC) Create(ctx context.Context, namespace *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	opts.TypeMeta.Kind = "namespace"

	var event = model.Event{
		Kind:      model.NamespaceKind,
		EventKind: model.CreateEventKind,
		Owner:     model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	return rc.namespaceRepo.Create(ctx, namespace, opts)
}

func (rc *NamespaceUC) Update(ctx context.Context, id string, request *model.NamespaceUpdateRequest) (*corev1.Namespace, error) {
	existNamespace, err := rc.GetByNameOrUID(ctx, id, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	request.Namespace.TypeMeta.Kind = "namespace"

	var event = model.Event{
		Kind:      model.NamespaceKind,
		EventKind: model.UpdateEventKind,
		Owner:     model.User{},
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	var kubeNamespace = rc.overwriteNamespace(request, existNamespace)

	return rc.namespaceRepo.Update(ctx, kubeNamespace, request.Opts)
}

func (rc *NamespaceUC) List(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error) {
	opts.TypeMeta.Kind = "namespace"

	return rc.namespaceRepo.List(ctx, opts)
}

func (rc *NamespaceUC) GetByNameOrUID(ctx context.Context, nameOrUID string, opts metav1.ListOptions) (*corev1.Namespace, error) {
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

func (rc *NamespaceUC) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	opts.TypeMeta.Kind = "namespace"

	var event = model.Event{
		Kind:      model.NamespaceKind,
		EventKind: model.DeleteEventKind,
		Owner:     model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return fmt.Errorf("failed to create event for %s: %w", event.EventKind, err)
	}

	return rc.namespaceRepo.Delete(ctx, name, opts)
}

func (rc *NamespaceUC) overwriteNamespace(newNamespace *model.NamespaceUpdateRequest, existNamespace *corev1.Namespace) *corev1.Namespace {
	existNamespace.Name = newNamespace.Namespace.Name

	existNamespace.Kind = "namespace"
	existNamespace.APIVersion = "v1"

	return existNamespace
}
