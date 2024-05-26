#!/usr/bin/env bash

if linux; then
    mkdir -p \
        ~/.var/app/org.mozilla.firefox/.mozilla/firefox \
        ~/.local/share/applications \
        ~/.local/share/icons
    for app in apps/*/; do
        app=$(basename "$app")
        dotlink "apps/$app/shortcut.desktop" ~/.local/share/applications/"$app.desktop"
        dotlink "apps/$app/icon.png" ~/.local/share/icons/"$app.png"
        profile_dir=~/.var/app/org.mozilla.firefox/.mozilla/firefox/"dotfiles.$app"
        # NOTE: Flathub apps can't use symlinks pointing outside of their file paths.
        if [[ ! -d "$profile_dir" ]]; then
            flatpak-spawn --host flatpak run org.mozilla.firefox -CreateProfile "$app $profile_dir" --no-remote
        fi
        mkdir -p "$profile_dir"
        cp "apps/$app/user.js" "$profile_dir"/user.js
        if [[ -d "apps/$app/chrome" ]]; then
            mkdir -p "$profile_dir"/chrome
            cp -rf "apps/$app/chrome/"* "$profile_dir"/chrome/
        fi
    done
fi
