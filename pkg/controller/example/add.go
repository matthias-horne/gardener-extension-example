package example

import (
	"context"

	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	// Type is the type of Extension resource.
	Type = "example"
	// ControllerName is the name of the example controller.
	ControllerName = "example"
	// FinalizerSuffix is the finalizer suffix for the example controller.
	FinalizerSuffix = "example"
)

var (
	// DefaultAddOptions are the default AddOptions for AddToManager.
	DefaultAddOptions = AddOptions{}
)

// AddOptions are options to apply when adding the example controller to the manager.
type AddOptions struct {
	// ControllerOptions contains options for the controller.
	ControllerOptions controller.Options

	// IgnoreOperationAnnotation specifies whether to ignore the operation annotation or not.
	IgnoreOperationAnnotation bool
	// ExtensionClasses defines the extension classes this extension is responsible for.
	ExtensionClasses []extensionsv1alpha1.ExtensionClass
}

// AddToManager adds a controller with the default Options to the given Controller Manager.
func AddToManager(ctx context.Context, mgr manager.Manager) error {
	return AddToManagerWithOptions(ctx, mgr, DefaultAddOptions)
}

// AddToManagerWithOptions adds a controller with the given Options to the given manager.
func AddToManagerWithOptions(ctx context.Context, mgr manager.Manager, opts AddOptions) error {
	return extension.Add(mgr, extension.AddArgs{
		Actuator:                  NewActuator(mgr),
		ControllerOptions:         opts.ControllerOptions,
		Name:                      ControllerName,
		FinalizerSuffix:           FinalizerSuffix,
		Predicates:                extension.DefaultPredicates(ctx, mgr, opts.IgnoreOperationAnnotation),
		IgnoreOperationAnnotation: opts.IgnoreOperationAnnotation,
		Type:                      Type,
		ExtensionClasses:          opts.ExtensionClasses,
	})
}
