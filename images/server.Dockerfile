FROM docker.io/library/archlinux

RUN pacman -Syu --noconfirm go npm gcc

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY ./web/frontend /build/web/frontend

WORKDIR /build/web/frontend

RUN npm i

RUN npm run build

WORKDIR /build/cmd/aurer/

COPY . /build

RUN go build -tags prod

FROM docker.io/library/archlinux

COPY --from=0 /build/cmd/aurer/aurer .

CMD ./aurer

EXPOSE 8008

LABEL aurer.server=1