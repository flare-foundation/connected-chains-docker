# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG VERSION=v29.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

# https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md#linux-distribution-specific-instructions
RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    git \
    build-essential \
    cmake \
    pkgconf \
    python3 \
    libevent-dev \
    libboost-dev \
    libzmq3-dev
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/bitcoin/bitcoin.git /opt/bitcoin

WORKDIR /opt/bitcoin

RUN <<-EOF
   cmake -B build -DENABLE_WALLET=OFF -DWITH_ZMQ=ON
   cmake --build build -j 4
   cmake --install build
EOF

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    libevent-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libboost-chrono-dev \
    libboost-test-dev \
    libboost-thread-dev \
    libzmq3-dev \
    net-tools \
    curl \
    jq \
    netcat
EOF

RUN groupadd --gid "10001" "bitcoin" && \
    useradd --uid "10001" --gid "10001" --shell /bin/bash --create-home "bitcoin" && \
    mkdir -p /opt/bitcoin/.bitcoin/db && \
    chown -R bitcoin:bitcoin /opt/bitcoin/

COPY --from=build /usr/local/bin /usr/local/bin

ENTRYPOINT ["bitcoind", "-conf=/opt/bitcoin/.bitcoin/bitcoin.conf", "-datadir=/opt/bitcoin/.bitcoin/db", "-txindex=1"]
