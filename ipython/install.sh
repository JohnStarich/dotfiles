#!/bin/bash

mkdir -p ~/.ipython/profile_default
dotlink ipython_config.py ~/.ipython/profile_default/ipython_config.py
pip3 install ipython3 lazy astunparse
