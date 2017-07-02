#!/bin/bash

function gitff-all {
	args=$@
	(	# Create subshell to prevent variables from influencing current shell
		# Check if there are git repos in here
		shopt -s nullglob
		set -- */.git
		if [ "$#" -le 0 ]; then
			echo 'Error finding git repo directories in this directory.' >&2
			return 2
		fi
		for dir in */.git; do
			(
				dir=$(dirname $dir)
				echo "Updating $dir... "
				cd $dir
				if [ -z "$(git remote)" ]; then
					echo "Repo $dir does not have a remote, skipping update" >&2
				elif [[ "$(git rev-parse --abbrev-ref HEAD)" != 'master' ]]; then
					echo "Repo $dir is not on master, skipping update" >&2
				else
					gitff $args
					if [[ $? == 0 ]]; then
						echo "Updated $dir."
					else
						echo "Failed to update $dir"
					fi
				fi
				cd ..
			) &
		done
		wait
	)
}

function glog {
	local since=--since="$1"
	local args=${@:2}
	if [ -z "$1" ]; then
		unset since
		args=$@
	fi
    git log --oneline --graph $since $args
}
