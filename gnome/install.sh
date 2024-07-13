#!/usr/bin/env bash

if ! linux; then
    exit
fi

dotlink applications ~/.local/share/applications/dotfiles
gsettings set org.gnome.desktop.input-sources xkb-options "['caps:escape']"

mkdir -p ~/.config/systemd/user
for service in systemd/*; do
    name="$(basename "$service")"
    dotlink "$service" ~/.config/systemd/user/"$name"
    systemctl --user enable --now "$name"
done
