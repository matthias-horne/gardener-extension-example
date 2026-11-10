package cmd

import (
	"errors"
	"fmt"
	"os"

	configapi "github.com/gardener/gardener-extension-example/pkg/apis/config"
	"github.com/gardener/gardener-extension-example/pkg/apis/config/v1alpha1"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	scheme  = runtime.NewScheme()
	decoder runtime.Decoder
)

func init() {
	utilruntime.Must(configapi.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
}

// ExampleOptions holds options for the example controller.
type ExampleOptions struct {
	ConfigLocation string
	config         *ControllerConfig
}

// AddFlags implements Flagger.AddFlags.
func (o *ExampleOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigLocation, "config", "", "Path to controller configuration")
}

// Complete implements Completer.Complete.
func (o *ExampleOptions) Complete() error {
	if o.ConfigLocation == "" {
		return errors.New("config location is not set")
	}
	data, err := os.ReadFile(o.ConfigLocation)
	if err != nil {
		return err
	}

	config := configapi.ControllerConfiguration{}
	if _, _, err := decoder.Decode(data, nil, &config); err != nil {
		return fmt.Errorf("failed to decode controller configuration: %w", err)
	}
	o.config = &ControllerConfig{config: config}
	return nil
}

// Completed returns the configuration after Complete succeeds.
func (o *ExampleOptions) Completed() *ControllerConfig {
	return o.config
}

// ControllerConfig contains the completed controller configuration.
type ControllerConfig struct {
	config configapi.ControllerConfiguration
}

// Config returns the completed controller configuration.
func (c *ControllerConfig) Config() *configapi.ControllerConfiguration {
	return &c.config
}
