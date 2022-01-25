package blueprint

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterResourceRef struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type InputRef struct {
	Name     string `json:"name"`
	Resource string `json:"resource"`
}

type Criteria struct {
	Selector    *metav1.LabelSelector `json:"selector"`
	Sources     []InputRef           `json:"sources"`
	Images      []InputRef           `json:"images"`
	Configs     []InputRef           `json:"configs"`
}

type Option struct {
	Name string `json:"name"`
	Selector *metav1.LabelSelector `json:"selector"`
	Criteria
}

type Resource struct {
	Name        string               `json:"name"`
	TemplateRef *ClusterResourceRef  `json:"templateRef"`
	Kind        string               `json:"kind"`
	Options     []Option             `json:"options"`
	Criteria
}

type Spec struct {
	Selector  *metav1.LabelSelector `json:"selector"`
	Resources []Resource            `json:"resources"`
}

type Blueprint struct {
	Spec Spec `json:"spec"`
}
