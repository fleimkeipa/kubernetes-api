package model

type NamespaceCreateRequest struct {
	Opts      CreateOptions `json:"opts"`
	Namespace Namespace     `json:"namespace"`
}

type NamespaceUpdateRequest struct {
	Opts      UpdateOptions `json:"opts"`
	Namespace Namespace     `json:"namespace"`
}
