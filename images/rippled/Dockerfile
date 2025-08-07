# syntax=docker/dockerfile:1.3-labs
FROM debian:12@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as build

ARG VERSION=2.5.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update &&
    apt-get -y install \
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
    python3 -m venv /opt/conan_env
    /opt/conan_env/bin/python -m pip install -r /tmp/requirements.txt --require-hashes
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/ripple/rippled.git /opt/ripple

WORKDIR /opt/ripple

ENV PATH="/opt/conan_env/bin:$PATH"

# New and improved build process, but is not in latest (2.5.0) release yet
# Move from copying the default profile to using one in 
# "/opt/ripple/conan/profiles/default" after it gets integrated
# https://github.com/XRPLF/rippled/blob/develop/BUILD.md
COPY conan_default conan_default

RUN <<-EOF
    git config --global http.postBuffer 1048576000
    conan profile detect
    conan config install conan_default -tf $(conan config home)/profiles 
    conan remote add --index 0 xrplf "https://conan.ripplex.io"
    conan config set tools.build:jobs=8
    mkdir cmake_build
EOF

RUN <<-EOF
    cd cmake_build
    conan install \
        .. \
        --output-folder . \
        --build missing \
        --settings build_type=Release
EOF

RUN <<-EOF
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
    cd cmake_build
    cmake --build . --target install -- -j 8
    find /opt/ripple/build/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
    mkdir -p /opt/ripple/.ripple/db
EOF

COPY rippled.conf validators.txt /opt/ripple/.ripple/

RUN chown -R 65532:65532 /opt/ripple/

FROM gcr.io/distroless/base-debian12:nonroot@sha256:dca858878977f01b80685240ca2cd9c1bb38958e2ab6f4ed5698c0ea23307143 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build /opt/ripple/build/bin/ /opt/ripple/bin/
COPY --from=build /opt/ripple/.ripple/ /opt/ripple/.ripple/

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["/opt/ripple/bin/rippled", "--conf=/opt/ripple/.ripple/rippled.conf", "--fg"]
