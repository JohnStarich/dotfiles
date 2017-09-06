#!/bin/zsh

autoload -Uz colors && colors

function js_host_prompt() {
    local hostname_color="%{$fg_bold[yellow]%}"
    local reset_color="%{$reset_color%}"
    if ! [[ -z "$SSH_TTY" && -z "$SSH_CLIENT" ]]; then
        print "${hostname_color}${HOST}${reset_color} "
    fi
}
