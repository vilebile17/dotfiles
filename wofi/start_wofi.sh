#!/bin/bash

if ! pgrep -x "wofi" >/dev/null; then
  wofi --show drun -I &
else
  echo "Wofi is already running."
fi
