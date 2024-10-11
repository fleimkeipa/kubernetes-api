package model

type (
	PodsCreateRequest struct {
		Opts CreateOptions `json:"opts"`
		Pod  Pod           `json:"pod"`
	}
)

type (
	PodsUpdateRequest struct {
		Opts UpdateOptions `json:"opts"`
		Pod  PodUpdate     `json:"pod"`
	}

	PodUpdate struct {
		Spec SpecRequest `json:"spec,omitempty"`
	}

	SpecRequest struct {
		InitContainers                []ContainerRequest `json:"initContainers,omitempty"`
		Containers                    []ContainerRequest `json:"containers"`
		ActiveDeadlineSeconds         *int64             `json:"activeDeadlineSeconds,omitempty"`
		TerminationGracePeriodSeconds *int64             `json:"terminationGracePeriodSeconds,omitempty"`
		Tolerations                   []Toleration       `json:"tolerations,omitempty"` // allow it to be set to 1 if it was previously negative
	}

	ContainerRequest struct {
		Name  string `json:"name"` // cannot changable
		Image string `json:"image,omitempty"`
	}
)
