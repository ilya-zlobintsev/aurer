FROM docker.io/archlinux

RUN pacman -Syu base-devel git --noprogressbar --noconfirm --needed

RUN echo "builder ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

COPY images/scripts/* /usr/local/bin/

RUN mkdir /repo

RUN mkdir /work

ENV PKGDEST /repo

WORKDIR /work

CMD setup.sh