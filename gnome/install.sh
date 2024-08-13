#!/usr/bin/env bash

if ! linux; then
    exit
fi

dotlink applications ~/.local/share/applications/dotfiles
gsettings set org.gnome.desktop.input-sources xkb-options "['caps:escape']"  # Set Caps Lock to press Escape instead
gsettings set org.gnome.desktop.notifications.application:/org/gnome/desktop/notifications/application/org-gnome-evolution-alarm-notify/ enable-sound-alerts false  # Disable high-pitch, annoying sounds for calendar alerts

mkdir -p ~/.config/systemd/user
for service in systemd/*; do
    name="$(basename "$service")"
    dotlink "$service" ~/.config/systemd/user/"$name"
    systemctl --user enable --now "$name"
done
