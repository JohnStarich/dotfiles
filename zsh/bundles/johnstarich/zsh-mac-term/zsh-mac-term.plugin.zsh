#!/bin/zsh

mac_term_light_background=false
function mac_term_check_light_background() {
    mac_term_light_background=$(mac-term-background-ternary false true)
}

autoload -Uz add-zsh-hook 
PERIOD=${PERIOD:-10} # seconds between calls
add-zsh-hook periodic mac_term_check_light_background
