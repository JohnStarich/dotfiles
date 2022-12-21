#!/bin/bash

if ! macos; then
    exit 0
fi

if ! which brew >/dev/null; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

brew_taps=(
    jcgay/jcgay # for maven-deluxe
    universal-ctags/universal-ctags
    homebrew/cask-fonts
)
brew_formulae=(
    bash-completion@2
    cmake
    coreutils
    delta
    docker
    dos2unix
    git
    gnupg
    go
    htop
    jq
    kubectl
    maven-deluxe
    mobile-shell
    ncdu
    node
    nvim
    pandoc
    pinentry-mac
    python
    python3
    reattach-to-user-namespace
    ruby
    sbt
    scala
    the_silver_searcher
    thefuck
    tmate
    tmux
    trash
    vim
    watch
    watchman
    wget
    zsh
    zsh-completions
)
brew_cask_formulae=(
    corelocationcli
    docker
    font-sf-mono-for-powerline
    graphql-playground
    qlstephen
    virtualbox
)
brew_head_only_formulae=(
    universal-ctags
)

shopt -s nullglob

set +e
err=0

for tap in "${brew_taps[@]}"; do
    if ! brew tap "$tap"; then
        echo "Error installing brew tap: '$tap'"
        err=1
    fi
done

# Brew formulae
if ! brew ls --versions "${brew_formulae[@]}" >/dev/null; then
    for formula in "${brew_formulae[@]}"; do
        if brew ls --versions "$formula" >/dev/null; then
            continue
        fi
        if ! brew install "$formula"; then
            echo "Error installing brew formula: '$formula'"
            err=1
        fi
    done
fi

# Brew head-only formulae
if ! brew ls --versions "${brew_head_only_formulae[@]}" >/dev/null; then
    for formula in "${brew_formulae[@]}"; do
        if brew ls --versions "$formula" >/dev/null; then
            continue
        fi
        if ! brew install --HEAD "$formula"; then
            echo "Error installing brew (head-only) formula: '$formula'"
            err=1
        fi
    done
fi

# Brew cask formulae
if ! brew ls --cask --versions "${brew_cask_formulae[@]}" >/dev/null; then
    for formula in "${brew_cask_formulae[@]}"; do
        if brew ls --cask --versions "$formula" >/dev/null; then
            continue
        fi
        if ! brew install --cask "$formula"; then
            echo "Error installing brew cask formula: '$formula'"
            err=1
        fi
    done
fi

exit $err
