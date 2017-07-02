#!/bin/bash

function _dmc_completer() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
    if [ $COMP_CWORD -gt 1 ]; then
        prev2="${COMP_WORDS[COMP_CWORD-2]}"
    fi

	if [[ ${prev} == 'dmc' ]]; then
		COMPREPLY=( $(compgen -W "console create restart shutdown start" -- ${cur}) )
		return 0
	elif [[ ${prev2} == 'dmc' && ${prev} != 'create' ]]; then
        local opts1=$(docker ps --all --filter=ancestor=johnstarich/sponge-vanilla:test --format={{.Names}})
		COMPREPLY=( $(compgen -W "${opts1}" -- ${cur}) )
		return 0
	fi
}

if [[ "$SHELL" == *bash* ]]; then
	complete -F _dmc_completer dmc
fi
