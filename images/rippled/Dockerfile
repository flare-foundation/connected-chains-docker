# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:24.04 as build

ARG VERSION=2.4.0

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
pkg-config \
python3-pip \
python3.12-dev \
python3.12-venv
EOF

RUN <<-EOF
    python3.12 -m venv /opt/conan_env && \
    /opt/conan_env/bin/python -m pip install --upgrade pip setuptools && \
    /opt/conan_env/bin/python -m pip install "conan<2"
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

ENV PATH="/opt/conan_env/bin:$PATH"

RUN <<-EOF
    git config --global http.postBuffer 1048576000 && \
    conan profile new default --detect && \
    conan profile update settings.compiler.cppstd=20 default && \
    conan export external/snappy snappy/1.1.10@ && \
    mkdir cmake_build && \
    cd cmake_build && \
    conan install .. --output-folder . --build missing --settings build_type=Release && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_TOOLCHAIN_FILE:FILEPATH=build/generators/conan_toolchain.cmake -DCMAKE_INSTALL_PREFIX=/opt/ripple/build -Dstatic=OFF -Dxrpld=ON -Dtests=ON -B./ ..
EOF

RUN <<-EOF
    cd cmake_build && \
    cmake --build . --target install -- -j 8 && \
    ./rippled --unittest
EOF

FROM ubuntu:24.04
COPY --from=build /opt/ripple/build /opt/ripple/

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/ripple/bin:${PATH}"

RUN <<-EOF
apt-get -y update && \
apt-get -y install \
libboost-all-dev \
net-tools \
curl \
jq \
netcat-openbsd
EOF

RUN <<-EOF
    mkdir -p /opt/ripple/.ripple/db && \
    chmod 755 /opt/ripple/.ripple/db
EOF

COPY ./rippled.conf /opt/ripple/.ripple/rippled.conf
COPY ./validators.txt /opt/ripple/.ripple/validators.txt

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]