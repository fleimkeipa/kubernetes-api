package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodsRepository struct {
	client *kubernetes.Clientset
}

func NewPodsRepository(client *kubernetes.Clientset) *PodsRepository {
	return &PodsRepository{client}
}

func (rc *PodsRepository) Create(ctx context.Context, pod *model.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	var kubePod = fillPod(pod)
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, kubePod, opts)
}

func (rc *PodsRepository) Get(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	return rc.client.CoreV1().Pods(namespace).List(ctx, opts)
}

func (rc *PodsRepository) Delete(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return rc.client.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (rc *PodsRepository) Update(ctx context.Context, pod *model.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	var kubePod = fillPod(pod)
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Update(ctx, kubePod, opts)
}

func fillPod(pod *model.Pod) *corev1.Pod {
	var volumes = make([]corev1.Volume, 0)
	for _, v := range pod.Spec.Volumes {
		volumes = append(volumes, corev1.Volume{
			Name:         v.Name,
			VolumeSource: corev1.VolumeSource{},
		})
	}

	var containers = make([]corev1.Container, 0)
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

	var conditions = make([]corev1.PodCondition, 0)
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
