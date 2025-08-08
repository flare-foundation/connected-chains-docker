# syntax=docker/dockerfile:1.3-labs
FROM debian:12@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as build

ARG VERSION=v29.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

# https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md#linux-distribution-specific-instructions
RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    git \
    build-essential \
    cmake \
    curl
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/bitcoin/bitcoin.git /opt/bitcoin

WORKDIR /opt/bitcoin

RUN <<-EOF
    make -C depends/ -j $(getconf _NPROCESSORS_ONLN) NO_QT=1 NO_QR=1 NO_WALLET=1
    cmake -B build \
        --toolchain /opt/bitcoin/depends/x86_64-pc-linux-gnu/toolchain.cmake \
        -DCMAKE_INSTALL_PREFIX=/opt/bitcoin/build/
    cmake --build build -j $(getconf _NPROCESSORS_ONLN) 
    cmake --install build
    find /opt/bitcoin/build/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
EOF

RUN mkdir -p /opt/bitcoin/.bitcoin/db && \
    chown -R 65532:65532 /opt/bitcoin/

FROM gcr.io/distroless/cc-debian12:nonroot@sha256:d1b8e4c52be1111aa108e959ef2a822299bb70fd1819dd250871a2601ca1e4b6 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build /opt/bitcoin/build/bin /opt/bitcoin/bin
COPY --from=build /opt/bitcoin/.bitcoin/db /opt/bitcoin/.bitcoin/db
COPY bitcoin.conf /opt/bitcoin/.bitcoin/bitcoin.conf

ENTRYPOINT ["/opt/bitcoin/bin/bitcoind", "-conf=/opt/bitcoin/.bitcoin/bitcoin.conf", "-datadir=/opt/bitcoin/.bitcoin/db", "-txindex=1"]
