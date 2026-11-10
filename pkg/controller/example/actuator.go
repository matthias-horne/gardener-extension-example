package example

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-example/pkg/apis/example/v1alpha1"
)

// NewActuator returns an actuator responsible for example Extension resources.
func NewActuator(mgr manager.Manager) extension.Actuator {
	return &actuator{
		decoder: serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder(),
	}
}

type actuator struct {
	decoder runtime.Decoder
}

// Reconcile validates the provider configuration and completes successfully.
func (a *actuator) Reconcile(_ context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	if ex.Spec.ProviderConfig != nil {
		config := &v1alpha1.ExampleConfig{}
		if _, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, config); err != nil {
			return fmt.Errorf("failed to decode provider config: %w", err)
		}
	}

	// Add logic to create or update the resources managed by your extension here.
	log.Info("Reconciling example extension", "namespace", ex.Namespace, "name", ex.Name)
	return nil
}

// Delete cleans up resources managed by the extension.
func (a *actuator) Delete(_ context.Context, _ logr.Logger, _ *extensionsv1alpha1.Extension) error {
	// Add cleanup logic here when the extension manages resources.
	return nil
}

// ForceDelete performs cleanup during forceful shoot deletion.
func (a *actuator) ForceDelete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, log, ex)
}

// Restore reconciles the extension after a shoot migration.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate cleans up the extension before a shoot migration.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, log, ex)
}
