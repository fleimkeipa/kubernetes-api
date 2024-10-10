package model

import (
	"time"

	"k8s.io/apimachinery/pkg/types"
)

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

// ListMeta describes metadata that synthetic resources must have, including lists and
// various status objects. A resource may have only one of {ObjectMeta, ListMeta}.
type ListMeta struct {
	SelfLink           string `json:"selfLink,omitempty"`
	ResourceVersion    string `json:"resourceVersion,omitempty"`
	Continue           string `json:"continue,omitempty"`
	RemainingItemCount *int64 `json:"remainingItemCount,omitempty"`
}

type ObjectMeta struct {
	UID                        string            `json:"uid,omitempty"`
	CreationTimestamp          time.Time         `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          *time.Time        `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds *int64            `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty"`
	Name                       string            `json:"name,omitempty"`
	GenerateName               string            `json:"generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	OwnerReferences            []OwnerReference  `json:"ownerReferences,omitempty"`
	Finalizers                 []string          `json:"finalizers,omitempty"`
	Generation                 int64             `json:"generation,omitempty"`
}

type OwnerReference struct {
	Controller         *bool  `json:"controller,omitempty"`
	BlockOwnerDeletion *bool  `json:"blockOwnerDeletion,omitempty"`
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
}

type Container struct {
	Name                   string          `json:"name"`
	Image                  string          `json:"image,omitempty"`
	WorkingDir             string          `json:"workingDir,omitempty"`
	TerminationMessagePath string          `json:"terminationMessagePath,omitempty"`
	Command                []string        `json:"command,omitempty"`
	Args                   []string        `json:"args,omitempty"`
	Ports                  []ContainerPort `json:"ports,omitempty"`
	Env                    []EnvVar        `json:"env,omitempty"`
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
	Protocol      Protocol `json:"protocol,omitempty"`
	HostIP        string   `json:"hostIP,omitempty"`
	HostPort      int32    `json:"hostPort,omitempty"`
	ContainerPort int32    `json:"containerPort"`
}
type Volume struct {
	VolumeSource `json:",inline"`
	Name         string `json:"name"`
}

type VolumeSource struct{}

type TaintEffect string

type TolerationOperator string

type Toleration struct {
	TolerationSeconds *int64             `json:"tolerationSeconds,omitempty"`
	Key               string             `json:"key,omitempty"`
	Operator          TolerationOperator `json:"operator,omitempty"`
	Value             string             `json:"value,omitempty"`
	Effect            TaintEffect        `json:"effect,omitempty"`
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

// CreateOptions may be provided when creating an API object.
type CreateOptions struct {
	TypeMeta        `json:",inline"`
	FieldManager    string   `json:"fieldManager,omitempty"`
	FieldValidation string   `json:"fieldValidation,omitempty"`
	DryRun          []string `json:"dryRun,omitempty"`
}

// UpdateOptions may be provided when updating an API object.
// All fields in UpdateOptions should also be present in PatchOptions.
type UpdateOptions struct {
	TypeMeta        `json:",inline"`
	FieldManager    string   `json:"fieldManager,omitempty"`
	FieldValidation string   `json:"fieldValidation,omitempty"`
	DryRun          []string `json:"dryRun,omitempty"`
}

// resourceVersionMatch specifies how the resourceVersion parameter is applied. resourceVersionMatch
// may only be set if resourceVersion is also set.
//
// "NotOlderThan" matches data at least as new as the provided resourceVersion.
// "Exact" matches data at the exact resourceVersion provided.
//
// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
// details.
type ResourceVersionMatch string

const (
	// ResourceVersionMatchNotOlderThan matches data at least as new as the provided
	// resourceVersion.
	ResourceVersionMatchNotOlderThan ResourceVersionMatch = "NotOlderThan"
	// ResourceVersionMatchExact matches data at the exact resourceVersion
	// provided.
	ResourceVersionMatchExact ResourceVersionMatch = "Exact"
)

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	TimeoutSeconds       *int64 `json:"timeoutSeconds,omitempty"`
	SendInitialEvents    *bool  `json:"sendInitialEvents,omitempty"`
	TypeMeta             `json:",inline"`
	LabelSelector        string               `json:"labelSelector,omitempty"`
	FieldSelector        string               `json:"fieldSelector,omitempty"`
	ResourceVersion      string               `json:"resourceVersion,omitempty"`
	ResourceVersionMatch ResourceVersionMatch `json:"resourceVersionMatch,omitempty"`
	Continue             string               `json:"continue,omitempty"`
	Limit                int64                `json:"limit,omitempty"`
	Watch                bool                 `json:"watch,omitempty"`
	AllowWatchBookmarks  bool                 `json:"allowWatchBookmarks,omitempty"`
}

// DeletionPropagation decides if a deletion will propagate to the dependents of
// the object, and how the garbage collector will handle the propagation.
type DeletionPropagation string

const (
	// Orphans the dependents.
	DeletePropagationOrphan DeletionPropagation = "Orphan"
	// Deletes the object from the key-value store, the garbage collector will
	// delete the dependents in the background.
	DeletePropagationBackground DeletionPropagation = "Background"
	// The object exists in the key-value store until the garbage collector
	// deletes all the dependents whose ownerReference.blockOwnerDeletion=true
	// from the key-value store.  API sever will put the "foregroundDeletion"
	// finalizer on the object, and sets its deletionTimestamp.  This policy is
	// cascading, i.e., the dependents will be deleted with Foreground.
	DeletePropagationForeground DeletionPropagation = "Foreground"
)

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
type Preconditions struct {
	// Specifies the target UID.
	// +optional
	UID *types.UID `json:"uid,omitempty"`
	// Specifies the target ResourceVersion
	// +optional
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeleteOptions may be provided when deleting an API object.
type DeleteOptions struct {
	TypeMeta `json:",inline"`

	// The duration in seconds before the object should be deleted. Value must be non-negative integer.
	// The value zero indicates delete immediately. If this value is nil, the default grace period for the
	// specified type will be used.
	// Defaults to a per object value if not specified. zero means delete immediately.
	// +optional
	GracePeriodSeconds *int64 `json:"gracePeriodSeconds,omitempty"`

	// Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status will be
	// returned.
	// +k8s:conversion-gen=false
	// +optional
	Preconditions *Preconditions `json:"preconditions,omitempty"`

	// Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7.
	// Should the dependent objects be orphaned. If true/false, the "orphan"
	// finalizer will be added to/removed from the object's finalizers list.
	// Either this field or PropagationPolicy may be set, but not both.
	// +optional
	OrphanDependents *bool `json:"orphanDependents,omitempty"`

	// Whether and how garbage collection will be performed.
	// Either this field or OrphanDependents may be set, but not both.
	// The default policy is decided by the existing finalizer set in the
	// metadata.finalizers and the resource-specific default policy.
	// Acceptable values are: 'Orphan' - orphan the dependents; 'Background' -
	// allow the garbage collector to delete the dependents in the background;
	// 'Foreground' - a cascading policy that deletes all dependents in the
	// foreground.
	// +optional
	PropagationPolicy *DeletionPropagation `json:"propagationPolicy,omitempty"`

	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	// +listType=atomic
	DryRun []string `json:"dryRun,omitempty"`
}
