#!/bin/bash

function _notes_completer() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	if [[ ${prev} == 'notes' ]]; then
        local class_names=$(ls -d ~/notes/* | sbasename | sed 's/ /\\\\ /g')
        local IFS=$'\n'
		COMPREPLY=( $(compgen -W "${class_names}" -- ${cur}) )
		return 0
	fi
}

if [[ "$SHELL" == *bash* ]]; then
	complete -F _notes_completer notes
fi

#complete -G ~/school/* -d notes
