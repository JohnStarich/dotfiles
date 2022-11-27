#!/bin/zsh

autoload -Uz colors && colors
zmodload zsh/stat

js_kubectl_prompt=''
declare -i js_kubectl_mtime=0

function js_kubectl_gen_prompt() {
    local kube_color="%{$fg_bold[magenta]%}"
    local context_color="%{$fg_bold[red]%}"
    local reset_color="%{$reset_color%}"
    if [[ -z "$KUBECONFIG" || ! -e "$KUBECONFIG" ]]; then
        js_kubectl_prompt=''
        js_kubectl_mtime=0
        return 0
    fi
    declare -i current_mtime=$(stat +mtime "$KUBECONFIG" 2>&- || echo 0)
    if (( js_kubectl_mtime < current_mtime )); then
        js_kubectl_mtime=$current_mtime
        js_kubectl_prompt="${kube_color}kube:(${context_color}$(kubectl config current-context)${kube_color})${reset_color} "
    fi
}

autoload -Uz add-zsh-hook 
add-zsh-hook precmd js_kubectl_gen_prompt

function js_kubectl_toggle() {
    if [[ -n "$KUBECONFIG" ]]; then
        unset KUBECONFIG
    else
        export KUBECONFIG=$HOME/.kube/config
    fi
}
