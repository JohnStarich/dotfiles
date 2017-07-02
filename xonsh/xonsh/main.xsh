#!/usr/bin/env xonsh

def main(args, stdin=None) -> None:
    if not !(which reattach-to-user-namespace):
        if os.isatty(0):
            input('Press enter to install tmux.conf dependency `reattach-to-user-namespace`')
            $[brew install reattach-to-user-namespace]
        else:
            print('Dependency reattach-to-user-namespace is not installed')
            return 1
    if not !(tmux has-session -t main):
        $[tmux new-session -s main]
    else:
        $[tmux attach -t main]
