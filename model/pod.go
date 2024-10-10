package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	TypeMeta   `json:",inline"`
	Spec       PodSpec   `json:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty"`
	ObjectMeta `json:"metadata,omitempty"`
}

// PodList is a list of Pods.
type PodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pod `json:"items"`
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
		LastProbeTime      metav1.Time      `json:"lastProbeTime,omitempty"`
		LastTransitionTime metav1.Time      `json:"lastTransitionTime,omitempty"`
		Reason             string           `json:"reason,omitempty"`
		Message            string           `json:"message,omitempty"`
	}
)
