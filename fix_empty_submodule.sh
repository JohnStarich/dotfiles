#!/bin/bash

path=$1

function usage() {
    echo "Usage: $(basename "$0") REPO_PATH" >&2
}

if [[ -z "$path" ]]; then
    usage
    exit 2
fi
if [[ ! -e ./.git ]]; then
    echo "Repo path is not a git repository: $path" >&2
    echo 'Try a submodule init first.' >&2
    usage
    exit 2
fi

if ! ls -A | grep -q -v ".git"; then
    git reset HEAD -- \*
    git checkout -- \*
fi
