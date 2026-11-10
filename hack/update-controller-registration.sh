#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Gardener's registration generator uses GNU tar options.
if command -v gtar >/dev/null 2>&1; then
  TASK_TOOLS_DIR="$(mktemp -d)"
  trap 'rm -rf "$TASK_TOOLS_DIR"' EXIT
  ln -s "$(command -v gtar)" "$TASK_TOOLS_DIR/tar"
  export PATH="$TASK_TOOLS_DIR:$PATH"
fi

bash "$GARDENER_HACK_DIR/generate-controller-registration.sh" \
  example . "$IMAGE_TAG" ../../example/controller-registration.yaml Extension:example
