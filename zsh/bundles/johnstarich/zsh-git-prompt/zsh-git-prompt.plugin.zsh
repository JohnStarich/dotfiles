#!/bin/zsh

autoload -Uz colors && colors

js_git_prompt=''
js_git_prev_dir=''
js_git_cache_main_ref=''
function js_git_gen_prompt() {
    local current_git_dir
    if ! current_git_dir=$(git rev-parse --show-toplevel --quiet 2>&-); then
        js_git_prompt=''
        return 0
    fi
    if [[ "$js_git_prev_dir" != "$current_git_dir" || "$js_git_cache_main_ref" == '' ]]; then
        js_git_cache_main_ref=$(git symbolic-ref HEAD 2> /dev/null) || js_git_cache_main_ref=$(git rev-parse --short HEAD 2> /dev/null) || true
        js_git_cache_main_ref=${js_git_cache_main_ref#refs/heads/}
        js_git_prev_dir=$current_git_dir
    fi

    js_git_prompt="%{${fg[blue]}%}git:(%{${fg[red]}%}$js_git_cache_main_ref%{${fg[blue]}%})"

    if ! git diff --no-ext-diff --quiet --exit-code >&- 2>&-; then
        # Has modifications
        js_git_prompt+=" %{${fg[yellow]}%}✗"
    elif git ls-files --others --exclude-standard --error-unmatch . >&- 2>&-; then
        # Has untracked files
        js_git_prompt+=" %{${fg[green]}%}✔"
    fi
    js_git_prompt+="%{$reset_color%} "
}

autoload -Uz add-zsh-hook 
add-zsh-hook precmd js_git_gen_prompt
