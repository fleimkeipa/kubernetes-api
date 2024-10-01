package model

type DeploymentCreateRequest struct {
	Opts       CreateOptions
	Deployment Deployment
}

type DeploymentUpdateRequest struct {
	Opts       UpdateOptions
	Deployment Deployment
}
