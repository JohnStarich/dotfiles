#!/bin/bash

mkdir -p ~/.config
dotlink config ~/.config/powerline

packages=(
    dnspython
    maxminddb
    git+https://github.com/powerline/powerline
    requests
)
dotpip3 "${packages[@]}"

if linux; then
    db=~/.local/lib/johnstarich-powerline/dbip-city.mmdb
    if [[ ! -f "$db" ]]; then
        mkdir -p "$(dirname "$db")"
        curl -s https://download.db-ip.com/free/dbip-city-lite-2023-10.mmdb.gz | gunzip > "$db"
    fi
fi
