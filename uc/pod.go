package uc

import (
	"context"
	"fmt"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories/interfaces"
)

type PodUC struct {
	podsRepo interfaces.PodInterfaces
	eventUC  *EventUC
}

func NewPodUC(podsRepo interfaces.PodInterfaces, eventUC *EventUC) *PodUC {
	return &PodUC{
		podsRepo: podsRepo,
		eventUC:  eventUC,
	}
}

func (rc *PodUC) Create(ctx context.Context, request model.PodsCreateRequest) (*model.Pod, error) {
	request.Pod.TypeMeta.Kind = "pod"
	if request.Pod.ObjectMeta.Namespace == "" {
		request.Pod.ObjectMeta.Namespace = "default"
	}

	event := model.Event{
		Category: model.PodCategory,
		Type:     model.CreateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	return rc.podsRepo.Create(ctx, &request.Pod, request.Opts)
}

func (rc *PodUC) Update(ctx context.Context, namespace, id string, request *model.PodsUpdateRequest) (*model.Pod, error) {
	event := model.Event{
		Category: model.PodCategory,
		Type:     model.UpdateEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.Type, err)
	}

	kubePod := rc.fillPod(request)
	kubePod.Namespace = namespace

	return rc.podsRepo.Update(ctx, id, kubePod, request.Opts)
}

func (rc *PodUC) List(ctx context.Context, namespace string, opts model.ListOptions) (*model.PodList, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	return rc.podsRepo.List(ctx, namespace, opts)
}

func (rc *PodUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Pod, error) {
	return rc.podsRepo.GetByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *PodUC) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	event := model.Event{
		Category: model.PodCategory,
		Type:     model.DeleteEventType,
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return fmt.Errorf("failed to create event for %s: %w", event.Type, err)
	}

	return rc.podsRepo.Delete(ctx, namespace, name, opts)
}

func (rc *PodUC) fillPod(podRequest *model.PodsUpdateRequest) *model.Pod {
	newPods := new(model.Pod)

	newPods.Spec.Containers = rc.fillContainers(podRequest.Pod.Spec.Containers)
	newPods.Spec.InitContainers = rc.fillContainers(podRequest.Pod.Spec.InitContainers)

	newPods.Spec.Tolerations = rc.fillTolerations(podRequest.Pod.Spec.Tolerations)

	newPods.Spec.ActiveDeadlineSeconds = podRequest.Pod.Spec.ActiveDeadlineSeconds

	graceSeconds := newPods.Spec.TerminationGracePeriodSeconds
	if graceSeconds == nil { // (allow it to be set to 1 if it was previously negative)
		newPods.Spec.TerminationGracePeriodSeconds = podRequest.Pod.Spec.TerminationGracePeriodSeconds
	}

	return newPods
}

// fillContainers change only images
func (rc *PodUC) fillContainers(newContainers []model.ContainerRequest) []model.Container {
	conts := make([]model.Container, 0)
	for _, v := range newContainers {
		conts = append(conts, model.Container{
			Name:  v.Name,
			Image: v.Image,
		})
	}

	return conts
}

// addToleration only additions to existing tolerations
func (rc *PodUC) fillTolerations(newTolerations []model.Toleration) []model.Toleration {
	additions := make([]model.Toleration, 0)
	for _, v := range newTolerations {
		additions = append(additions, model.Toleration{
			Key:               v.Key,
			Operator:          model.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            model.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		})
	}

	return additions
}
