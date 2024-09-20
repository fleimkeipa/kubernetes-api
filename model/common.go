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
	OwnerReferences            []OwnerReference  `json:"ownerReferences,omitempty"`
	Finalizers                 []string          `json:"finalizers,omitempty"`
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
	Ports                  []ContainerPort `json:"ports,omitempty"`
	Env                    []EnvVar        `json:"env,omitempty"`
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

// Note:
// There are two different styles of label selectors used in versioned types:
// an older style which is represented as just a string in versioned types, and a
// newer style that is structured.  LabelSelector is an internal representation for the
// latter style.

// A label selector is a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed. An empty label selector matches all objects. A null
// label selector matches no objects.
// +structType=atomic
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	// +listType=atomic
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty"`
}

// A label selector requirement is a selector that contains values, a key, and an operator that
// relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	Key string `json:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic
	// merge patch.
	// +optional
	// +listType=atomic
	Values []string `json:"values,omitempty"`
}

// A label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

const (
	LabelSelectorOpIn           LabelSelectorOperator = "In"
	LabelSelectorOpNotIn        LabelSelectorOperator = "NotIn"
	LabelSelectorOpExists       LabelSelectorOperator = "Exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)
