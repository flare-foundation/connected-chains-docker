# syntax=docker/dockerfile:1.3-labs
FROM debian:12@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as build

ARG VERSION=v1.14.9

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    set -e
    apt-get -y update && apt-get -y install \
        git \
        build-essential \
        cmake \
        libtool \
        pkg-config \
        curl
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/dogecoin/dogecoin.git /opt/dogecoin

WORKDIR /opt/dogecoin

RUN <<-EOF
    set -e
    make -C /opt/dogecoin/depends/ -j $(getconf _NPROCESSORS_ONLN) NO_QT=1 NO_QR=1 NO_WALLET=1 NO_BDB=1
    ./autogen.sh
    ./configure \
        --disable-wallet \
        --without-gui \
        --disable-tests \
        --disable-bench \
        --prefix=/opt/dogecoin/depends/x86_64-pc-linux-gnu/
    make -j $(getconf _NPROCESSORS_ONLN)
    make install
    find /opt/dogecoin/depends/x86_64-pc-linux-gnu/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
EOF

COPY dogecoin.conf /opt/dogecoin/.dogecoin/dogecoin.conf

RUN <<-EOF
    set -e
    mkdir -p /opt/dogecoin/.dogecoin/db
EOF

FROM gcr.io/distroless/cc-debian12:nonroot@sha256:d1b8e4c52be1111aa108e959ef2a822299bb70fd1819dd250871a2601ca1e4b6 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build --chown=65532:65532 /opt/dogecoin/depends/x86_64-pc-linux-gnu/bin/ /opt/dogecoin/bin/
COPY --from=build --chown=65532:65532 /opt/dogecoin/.dogecoin/ /opt/dogecoin/.dogecoin/

ENTRYPOINT ["/opt/dogecoin/bin/dogecoind", "-conf=/opt/dogecoin/.dogecoin/dogecoin.conf", "-datadir=/opt/dogecoin/.dogecoin/db", "-printtoconsole", "-txindex=1"]
