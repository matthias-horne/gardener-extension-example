package app

import (
	"context"
	"fmt"

	"github.com/gardener/gardener-extension-example/pkg/apis/example/install"
	"github.com/gardener/gardener-extension-example/pkg/controller/example"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	heartbeatcontroller "github.com/gardener/gardener/extensions/pkg/controller/heartbeat"
	"github.com/gardener/gardener/extensions/pkg/util"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"

	"github.com/spf13/cobra"

	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NewControllerManagerCommand creates a new command that is used to start the controller.
func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	options := NewOptions()

	cmd := &cobra.Command{
		Use:           "gardener-extension-example",
		Short:         "Runs a minimal Gardener extension for shoot clusters.",
		SilenceErrors: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.optionAggregator.Complete(); err != nil {
				return fmt.Errorf("error completing options: %w", err)
			}
			if err := options.heartbeatOptions.Validate(); err != nil {
				return err
			}

			cmd.SilenceUsage = true
			return options.run(ctx)
		},
	}

	options.optionAggregator.AddFlags(cmd.Flags())

	return cmd
}

func (o *Options) run(ctx context.Context) error {
	ctrlConfig := o.exampleOptions.Completed()
	util.ApplyClientConnectionConfigurationToRESTConfig(ctrlConfig.Config().ClientConnection, o.restOptions.Completed().Config)

	mgrOpts := o.managerOptions.Completed().Options()

	mgr, err := manager.New(o.restOptions.Completed().Config, mgrOpts)
	if err != nil {
		return fmt.Errorf("could not instantiate controller-manager: %w", err)
	}

	if err := extensionscontroller.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("could not update manager scheme: %w", err)
	}

	if err := install.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("could not update manager scheme: %w", err)
	}

	o.controllerOptions.Completed().Apply(&example.DefaultAddOptions.ControllerOptions)
	o.reconcileOptions.Completed().Apply(&example.DefaultAddOptions.IgnoreOperationAnnotation)
	example.DefaultAddOptions.ExtensionClasses = o.generalOptions.Completed().ExtensionClasses
	o.heartbeatOptions.Completed().Apply(&heartbeatcontroller.DefaultAddOptions)

	if err := o.controllerSwitches.Completed().AddToManager(ctx, mgr); err != nil {
		return fmt.Errorf("could not add controllers to manager: %w", err)
	}

	if err := mgr.AddReadyzCheck("informer-sync", gardenerhealthz.NewCacheSyncHealthz(mgr.GetCache())); err != nil {
		return fmt.Errorf("could not add ready check for informers: %w", err)
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		return fmt.Errorf("could not add health check to manager: %w", err)
	}

	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("error running manager: %w", err)
	}

	return nil
}
