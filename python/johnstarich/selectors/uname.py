#!/usr/bin/env python3

import subprocess

def is_kernel(kernel):
    uname = subprocess.check_output('uname').decode().strip()
    return uname == kernel

def is_macos(pl, segment_info, mode):
    return is_kernel('Darwin')

def is_linux(pl, segment_info, mode):
    return is_kernel('Linux')
