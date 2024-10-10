package model

import (
	"time"
)

type Pod struct {
	TypeMeta   `json:",inline"`
	Spec       PodSpec   `json:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty"`
	ObjectMeta `json:"metadata,omitempty"`
}

// PodList is a list of Pods.
type PodList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`
	Items    []Pod `json:"items"`
}

type PodSpec struct {
	Volumes                       []Volume     `json:"volumes,omitempty"`
	InitContainers                []Container  `json:"initContainers,omitempty"`
	Containers                    []Container  `json:"containers"`
	ActiveDeadlineSeconds         *int64       `json:"activeDeadlineSeconds,omitempty"`
	TerminationGracePeriodSeconds *int64       `json:"terminationGracePeriodSeconds,omitempty"`
	Tolerations                   []Toleration `json:"tolerations,omitempty"`
}

type ConditionStatus string

type (
	PodPhase string

	PodStatus struct {
		Phase             PodPhase       `json:"phase,omitempty"`
		Message           string         `json:"message,omitempty"`
		Reason            string         `json:"reason,omitempty"`
		NominatedNodeName string         `json:"nominatedNodeName,omitempty"`
		HostIP            string         `json:"hostIP,omitempty"`
		Conditions        []PodCondition `json:"conditions,omitempty"`
	}
)

type (
	PodConditionType string

	PodCondition struct {
		Type               PodConditionType `json:"type"`
		Status             ConditionStatus  `json:"status"`
		LastProbeTime      time.Time        `json:"lastProbeTime,omitempty"`
		LastTransitionTime time.Time        `json:"lastTransitionTime,omitempty"`
		Reason             string           `json:"reason,omitempty"`
		Message            string           `json:"message,omitempty"`
	}
)
