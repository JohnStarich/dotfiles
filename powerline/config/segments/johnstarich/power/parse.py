#!/usr/bin/env python3

import subprocess
import sys
import re


power_regex = re.compile(
    r"""
    (?P<percentage> \d+)%; \s
    (?P<status> [^;]+); \s
    (?P<time> (?: \d+:\d+ \s)? .+) \s
    present .*
    """,
    flags=re.X,
)


def raw_power_info() -> dict:
    """
    Runs pmstat utility to get power information.
    Note: calls subprocess, so use this call sparingly
    """
    process = subprocess.run(['pmset', '-g', 'ps'], stdout=subprocess.PIPE,
                             encoding='utf8')
    output = process.stdout
    match = power_regex.search(output)
    if match is None:
        return None
    return {
        'percentage': match.group("percentage"),
        'status': match.group("status"),
        'time': match.group("time"),
    }


def power_info(format_string: str='{percentage}% {status}; {time}') -> str:
    return format_power_info(raw_power_info(), format_string)


def format_power_info(info: dict, format_string: str) -> str:
    return format_string.format(
        p=info['percentage'],
        s=info['status'],
        t=info['time'],
        **info
    )


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print(power_info())
    else:
        format_string = ' '.join(sys.argv[1:])
        print(power_info(format_string))
