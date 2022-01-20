#!/bin/bash

git clone https://aur.archlinux.org/"$PACKAGE".git /work

makepkg -s -c --noconfirm --noprogressbar

repo-add -R /repo/aurer.db.tar "$(makepkg --packagelist)"