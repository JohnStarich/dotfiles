#!/usr/bin/env xonsh

def time_xonsh(*args, **kwargs) -> float:
    time_output = $(time xonsh -ilc '' err>out)
    total_time = time_output.lstrip().split(sep=' ', maxsplit=1)[0]
    return float(total_time)

def avg(l: list, ndigits=None) -> float:
    actual_list = list(l)
    average = sum(actual_list) / len(actual_list)
    return round(average, ndigits=ndigits)

def time_xonsh_n(n: int) -> float:
    return avg(map(time_xonsh, range(n)), ndigits=3)

def timesh(args, stdin=None) -> str:
    if len(args) > 1:
        raise Error("Too many arguments. Usage: timesh [trials]")
    trials = 1
    if len(args) == 1:
        if not args[0].isdecimal():
            raise TypeError("Trials must be an integer. Usage: timesh [trials]")
        trials = int(args[0])
    print(time_xonsh_n(trials))
