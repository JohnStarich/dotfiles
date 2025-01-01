#!/bin/zsh

autoload -Uz colors && colors

js_toolbox_prompt=''
if [[ -f /run/.containerenv && $(< /run/.containerenv) =~ 'name="([^"]*)"' ]]; then
    local toolbox_name=$match
    if [[ "$toolbox_name" == *"$USER"* ]]; then
        toolbox_name='🛠'
    fi
    local toolbox_color="%{$fg_bold[yellow]%}"
    local reset_color="%{$reset_color%}"
    if [[ -n "$toolbox_name" ]]; then
        js_toolbox_prompt="${toolbox_color}${toolbox_name}${reset_color} "
    fi
fi
