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
        sudo dnf config-manager addrepo --from-repofile="$url"
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
    gh
    git-delta
    golang
    htop
    jq
    kde-connect
    neovim
    openssh-server
    powerline-fonts
    ripgrep
    the_silver_searcher
    tmux
    zeal
    zsh
)

ensure_installed "${packages[@]}"

# Build and install Bitwarden polkit fix for Fedora Silverblue
command_prefix=()
if [[ -f /run/.containerenv ]]; then
    command_prefix=(flatpak-spawn --host)
fi
if ! "${command_prefix[@]}" which podman rpm-ostree >/dev/null 2>/dev/null; then
    exit 0
fi

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
