package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceRepository struct {
	client *kubernetes.Clientset
}

func NewNamespaceRepository(client *kubernetes.Clientset) *NamespaceRepository {
	return &NamespaceRepository{
		client: client,
	}
}

func (rc *NamespaceRepository) Create(ctx context.Context, namespace *corev1.Namespace, opts model.CreateOptions) (*corev1.Namespace, error) {
	metaOpts := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.CoreV1().Namespaces().Create(ctx, namespace, metaOpts)
}

func (rc *NamespaceRepository) Update(ctx context.Context, namespace *corev1.Namespace, opts model.UpdateOptions) (*corev1.Namespace, error) {
	metaOpts := metav1.UpdateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return rc.client.CoreV1().Namespaces().Update(ctx, namespace, metaOpts)
}

func (rc *NamespaceRepository) List(ctx context.Context, opts model.ListOptions) (*corev1.NamespaceList, error) {
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
	return rc.client.CoreV1().Namespaces().List(ctx, metaOpts)
}

func (rc *NamespaceRepository) Delete(ctx context.Context, name string, opts model.DeleteOptions) error {
	metaOpts := metav1.DeleteOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		GracePeriodSeconds: opts.GracePeriodSeconds,
		Preconditions: &metav1.Preconditions{
			UID:             opts.Preconditions.UID,
			ResourceVersion: opts.Preconditions.ResourceVersion,
		},
		OrphanDependents: opts.OrphanDependents,
		DryRun:           opts.DryRun,
	}
	return rc.client.CoreV1().Namespaces().Delete(ctx, name, metaOpts)
}
