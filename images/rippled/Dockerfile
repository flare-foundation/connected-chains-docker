# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:24.04 as build

ARG VERSION=2.5.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    git \
    build-essential \
    cmake \
    pkg-config \
    python3-pip \
    python3.12-dev \
    python3.12-venv \
    net-tools \
    curl \
    jq \
    netcat-openbsd
EOF

RUN <<-EOF
    python3.12 -m venv /opt/conan_env && \
    /opt/conan_env/bin/python -m pip install --upgrade pip setuptools && \
    /opt/conan_env/bin/python -m pip install "conan<2"
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/ripple/rippled.git /opt/ripple

WORKDIR /opt/ripple

ENV PATH="/opt/conan_env/bin:$PATH"

# Increasing job count increases risk of race conditions occuring when installing OpenSSL (conan bug?)
RUN <<-EOF
    git config --global http.postBuffer 1048576000 && \
    conan profile new default --detect && \
    conan profile update settings.compiler.cppstd=20 default && \
    conan profile update conf.tools.build:jobs=4 default && \
    conan export external/snappy snappy/1.1.10@ && \
    mkdir cmake_build && \
    cd cmake_build && \
    conan install .. --output-folder . --build missing --settings build_type=Release && \
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_TOOLCHAIN_FILE:FILEPATH=build/generators/conan_toolchain.cmake -DCMAKE_INSTALL_PREFIX=/opt/ripple/build -Dstatic=OFF -Dxrpld=ON -Dtests=OFF -B./ ..
EOF

# Job count of 8 consumes ~10GB of memory at peak
RUN <<-EOF
    cd cmake_build && \
    cmake --build . --target install -- -j 8 && \
    chown -R 65532:65532 /opt/ripple/build
EOF

RUN <<-EOF
    mkdir -p /opt/ripple/.ripple/db && \
    chown 65532:65532 /opt/ripple/.ripple/db
EOF

FROM gcr.io/distroless/cc-debian12:nonroot as final

COPY --from=build /opt/ripple/build /opt/ripple/

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"
ENV PATH="/opt/ripple/bin:${PATH}"

# Tools
COPY --from=build /usr/bin/curl /usr/bin/curl
COPY --from=build /usr/bin/jq /usr/bin/jq
COPY --from=build /usr/bin/netcat /usr/bin/netcat
COPY --from=build /usr/bin/netstat /usr/bin/netstat

# Ripple data
COPY --from=build /opt/ripple/.ripple/db /opt/ripple/.ripple/db 

# Ripple configuration
COPY ./rippled.conf /opt/ripple/.ripple/rippled.conf
COPY ./validators.txt /opt/ripple/.ripple/validators.txt

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]