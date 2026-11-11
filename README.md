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
| Helm chart | `ghcr.io/gardener/charts/gardener-extension-example:vX.Y.Z` |
| Extension resource (`extension.yaml`) | `ghcr.io/gardener/extensions/gardener-extension-example:vX.Y.Z` |
| Compatibility information (`compatibility.yaml`) | OCI referrer attached to the extension resource by digest |

The image is also tagged `latest` on releases. The chart's `appVersion` and `image.tag` match the image tag for the triggering event. For forks, replace `gardener` with the repository owner.

Pull requests from forks are built and tested, but artifacts are only published for pull requests from branches in this repository because fork pull requests receive a read-only token.

The published `extension.yaml` is generated from [example/extension.yaml](example/extension.yaml), with its Helm chart and manifest references set to the artifacts for the triggering event. [compatibility.yaml](compatibility.yaml) tracks manually maintained `requires` and `supportedFeatures`. The workflow adds the Gardener dependency version from [go.mod](go.mod) as `builtWith."gardener/gardener"` and preserves any `testedWith` information added after compatibility testing. The extension artifact and its YAML file use media type `application/gardener.extension+yaml`; the compatibility artifact and its YAML file use `application/gardener.compatibility+yaml`.

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

For a Gardener local setup, use `make push-to-gardener-local` to build and load the image, then register the extension with `kubectl apply -k example/`. Adapt and apply [example/shoot.yaml](example/shoot.yaml) to enable it for a shoot.

## Adapt the example

Implement reconciliation and cleanup in [pkg/controller/example/actuator.go](pkg/controller/example/actuator.go). Add provider configuration fields to the internal and versioned types under [pkg/apis/example](pkg/apis/example), and controller configuration fields under [pkg/apis/config](pkg/apis/config).

After changing API types or the Helm chart, run `make generate` to regenerate API helpers and the embedded chart in [example/controller-registration.yaml](example/controller-registration.yaml). Generation requires Helm and `yq`; Gardener's registration generator also requires GNU tar (available as `gtar` on macOS).
