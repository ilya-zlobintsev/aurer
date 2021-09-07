#!/bin/bash

sudo pacman -Syu --noconfirm

sudo chown build:build /work

git clone https://aur.archlinux.org/"$PACKAGE".git /work

makepkg -s -c --noconfirm --noprogressbar