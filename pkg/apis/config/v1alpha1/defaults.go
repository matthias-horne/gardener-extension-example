package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	configv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_ControllerConfiguration sets defaults for the ControllerConfiguration.
func SetDefaults_ControllerConfiguration(obj *ControllerConfiguration) {
	if obj.ClientConnection == nil {
		obj.ClientConnection = &configv1alpha1.ClientConnectionConfiguration{
			QPS:   100,
			Burst: 130,
		}
	}
}
