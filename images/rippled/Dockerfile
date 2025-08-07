# syntax=docker/dockerfile:1.3-labs
# debian:12
FROM debian@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as build

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

# to generate requirements.txt, run:
# echo "conan<2" > requirements.in
# pip-compile --generate-hashes requirements.in
COPY requirements.txt /tmp/requirements.txt

RUN <<-EOF
    python3 -m venv /opt/conan_env && \
    /opt/conan_env/bin/python -m pip install -r /tmp/requirements.txt --require-hashes
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

# gcr.io/distroless/base-debian12:nonroot
FROM gcr.io/distroless/base-debian12@sha256:dca858878977f01b80685240ca2cd9c1bb38958e2ab6f4ed5698c0ea23307143 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build /opt/ripple/build/bin/ /opt/ripple/bin/
COPY --from=build /opt/ripple/.ripple/ /opt/ripple/.ripple/

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["/opt/ripple/bin/rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]
