#!/usr/bin/env bash

if ! which dnf >/dev/null 2>/dev/null; then
    exit 0
fi

installed_packages=$(dnf list --installed)
function is_installed() {
    local pkg=$1
    [[ "$installed_packages" =~ (^|$'\n')"$pkg"\. ]]
}
function ensure_installed() {
    local packages=("$@")
    for pkg in "${packages[@]}"; do
        if ! is_installed "$pkg"; then
            echo "Package $pkg not installed. Running bulk install..."
            sudo dnf install -y "${packages[@]}"
            return
        fi
    done
}
function ensure_repo() {
    local name=$1
    local url=$2
    if [[ ! -f "/etc/yum.repos.d/$name" ]]; then
        sudo dnf config-manager --add-repo "$repo"
    fi
}

ensure_installed dnf-plugins-core
ensure_repo docker-ce.repo https://download.docker.com/linux/fedora/docker-ce.repo

packages=(
    containerd.io
    docker-buildx-plugin
    docker-ce
    docker-ce-cli
    docker-compose-plugin
    git-delta
    neovim
    powerline-fonts
    zeal
)

ensure_installed "${packages[@]}"
