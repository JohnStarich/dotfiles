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
            rpm-ostree install --idempotent -y "${packages[@]}"
            return
        fi
    done
}

if ! is_installed rpmfusion-free-release; then
    rpm-ostree install -y https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm
    echo A reboot is required to finish installing rpmfusion. >&2
fi
if ! is_installed rpmfusion-nonfree-release; then
    rpm-ostree install -y https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm
    echo A reboot is required to finish installing rpmfusion. >&2
fi

packages=(
    libva-utils
    simple-scan
    zsh
)

cpu_vendor=$(lscpu --json | jq -r '.lscpu[] | select(.field == "Vendor ID:") | .data')
if [[ "$cpu_vendor" == GenuineIntel ]]; then
    packages+=(intel-media-driver)
fi

ensure_installed "${packages[@]}"
