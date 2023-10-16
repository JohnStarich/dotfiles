#!/usr/bin/env bash

if ! which dnf >/dev/null 2>/dev/null; then
    exit 0
fi

packages=(
    git-delta
    powerline-fonts
    zeal
)

installed_packages=$(dnf list --installed)
for pkg in "${packages[@]}"; do
    if ! [[ "$installed_packages" =~ (^|$'\n')"$pkg"\. ]]; then
        echo "Package $pkg not installed. Running bulk install...";
        sudo dnf install -y "${packages[@]}"
        break
    fi
done
