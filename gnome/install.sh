#!/usr/bin/env bash

if ! linux; then
    exit
fi

dotlink applications ~/.local/share/applications/dotfiles
gsettings set org.gnome.desktop.input-sources xkb-options "['caps:escape']"

mkdir -p ~/.config/systemd/user
for service in systemd/*; do
    dotlink "$service" ~/.config/systemd/user/"$(basename "$service")"
    systemctl --user enable --now power-monitor
done
