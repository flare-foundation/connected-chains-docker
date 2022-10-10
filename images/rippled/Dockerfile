# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:20.04 as build

ARG VERSION=1.9.4

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
apt-get -y update && \
apt-get -y install \
git \
build-essential \
cmake \
gcc \
g++ \
protobuf-compiler \
libprotobuf-dev \
flex \
bison \
libssl-dev \
libboost-all-dev \
pkg-config
EOF

RUN <<-EOF
    echo "* soft nofile 1048576" >> /etc/security/limits.conf && \
    echo "* hard nofile 1048576" >> /etc/security/limits.conf && \
    echo "{{ ops_user }} soft nofile 1048576" >> /etc/security/limits.conf && \
    echo "{{ ops_user }} hard nofile 1048576" >> /etc/security/limits.conf && \
    cat /etc/security/limits.conf && \
    ulimit -Sn
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/ripple/rippled.git /opt/ripple

WORKDIR /opt/ripple

RUN <<-EOF
    git config --global http.postBuffer 1048576000 && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/opt/ripple/build -Dstatic=OFF -B./cmake_build
EOF

RUN <<-EOF
    cd cmake_build && \
    cmake --build . --target install -- -j 1 && \
    mkdir -p /opt/ripple/.ripple/db && \
    chmod 755 /opt/ripple/.ripple/db
EOF

RUN rm -rf .git src .nih_c cmake_build

COPY ./rippled.conf /opt/ripple/.ripple/rippled.conf
COPY ./validators.txt /opt/ripple/.ripple/validators.txt

FROM ubuntu:20.04
COPY --from=build /opt/ripple/ /opt/ripple/

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/ripple/build/bin:${PATH}"

RUN <<-EOF
apt-get -y update && \
apt-get -y install \
libboost-all-dev \
net-tools \
curl \
jq \
netcat
EOF

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]