#!/bin/bash

if ! pgrep -x "wofi" >/dev/null; then
  wofi --show drun -I &
else
  killall wofi
fi
