#!/usr/bin/env bash

set -e -o pipefail

# Build and install Bitwarden polkit fix for Fedora Silverblue
command_prefix=()
if [[ -f /run/.containerenv ]]; then
    command_prefix=(flatpak-spawn --host)
fi
if ! "${command_prefix[@]}" which podman rpm-ostree >/dev/null 2>/dev/null; then
    exit 0
fi

if "${command_prefix[@]}" test ! -f /usr/share/polkit-1/actions/com.bitwarden.Bitwarden.policy; then
    # Set up Bitwarden polkit policy
    cd "$(dirname "$0")/bitwarden-ostree-policy-rpm"

    chcon system_u:object_r:usr_t:s0 ./buildroot/usr/share/polkit-1/actions/com.bitwarden.Bitwarden.policy  # This can't be run inside a (rootless) container build.
    image_tag=localhost/bitwarden-ostree-policy-rpm:latest
    podman build -t "$image_tag" .

    rpm_dir=$PWD/out
    rm -rf "$rpm_dir"
    mkdir -p "$rpm_dir"
    podman run \
        --rm \
        --name "bitwarden-ostree-policy-rpm-$RANDOM" \
        -v "./buildroot:/buildroot:ro" \
        -v "$rpm_dir:/data:Z" \
        "$image_tag"
    rpm_file=./out/bitwarden-polkit-policy.noarch.rpm
    ls "$rpm_file"

    "${command_prefix[@]}" sudo rpm-ostree install \
        --assumeyes \
        --idempotent \
        --apply-live \
        --uninstall bitwarden-polkit-policy \
        "$rpm_file"
    "${command_prefix[@]}" sudo systemctl restart polkit
    cd -
fi

desired_toolbox_image_id=$(./toolbox/build.sh --quiet)
actual_toolbox_image_id=$(podman inspect "$USER" | jq -r '.[].Image')
if [[ "$actual_toolbox_image_id" != "$desired_toolbox_image_id" ]]; then
    echo "Toolbox image is out of date, migrating to new one. Current toolbox sessions will exit." >&2
    sleep 1
    ./toolbox/migrate.sh
fi
