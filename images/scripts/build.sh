#!/bin/bash
set -e 

sudo pacman -Syu --noconfirm --noprogressbar

git clone https://aur.archlinux.org/"$PACKAGE".git /work

makepkg -f -s -c --noconfirm --noprogressbar

repo-add -R /repo/aurer.db.tar "$(makepkg --packagelist)"
