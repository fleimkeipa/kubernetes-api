package model

type DeploymentCreateRequest struct {
	Opts       CreateOptions `json:"opts"`
	Deployment Deployment    `json:"deployment"`
}

type (
	DeploymentUpdateRequest struct {
		Opts       UpdateOptions    `json:"opts"`
		Deployment DeploymentUpdate `json:"deployment"`
	}

	DeploymentUpdate struct {
		DeploymentObjectMetaUpdateRequest `json:"metadata,omitempty"`
		Spec                              DeploymentSpecUpdateRequest `json:"spec,omitempty"`
	}

	DeploymentObjectMetaUpdateRequest struct {
		Labels      map[string]string `json:"labels,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}

	DeploymentSpecUpdateRequest struct {
		Replicas                *int32             `json:"replicas,omitempty"`
		Template                PodTemplateSpec    `json:"template"`
		Strategy                DeploymentStrategy `json:"strategy,omitempty"`
		MinReadySeconds         int32              `json:"minReadySeconds,omitempty"`
		ProgressDeadlineSeconds *int32             `json:"progressDeadlineSeconds,omitempty"`
	}
)
