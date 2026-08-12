# syntax=docker/dockerfile:1.3-labs
FROM gcc:15-trixie@sha256:486d5379b2d679f18bca957901d628952aa0ebc2d82817047bbef1dbfbc13d70 as build

ARG VERSION=3.3.0

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

ENV PATH="/opt/conan_env/bin:$PATH"

RUN <<-EOF
    set -e
    echo "conan" > /tmp/requirements.in
    python3 -m venv /opt/conan_env
    /opt/conan_env/bin/python -m pip install pip-tools
    /opt/conan_env/bin/pip-compile --generate-hashes --output-file=/tmp/requirements.txt /tmp/requirements.in
    /opt/conan_env/bin/python -m pip install -r /tmp/requirements.txt --require-hashes
EOF

RUN git clone --single-branch --branch "${VERSION}" https://github.com/ripple/rippled.git /opt/xrpld

WORKDIR /opt/xrpld

# https://github.com/XRPLF/rippled/blob/develop/BUILD.md
RUN <<-EOF
    set -e
    git config --global http.postBuffer 1048576000
    conan config install /opt/xrpld/conan/profiles/default -tf $(conan config home)/profiles
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
        -DCMAKE_INSTALL_PREFIX=/opt/xrpld/build \
        -Dstatic=ON \
        -Dxrpld=ON \
        -Dtests=OFF \
        -B./ ..
EOF

RUN <<-EOF
    set -e
    cd cmake_build
    cmake --build . --target install -- -j 8
    find /opt/xrpld/build/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
    mkdir -p /opt/xrpld/.xrpld/db
EOF

WORKDIR /opt/validator-keys-tool

RUN git clone --single-branch --branch "master" https://github.com/ripple/validator-keys-tool.git /opt/validator-keys-tool && \
    git checkout 05742108d4b7c0b7d9a97c073acd690d2c75a2b4

RUN <<-EOF
    set -e
    mkdir .build
    cd .build
    conan install \
       .. \
       --profile:host=default \
       --profile:build=default \
       --output-folder . \
       --build missing \
       --settings build_type=Release
    cmake -DCMAKE_POLICY_DEFAULT_CMP0091=NEW \
        -DCMAKE_TOOLCHAIN_FILE:FILEPATH=conan_toolchain.cmake \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_EXE_LINKER_FLAGS="-static" \
        ..
EOF

RUN <<-EOF
    set -e
    cd .build
    cmake --build .
EOF

COPY xrpld.cfg validators.txt /opt/xrpld/.xrpld/

FROM gcr.io/distroless/base-debian13:nonroot@sha256:ab7554b6d07ad354fad31957f8a1a813e65dfb93a8ad160568c79c3f2be6884f as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

COPY --from=build --chown=65532:65532 /opt/xrpld/build/bin/ /opt/xrpld/bin/
COPY --from=build --chown=65532:65532 /opt/xrpld/.xrpld/ /opt/xrpld/.xrpld/
COPY --from=build /opt/validator-keys-tool/.build/validator-keys /opt/xrpld/bin/

# https://xrpl.org/commandline-usage.html
ENTRYPOINT ["/opt/xrpld/bin/xrpld", "--conf=/opt/xrpld/.xrpld/xrpld.cfg", "--fg"]
