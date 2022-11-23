#!/bin/bash

if ! which apt-get >/dev/null; then
    exit 0
fi

apt_repos=(
    ppa:masterminds/glide
)
apt_packages=(
    bash-completion
    cmake
    coreutils
    ctags
    dos2unix
    git
    glide
    golang
    htop
    jq
    maven
    mongodb
    mosh
    ncdu
    neovim
    nodejs
    pandoc
    python
    python3
    ruby
    sbt
    scala
    silversearcher-ag
    thefuck
    tmate
    tmux
    trash-cli
    vim
    watch
    wget
    openjdk-8-jdk
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
    # sbt
    if [[ ! -f /etc/apt/sources.list.d/sbt.list ]]; then
        echo "deb https://dl.bintray.com/sbt/debian /" | sudo tee -a /etc/apt/sources.list.d/sbt.list
        apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 2EE0EA64E40A89B84B2DF73499E82A75642AC823
    fi

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
