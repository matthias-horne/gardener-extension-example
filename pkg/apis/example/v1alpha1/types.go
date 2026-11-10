package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExampleConfig defines the configuration for an individual example extension.
// Add extension-specific fields here when implementing your own extension.
type ExampleConfig struct {
	metav1.TypeMeta `json:",inline"`
}
