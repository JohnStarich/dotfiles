#!/bin/bash

definitely_me dotlink "$HOME/Library/Mobile Documents/com~apple~Automator/Documents" ~/ibin
if ! definitely_me dotlink /usr/local/bin/pinentry-mac ~/bin/pinentry-mac >/dev/null; then
    definitely_me dotlink /opt/homebrew/bin/pinentry-mac ~/bin/pinentry-mac
fi
definitely_me dotlink ~/ibin/gnupg ~/.gnupg
definitely_me chmod 700 ~/.gnupg
definitely_me dotlink gitconfig ~/.gitconfig
macos dotlink gitconfig_macos ~/.gitconfig_macos

dotlink gitignore ~/.gitignore
dotlink hooks ~/.githooks
