package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type PodRepository struct {
	client *kubernetes.Clientset
}

func NewPodRepository(client *kubernetes.Clientset) *PodRepository {
	return &PodRepository{client}
}

func (rc *PodRepository) Create(ctx context.Context, pod *model.Pod, opts model.CreateOptions) (*model.Pod, error) {
	createOptions := convertCreateOptsToKube(opts)

	kubePod := rc.fillRequestPod(pod)

	createdPod, err := rc.client.CoreV1().Pods(pod.Namespace).Create(ctx, kubePod, createOptions)
	if err != nil {
		return nil, err
	}

	return rc.fillResponsePod(createdPod), nil
}

func (rc *PodRepository) Update(ctx context.Context, podID string, pod *model.Pod, opts model.UpdateOptions) (*model.Pod, error) {
	updateOptions := convertUpdateOptsToKube(opts)

	existPod, err := rc.getByNameOrUID(ctx, pod.Namespace, podID, model.ListOptions{})
	if err != nil {
		return nil, err
	}

	kubePod := rc.overwriteOnKubePod(pod, existPod)

	updatedPod, err := rc.client.CoreV1().Pods(pod.Namespace).Update(ctx, kubePod, updateOptions)
	if err != nil {
		return nil, err
	}

	return rc.fillResponsePod(updatedPod), nil
}

func (rc *PodRepository) List(ctx context.Context, namespace string, opts model.ListOptions) (*model.PodList, error) {
	kubePods, err := rc.list(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}

	podList := model.PodList{}
	for _, kubePod := range kubePods.Items {
		podList.Items = append(podList.Items, *rc.fillResponsePod(&kubePod))
	}

	podList.ListMeta = model.ListMeta{
		RemainingItemCount: kubePods.ListMeta.RemainingItemCount,
		ResourceVersion:    kubePods.ListMeta.ResourceVersion,
		Continue:           kubePods.ListMeta.Continue,
	}
	podList.TypeMeta = model.TypeMeta(kubePods.TypeMeta)

	return &podList, nil
}

func (rc *PodRepository) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	deleteOpts := convertDeleteOptsToKube(opts)

	return rc.client.CoreV1().Pods(namespace).Delete(ctx, name, deleteOpts)
}

func (rc *PodRepository) GetByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*model.Pod, error) {
	pod, err := rc.getByNameOrUID(ctx, namespace, nameOrUID, opts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponsePod(pod), nil
}

func (rc *PodRepository) getByNameOrUID(ctx context.Context, namespace, nameOrUID string, opts model.ListOptions) (*corev1.Pod, error) {
	opts.TypeMeta.Kind = "pod"
	if namespace == "" {
		namespace = "default"
	}

	opts.Limit = 100
	pods, err := rc.list(ctx, namespace, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range pods.Items {
		if v.Name == nameOrUID || v.UID == types.UID(nameOrUID) {
			return &v, nil
		}
	}

	if pods.ListMeta.Continue == "" {
		return nil, fmt.Errorf("pod %s not found", nameOrUID)
	}

	opts.Continue = pods.ListMeta.Continue
	return rc.getByNameOrUID(ctx, namespace, nameOrUID, opts)
}

func (rc *PodRepository) list(ctx context.Context, namespace string, opts model.ListOptions) (*corev1.PodList, error) {
	listOpts := convertListOptsToKube(opts)

	return rc.client.CoreV1().Pods(namespace).List(ctx, listOpts)
}

func (rc *PodRepository) overwriteOnKubePod(newPod *model.Pod, existPod *corev1.Pod) *corev1.Pod {
	existPod.Spec.Containers = overwriteKubeContainers(newPod.Spec.Containers, existPod.Spec.Containers)
	existPod.Spec.InitContainers = overwriteKubeContainers(newPod.Spec.InitContainers, existPod.Spec.InitContainers)

	existPod.Spec.Tolerations = addKubeTolerations(newPod.Spec.Tolerations, existPod.Spec.Tolerations)

	existPod.Spec.ActiveDeadlineSeconds = newPod.Spec.ActiveDeadlineSeconds

	graceSeconds := existPod.Spec.TerminationGracePeriodSeconds
	if graceSeconds == nil { // (allow it to be set to 1 if it was previously negative)
		existPod.Spec.TerminationGracePeriodSeconds = newPod.Spec.TerminationGracePeriodSeconds
	}

	return existPod
}

func (rc *PodRepository) fillRequestPod(pod *model.Pod) *corev1.Pod {
	containers := convertContainersToKube(pod.Spec.Containers)

	initContainers := convertContainersToKube(pod.Spec.InitContainers)

	tolerations := convertTolerationsToKube(pod.Spec.Tolerations)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.Name,
		},
		Spec: corev1.PodSpec{
			InitContainers:        initContainers,
			Containers:            containers,
			Tolerations:           tolerations,
			ActiveDeadlineSeconds: pod.Spec.ActiveDeadlineSeconds,
		},
	}
}

