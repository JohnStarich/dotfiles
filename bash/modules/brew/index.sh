#!/bin/bash

if [[ "$SHELL" == *bash* && -f "$BREW_PREFIX"/etc/bash_completion ]]; then
	source "$BREW_PREFIX"/etc/bash_completion
fi

