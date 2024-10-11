package model

// MiniDeployment is a Deployment with only the information needed for the UI.
type MiniDeployment struct {
	MiniObjectMeta `json:"metadata,omitempty"`
}

// MiniDeploymentList is a list of Deployments with only the information needed for the UI.
type MiniDeploymentList struct {
	ListMeta `json:"metadata,omitempty"`
	Items    []MiniDeployment `json:"items"`
}

// ConvertMini converts a DeploymentList object into a MiniDeploymentList object.
func (rc *DeploymentList) ConvertMini() MiniDeploymentList {
	return MiniDeploymentList{
		ListMeta: ListMeta(rc.ListMeta),
		Items:    rc.convertDeploymentsToMini(),
	}
}

// convertDeploymentsToMini converts a slice of Deployment objects into a slice of MiniDeployment objects.
func (rc *DeploymentList) convertDeploymentsToMini() []MiniDeployment {
	deployments := make([]MiniDeployment, len(rc.Items))
	for i, deployment := range rc.Items {
		deployments[i] = MiniDeployment{
			MiniObjectMeta: MiniObjectMeta{
				UID:               deployment.UID,
				CreationTimestamp: deployment.CreationTimestamp,
				Name:              deployment.Name,
				GenerateName:      deployment.GenerateName,
				Namespace:         deployment.Namespace,
			},
		}
	}

	return deployments
}