func (rc *PodRepository) fillResponsePod(pod *corev1.Pod) *model.Pod {
	volumes := make([]model.Volume, 0, len(pod.Spec.Volumes))
	for _, v := range pod.Spec.Volumes {
		volumes = append(volumes, model.Volume{
			Name:         v.Name,
			VolumeSource: model.VolumeSource{},
		})
	}

	containers := make([]model.Container, 0, len(pod.Spec.Containers))
	for _, v := range pod.Spec.Containers {
		containers = append(containers, model.Container{
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

	initContainers := make([]model.Container, 0, len(pod.Spec.InitContainers))
	for _, v := range pod.Spec.InitContainers {
		initContainers = append(initContainers, model.Container{
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

	conditions := make([]model.PodCondition, 0, len(pod.Status.Conditions))
	for _, v := range pod.Status.Conditions {
		conditions = append(conditions, model.PodCondition{
			Type:               model.PodConditionType(v.Type),
			Status:             model.ConditionStatus(v.Status),
			LastProbeTime:      v.LastProbeTime.Time,
			LastTransitionTime: v.LastTransitionTime.Time,
			Reason:             v.Reason,
			Message:            v.Message,
		})
	}

	tolerations := make([]model.Toleration, 0, len(pod.Spec.Tolerations))
	for _, v := range pod.Spec.Tolerations {
		tolerations = append(tolerations, model.Toleration{
			Key:               v.Key,
			Operator:          model.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            model.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		})
	}

	ownerReferences := make([]model.OwnerReference, 0, len(pod.OwnerReferences))
	for _, v := range pod.OwnerReferences {
		ownerReferences = append(ownerReferences, model.OwnerReference{
			Controller:         v.Controller,
			BlockOwnerDeletion: v.BlockOwnerDeletion,
			APIVersion:         v.APIVersion,
			Kind:               v.Kind,
			Name:               v.Name,
		})
	}

	deletionTimestamp := new(time.Time)
	if deletionTime := pod.DeletionTimestamp; deletionTime != nil {
		deletionTimestamp = &deletionTime.Time
	}

	return &model.Pod{
		TypeMeta: model.TypeMeta(pod.TypeMeta),
		ObjectMeta: model.ObjectMeta{
			UID:                        string(pod.UID),
			CreationTimestamp:          pod.CreationTimestamp.Time,
			DeletionTimestamp:          deletionTimestamp,
			DeletionGracePeriodSeconds: pod.DeletionGracePeriodSeconds,
			Labels:                     pod.Labels,
			Annotations:                pod.Annotations,
			Name:                       pod.Name,
			GenerateName:               pod.GenerateName,
			Namespace:                  pod.Namespace,
			ResourceVersion:            pod.ResourceVersion,
			OwnerReferences:            ownerReferences,
			Finalizers:                 pod.Finalizers,
			Generation:                 pod.Generation,
		},
		Spec: model.PodSpec{
			Volumes:                       volumes,
			InitContainers:                initContainers,
			Containers:                    containers,
			ActiveDeadlineSeconds:         pod.DeletionGracePeriodSeconds,
			TerminationGracePeriodSeconds: pod.DeletionGracePeriodSeconds,
			Tolerations:                   tolerations,
		},
		Status: model.PodStatus{
			Phase:             model.PodPhase(pod.Status.Phase),
			Conditions:        conditions,
			Message:           pod.Status.Message,
			Reason:            pod.Status.Reason,
			NominatedNodeName: pod.Status.NominatedNodeName,
			HostIP:            pod.Status.HostIP,
		},
	}
}
