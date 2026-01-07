# syntax=docker/dockerfile:1.3-labs
FROM debian:bookworm@sha256:93492d1405a072c9d3b89110490e330e0b7eb37754cafc9520908b7fe6873f5f as build

ARG VERSION=v30.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

# https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md#linux-distribution-specific-instructions
RUN <<-EOF
    set -e
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
    set -e
    make -C depends/ -j $(getconf _NPROCESSORS_ONLN) NO_QT=1 NO_QR=1 NO_WALLET=1 NO_IPC=1
    cmake -B build \
        --toolchain /opt/bitcoin/depends/x86_64-pc-linux-gnu/toolchain.cmake \
        -DCMAKE_INSTALL_PREFIX=/opt/bitcoin/build/
    cmake --build build -j $(getconf _NPROCESSORS_ONLN)
    cmake --install build
    find /opt/bitcoin/build/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
EOF

RUN mkdir -p /opt/bitcoin/.bitcoin/db

FROM gcr.io/distroless/cc-debian12:nonroot@sha256:189bd2ce1f7750193c2c10220d9201ba38c11e30fbb75b036606829fadbc81b1 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build --chown=65532:65532 /opt/bitcoin/build/bin /opt/bitcoin/bin
COPY --from=build --chown=65532:65532 /opt/bitcoin/.bitcoin/db /opt/bitcoin/.bitcoin/db
COPY --chown=65532:65532 bitcoin.conf /opt/bitcoin/.bitcoin/bitcoin.conf

ENTRYPOINT ["/opt/bitcoin/bin/bitcoind", "-conf=/opt/bitcoin/.bitcoin/bitcoin.conf", "-datadir=/opt/bitcoin/.bitcoin/db", "-txindex=1"]
