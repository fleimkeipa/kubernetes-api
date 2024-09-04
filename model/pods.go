package model

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Pods struct {
	Name string
}

type PodsRequest struct {
	TypeMeta   TypeMeta             `json:"type_meta"`
	ObjectMeta ObjectMeta           `json:"object_meta"`
	Spec       Spec                 `json:"spec"`
	Opts       metav1.CreateOptions `json:"opts"`
}

type TypeMeta struct {
	Kind string `json:"kind"`
}

type ObjectMeta struct {
	Name      string `json:"name"`
	NameSpace string `json:"namespace"`
}

type Spec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
