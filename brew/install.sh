#!/bin/bash

if ! macos; then
    exit 0
fi

if ! which brew >/dev/null; then
    /usr/bin/ruby -e "`curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install`"
fi

brew_taps=(
    jcgay/jcgay # for maven-deluxe
)
brew_formulae=(
    bash-completion2
    cmake
    coreutils
    ctags
    dos2unix
    git
    glide
    go
    htop
    javarepl
    jq
    maven-deluxe
    mongodb
    ncdu
    node
    pandoc
    python
    python3
    reattach-to-user-namespace
    ruby
    sbt
    scala
    the_silver_searcher
    tmate
    tmux
    trash
    vim
    watch
    watchman
    wget
    zsh-completions
)
brew_cask_formulae=(
    corelocationcli
    java
    qlstephen
    virtualbox
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

if ! brew cask ls --versions "${brew_cask_formulae[@]}" >/dev/null; then
    for formula in "${brew_cask_formulae[@]}"; do
        if brew cask ls --versions "$formula" >/dev/null; then
            continue
        fi
        if ! brew cask install "$formula"; then
            echo "Error installing brew cask formula: '$formula'"
            err=1
        fi
    done
fi

exit $err

