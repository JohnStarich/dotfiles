#!/bin/zsh

autoload -Uz colors && colors
zmodload zsh/stat

js_docker_prompt=''
declare -i js_docker_mtime=0

function js_docker_gen_prompt() {
    if [[ -n "$DOCKER_MACHINE_NAME" ]]; then
        js_docker_prompt=$(js_docker_prompt_text "$DOCKER_MACHINE_NAME")
        return 0
    fi
    local docker_config=~/.docker/config.json
    if [[ ! -f "$docker_config" ]]; then
        js_docker_prompt=''
        return 0
    fi
    declare -i current_mtime=$(stat +mtime "$docker_config" 2>&- || echo 0)
    if (( js_docker_mtime >= current_mtime )); then
        return 0
    fi
    js_docker_mtime=$current_mtime

    local context_name=$(jq -r .currentContext "$docker_config")
    if [[ -n "$context_name" && "$context_name" != null ]]; then
        js_docker_prompt=$(js_docker_prompt_text "$context_name")
    fi
}

autoload -Uz add-zsh-hook 
add-zsh-hook precmd js_docker_gen_prompt

function js_docker_prompt_text() {
    print "%{${fg[green]}%}docker:(%{${fg[yellow]}%}$*%{${fg[green]}%})%{$reset_color%} "
}
