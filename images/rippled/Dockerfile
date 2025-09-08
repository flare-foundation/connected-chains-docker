# syntax=docker/dockerfile:1.3-labs
FROM debian:13@sha256:6d87375016340817ac2391e670971725a9981cfc24e221c47734681ed0f6c0f5 as build

ARG VERSION=2.6.0

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

# To generate requirements.txt, run:
# echo "conan" > requirements.in
# pip-compile --generate-hashes --output-file=requirements.txt requirements.in
COPY requirements.txt /tmp/requirements.txt

RUN <<-EOF
    set -e
    python3 -m venv /opt/conan_env
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
