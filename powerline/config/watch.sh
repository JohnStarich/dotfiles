#!/bin/bash

cd "$(dirname "$0")"
cd "$(pwd -P)/../.."

watch-dir . py,json powerline-daemon --replace
