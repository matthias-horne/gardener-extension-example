package example

import (
	"context"
	"testing"

	"github.com/gardener/gardener-extension-example/pkg/apis/example/install"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func TestReconcileProviderConfig(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := install.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	a := &actuator{decoder: serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()}

	for _, tc := range []struct {
		name    string
		config  string
		wantErr bool
	}{
		{name: "no provider configuration"},
		{name: "empty example configuration", config: `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig"}`},
		{name: "unknown field", config: `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig","unknownField":true}`, wantErr: true},
		{name: "unregistered kind", config: `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"OtherConfig"}`, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ex := &extensionsv1alpha1.Extension{}
			if tc.config != "" {
				ex.Spec.ProviderConfig = &runtime.RawExtension{Raw: []byte(tc.config)}
			}
			err := a.Reconcile(context.Background(), logr.Discard(), ex)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Reconcile() error = %v, want error: %t", err, tc.wantErr)
			}
		})
	}
}
