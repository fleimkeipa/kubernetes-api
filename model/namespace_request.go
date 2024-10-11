package model

type NamespaceCreateRequest struct {
	Opts      CreateOptions `json:"opts"`
	Namespace Namespace     `json:"namespace"`
}

type (
	NamespaceUpdateRequest struct {
		Opts      UpdateOptions   `json:"opts"`
		Namespace NamespaceUpdate `json:"namespace"`
	}

	NamespaceUpdate struct {
		NamespaceObjectMetaUpdateRequest `json:"metadata,omitempty"`
		Spec                             NamespaceSpec `json:"spec,omitempty"`
	}
	NamespaceObjectMetaUpdateRequest struct {
		Labels      map[string]string `json:"labels,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
	}
)
