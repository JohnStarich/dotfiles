#!/bin/bash
# Configure and set up the dotfiles in this repository.
#
# This installer was designed to be safe to run even if you
# already have your own dotfiles. However, keep in mind this
# script has the potential to mess up some stuff. I am not
# responsible for any damages this script may cause.

cd "$(dirname "$0")"
DOTFILES_DIR=$PWD

mkdir -p "$DOTFILES_DIR"/bin

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
    ln -sf "$PWD/$src" "$DOTFILES_DIR"/bin/"$src"
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
    ln -sf "$PWD/$src" "$dest"
}

shopt -s nullglob

for installer in */install.sh; do
    module_dir=$(dirname "$installer")
    if [[ -f "$installer" ]]; then
        (
            set -e
            cd "$module_dir"
            source "install.sh"
        )
        if [[ $? != 0 ]]; then
            echo "Installation failed for $module_dir"
        fi
    fi
done
echo 'Done!'
