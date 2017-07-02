#!/bin/bash

BASH_PROMPT="${USER/johnstarich/js}@${HOSTNAME:0:9}:\W$ "
# used to be PS1_TMP=...
#PS1="${PS1/\S+\s$/$PS1_TMP/}"

function pcd {
	if [[ ! -z "$@" ]]; then
		cd "$@"
	fi
	cd "$(pwd -P)"
}
