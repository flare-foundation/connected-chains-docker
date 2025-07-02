# syntax=docker/dockerfile:1.3-labs
FROM debian:12 as build

ARG VERSION=2.5.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    git \
    build-essential \
    cmake \
    pkg-config \
    python3 \
    python3.11-venv \
EOF

RUN <<-EOF
    python3 -m venv /opt/conan_env && \
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
    conan profile update conf.tools.build:jobs=1 default && \
    conan export external/snappy snappy/1.1.10@ && \
    mkdir cmake_build && \
    cd cmake_build && \
    conan install \
        .. \
        --output-folder . \
        --build missing \
        --settings build_type=Release && \
    cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_TOOLCHAIN_FILE:FILEPATH=build/generators/conan_toolchain.cmake \
        -DCMAKE_INSTALL_PREFIX=/opt/ripple/build \
        -Dstatic=ON \
        -Dxrpld=ON \
        -Dtests=OFF \
        -B./ ..
EOF

# Job count of 8 consumes ~10GB of memory at peak
RUN <<-EOF
    cd cmake_build && \
    cmake --build . --target install -- -j 8 && \
    strip -s /opt/ripple/build/bin/* && \
    mkdir -p /opt/ripple/.ripple/db
EOF

COPY rippled.conf validators.txt /opt/ripple/.ripple/

RUN chown -R 65532:65532 /opt/ripple/

FROM gcr.io/distroless/base-debian12:nonroot as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build /opt/ripple/build/bin/ /opt/ripple/bin/
COPY --from=build /opt/ripple/.ripple/ /opt/ripple/.ripple/

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["/opt/ripple/bin/rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]
