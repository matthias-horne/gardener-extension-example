package app

import (
	"context"
	"fmt"
	"os"

	exampleinstall "github.com/gardener/gardener-extension-example/pkg/apis/example/install"
	"github.com/gardener/gardener/extensions/pkg/util"
	"github.com/gardener/gardener/pkg/apis/core/install"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NewAdmissionCommand creates a command for running the example admission webhook.
func NewAdmissionCommand(ctx context.Context) *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:           "gardener-extension-admission-example",
		Short:         "Validates example extension configuration on Gardener Shoots.",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if gardenKubeconfig := os.Getenv("GARDEN_KUBECONFIG"); gardenKubeconfig != "" {
				o.restOptions.Kubeconfig = gardenKubeconfig
			}
			if err := o.optionAggregator.Complete(); err != nil {
				return fmt.Errorf("error completing options: %w", err)
			}
			cmd.SilenceUsage = true
			return o.run(ctx)
		},
	}
	o.optionAggregator.AddFlags(cmd.Flags())
	return cmd
}

func (o *Options) run(ctx context.Context) error {
	util.ApplyClientConnectionConfigurationToRESTConfig(&v1alpha1.ClientConnectionConfiguration{
		QPS: 100, Burst: 130,
	}, o.restOptions.Completed().Config)
	mgrOpts := o.managerOptions.Completed().Options()
	namespace := o.webhookOptions.Server.Completed().Namespace

	// Certificates and leader election may live in the runtime garden while
	// webhook registrations are managed in the virtual garden.
	var sourceConfig *rest.Config
	if os.Getenv("SOURCE_CLUSTER") != "" {
		var err error
		sourceConfig, err = clientcmd.BuildConfigFromFlags("", os.Getenv("SOURCE_KUBECONFIG"))
		if err != nil {
			return fmt.Errorf("could not configure source cluster: %w", err)
		}
		mgrOpts.LeaderElectionConfig = sourceConfig
	} else if namespace != "" {
		// Certificate management only needs secrets in the webhook namespace.
		mgrOpts.Cache.ByObject = map[client.Object]cache.ByObject{
			&corev1.Secret{}: {Namespaces: map[string]cache.Config{namespace: {}}},
		}
	}

	mgr, err := manager.New(o.restOptions.Completed().Config, mgrOpts)
	if err != nil {
		return fmt.Errorf("could not instantiate admission manager: %w", err)
	}
	install.Install(mgr.GetScheme())
	if err := exampleinstall.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("could not update manager scheme: %w", err)
	}

	var sourceCluster cluster.Cluster
	if sourceConfig != nil {
		sourceCluster, err = cluster.New(sourceConfig, func(opts *cluster.Options) {
			if namespace != "" {
				opts.Cache.DefaultNamespaces = map[string]cache.Config{namespace: {}}
			}
		})
		if err != nil {
			return fmt.Errorf("could not instantiate source cluster: %w", err)
		}
		if err := mgr.Add(sourceCluster); err != nil {
			return fmt.Errorf("could not add source cluster: %w", err)
		}
		if err := mgr.AddReadyzCheck("source-informer-sync", gardenerhealthz.NewCacheSyncHealthz(sourceCluster.GetCache())); err != nil {
			return fmt.Errorf("could not add source cluster ready check: %w", err)
		}
	}

	if _, err := o.webhookOptions.Completed().AddToManager(ctx, mgr, sourceCluster); err != nil {
		return fmt.Errorf("could not add webhooks to manager: %w", err)
	}
	if err := mgr.AddReadyzCheck("informer-sync", gardenerhealthz.NewCacheSyncHealthz(mgr.GetCache())); err != nil {
		return fmt.Errorf("could not add informer ready check: %w", err)
	}
	if err := mgr.AddReadyzCheck("webhook-server", mgr.GetWebhookServer().StartedChecker()); err != nil {
		return fmt.Errorf("could not add webhook ready check: %w", err)
	}
	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		return fmt.Errorf("could not add health check: %w", err)
	}
	return mgr.Start(ctx)
}
