#!/bin/bash

sudo pacman -Syu --noconfirm

sudo chown -R build:build /work /output

git clone https://aur.archlinux.org/"$PACKAGE".git /work

makepkg -s -c --noconfirm

repo-add /output/aurer.db.tar /output/"$PACKAGE"*.pkg.tar.zst