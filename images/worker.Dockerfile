FROM docker.io/archlinux

RUN pacman -Syu base-devel git --noprogressbar --noconfirm --needed

RUN rm /var/cache/pacman/pkg/*

RUN echo "builder ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN echo -e '\
[aurer]\n\
Server = file:///repo\n\
SigLevel = Optional TrustAll'\
>> /etc/pacman.conf

COPY images/scripts/* /usr/local/bin/

RUN mkdir /repo

RUN mkdir /work

ENV PKGDEST /repo

WORKDIR /work

CMD setup.sh