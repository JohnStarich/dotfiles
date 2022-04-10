#!/bin/zsh

autoload -Uz colors && colors

function js_docker_prompt_text() {
    print "%{${fg[green]}%}docker:(%{${fg[yellow]}%}$*%{${fg[green]}%})%{$reset_color%} "
}

function js_docker_prompt() {
    if [[ -f ~/.docker/config.json ]]; then
        contextName=$(jq -r .currentContext ~/.docker/config.json)
        if [[ -n "$contextName" && "$contextName" != null ]]; then
            js_docker_prompt_text "$contextName"
        fi
    fi
    if [[ -n "$DOCKER_MACHINE_NAME" ]]; then
        js_docker_prompt_text "$DOCKER_MACHINE_NAME"
    fi
}
