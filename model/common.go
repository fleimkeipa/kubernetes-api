package model

import (
	"time"
)

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ObjectMeta struct {
	Name                       string            `json:"name,omitempty"`
	GenerateName               string            `json:"generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	Generation                 int64             `json:"generation,omitempty"`
	CreationTimestamp          time.Time         `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          *time.Time        `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds *int64            `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty"`
	OwnerReferences            []OwnerReference  `json:"ownerReferences,omitempty" patchStrategy:"merge" patchMergeKey:"uid"`
	Finalizers                 []string          `json:"finalizers,omitempty" patchStrategy:"merge"`
}

type OwnerReference struct {
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	Controller         *bool  `json:"controller,omitempty"`
	BlockOwnerDeletion *bool  `json:"blockOwnerDeletion,omitempty"`
}

type Container struct {
	Name                   string          `json:"name"`
	Image                  string          `json:"image,omitempty"`
	Command                []string        `json:"command,omitempty"`
	Args                   []string        `json:"args,omitempty"`
	WorkingDir             string          `json:"workingDir,omitempty"`
	Ports                  []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort"`
	Env                    []EnvVar        `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
	TerminationMessagePath string          `json:"terminationMessagePath,omitempty"`
	Stdin                  bool            `json:"stdin,omitempty"`
	StdinOnce              bool            `json:"stdinOnce,omitempty"`
	TTY                    bool            `json:"tty,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type Protocol string

type ContainerPort struct {
	Name          string   `json:"name,omitempty"`
	HostPort      int32    `json:"hostPort,omitempty"`
	ContainerPort int32    `json:"containerPort"`
	Protocol      Protocol `json:"protocol,omitempty"`
	HostIP        string   `json:"hostIP,omitempty"`
}
type Volume struct {
	Name         string `json:"name"`
	VolumeSource `json:",inline"`
}

type VolumeSource struct {
}

type TaintEffect string

type TolerationOperator string

type Toleration struct {
	Key               string             `json:"key,omitempty"`
	Operator          TolerationOperator `json:"operator,omitempty"`
	Value             string             `json:"value,omitempty"`
	Effect            TaintEffect        `json:"effect,omitempty"`
	TolerationSeconds *int64             `json:"tolerationSeconds,omitempty"`
}
