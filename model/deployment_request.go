package model

type DeploymentCreateRequest struct {
	Opts       CreateOptions `json:"opts"`
	Deployment Deployment    `json:"deployment"`
}

type DeploymentUpdateRequest struct {
	Opts       UpdateOptions `json:"opts"`
	Deployment Deployment    `json:"deployment"`
}
