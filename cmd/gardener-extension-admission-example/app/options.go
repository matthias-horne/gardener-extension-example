package app

import (
	"os"

	admissioncmd "github.com/gardener/gardener-extension-example/pkg/admission/cmd"
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	webhookcmd "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
)

// AdmissionName is the name of the admission component.
const AdmissionName = "admission-example"

// Options holds configuration for the admission controller.
type Options struct {
	restOptions      *controllercmd.RESTOptions
	managerOptions   *controllercmd.ManagerOptions
	webhookOptions   *webhookcmd.AddToManagerOptions
	optionAggregator controllercmd.OptionAggregator
}

// NewOptions creates a new Options instance.
func NewOptions() *Options {
	o := &Options{
		restOptions: &controllercmd.RESTOptions{},
		managerOptions: &controllercmd.ManagerOptions{
			LeaderElection:          true,
			LeaderElectionID:        controllercmd.LeaderElectionNameID(AdmissionName),
			LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
			WebhookServerPort:       10250,
			WebhookCertDir:          "/tmp/admission-example-cert",
			MetricsBindAddress:      ":8080",
			HealthBindAddress:       ":8081",
		},
		webhookOptions: webhookcmd.NewAddToManagerOptions(
			AdmissionName, "", nil, nil,
			&webhookcmd.ServerOptions{Namespace: os.Getenv("WEBHOOK_CONFIG_NAMESPACE")},
			admissioncmd.GardenWebhookSwitchOptions(),
		),
	}
	o.optionAggregator = controllercmd.NewOptionAggregator(o.restOptions, o.managerOptions, o.webhookOptions)
	return o
}
