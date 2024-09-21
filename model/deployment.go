package model

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Deployment struct {
	TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Deployment.
	// +optional
	Spec DeploymentSpec `json:"spec,omitempty"`

	// Most recently observed status of the Deployment.
	// +optional
	Status DeploymentStatus `json:"status,omitempty"`
}

// DeploymentSpec is the specification of the desired behavior of the Deployment.
type DeploymentSpec struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	// It must match the pod template's labels.
	Selector *LabelSelector `json:"selector"`

	// Template describes the pods that will be created.
	// The only allowed template.spec.restartPolicy value is "Always".
	Template v1.PodTemplateSpec `json:"template"`

	// The deployment strategy to use to replace existing pods with new ones.
	// +optional
	// s
	Strategy DeploymentStrategy `json:"strategy,omitempty"`

	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty"`

	// The number of old ReplicaSets to retain to allow rollback.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 10.
	// +optional
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`

	// Indicates that the deployment is paused.
	// +optional
	Paused bool `json:"paused,omitempty"`

	// The maximum time in seconds for a deployment to make progress before it
	// is considered to be failed. The deployment controller will continue to
	// process failed deployments and a condition with a ProgressDeadlineExceeded
	// reason will be surfaced in the deployment status. Note that progress will
	// not be estimated during the time a deployment is paused. Defaults to 600s.
	ProgressDeadlineSeconds *int32 `json:"progressDeadlineSeconds,omitempty"`
}

// DeploymentStrategy describes how to replace existing pods with new ones.
type DeploymentStrategy struct {
	// Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.
	// +optional
	Type DeploymentStrategyType `json:"type,omitempty"`

	// Rolling update config params. Present only if DeploymentStrategyType =
	// RollingUpdate.
	//---
	// TODO: Update this to follow our convention for oneOf, whatever we decide it
	// to be.
	// +optional
	RollingUpdate *RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

// +enum
type DeploymentStrategyType string

const (
	// Kill all existing pods before creating new ones.
	RecreateDeploymentStrategyType DeploymentStrategyType = "Recreate"

	// Replace the old ReplicaSets by new one using rolling update i.e gradually scale down the old ReplicaSets and scale up the new one.
	RollingUpdateDeploymentStrategyType DeploymentStrategyType = "RollingUpdate"
)

// Spec to control the desired behavior of rolling update.
type RollingUpdateDeployment struct {
	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// This can not be 0 if MaxSurge is 0.
	// Defaults to 25%.
	// Example: when this is set to 30%, the old ReplicaSet can be scaled down to 70% of desired pods
	// immediately when the rolling update starts. Once new pods are ready, old ReplicaSet
	// can be scaled down further, followed by scaling up the new ReplicaSet, ensuring
	// that the total number of pods available at all times during the update is at
	// least 70% of desired pods.
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// The maximum number of pods that can be scheduled above the desired number of
	// pods.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	// Absolute number is calculated from percentage by rounding up.
	// Defaults to 25%.
	// Example: when this is set to 30%, the new ReplicaSet can be scaled up immediately when
	// the rolling update starts, such that the total number of old and new pods do not exceed
	// 130% of desired pods. Once old pods have been killed,
	// new ReplicaSet can be scaled up further, ensuring that total number of pods running
	// at any time during the update is at most 130% of desired pods.
	// +optional
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`
}

// DeploymentStatus is the most recently observed status of the Deployment.
type DeploymentStatus struct {
	// The generation observed by the deployment controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Total number of non-terminated pods targeted by this deployment (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty"`

	// readyReplicas is the number of pods targeted by this Deployment with a Ready Condition.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`

	// Total number of unavailable pods targeted by this deployment. This is the total number of
	// pods that are still required for the deployment to have 100% available capacity. They may
	// either be pods that are running but not yet available or pods that still have not been created.
	// +optional
	UnavailableReplicas int32 `json:"unavailableReplicas,omitempty"`

	// Represents the latest available observations of a deployment's current state.
	Conditions []DeploymentCondition `json:"conditions,omitempty"`

	// Count of hash collisions for the Deployment. The Deployment controller uses this
	// field as a collision avoidance mechanism when it needs to create the name for the
	// newest ReplicaSet.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty"`
}

type DeploymentConditionType string

// These are valid conditions of a deployment.
const (
	// Available means the deployment is available, ie. at least the minimum available
	// replicas required are up and running for at least minReadySeconds.
	DeploymentAvailable DeploymentConditionType = "Available"
	// Progressing means the deployment is progressing. Progress for a deployment is
	// considered when a new replica set is created or adopted, and when new pods scale
	// up or old pods scale down. Progress is not estimated for paused deployments or
	// when progressDeadlineSeconds is not specified.
	DeploymentProgressing DeploymentConditionType = "Progressing"
	// ReplicaFailure is added in a deployment when one of its pods fails to be created
	// or deleted.
	DeploymentReplicaFailure DeploymentConditionType = "ReplicaFailure"
)

// DeploymentCondition describes the state of a deployment at a certain point.
type DeploymentCondition struct {
	// Type of deployment condition.
	Type DeploymentConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime time.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}
