# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG VERSION=v1.14.8

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
    libboost-program-options-dev \
    libssl-dev \
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

RUN git clone --single-branch --branch "${VERSION}" https://github.com/dogecoin/dogecoin.git /opt/dogecoin

WORKDIR /opt/dogecoin

RUN <<-EOF
    bash ./autogen.sh && \
    bash ./configure --enable-hardening --disable-wallet --without-gui --disable-tests --disable-bench --prefix=/opt/dogecoin/build && \
    echo $(getconf _NPROCESSORS_ONLN) && \
    make -j $(getconf _NPROCESSORS_ONLN) && \
    make install && \
    mkdir -p /opt/dogecoin/.dogecoin/db && \
    chmod 755 /opt/dogecoin/.dogecoin/db
EOF

COPY ./dogecoin.conf /opt/dogecoin/.dogecoin/dogecoin.conf

FROM ubuntu:22.04
COPY --from=build /opt/dogecoin/ /opt/dogecoin/

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/dogecoin/build/bin:${PATH}"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    libevent-dev \
    libboost-system-dev \
    libboost-filesystem-dev \
    libboost-chrono-dev \
    libboost-test-dev \
    libboost-thread-dev \
    libboost-program-options-dev \
    net-tools \
    curl \
    jq \
    netcat
EOF

ENTRYPOINT ["dogecoind", "-conf=/opt/dogecoin/.dogecoin/dogecoin.conf", "-datadir=/opt/dogecoin/.dogecoin/db", "-printtoconsole", "-txindex=1"]