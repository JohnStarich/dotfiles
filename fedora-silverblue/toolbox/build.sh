#!/usr/bin/env bash

cd "$(dirname "$0")"
podman build -t localhost/johnstarich-dev:latest -f Containerfile "$@" .
