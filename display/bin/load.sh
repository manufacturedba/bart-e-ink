#!/bin/bash
# Dumb script to reload files to lib

ampy --port $1 rmdir lib
ampy --port $1 put lib