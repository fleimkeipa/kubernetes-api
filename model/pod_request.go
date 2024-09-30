package model

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type (
	PodsCreateRequest struct {
		Opts metav1.CreateOptions `json:"opts"`
		Pod  Pod                  `json:"pod"`
	}
)

type (
	PodsUpdateRequest struct {
		Opts metav1.UpdateOptions `json:"opts"`
		Pod  PodUpdate            `json:"pod"`
	}

	PodUpdate struct {
		ID        string
		Namespace string
		Spec      SpecRequest `json:"spec,omitempty"`
	}

	SpecRequest struct {
		InitContainers                []ContainerRequest `json:"initContainers,omitempty"`
		Containers                    []ContainerRequest `json:"containers"`
		ActiveDeadlineSeconds         *int64             `json:"activeDeadlineSeconds,omitempty"`
		TerminationGracePeriodSeconds *int64             `json:"terminationGracePeriodSeconds,omitempty"`
		Tolerations                   []Toleration       `json:"tolerations,omitempty"`
	}

	ContainerRequest struct {
		Name  string `json:"name"`
		Image string `json:"image,omitempty"`
	}
)
