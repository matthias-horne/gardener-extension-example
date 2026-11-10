package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerConfiguration defines the configuration for the example controller.
type ControllerConfiguration struct {
	metav1.TypeMeta

	// ClientConnection configures the controller's connection to the seed API server.
	// +optional
	ClientConnection *configv1alpha1.ClientConnectionConfiguration
}
