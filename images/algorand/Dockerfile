# syntax=docker/dockerfile:1.3-labs
FROM debian:12@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as build

ARG ALGORAND_VERSION=v4.1.2-stable
ARG GO_VERSION=1.23.9

WORKDIR /opt/algorand
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/opt/algorand/go

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    set -e
    apt-get -y update && apt-get -y install \
        curl \
        git \
        make \
        gcc \
        g++ \
        pkg-config \
        libboost-dev \
        autoconf \
        automake \
        python3
    curl -L -o /tmp/go-linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    tar -C /usr/local -xf /tmp/go-linux-amd64.tar.gz
EOF
    
RUN <<-EOF
    set -e
    git clone --branch "$ALGORAND_VERSION" https://github.com/algorand/go-algorand.git
    cd go-algorand
    ./scripts/configure_dev.sh
    ./scripts/buildtools/install_buildtools.sh
    make build
    find /opt/algorand/go/bin -type f -executable -exec strip -s {} + 2>/dev/null || true
    ./scripts/local_install.sh -c stable -p ${GOPATH}/bin -d /opt/algorand/algorand-defaults -f -s
    rm /opt/algorand/algorand-defaults/genesis.json
    cp installer/genesis/mainnet/genesis.json /opt/algorand/algorand-defaults/genesis-mainnet.json
    cp installer/genesis/betanet/genesis.json /opt/algorand/algorand-defaults/genesis-betanet.json
    cp installer/genesis/devnet/genesis.json /opt/algorand/algorand-defaults/genesis-devnet.json
    cp installer/genesis/testnet/genesis.json /opt/algorand/algorand-defaults/genesis-testnet.json
EOF

RUN <<-EOF
    set -e
    mkdir -p /opt/algorand/.algorand
    touch /opt/algorand/.algorand/algod.token
    touch /opt/algorand/.algorand/algod.admin.token
EOF

COPY scripts/entrypoint.go /tmp/entrypoint.go
RUN go build -o /tmp/entrypoint /tmp/entrypoint.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:cdf4daaf154e3e27cfffc799c16f343a384228f38646928a1513d925f473cb46 as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

ENV GOPATH=/opt/algorand/go
ENV ALGORAND_DATA=/opt/algorand/.algorand

# Binaries for core functionality
COPY --from=build --chown=65532:65532 \
    /opt/algorand/go/bin/algod \
    /opt/algorand/go/bin/goal \
    /opt/algorand/go/bin/kmd \
    /opt/algorand/go/bin/
COPY --from=build --chown=65532:65532 /opt/algorand/algorand-defaults/ /opt/algorand/algorand-defaults/
COPY --from=build --chown=65532:65532 /opt/algorand/.algorand/ /opt/algorand/.algorand/

COPY --from=build --chown=65532:65532 /tmp/entrypoint /opt/algorand/scripts/entrypoint
COPY --chown=65532:65532 default-config/ /opt/algorand/default-config/

ENTRYPOINT ["/opt/algorand/scripts/entrypoint"]
CMD ["/opt/algorand/go/bin/algod"]
