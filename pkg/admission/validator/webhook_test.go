package validator

import (
	"context"
	"fmt"
	"strings"
	"testing"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core/install"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type webhookManager struct {
	manager.Manager
	scheme *runtime.Scheme
	client client.Client
}

func (m *webhookManager) GetScheme() *runtime.Scheme { return m.scheme }
func (m *webhookManager) GetClient() client.Client   { return m.client }

func TestWebhookAdmissionAndRegistration(t *testing.T) {
	scheme := runtime.NewScheme()
	install.Install(scheme)
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{corev1beta1.SchemeGroupVersion})
	mapper.Add(corev1beta1.SchemeGroupVersion.WithKind("Shoot"), meta.RESTScopeNamespace)
	mgr := &webhookManager{scheme: scheme, client: fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).Build()}
	webhook, err := New(mgr)
	if err != nil {
		t.Fatal(err)
	}
	version := corev1beta1.SchemeGroupVersion.Version
	for _, operation := range []admissionv1.Operation{admissionv1.Create, admissionv1.Update} {
		for _, tc := range []struct {
			name    string
			config  string
			allowed bool
		}{
			{name: "valid", config: `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig"}`, allowed: true},
			{name: "unknown field", config: `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig","unknownField":true}`},
		} {
			t.Run(fmt.Sprintf("%s/%s/%s", version, operation, tc.name), func(t *testing.T) {
				raw := fmt.Sprintf(`{"apiVersion":"core.gardener.cloud/%s","kind":"Shoot","metadata":{"name":"example"},"spec":{"extensions":[{"type":"example","providerConfig":%s}]}}`, version, tc.config)
				req := admission.Request{AdmissionRequest: admissionv1.AdmissionRequest{
					Kind:      metav1.GroupVersionKind{Group: "core.gardener.cloud", Version: version, Kind: "Shoot"},
					Operation: operation,
					Object:    runtime.RawExtension{Raw: []byte(raw)},
				}}
				if operation == admissionv1.Update {
					req.OldObject = runtime.RawExtension{Raw: []byte(fmt.Sprintf(`{"apiVersion":"core.gardener.cloud/%s","kind":"Shoot"}`, version))}
				}
				response := webhook.Webhook.Handle(context.Background(), req)
				if response.Allowed != tc.allowed {
					t.Fatalf("Allowed = %t, want %t; result: %v", response.Allowed, tc.allowed, response.Result)
				}
				if !tc.allowed && (response.Result == nil || !strings.Contains(response.Result.Message, "spec.extensions[0].providerConfig")) {
					t.Fatalf("expected configuration field in denial: %v", response.Result)
				}
			})
		}
	}

	for _, mode := range []string{extensionswebhook.ModeService, extensionswebhook.ModeURL} {
		t.Run("registration/"+mode, func(t *testing.T) {
			configs, _, err := extensionswebhook.BuildWebhookConfigs(
				[]*extensionswebhook.Webhook{webhook}, mgr.client, "garden", "admission-example", false,
				443, mode, "gardener-extension-admission-example.garden", nil, false,
			)
			if err != nil {
				t.Fatal(err)
			}
			config := configs.ValidatingWebhookConfig
			if config == nil || config.Name != "gardener-extension-admission-example" || len(config.Webhooks) != 1 {
				t.Fatalf("unexpected registration: %+v", config)
			}
			registered := config.Webhooks[0]
			if len(registered.Rules) != 1 || len(registered.Rules[0].Resources) != 1 || registered.Rules[0].Resources[0] != "shoots" {
				t.Fatalf("unexpected webhook rules: %+v", registered.Rules)
			}
			if registered.ObjectSelector.MatchLabels["extensions.extensions.gardener.cloud/example"] != "true" {
				t.Fatalf("unexpected selector: %+v", registered.ObjectSelector)
			}
			if mode == extensionswebhook.ModeService {
				service := registered.ClientConfig.Service
				if service == nil || service.Name != config.Name || service.Namespace != "garden" || *service.Port != 443 || *service.Path != "/webhooks/validate" {
					t.Fatalf("unexpected service: %+v", service)
				}
			} else if registered.ClientConfig.URL == nil || *registered.ClientConfig.URL != "https://gardener-extension-admission-example.garden/webhooks/validate" {
				t.Fatalf("unexpected URL: %+v", registered.ClientConfig.URL)
			}
		})
	}
}
