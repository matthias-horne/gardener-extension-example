# Gardener Extension Example

This repository is a minimal [Gardener](https://gardener.cloud/) extension that reconciles `Extension` resources of `type: example`.

The controller accepts an optional, empty `ExampleConfig`, logs reconciliation, and completes the extension lifecycle without creating additional resources. Gardener's generic controller handles operation status and finalizers. The controller manager includes leader election, heartbeat, metrics, and health and readiness probes.

The admission controller validates `ExampleConfig` in a Shoot's `spec.extensions` on creation and updates. It accepts an omitted configuration and rejects malformed configurations, unknown fields, and unsupported kinds or API versions. Gardener's webhook framework manages TLS certificates, CA rotation, and the `ValidatingWebhookConfiguration`. The webhook selects Shoots labelled `extensions.extensions.gardener.cloud/example: "true"`, which Gardener sets for Shoots using this extension.

## Build and test

Use the Go version specified in [go.mod](go.mod):

```sh
make build
make test
```

`make build` builds both controller binaries. Build both container images with `make docker-image`, or use `make docker-image-extension` and `make docker-image-admission` individually. The default images are `ghcr.io/gardener/gardener-extension-example:latest` and `ghcr.io/gardener/gardener-extension-admission-example:latest`; override `REGISTRY` and `IMAGE_TAG` as needed.

## Published artifacts

The pipeline publishes artifacts for GitHub releases, pushes to `main`, and pull requests. Tags are derived from the event that triggered the workflow:

| Event | Image, extension, compatibility | Helm chart (OCI) |
| --- | --- | --- |
| Release `<tag>` | `<tag>` | `<tag>` |
| Push to `main` | `branch-main` | `v0.0.0-main` |
| Pull request `#<number>` from branch `<branch>` | `pr-<number>-<branch>` | `v0.0.0-pr.<number>` |

Helm chart tags use a stricter semver-compatible format because Helm parses them as versions. For a release `vX.Y.Z`, the artifacts are published to GHCR as:

| Artifact | Reference for release `vX.Y.Z` |
| --- | --- |
| Extension image | `ghcr.io/gardener/gardener-extension-example:vX.Y.Z` |
| Admission image | `ghcr.io/gardener/gardener-extension-admission-example:vX.Y.Z` |
| Extension chart | `ghcr.io/gardener/charts/gardener-extension-example:vX.Y.Z` |
| Admission runtime chart | `ghcr.io/gardener/charts/admission-example-runtime:vX.Y.Z` |
| Admission application chart | `ghcr.io/gardener/charts/admission-example-application:vX.Y.Z` |
| Extension resource (`extension.yaml`) | `ghcr.io/gardener/extensions/gardener-extension-example:vX.Y.Z` |
| Compatibility information (`compatibility.yaml`) | OCI referrer attached to the extension resource by digest |

Both images are also tagged `latest` on releases. The chart `appVersion` values and image tags match the image tag for the triggering event. For forks, replace `gardener` with the repository owner.

Pull requests from forks are built and tested, but artifacts are only published for pull requests from branches in this repository because fork pull requests receive a read-only token.

The published `extension.yaml` is generated from [example/extension.yaml](example/extension.yaml), with its extension and admission chart refs, manifest annotation, and image values set to the artifacts for the triggering event. [compatibility.yaml](compatibility.yaml) tracks manually maintained `requires` and `supportedFeatures`. The workflow adds the Gardener dependency version from [go.mod](go.mod) as `builtWith."gardener/gardener"` and preserves any `testedWith` information added after compatibility testing. The extension artifact and its YAML file use media type `application/gardener.extension+yaml`; the compatibility artifact and its YAML file use `application/gardener.compatibility+yaml`.

### Download artifacts

Install [ORAS](https://oras.land/docs/installation/) to download the extension resource. Replace `v0.0.2` with the desired release tag:

```sh
release_tag=v0.0.2
image=ghcr.io/gardener/gardener-extension-example
chart=ghcr.io/gardener/charts/gardener-extension-example
extension=ghcr.io/gardener/extensions/gardener-extension-example

docker pull "${image}:${release_tag}"
helm pull "oci://${chart}" --version "${release_tag}"
oras pull "${extension}:${release_tag}" --output ./release-extension
```

The last command writes `release-extension/extension.yaml`.

### Discover attached compatibility information

The pipeline uses the OCI referrers tag fallback for GHCR. Discover compatibility information attached to the extension resource with:

```sh
oras discover "${extension}:${release_tag}" \
  --distribution-spec v1.1-referrers-tag \
  --artifact-type application/gardener.compatibility+yaml
```

Copy the attachment digest from the discovery output to download its `compatibility.yaml`:

```sh
oras pull "${extension}@<compatibility-attachment-digest>" --output ./release-compatibility
```

## Run locally

With a kubeconfig pointing to a seed cluster:

```sh
bin/gardener-extension-example \
  --config=example/00-componentconfig.yaml \
  --kubeconfig=/path/to/seed-kubeconfig \
  --leader-election=false \
  --heartbeat-namespace=garden
```

For a Gardener local setup, use `make push-to-gardener-local` to build and load both images, then register the extension with `kubectl apply -k example/`. Adapt and apply [example/shoot.yaml](example/shoot.yaml) to enable it for a shoot.

## Admission deployment

The admission charts follow the [GCP extension](https://github.com/gardener/gardener-extension-provider-gcp) layout:

- [Runtime chart](charts/gardener-extension-admission-example/charts/runtime) deploys the controller, Service, ServiceAccount, namespace RBAC for certificate secrets and leader election, and PodDisruptionBudget.
- [Application chart](charts/gardener-extension-admission-example/charts/application) grants access to webhook configurations and namespaces in the garden cluster. Webhook registrations are created by the controller at startup.

[example/extension.yaml](example/extension.yaml) references both charts for deployment through gardener-operator. The operator injects the garden kubeconfig; certificates and leader election use the runtime cluster's credentials. `GARDEN_KUBECONFIG` selects the garden kubeconfig, and `SOURCE_CLUSTER` enables a separate cluster for certificates and leader election, using `SOURCE_KUBECONFIG` or in-cluster credentials.

For a single garden cluster, install both charts into the same namespace, with your kubeconfig pointing to that garden:

```sh
helm upgrade --install example-admission-application \
  charts/gardener-extension-admission-example/charts/application \
  --namespace garden --create-namespace \
  --set gardener.virtualCluster.enabled=false
helm upgrade --install example-admission \
  charts/gardener-extension-admission-example/charts/runtime \
  --namespace garden \
  --set gardener.virtualCluster.enabled=false
```

The Service exposes HTTPS on port 443 and forwards to the non-root webhook server on port 10250. `make push-to-gardener-local` loads both images; `kubectl apply -k example/` registers the seed extension, and the Helm commands above install admission in the garden cluster.

## Adapt the example

Implement reconciliation and cleanup in [pkg/controller/example/actuator.go](pkg/controller/example/actuator.go). Add provider configuration fields to the internal and versioned types under [pkg/apis/example](pkg/apis/example), and controller configuration fields under [pkg/apis/config](pkg/apis/config).

Add semantic admission checks in [pkg/admission/validator/shoot.go](pkg/admission/validator/shoot.go) when extending `ExampleConfig`. Register additional webhooks in [pkg/admission/cmd/options.go](pkg/admission/cmd/options.go).

After changing API types or the Helm chart, run `make generate` to regenerate API helpers and the embedded chart in [example/controller-registration.yaml](example/controller-registration.yaml). Generation requires Helm and `yq`; Gardener's registration generator also requires GNU tar (available as `gtar` on macOS).
