FROM docker.io/archlinux

RUN echo "NoProgressBar" >> /etc/pacman.conf

RUN pacman -Syu base-devel git --noconfirm --needed

RUN useradd -m build

RUN echo "build ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers


COPY build.sh /usr/local/bin/build.sh

RUN mkdir /output

RUN mkdir /work

ENV PKGDEST /output

WORKDIR /work

USER build

CMD build.sh