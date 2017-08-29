#!/bin/bash

# Find your version by running `pip --version`
for lib in /usr/local/lib/python{3.7,3.6,3.5,2.7}/{site,dist}-packages; do
	if [[ -d "$lib" ]]; then
		export PYTHON_LIB="$lib"
		break
	fi
done
