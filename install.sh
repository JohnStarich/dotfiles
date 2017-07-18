#!/bin/bash
# Configure and set up the dotfiles in this repository.
#
# This installer was designed to be safe to run even if you
# already have your own dotfiles. However, keep in mind this
# script has the potential to mess up some stuff. I am not
# responsible for any damages this script may cause.

cd "$(dirname "$0")"
DOTFILES_DIR=$(pwd -P)

mkdir -p "$DOTFILES_DIR"/bin

source ./installer_functions.sh

dotlink "$DOTFILES_DIR" ~/.dotfiles

shopt -s nullglob

if [[ -z "$@" ]]; then
    installers=( */install.sh )
else
    installers=()
    for module in "$@"; do
        installers+=("$module/install.sh")
    done
fi

for installer in "${installers[@]}"; do
    module_dir=$(dirname "$installer")
    if [[ ! -f "$installer" ]]; then
        printr "Module installation failed: installer not found '$installer'\n"
    else
        (
            printr "Installing $module_dir... "
            set -e
            cd "$module_dir"
            source "install.sh"
        )
        if [[ $? != 0 ]]; then
            printr "Installation failed for ${module_dir}\n"
        fi
    fi
done
printr 'Done!\n'
