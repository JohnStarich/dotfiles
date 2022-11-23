#!/bin/bash

if ! which apt-get >/dev/null; then
    exit 0
fi

apt_repos=(
)
apt_packages=(
    bash-completion
    cmake
    coreutils
    dos2unix
    git
    golang
    htop
    jq
    mosh
    ncdu
    neovim
    python3
    silversearcher-ag
    thefuck
    tmux
    trash-cli
    vim
    watch
    wget
)

shopt -s nullglob

set +e
err=0

function is_installed() {
    installed=$(dpkg-query --show --showformat='${db:Status-Status}\n' "$@" 2>&1 | grep -v '^installed$' | wc -l)
    [[ $installed -ne 0 ]]
}

if is_installed "${apt_packages[@]}"; then
    export apt_packages
    sudo -u root bash <<EOT
    function is_installed() {
        installed=\$(dpkg-query --show --showformat='\${db:Status-Status}\n' "\$@" 2>&1 | grep -v '^installed$' | wc -l)
        [[ \$installed -ne 0 ]]
    }

    apt-get update
    for repo in ${apt_repos[@]}; do
        add-apt-repository \$repo
    done
    apt-get update
    for package in ${apt_packages[@]}; do
        if is_installed "$package" >/dev/null; then
            continue
        fi
        if ! apt-get install -y "\$package"; then
            echo "Error installing apt package: '\$package'"
            err=1
        fi
    done
EOT
fi

exit $err
