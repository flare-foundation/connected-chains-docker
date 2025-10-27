# syntax=docker/dockerfile:1.3-labs
FROM debian:bookworm@sha256:26f2a7cab45014541c65f9d140ccfa6aaefbb49686c6759bea9c6f7f5bb3d72f as build

ARG VERSION=2.6.1

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    set -e
    apt-get -y update && apt-get -y install \
        python3-venv \
        curl \
        git \
        build-essential \
        cmake
EOF

RUN <<-EOF
    set -e
    echo "conan" > /tmp/requirements.in
    python3 -m venv /opt/conan_env
    /opt/conan_env/bin/python -m pip install pip-tools
    /opt/conan_env/bin/pip-compile --generate-hashes --output-file=/tmp/requirements.txt /tmp/requirements.in
    /opt/conan_env/bin/python -m pip install -r /tmp/requirements.txt --require-hashes
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/ripple/rippled.git /opt/ripple

WORKDIR /opt/ripple

ENV PATH="/opt/conan_env/bin:$PATH"

# https://github.com/XRPLF/rippled/blob/develop/BUILD.md
RUN <<-EOF
    set -e
    git config --global http.postBuffer 1048576000
    conan config install /opt/ripple/conan/profiles/default -tf $(conan config home)/profiles 
    conan remote add --index 0 xrplf "https://conan.ripplex.io"
    mkdir cmake_build
EOF

RUN <<-EOF
    set -e
    cd cmake_build
    conan install \
        .. \
        --profile:host=default \
        --profile:build=default \
        --output-folder . \
        --build missing \
        --settings build_type=Release
EOF

RUN <<-EOF
    set -e
    cd cmake_build
    cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_TOOLCHAIN_FILE:FILEPATH=build/generators/conan_toolchain.cmake \
        -DCMAKE_INSTALL_PREFIX=/opt/ripple/build \
        -Dstatic=ON \
        -Dxrpld=ON \
        -Dtests=OFF \
        -B./ ..
EOF

RUN <<-EOF
    set -e
    cd cmake_build
    cmake --build . --target install -- -j 8
    find /opt/ripple/build/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
    mkdir -p /opt/ripple/.ripple/db
EOF

COPY rippled.conf validators.txt /opt/ripple/.ripple/

FROM gcr.io/distroless/base-debian12:nonroot@sha256:c1201b805d3a35a4e870f9ce9775982dd166a2b0772232638dd2440fbe0e0134 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build --chown=65532:65532 /opt/ripple/build/bin/ /opt/ripple/bin/
COPY --from=build --chown=65532:65532 /opt/ripple/.ripple/ /opt/ripple/.ripple/

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["/opt/ripple/bin/rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]
