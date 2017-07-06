# John's dotfiles
My config files and automated installation scripts I use on a day-to-day basis.

## Setup

To set up these dotfiles, simply run `./install.sh` from the root of this repository.

If you prefer, you can install just a few modules by typing them out like `./install.sh bash tmux`

This installation script will attempt to only install dotfile symlinks. If those symlinks would overwrite a file, then it skips installing that particular dotfile. However, I am not responsible for any potential damages this script may cause.
