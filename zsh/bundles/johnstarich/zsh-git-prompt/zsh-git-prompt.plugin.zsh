#!/bin/zsh

autoload -Uz colors && colors
js_ret_prompt="%(?:%{${fg[green]}%}➜:%{${fg[red]}%}➜)"
function js_git_prompt() {
    local ref
    ref=$(git symbolic-ref HEAD 2> /dev/null) || ref=$(git rev-parse --short HEAD 2> /dev/null) || return 0
    ref=${ref#refs/heads/}
    echo -n "%{${fg[blue]}%}git:(%{${fg[red]}%}$ref%{${fg[blue]}%})%{$reset_color%} "
    if ! git diff --no-ext-diff --quiet --exit-code; then
        echo -n "%{${fg[yellow]}%}✗%{$reset_color%} "
    fi
}
