#!/usr/bin/env xonsh

def glog(args, stdin=None) -> None:
    command = ['git', 'log', '--oneline', '--graph', '--decorate=full']
    if len(args) >= 1:
        command.append('--since={}'.format(args[0]))
    if len(args) > 1:
        command += args[1:]
    $[@(command)]
