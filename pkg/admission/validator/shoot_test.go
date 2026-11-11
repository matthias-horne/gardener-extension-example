package validator

import (
	"context"
	"strings"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestValidateShoot(t *testing.T) {
	validConfig := `{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig"}`
	validator := NewShootValidator()
	for _, tc := range []struct {
		name       string
		extensions []core.Extension
		wantError  string
	}{
		{name: "no extensions"},
		{name: "no provider configuration", extensions: []core.Extension{{Type: "example"}}},
		{name: "valid configuration", extensions: []core.Extension{extensionWithConfig("example", validConfig)}},
		{name: "other extension", extensions: []core.Extension{extensionWithConfig("other", `{invalid`)}},
		{
			name: "unknown field",
			extensions: []core.Extension{extensionWithConfig("example",
				`{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"ExampleConfig","unknownField":true}`)},
			wantError: `unknown field "unknownField"`,
		},
		{
			name: "wrong kind",
			extensions: []core.Extension{extensionWithConfig("example",
				`{"apiVersion":"example.extensions.gardener.cloud/v1alpha1","kind":"OtherConfig"}`)},
			wantError: `no kind "OtherConfig" is registered`,
		},
		{
			name: "wrong API version",
			extensions: []core.Extension{extensionWithConfig("example",
				`{"apiVersion":"example.extensions.gardener.cloud/v1beta1","kind":"ExampleConfig"}`)},
			wantError: `no kind "ExampleConfig" is registered`,
		},
		{name: "malformed JSON", extensions: []core.Extension{extensionWithConfig("example", `{invalid`)}, wantError: "failed to decode example configuration"},
		{name: "empty raw configuration", extensions: []core.Extension{extensionWithConfig("example", "")}},
		{
			name:       "error includes extension index",
			extensions: []core.Extension{extensionWithConfig("other", `{invalid`), extensionWithConfig("example", `{invalid`)},
			wantError:  "spec.extensions[1].providerConfig",
		},
		{
			name:       "validates all example entries",
			extensions: []core.Extension{extensionWithConfig("example", validConfig), extensionWithConfig("example", `{invalid`)},
			wantError:  "spec.extensions[1].providerConfig",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Workerless Shoots must be validated as well.
			shoot := &core.Shoot{Spec: core.ShootSpec{Extensions: tc.extensions}}
			for _, old := range []*core.Shoot{nil, shoot.DeepCopy()} {
				err := validator.Validate(context.Background(), shoot, old)
				if tc.wantError == "" {
					if err != nil {
						t.Fatalf("Validate() failed: %v", err)
					}
				} else if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("Validate() error = %v, want containing %q", err, tc.wantError)
				}
			}
		})
	}
	if err := validator.Validate(context.Background(), &core.CloudProfile{}, nil); err == nil {
		t.Fatal("expected error for a non-Shoot object")
	}
}

func extensionWithConfig(extensionType, raw string) core.Extension {
	return core.Extension{Type: extensionType, ProviderConfig: &runtime.RawExtension{Raw: []byte(raw)}}
}
