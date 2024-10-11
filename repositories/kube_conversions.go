package repositories

import (
	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/google/uuid"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func convertTolerationsToModel(tolerations []corev1.Toleration) []model.Toleration {
	newTolerations := make([]model.Toleration, 0, len(tolerations))
	for _, v := range tolerations {
		newTolerations = append(newTolerations, model.Toleration{
			Key:               v.Key,
			Operator:          model.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            model.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		})
	}
	return newTolerations
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

func convertContainersToModel(containers []corev1.Container) []model.Container {
	newContainers := make([]model.Container, 0, len(containers))
	for _, v := range containers {
		newContainers = append(newContainers, model.Container{
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

func convertTemplateToKube(template *model.PodTemplateSpec) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:        template.ObjectMeta.Name,
			Labels:      template.ObjectMeta.Labels,
			Annotations: template.ObjectMeta.Annotations,
		},
		Spec: corev1.PodSpec{
			Containers:            convertContainersToKube(template.Spec.Containers),
			InitContainers:        convertContainersToKube(template.Spec.InitContainers),
			Tolerations:           convertTolerationsToKube(template.Spec.Tolerations),
			ActiveDeadlineSeconds: template.Spec.ActiveDeadlineSeconds,
		},
	}
}

func convertTemplateToModel(deployment *v1.Deployment) model.PodTemplateSpec {
	template := model.PodTemplateSpec{
		ObjectMeta: model.ObjectMeta{
			Name:        deployment.Spec.Template.ObjectMeta.Name,
			Labels:      deployment.Spec.Template.ObjectMeta.Labels,
			Annotations: deployment.Spec.Template.ObjectMeta.Annotations,
		},
		Spec: model.PodSpec{
			Containers:            convertContainersToModel(deployment.Spec.Template.Spec.Containers),
			InitContainers:        convertContainersToModel(deployment.Spec.Template.Spec.InitContainers),
			Tolerations:           convertTolerationsToModel(deployment.Spec.Template.Spec.Tolerations),
			ActiveDeadlineSeconds: deployment.Spec.Template.Spec.ActiveDeadlineSeconds,
		},
	}
	return template
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

	metaOpts := metav1.DeleteOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       opts.TypeMeta.Kind,
			APIVersion: opts.TypeMeta.APIVersion,
		},
		GracePeriodSeconds: gracePeriodSeconds,
		DryRun:             opts.DryRun,
	}

	if opts.Preconditions == nil || opts.Preconditions.UID == nil {
		return metaOpts
	}

	uidStr := *opts.Preconditions.UID
	if err := uuid.Validate(string(uidStr)); err != nil {
		return metaOpts
	}

	opts.Preconditions.UID = &uidStr
	uidTypes := types.UID(uidStr)
	metaOpts.Preconditions = &metav1.Preconditions{
		UID: &uidTypes,
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
