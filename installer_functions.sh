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
    dest="$DOTFILES_DIR/bin/$src"
    dest_dir=$(dirname "$dest")
    src="$PWD/$src"
    ln -sf "$(relpath "$src" "$dest_dir")" "$dest"
}

# Only perform the specified command if DEFINITELY_ME is 'yes'.
# This prevents "John Starich"-specific configuration files
# from showing up in someone else's git commits or something.
definitely_me() {
    if [[ "$DEFINITELY_ME" == 'yes' ]]; then
        eval "$@"
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

function macos() {
    if [[ "`uname`" == 'Darwin' ]]; then
        if [[ ! -z "$@" ]]; then
            eval "$@"
            return $?
        else
            return 0
        fi
    fi
    return 1
}

# Clear the line and print the string
function printr() {
    echo -en "\r\033[K$*"
}
