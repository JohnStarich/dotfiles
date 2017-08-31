#!/bin/zsh

autoload -Uz colors && colors

function js_docker_prompt() {
    if [[ -n "$DOCKER_MACHINE_NAME" ]]; then
        print "%{${fg[green]}%}docker:(%{${fg[yellow]}%}$DOCKER_MACHINE_NAME%{${fg[green]}%})%{$reset_color%} "
    fi
}
