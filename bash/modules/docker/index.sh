#!/bin/bash

if which docker-machine 2>&1 > /dev/null; then
	DM_COLOR_RESET='\033[0;00m'
	function _docker_machine_name {
		if [ "$DOCKER_MACHINE_NAME" ]; then
			echo -n "$DOCKER_MACHINE_NAME"
		elif [ "$DOCKER_HOST" ]; then
			# Remove host protocol and port number to clean up the look
			#echo -n "$DOCKER_HOST" | sed -E 's/^.+:\/\///' | sed -E 's/:[0-9]+//'
			# Remove all but first part of host name
			echo -n "$DOCKER_HOST" | sed -E 's/^.+:\/\///' | sed -E 's/\..*//' | tr -d '\n'
		fi
	}
	function _docker_machine_wrap_color {
		echo -en "\001$@\002"
	}
	function _docker_machine_color {
		DM_RAND="$(echo "$(_docker_machine_name)" | md5 | cut -c1-4)"
		DM_RAND=$((31 + ( (16#${DM_RAND}) % 7) ))
		DM_COLOR="\033[1;${DM_RAND}m"
		_docker_machine_wrap_color "$DM_COLOR"
    }
	function _docker_machine_ps1 {
		if [ "$(_docker_machine_name)" ]; then
			echo -n '['
			_docker_machine_color
			_docker_machine_name
			_docker_machine_wrap_color "\033[m"
			echo -n '] '
		fi
	}
	BASH_PROMPT_PREFIX+="\$([[ -n \"\$(type -t _docker_machine_ps1)\" ]] && _docker_machine_ps1)"
fi
