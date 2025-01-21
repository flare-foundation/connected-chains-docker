# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG VERSION=v28.1

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    git \
    build-essential \
    libtool \
    autotools-dev \
    automake \
    pkg-config \
    bsdmainutils \
    python3 \
    python3-pip \
    libevent-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libboost-chrono-dev \
    libboost-test-dev \
    libboost-thread-dev \
    libzmq3-dev \
    acl
EOF

RUN <<-EOF
    echo "* soft nofile 1048576" >> /etc/security/limits.conf && \
    echo "* hard nofile 1048576" >> /etc/security/limits.conf && \
    echo "{{ ops_user }} soft nofile 1048576" >> /etc/security/limits.conf && \
    echo "{{ ops_user }} hard nofile 1048576" >> /etc/security/limits.conf && \
    cat /etc/security/limits.conf && \
    ulimit -Sn
EOF

RUN pip install --upgrade pip && pip3 install pexpect

RUN git clone --single-branch --branch "${VERSION}" https://github.com/bitcoin/bitcoin.git /opt/bitcoin

WORKDIR /opt/bitcoin

RUN <<-EOF
    bash ./autogen.sh && \
    bash ./configure --disable-wallet --without-gui --disable-tests --disable-bench --prefix=/opt/bitcoin/build && \
    echo $(getconf _NPROCESSORS_ONLN) && \
    make -j $(getconf _NPROCESSORS_ONLN) && \
    make install && \
    mkdir -p /opt/bitcoin/.bitcoin/db && \
    chmod 755 /opt/bitcoin/.bitcoin/db
EOF

COPY ./bitcoin.conf /opt/bitcoin/.bitcoin/bitcoin.conf

FROM ubuntu:22.04
COPY --from=build /opt/bitcoin/build /opt/bitcoin/build

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/bitcoin/build/bin:${PATH}"

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

ENTRYPOINT ["bitcoind", "-conf=/opt/bitcoin/.bitcoin/bitcoin.conf", "-datadir=/opt/bitcoin/.bitcoin/db", "-txindex=1"]