package repositories

import (
	"github.com/fleimkeipa/kubernetes-api/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// overwriteKubeContainers change only images
func overwriteKubeContainers(newContainers []model.Container, existContainers []corev1.Container) []corev1.Container {
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

// addKubeTolerations only additions to existing tolerations
func addKubeTolerations(newTolerations []model.Toleration, existTolerations []corev1.Toleration) []corev1.Toleration {
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

func convertTolerationsToKube(tolerations []model.Toleration) []corev1.Toleration {
	newTolerations := make([]corev1.Toleration, 0, len(tolerations))
	for _, v := range tolerations {
		newTolerations = append(newTolerations, corev1.Toleration{
			Key:               v.Key,
			Operator:          corev1.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            corev1.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		})
	}
	return newTolerations
}

func convertPodConditionsToKube(conditions []model.PodCondition) []corev1.PodCondition {
	kubeConditions := make([]corev1.PodCondition, 0, len(conditions))
	for _, condition := range conditions {
		kubeConditions = append(kubeConditions, corev1.PodCondition{
			Type:               corev1.PodConditionType(condition.Type),
			Status:             corev1.ConditionStatus(condition.Status),
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}
	return kubeConditions
}

func convertContainersToKube(containers []model.Container) []corev1.Container {
	newContainers := make([]corev1.Container, 0, len(containers))
	for _, v := range containers {
		newContainers = append(newContainers, corev1.Container{
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
	return newContainers
}

/*---------------- Kube Options Conversions ----------------*/

func convertCreateOptsToKube(opts model.CreateOptions) metav1.CreateOptions {
	createOptions := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return createOptions
}

func convertUpdateOptsToKube(opts model.UpdateOptions) metav1.UpdateOptions {
	updateOptions := metav1.UpdateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		DryRun:          opts.DryRun,
		FieldManager:    opts.FieldManager,
		FieldValidation: opts.FieldValidation,
	}
	return updateOptions
}

func convertDeleteOptsToKube(opts model.DeleteOptions) metav1.DeleteOptions {
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

	return metaOpts
}

func convertListOptsToKube(opts model.ListOptions) metav1.ListOptions {
	listOpts := metav1.ListOptions{
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
	return listOpts
}
