#!/bin/zsh

CURRENT_DIR="$(cd "$(dirname $0)" && pwd -P)"
CURRENT_FILE="$(basename $0)"

fpath=("$CURRENT_DIR" $fpath)

autoload -Uz $(
	ls "$CURRENT_DIR" |
	grep -v '\.zwc$' |
	grep -v "^$CURRENT_FILE$"
)
