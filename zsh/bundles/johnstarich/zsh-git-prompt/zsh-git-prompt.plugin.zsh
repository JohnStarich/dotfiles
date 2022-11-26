#!/bin/zsh

function js_git_prompt_ternary_string() {
    local success=true
    local failure=false
    if [[ "$1" == '!' ]]; then
        shift
        success=false
        failure=true
    fi
    if "$@" >&- 2>&-; then
        echo "$success"
    else
        echo "$failure"
    fi
}

autoload -Uz colors && colors
js_ret_prompt="%(?:%{${fg[green]}%}➜:%{${fg[red]}%}➜)"
function js_git_prompt() {
    local ref
    ref=$(git symbolic-ref HEAD 2> /dev/null) || ref=$(git rev-parse --short HEAD 2> /dev/null) || return 0
    ref=${ref#refs/heads/}
    echo -n "%{${fg[blue]}%}git:(%{${fg[red]}%}$ref%{${fg[blue]}%})%{$reset_color%} "

    local has_modifications=$(js_git_prompt_ternary_string ! git diff --no-ext-diff --quiet --exit-code)
    local has_untracked=$(js_git_prompt_ternary_string git ls-files --others --exclude-standard --error-unmatch .)
    if [[ "$has_modifications" == true || "$has_untracked" == true ]]; then
        echo -n "%{${fg[yellow]}%}✗%{$reset_color%} "
    fi
}
