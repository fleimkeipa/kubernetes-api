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

func (rc *PodUC) Create(ctx context.Context, pod *model.Pod, opts model.CreateOptions) (*corev1.Pod, error) {
	pod.TypeMeta.Kind = "pod"
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = "default"
	}

	event := model.Event{
		Category: model.PodKind,
		Type:     model.CreateEventKind,
		Owner:    model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, err
	}

	kubePod := rc.fillPod(pod)

	return rc.podsRepo.Create(ctx, kubePod, opts)
}

func (rc *PodUC) Update(ctx context.Context, id string, request *model.PodsUpdateRequest) (*corev1.Pod, error) {
	existPod, err := rc.GetByNameOrUID(ctx, request.Pod.Namespace, id, model.ListOptions{})
	if err != nil {
		return nil, err
	}

	event := model.Event{
		Category: model.PodKind,
		Type:     model.UpdateEventKind,
		Owner:    model.User{},
	}
	_, err = rc.eventUC.Create(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event for %s: %w", event.Type, err)
	}

	kubePod := rc.overwritePod(request, existPod)

	return rc.podsRepo.Update(ctx, id, kubePod, request.Opts)
}

func (rc *PodUC) List(ctx context.Context, namespace string, opts model.ListOptions) (*corev1.PodList, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	return rc.podsRepo.List(ctx, namespace, opts)
}

func (rc *PodUC) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*corev1.Pod, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	pods, err := rc.podsRepo.List(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range pods.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if pods.ListMeta.Continue == "" {
		return &corev1.Pod{}, nil
	}

	opts.Continue = pods.ListMeta.Continue
	return rc.GetByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *PodUC) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	event := model.Event{
		Category: model.PodKind,
		Type:     model.DeleteEventKind,
		Owner:    model.User{},
	}
	_, err := rc.eventUC.Create(ctx, &event)
	if err != nil {
		return fmt.Errorf("failed to create event for %s: %w", event.Type, err)
	}

	return rc.podsRepo.Delete(ctx, namespace, name, opts)
}

func (rc *PodUC) fillPod(pod *model.Pod) *corev1.Pod {
	volumes := make([]corev1.Volume, 0)
	for _, v := range pod.Spec.Volumes {
		volumes = append(volumes, corev1.Volume{
			Name:         v.Name,
			VolumeSource: corev1.VolumeSource{},
		})
	}

	containers := make([]corev1.Container, 0)
	for _, v := range pod.Spec.Containers {
		containers = append(containers, corev1.Container{
			Name:                   v.Name,
			Image:                  v.Image,
			Command:                v.Command,
			Args:                   v.Args,
			WorkingDir:             v.WorkingDir,
			TerminationMessagePath: v.TerminationMessagePath,
			Stdin:                  v.Stdin,
			StdinOnce:              v.StdinOnce,
			TTY:                    v.TTY,
		})
	}

	conditions := make([]corev1.PodCondition, 0)
	for _, v := range pod.Status.Conditions {
		conditions = append(conditions, corev1.PodCondition{
			Type:               corev1.PodConditionType(v.Type),
			Status:             corev1.ConditionStatus(v.Status),
			LastProbeTime:      v.LastProbeTime,
			LastTransitionTime: v.LastTransitionTime,
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}

	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta(pod.TypeMeta),
		ObjectMeta: metav1.ObjectMeta{
			Name:                       pod.Name,
			GenerateName:               pod.GenerateName,
			Namespace:                  pod.Namespace,
			ResourceVersion:            pod.ResourceVersion,
			Generation:                 pod.Generation,
			DeletionGracePeriodSeconds: pod.DeletionGracePeriodSeconds,
			Labels:                     pod.Labels,
			Annotations:                pod.Annotations,
			Finalizers:                 pod.Finalizers,
		},
		Spec: corev1.PodSpec{
			Volumes:    volumes,
			Containers: containers,
		},
		Status: corev1.PodStatus{
			Phase:             corev1.PodPhase(pod.Status.Phase),
			Conditions:        conditions,
			Message:           pod.Status.Message,
			Reason:            pod.Status.Reason,
			NominatedNodeName: pod.Status.NominatedNodeName,
			HostIP:            pod.Status.HostIP,
		},
	}
}

func (rc *PodUC) overwritePod(newPod *model.PodsUpdateRequest, existPod *corev1.Pod) *corev1.Pod {
	existPod.Spec.Containers = rc.overwriteContainers(newPod.Pod.Spec.Containers, existPod.Spec.Containers)
	existPod.Spec.InitContainers = rc.overwriteContainers(newPod.Pod.Spec.InitContainers, existPod.Spec.InitContainers)

	existPod.Spec.Tolerations = rc.addTolerations(newPod.Pod.Spec.Tolerations, existPod.Spec.Tolerations)

	existPod.Spec.ActiveDeadlineSeconds = newPod.Pod.Spec.ActiveDeadlineSeconds

	graceSeconds := existPod.Spec.TerminationGracePeriodSeconds
	if graceSeconds == nil { // (allow it to be set to 1 if it was previously negative)
		existPod.Spec.TerminationGracePeriodSeconds = newPod.Pod.Spec.TerminationGracePeriodSeconds
	}

	return existPod
}

// overwriteContainers change only images
func (rc *PodUC) overwriteContainers(newContainers []model.ContainerRequest, existContainers []corev1.Container) []corev1.Container {
	containersMap := make(map[string]model.Container, 0)
	for _, v := range newContainers {
		containersMap[v.Name] = model.Container{
			Name:  v.Name,
			Image: v.Image,
		}
	}

	for i, v := range existContainers {
		matched, ok := containersMap[v.Name]
		if !ok {
			continue
		}

		existContainers[i].Image = matched.Image
		existContainers[i].WorkingDir = matched.WorkingDir
	}

	return existContainers
}

// addToleration only additions to existing tolerations
func (rc *PodUC) addTolerations(newTolerations []model.Toleration, existTolerations []corev1.Toleration) []corev1.Toleration {
	additions := make([]corev1.Toleration, 0)
	for _, v := range newTolerations {
		additions = append(additions, corev1.Toleration{
			Key:               v.Key,
			Operator:          corev1.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            corev1.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		})
	}

	existTolerations = append(existTolerations, additions...)

	return existTolerations
}
