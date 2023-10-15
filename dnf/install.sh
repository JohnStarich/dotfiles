#!/usr/bin/env bash

if ! which dnf >/dev/null 2>/dev/null; then
    exit 0
fi

packages=(
    powerline-fonts
)

sudo dnf install -y "${packages[@]}"
