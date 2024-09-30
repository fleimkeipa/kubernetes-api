package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	Spec       PodSpec   `json:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty"`
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
		Conditions        []PodCondition `json:"conditions,omitempty"`
		Message           string         `json:"message,omitempty"`
		Reason            string         `json:"reason,omitempty"`
		NominatedNodeName string         `json:"nominatedNodeName,omitempty"`
		HostIP            string         `json:"hostIP,omitempty"`
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
