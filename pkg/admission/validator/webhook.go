package validator

import (
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Name is the name of the validation webhook.
const Name = "validator"

// New creates a webhook that validates Shoots using the example extension.
func New(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Name: Name,
		Path: "/webhooks/validate",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewShootValidator(): {{Obj: &core.Shoot{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{v1beta1constants.LabelExtensionExtensionTypePrefix + "example": "true"},
		},
	})
}
