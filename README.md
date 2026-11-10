# Gardener Extension Example

This repository is a minimal [Gardener](https://gardener.cloud/) extension that reconciles `Extension` resources of `type: example`.

The controller accepts an optional, empty `ExampleConfig`, logs reconciliation, and completes the extension lifecycle without creating additional resources. Gardener's generic controller handles operation status and finalizers. The controller manager includes leader election, heartbeat, metrics, and health and readiness probes.

## Build and test

Use the Go version specified in [go.mod](go.mod):

```sh
make build
make test
```

Build the container image with `make docker-image`. The default image is `ghcr.io/gardener/gardener-extension-example:latest`; override `REGISTRY` and `IMAGE_TAG` as needed.

## Run locally

With a kubeconfig pointing to a seed cluster:

```sh
bin/gardener-extension-example \
  --config=example/00-componentconfig.yaml \
  --kubeconfig=/path/to/seed-kubeconfig \
  --leader-election=false \
  --heartbeat-namespace=garden
```

For a Gardener local setup, use `make push-to-gardener-local` to build and load the image, then register the extension with `kubectl apply -k example/`. Adapt and apply [example/shoot.yaml](example/shoot.yaml) to enable it for a shoot.

## Adapt the example

Implement reconciliation and cleanup in [pkg/controller/example/actuator.go](pkg/controller/example/actuator.go). Add provider configuration fields to the internal and versioned types under [pkg/apis/example](pkg/apis/example), and controller configuration fields under [pkg/apis/config](pkg/apis/config).

After changing API types or the Helm chart, run `make generate` to regenerate API helpers and the embedded chart in [example/controller-registration.yaml](example/controller-registration.yaml). Generation requires Helm and `yq`; Gardener's registration generator also requires GNU tar (available as `gtar` on macOS).
