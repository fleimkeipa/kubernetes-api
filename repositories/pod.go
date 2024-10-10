package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodRepository struct {
	client *kubernetes.Clientset
}

func NewPodRepository(client *kubernetes.Clientset) *PodRepository {
	return &PodRepository{client}
}

func (rc *PodRepository) Create(ctx context.Context, pod *corev1.Pod, opts model.CreateOptions) (*corev1.Pod, error) {
	metaOpts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, pod, metaOpts)
}

func (rc *PodRepository) Update(ctx context.Context, id string, pod *corev1.Pod, opts model.UpdateOptions) (*corev1.Pod, error) {
	metaOpts := metav1.UpdateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.CoreV1().Pods(pod.ObjectMeta.Namespace).Update(ctx, pod, metaOpts)
}

func (rc *PodRepository) List(ctx context.Context, namespace string, opts model.ListOptions) (*corev1.PodList, error) {
	metaOpts := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		LabelSelector:        opts.LabelSelector,
		FieldSelector:        opts.FieldSelector,
		Watch:                opts.Watch,
		AllowWatchBookmarks:  opts.AllowWatchBookmarks,
		ResourceVersion:      opts.ResourceVersion,
		ResourceVersionMatch: metav1.ResourceVersionMatch(opts.ResourceVersionMatch),
		TimeoutSeconds:       opts.TimeoutSeconds,
		Limit:                opts.Limit,
		Continue:             opts.Continue,
		SendInitialEvents:    opts.SendInitialEvents,
	}
	return rc.client.CoreV1().Pods(namespace).List(ctx, metaOpts)
}

func (rc *PodRepository) Delete(ctx context.Context, namespace, name string, opts model.DeleteOptions) error {
	gracePeriodSeconds := new(int64)
	if opts.GracePeriodSeconds != nil {
		gracePeriodSeconds = opts.GracePeriodSeconds
	}

	preconditions := new(model.Preconditions)
	if opts.Preconditions != nil {
		preconditions = opts.Preconditions
	}

	metaOpts := metav1.DeleteOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		GracePeriodSeconds: gracePeriodSeconds,
		Preconditions: &metav1.Preconditions{
			UID:             preconditions.UID,
			ResourceVersion: preconditions.ResourceVersion,
		},
		DryRun: opts.DryRun,
	}

	return rc.client.CoreV1().Pods(namespace).Delete(ctx, name, metaOpts)
}
