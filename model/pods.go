package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsRequest struct {
	Pod  Pod                  `json:"pod"`
	Opts metav1.CreateOptions `json:"opts"`
}

type Pod struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`
	Spec       PodSpec   `json:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty"`
}

type PodSpec struct {
	Volumes        []Volume    `json:"volumes,omitempty"`
	InitContainers []Container `json:"initContainers,omitempty"`
	Containers     []Container `json:"containers"`
}
type PodPhase string
type PodStatus struct {
	Phase             PodPhase       `json:"phase,omitempty"`
	Conditions        []PodCondition `json:"conditions,omitempty"`
	Message           string         `json:"message,omitempty"`
	Reason            string         `json:"reason,omitempty"`
	NominatedNodeName string         `json:"nominatedNodeName,omitempty"`
	HostIP            string         `json:"hostIP,omitempty"`
}

type PodConditionType string

type ConditionStatus string

type PodCondition struct {
	Type               PodConditionType `json:"type"`
	Status             ConditionStatus  `json:"status"`
	LastProbeTime      metav1.Time      `json:"lastProbeTime,omitempty"`
	LastTransitionTime metav1.Time      `json:"lastTransitionTime,omitempty"`
	Reason             string           `json:"reason,omitempty"`
	Message            string           `json:"message,omitempty"`
}
