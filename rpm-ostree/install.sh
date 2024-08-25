#!/usr/bin/env bash

if ! which rpm-ostree >/dev/null 2>/dev/null; then
    exit 0
fi
if ! which jq >/dev/null 2>/dev/null; then
    echo 'jq is required to interact with rpm-ostree' >&2
    exit 1
fi

installed_packages=$(rpm-ostree status --json | jq -r '.deployments[] | select(.booted) | .packages[]')
function is_installed() {
    local pkg=$1
    [[ "$installed_packages" =~ (^|$'\n')"$pkg"($|$'\n') ]]
}
function ensure_installed() {
    local packages=("$@")
    for pkg in "${packages[@]}"; do
        if ! is_installed "$pkg"; then
            echo "Package $pkg not installed. Running bulk install..."
            rpm-ostree install -y "${packages[@]}"
            return
        fi
    done
}

packages=(
    libva-utils
    rpmfusion-free-release
    rpmfusion-nonfree-release
    simple-scan
    zsh
)

cpu_vendor=$(lscpu --json | jq -r '.lscpu[] | select(.field == "Vendor ID:") | .data')
if [[ "$cpu_vendor" == GenuineIntel ]]; then
    packages+=(intel-media-driver)
fi

ensure_installed "${packages[@]}"
