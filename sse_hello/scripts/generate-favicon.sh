#!/bin/bash

# Install ImageMagick if not already installed
# On macOS: brew install imagemagick
# On Ubuntu: sudo apt-get install imagemagick

# Convert SVG to ICO
magick public/clock.svg -background transparent -define icon:auto-resize=32,16 public/favicon.ico 