#!/usr/bin/env bash

if ! linux; then
    exit
fi

dotlink applications ~/.local/share/applications/dotfiles
gsettings set org.gnome.desktop.input-sources xkb-options "['caps:escape']"
