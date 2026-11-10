package cmd

import (
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	extensionsheartbeatcontroller "github.com/gardener/gardener/extensions/pkg/controller/heartbeat"

	"github.com/gardener/gardener-extension-example/pkg/controller/example"
)

// ControllerSwitchOptions configures the example and heartbeat controllers.
func ControllerSwitchOptions() *controllercmd.SwitchOptions {
	return controllercmd.NewSwitchOptions(
		controllercmd.Switch(example.ControllerName, example.AddToManager),
		controllercmd.Switch(extensionsheartbeatcontroller.ControllerName, extensionsheartbeatcontroller.AddToManager),
	)
}
