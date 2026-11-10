package app

import (
	"os"

	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	heartbeatcmd "github.com/gardener/gardener/extensions/pkg/controller/heartbeat/cmd"

	examplecmd "github.com/gardener/gardener-extension-example/pkg/cmd"
)

// ExtensionName is the name of the extension.
const ExtensionName = "extension-example"

// Options holds configuration for the example controller manager.
type Options struct {
	generalOptions     *controllercmd.GeneralOptions
	exampleOptions     *examplecmd.ExampleOptions
	restOptions        *controllercmd.RESTOptions
	managerOptions     *controllercmd.ManagerOptions
	controllerOptions  *controllercmd.ControllerOptions
	heartbeatOptions   *heartbeatcmd.Options
	controllerSwitches *controllercmd.SwitchOptions
	reconcileOptions   *controllercmd.ReconcilerOptions
	optionAggregator   controllercmd.OptionAggregator
}

// NewOptions creates a new Options instance.
func NewOptions() *Options {
	options := &Options{
		generalOptions: &controllercmd.GeneralOptions{},
		exampleOptions: &examplecmd.ExampleOptions{},
		restOptions:    &controllercmd.RESTOptions{},
		managerOptions: &controllercmd.ManagerOptions{
			LeaderElection:          true,
			LeaderElectionID:        controllercmd.LeaderElectionNameID(ExtensionName),
			LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
			MetricsBindAddress:      ":8080",
			HealthBindAddress:       ":8081",
		},
		controllerOptions: &controllercmd.ControllerOptions{
			MaxConcurrentReconciles: 5,
		},
		heartbeatOptions: &heartbeatcmd.Options{
			ExtensionName:        ExtensionName,
			RenewIntervalSeconds: 30,
			Namespace:            os.Getenv("LEADER_ELECTION_NAMESPACE"),
		},
		controllerSwitches: examplecmd.ControllerSwitchOptions(),
		reconcileOptions:   &controllercmd.ReconcilerOptions{},
	}
	options.optionAggregator = controllercmd.NewOptionAggregator(
		options.generalOptions,
		options.restOptions,
		options.managerOptions,
		options.controllerOptions,
		options.exampleOptions,
		controllercmd.PrefixOption("heartbeat-", options.heartbeatOptions),
		options.controllerSwitches,
		options.reconcileOptions,
	)
	return options
}
