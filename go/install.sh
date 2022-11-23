#!/usr/bin/env bash

if ! which go; then
    echo Install Go to continue.
    exit 1
fi

GOBIN="$DOTFILES_DIR/bin" go install github.com/johnstarich/go/goop/cmd/goop@latest
