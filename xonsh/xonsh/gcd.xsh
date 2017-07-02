#!/usr/bin/env xonsh

def gcd(args, stdin=None) -> None:
    root_dir = $(git rev-parse --show-toplevel).strip()
    if len(args) == 0:
        os.chdir(root_dir)
    elif len(args) == 1:
        os.chdir(os.path.join(root_dir, ' '.join(args)))
