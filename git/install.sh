#!/bin/bash

macos definitely_me dotlink "$HOME/Library/Mobile Documents/com~apple~Automator/Documents" ~/ibin
if ! macos definitely_me dotlink /usr/local/bin/pinentry-mac ~/bin/pinentry-mac >/dev/null; then
    macos definitely_me dotlink /opt/homebrew/bin/pinentry-mac ~/bin/pinentry-mac
fi
macos definitely_me dotlink ~/ibin/gnupg ~/.gnupg
macos definitely_me chmod 700 ~/.gnupg
definitely_me dotlink gitconfig ~/.gitconfig
macos dotlink gitconfig_macos ~/.gitconfig_macos
linux dotlink gitconfig_linux ~/.gitconfig_linux

dotlink gitignore ~/.gitignore
dotlink hooks ~/.githooks
