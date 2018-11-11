#!/bin/bash

if [[ "$SHELL" == *bash* ]]; then
	complete -W "$([[ -e ~/.ssh/known_hosts ]] && echo $(cat ~/.ssh/known_hosts | cut -f 1 -d ' ' | sed 's/,.*//g' | uniq | grep -v "\["))" ssh
fi
