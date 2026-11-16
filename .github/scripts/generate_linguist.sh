#!/usr/bin/env bash
set -euo pipefail

# Only this scratch tree is writable by upstream-dependent generation. Never
# expose the runner environment, home, Docker socket, or Actions caches.
work=$(mktemp -d "${RUNNER_TEMP:-/tmp}/enry-linguist.XXXXXX")
modules=$(mktemp -d "${RUNNER_TEMP:-/tmp}/enry-modules.XXXXXX")
trap 'chmod -R u+w "$work" "$modules"; rm -rf "$work" "$modules"' EXIT
image='golang:1.26.8-bookworm@sha256:a688600ca24f8a4d3ca77f95b0dd40704a9fc787c826660eb7ba0b641b8b175d'
git archive HEAD | tar -x -C "$work"
(cd "$work" && python3 .github/scripts/sync_linguist.py update)

sandbox=(--rm --read-only --cap-drop=ALL --security-opt=no-new-privileges
  --user "$(id -u):$(id -g)" --tmpfs /tmp:rw,exec,mode=1777
  --mount "type=bind,src=$work,dst=/work" --workdir /work
  --env HOME=/tmp --env GOCACHE=/tmp/go-build
  --env GOPATH=/modules --env GOMODCACHE=/modules/pkg/mod
  --env GOTOOLCHAIN=local)

# Fetch dependencies and the validated revision without executing upstream code.
# No GitHub/Actions tokens or cache directories enter either container.
docker run "${sandbox[@]}" --mount "type=bind,src=$modules,dst=/modules" \
  --env LINGUIST_COMMIT "$image" sh -eu -c '
    go mod download
    git init -q .linguist
    git -C .linguist fetch --depth=1 https://github.com/github/linguist.git "$LINGUIST_COMMIT"
    git -C .linguist checkout --detach FETCH_HEAD
  '

# Generated Go is compiled by golden tests: run that work without networking.
docker run "${sandbox[@]}" --network=none \
  --mount "type=bind,src=$modules,dst=/modules,readonly" \
  "$image" make code-generate

# Use the trusted collector, not any script or Git configuration from the sandbox.
python3 .github/scripts/sync_linguist.py collect "$work"
git diff --binary -- data README.md internal/code-generator/generator > linguist.patch
