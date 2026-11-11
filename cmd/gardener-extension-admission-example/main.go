package main

import (
	"os"

	"github.com/gardener/gardener-extension-example/cmd/gardener-extension-admission-example/app"
	"github.com/gardener/gardener/pkg/logger"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	runtimelog.SetLogger(logger.MustNewZapLogger(logger.InfoLevel, logger.FormatJSON))
	if err := app.NewAdmissionCommand(signals.SetupSignalHandler()).Execute(); err != nil {
		runtimelog.Log.Error(err, "error executing the admission command")
		os.Exit(1)
	}
}
