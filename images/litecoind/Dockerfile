# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG VERSION=v0.21.3

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

# Remove libdb++-dev after upstream fixes the problem with --disable-wallet requiring it
RUN <<-EOF
    apt-get -y update && \
    apt-get -y install software-properties-common && \
    add-apt-repository universe && \
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
    libssl-dev \
    libuv1-dev \
    acl \
    libfmt-dev \
    libdb++-dev
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

RUN git clone --single-branch --branch "${VERSION}" https://github.com/litecoin-project/litecoin.git /opt/litecoin

WORKDIR /opt/litecoin

RUN <<-EOF
    bash ./autogen.sh && \
    bash ./configure --enable-hardening --disable-wallet --without-gui --disable-tests --disable-bench --prefix=/opt/litecoin/build && \
    echo $(getconf _NPROCESSORS_ONLN) && \
    make -j $(getconf _NPROCESSORS_ONLN) && \
    make install && \
    mkdir -p /opt/litecoin/.litecoin/db && \
    chmod 755 /opt/litecoin/.litecoin/db
EOF

COPY ./litecoin.conf /opt/litecoin/.litecoin/litecoin.conf

FROM ubuntu:22.04
COPY --from=build /opt/litecoin/ /opt/litecoin/

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/litecoin/build/bin:${PATH}"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    libevent-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libboost-chrono-dev \
    libboost-test-dev \
    libboost-thread-dev \
    libfmt-dev \
    net-tools \
    curl \
    jq \
    netcat
EOF

RUN ls -la /opt/litecoin

ENTRYPOINT ["litecoind", "-conf=/opt/litecoin/.litecoin/litecoin.conf", "-datadir=/opt/litecoin/.litecoin/db", "-txindex=1"]