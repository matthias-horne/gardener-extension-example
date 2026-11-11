package cmd

import (
	"github.com/gardener/gardener-extension-example/pkg/admission/validator"
	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
)

// GardenWebhookSwitchOptions configures the admission webhooks.
func GardenWebhookSwitchOptions() *webhookcmd.SwitchOptions {
	return webhookcmd.NewSwitchOptions(webhookcmd.Switch(validator.Name, validator.New))
}
