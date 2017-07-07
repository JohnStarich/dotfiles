#!/bin/bash

if [ ! -f ~/.ssh/known_hosts ]; then
    echo '~/.ssh/known_hosts does not exist'
    return 1
fi

if [[ "$SHELL" == *bash* ]]; then
	complete -W "$(echo $(cat ~/.ssh/known_hosts | cut -f 1 -d ' ' | sed 's/,.*//g' | uniq | grep -v "\["))" ssh
fi
