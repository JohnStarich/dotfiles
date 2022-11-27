#!/bin/zsh

autoload -Uz colors && colors

js_git_prompt=''
function js_git_gen_prompt() {
    local ref
    if ! { ref=$(git symbolic-ref HEAD 2> /dev/null) || ref=$(git rev-parse --short HEAD 2> /dev/null) }; then
        js_git_prompt=''
        return 0
    fi
    ref=${ref#refs/heads/}

    js_git_prompt="%{${fg[blue]}%}git:(%{${fg[red]}%}$ref%{${fg[blue]}%})"

    if js_git_has_modifications; then
        js_git_prompt+=" %{${fg[yellow]}%}✗"
    elif js_git_has_untracked_files; then
        js_git_prompt+=" %{${fg[green]}%}✔"
    fi
    js_git_prompt+="%{$reset_color%} "
}

autoload -Uz add-zsh-hook 
add-zsh-hook precmd js_git_gen_prompt

function js_git_has_modifications() {
    if ! git diff --no-ext-diff --quiet --exit-code >&- 2>&-; then
        return 0
    fi
    return 1
}

function js_git_has_untracked_files() {
    git ls-files --others --exclude-standard --error-unmatch . >&- 2>&-
}

