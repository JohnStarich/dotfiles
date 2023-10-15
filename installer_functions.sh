#!/bin/bash

# Original source: https://stackoverflow.com/a/31236568/1530494
function relpath() {
    python -c "import os, sys; print(os.path.relpath(*sys.argv[1:]))" "$@"
}

# Link the given src file into the dotfiles' bin directory using relative paths
function dotbinlink() {
    src="$1"
    if [[ -z "$src" ]]; then
        echo 'Usage: dotbinlink executable_src'
        return 2
    fi
    if [[ ! -e "$src" ]]; then
        echo "dotbinlink: Source does not exist: '$src'"
        return 2
    fi
    dest="$HOME/.dotfiles/bin/$src"
    dest_dir=$(dirname "$dest")
    src="$PWD/$src"
    ln -sf "$(relpath "$src" "$dest_dir")" "$dest"
}

# Only perform the specified command if DEFINITELY_ME is 'yes'.
# This prevents "John Starich"-specific configuration files
# from showing up in someone else's git commits or something.
definitely_me() {
    if [[ "$DEFINITELY_ME" == 'yes' ]]; then
        "$@"
    fi
}

function dotlink() {
    src="$1"
    dest="$2"
    if [[ -z "$src" || -z "$dest" ]]; then
        echo 'Usage: dotlink module_relative_src absolute_link_location'
        return 2
    fi
    if [[ ! -e "$src" ]]; then
        echo "dotlink: Source does not exist: '$src'"
        return 2
    fi
    if [[ -e "$dest" ]]; then
        if [[ -L "$dest" ]]; then
            rm "$dest"
        else
            echo "Skipping: Dotfiles link could not be made, non-link file exists: $dest"
            return 0
        fi
    fi
    if [[ ! "$src" =~ ^/ ]]; then
        src="$PWD/$src"
    fi
    ln -sf "$src" "$dest"
}

function dotpip() {
    local packages=$(pip freeze | sed -e 's/==//')
    local toinstall=()
    for arg in "$@"; do
        if ! grep -q "$arg" <<<"$packages"; then
            toinstall+=("$arg")
        fi
    done
    if (( ${#toinstall} > 0 )); then
        pip install "${toinstall[@]}"
    fi
}

function dotpip3() {
    local packages=$(pip3 freeze | sed -e 's/==//')
    local toinstall=()
    for arg in "$@"; do
        if ! grep -q "$arg" <<<"$packages"; then
            toinstall+=("$arg")
        fi
    done
    if (( ${#toinstall} > 0 )); then
        pip3 install "${toinstall[@]}"
    fi
}

function macos() {
    if [[ -z "$@" ]]; then
        [[ "`uname`" == 'Darwin' ]]
        return $?
    fi
    if [[ "`uname`" == 'Darwin' ]]; then
        eval "$@"
        return $?
    fi
    return 0
}

# Clear the line and print the string
function printr() {
    echo -en "\r\033[K$*"
}
