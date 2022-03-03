# syntax=docker/dockerfile:experimental


# Build ginit as an init
FROM golang:1.17-alpine as dev
RUN apk add --no-cache git ca-certificates gcc linux-headers musl-dev
COPY . /go/src/github.com/thebsdbox/ginit/
WORKDIR /go/src/github.com/thebsdbox/ginit
ENV GO111MODULE=on
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod/cache \
    --mount=type=cache,sharing=locked,id=goroot,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o init
    
#RUN go get; CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o init

# Build Busybox
FROM gcc:10.1.0 as Busybox
RUN apt-get update; apt-get install -y cpio
RUN curl -O https://busybox.net/downloads/busybox-1.31.1.tar.bz2
RUN tar -xf busybox*bz2
WORKDIR busybox-1.31.1
RUN make defconfig; make LDFLAGS=-static CONFIG_PREFIX=./initramfs install

WORKDIR initramfs 
COPY --from=dev /go/src/github.com/thebsdbox/ginit/init .

# Package initramfs
RUN find . -print0 | cpio --null -ov --format=newc > ../initramfs.cpio 
RUN gzip ../initramfs.cpio
RUN mv ../initramfs.cpio.gz /

FROM scratch
COPY --from=Busybox /initramfs.cpio.gz .
