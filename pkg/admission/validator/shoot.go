package validator

import (
	"context"
	"fmt"

	"github.com/gardener/gardener-extension-example/pkg/apis/example/install"
	"github.com/gardener/gardener-extension-example/pkg/apis/example/v1alpha1"
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewShootValidator creates a validator for the example provider configuration.
func NewShootValidator() extensionswebhook.Validator {
	scheme := runtime.NewScheme()
	install.Install(scheme)
	return &shootValidator{decoder: serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()}
}

type shootValidator struct {
	decoder runtime.Decoder
}

// Validate validates the optional configuration of each example extension.
func (s *shootValidator) Validate(_ context.Context, new, _ client.Object) error {
	shoot, ok := new.(*core.Shoot)
	if !ok {
		return fmt.Errorf("expected a Shoot, got %T", new)
	}
	var allErrs field.ErrorList
	for i, extension := range shoot.Spec.Extensions {
		if extension.Type != "example" || extension.ProviderConfig == nil {
			continue
		}
		config := &v1alpha1.ExampleConfig{}
		if _, _, err := s.decoder.Decode(extension.ProviderConfig.Raw, nil, config); err != nil {
			path := field.NewPath("spec", "extensions").Index(i).Child("providerConfig")
			allErrs = append(allErrs, field.Invalid(path, extension.ProviderConfig, fmt.Sprintf("failed to decode example configuration: %v", err)))
		}
		// Add extension-specific semantic validation here when ExampleConfig gains fields.
	}
	return allErrs.ToAggregate()
}
