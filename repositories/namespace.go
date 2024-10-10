package repositories

import (
	"context"

	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func (rc *NamespaceRepository) Create(ctx context.Context, namespace *model.Namespace, opts model.CreateOptions) (*model.Namespace, error) {
	metaOpts := convertCreateOptsToKube(opts)

	kubeNamespace := rc.fillRequestNamespace(namespace)

	createdNamespace, err := rc.client.CoreV1().Namespaces().Create(ctx, kubeNamespace, metaOpts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseNamespace(createdNamespace), nil
}

func (rc *NamespaceRepository) Update(ctx context.Context, namespace *model.Namespace, opts model.UpdateOptions) (*model.Namespace, error) {
	metaOpts := convertUpdateOptsToKube(opts)

	existNamespace, err := rc.getByNameOrUID(ctx, namespace.Name, model.ListOptions{})
	if err != nil {
		return nil, err
	}

	kubeNamespace := rc.overwriteOnKubeNamespace(namespace, existNamespace)

	createdNamespace, err := rc.client.CoreV1().Namespaces().Update(ctx, kubeNamespace, metaOpts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseNamespace(createdNamespace), nil
}

func (rc *NamespaceRepository) List(ctx context.Context, opts model.ListOptions) (*model.NamespaceList, error) {
	kubeNamespaces, err := rc.list(ctx, opts)
	if err != nil {
		return nil, err
	}

	namespaceList := model.NamespaceList{}
	for _, kubeNamespace := range kubeNamespaces.Items {
		namespaceList.Items = append(namespaceList.Items, *rc.fillResponseNamespace(&kubeNamespace))
	}

	namespaceList.ListMeta = model.ListMeta(kubeNamespaces.ListMeta)
	namespaceList.TypeMeta = model.TypeMeta(kubeNamespaces.TypeMeta)

	return &namespaceList, nil
}

func (rc *NamespaceRepository) GetByNameOrUID(ctx context.Context, nameOrUID string, opts model.ListOptions) (*model.Namespace, error) {
	namespace, err := rc.getByNameOrUID(ctx, nameOrUID, opts)
	if err != nil {
		return nil, err
	}

	return rc.fillResponseNamespace(namespace), nil
}

func (rc *NamespaceRepository) Delete(ctx context.Context, name string, opts model.DeleteOptions) error {
	metaOpts := convertDeleteOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().Delete(ctx, name, metaOpts)
}

func (rc *NamespaceRepository) list(ctx context.Context, opts model.ListOptions) (*corev1.NamespaceList, error) {
	metaOpts := convertListOptsToKube(opts)

	return rc.client.CoreV1().Namespaces().List(ctx, metaOpts)
}

func (rc *NamespaceRepository) getByNameOrUID(ctx context.Context, nameOrUID string, opts model.ListOptions) (*corev1.Namespace, error) {
	opts.TypeMeta.Kind = "namespace"

	opts.Limit = 100
	namespaces, err := rc.list(ctx, opts)
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
	return rc.getByNameOrUID(ctx, nameOrUID, opts)
}

func (rc *NamespaceRepository) fillRequestNamespace(namespace *model.Namespace) *corev1.Namespace {
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

func (rc *NamespaceRepository) fillResponseNamespace(namespace *corev1.Namespace) *model.Namespace {
	conditions := make([]model.NamespaceCondition, 0)
	for _, v := range namespace.Status.Conditions {
		conditions = append(conditions, model.NamespaceCondition{
			Type:   model.NamespaceConditionType(v.Type),
			Status: model.ConditionStatus(v.Status),
			// LastTransitionTime: metav1.Time{Time: v.LastTransitionTime},
			Reason:  v.Reason,
			Message: v.Message,
		})
	}

	finalizers := make([]model.FinalizerName, 0)
	for _, v := range namespace.Finalizers {
		finalizers = append(finalizers, model.FinalizerName(v))
	}

	return &model.Namespace{
		TypeMeta: model.TypeMeta(namespace.TypeMeta),
		ObjectMeta: model.ObjectMeta{
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
		Spec: model.NamespaceSpec{
			Finalizers: finalizers,
		},
		Status: model.NamespaceStatus{
			Phase:      model.NamespacePhase(namespace.Status.Phase),
			Conditions: conditions,
		},
	}
}

func (rc *NamespaceRepository) overwriteOnKubeNamespace(newNamespace *model.Namespace, existNamespace *corev1.Namespace) *corev1.Namespace {
	existNamespace.Name = newNamespace.Namespace

	existNamespace.Kind = "namespace"
	existNamespace.APIVersion = "v1"

	return existNamespace
}
