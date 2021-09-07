FROM docker.io/archlinux

RUN echo "NoProgressBar" >> /etc/pacman.conf

RUN pacman -Syu base-devel git --noconfirm --needed

RUN useradd -m build

RUN echo "build ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN mkdir /output

RUN mkdir /work

RUN chown build:build /work /output

ENV PKGDEST /output

USER build

WORKDIR /work

COPY build.sh /usr/local/bin/build.sh

CMD build.sh