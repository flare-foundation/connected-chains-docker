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
    libzmq3-dev \
    net-tools \
    curl \
    jq \
    netcat
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/bitcoin/bitcoin.git /opt/bitcoin

WORKDIR /opt/bitcoin

RUN <<-EOF
   cmake -B build -DENABLE_WALLET=OFF -DWITH_ZMQ=ON
   cmake --build build -j 4
   cmake --install build
EOF

RUN mkdir -p /opt/bitcoin/.bitcoin/db && \
    chown -R 65532:65532 /opt/bitcoin/

FROM gcr.io/distroless/cc-debian12:nonroot as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build /usr/local/bin /usr/local/bin

# Libraries (ldd /usr/local/bin/bitcoind)
COPY --from=build /usr/lib/x86_64-linux-gnu/libevent*.so* /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libzmq.so.5 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libbsd.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libsodium.so.23 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libpgm-5.3.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libnorm.so.1 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libmd.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libkrb5.so.3 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libk5crypto.so.3 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libcom_err.so.2 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libkrb5support.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=build /usr/lib/x86_64-linux-gnu/libkeyutils.so.1 /usr/lib/x86_64-linux-gnu/

# Tools
COPY --from=build /usr/bin/curl /usr/bin/curl
COPY --from=build /usr/bin/jq /usr/bin/jq
COPY --from=build /usr/bin/netcat /usr/bin/netcat
COPY --from=build /usr/bin/netstat /usr/bin/netstat

# Bitcoin Core
COPY --from=build /opt/bitcoin/.bitcoin/db /opt/bitcoin/.bitcoin/db

ENTRYPOINT ["bitcoind", "-conf=/opt/bitcoin/.bitcoin/bitcoin.conf", "-datadir=/opt/bitcoin/.bitcoin/db", "-txindex=1"]
