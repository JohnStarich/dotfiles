#!/usr/bin/env bash

function isodate() {
    date --rfc-3339=seconds -u \
        | tr ' ' T \
        | sed \
            -e 's/+00:00.*/Z/' \
            -e 's/:/-/g' \
        | tr -d '\n'
}

set -ex -o pipefail

old_toolbox="$USER-old-$(isodate)"
podman container rename "$USER" "$old_toolbox"
toolbox create --image localhost/johnstarich-dev:latest "$USER"
podman stop "$old_toolbox" || true
