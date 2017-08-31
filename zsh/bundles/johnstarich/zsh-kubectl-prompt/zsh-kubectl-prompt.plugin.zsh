#!/bin/zsh

_kubectl_prompt_oldconfig=''
_kubectl_prompt_context=''

function js_kubectl_prompt() {
    local kube_color="%{$fg_bold[magenta]%}"
    local context_color="%{$fg_bold[red]%}"
    local reset_color="%{$reset_color%}"
    if [[ ! -z "$KUBECONFIG" ]]; then
        if [[ "$_kubectl_prompt_oldconfig" != "$KUBECONFIG" ]]; then
            _kubectl_prompt_context=$(kubectl config current-context)
            _kubectl_prompt_oldconfig=$KUBECONFIG
        fi
        print "${kube_color}kube:(${context_color}${_kubectl_prompt_context}${kube_color})${reset_color} "
    fi
}
