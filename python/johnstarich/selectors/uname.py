#!/usr/bin/env python3

import subprocess

def get_kernel() -> str:
    return subprocess.check_output('uname').decode().strip()

def is_macos(pl, segment_info, mode):
    return get_kernel() == 'Darwin'

def is_linux(pl, segment_info, mode):
    return get_kernel() == 'Linux'
