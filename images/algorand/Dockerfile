# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG ALGORAND_VERSION=v3.25.0-stable
ARG GO_VERSION=1.21.10

WORKDIR /opt/algorand
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/opt/algorand/go

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install \
    jq \
    curl \
    git \
    make \
    gcc \
    g++ \
    pkg-config \
    libboost-dev \
    autoconf \
    automake \
    python3 && \
    curl -L -o /tmp/go-linux-amd64.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xf /tmp/go-linux-amd64.tar.gz
EOF

RUN <<-EOF
    git clone --branch "$ALGORAND_VERSION" https://github.com/algorand/go-algorand.git && \
    cd go-algorand && \
    ./scripts/configure_dev.sh && \
    ./scripts/buildtools/install_buildtools.sh && \
    make build && \
    ./scripts/local_install.sh -c stable -p ${GOPATH}/bin -d /opt/algorand/algorand-defaults -f -s && \
    rm /opt/algorand/algorand-defaults/genesis.json && \
    cp installer/genesis/mainnet/genesis.json /opt/algorand/algorand-defaults/genesis-mainnet.json && \
    cp installer/genesis/betanet/genesis.json /opt/algorand/algorand-defaults/genesis-betanet.json && \
    cp installer/genesis/devnet/genesis.json /opt/algorand/algorand-defaults/genesis-devnet.json && \
    cp installer/genesis/testnet/genesis.json /opt/algorand/algorand-defaults/genesis-testnet.json
EOF

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

WORKDIR /opt/algorand
ENV PATH=/opt/algorand/go/bin:/opt/algorand/scripts:/usr/local/go/bin:$PATH
ENV GOPATH=/opt/algorand/go

ARG USER_ID=10200
ARG GROUP_ID=10200

RUN groupadd -g ${GROUP_ID} algo && useradd -u ${USER_ID} -g algo -d /opt/algorand algo

COPY --from=build /tmp/go-linux-amd64.tar.gz /tmp
COPY --from=build /opt/algorand/go/ /opt/algorand/go/
COPY --from=build /opt/algorand/algorand-defaults/ /opt/algorand/algorand-defaults/

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install jq curl gosu && \
    cd /tmp && \
    tar -C /usr/local -xf go-linux-amd64.tar.gz && \
    rm go-linux-amd64.tar.gz && \
    apt-get clean
EOF

COPY scripts/ /opt/algorand/scripts/
COPY default-config/ /opt/algorand/default-config/

RUN <<-EOF
    mkdir -p /opt/algorand/.algorand && \
    chown -R algo:algo /opt/algorand && \
    touch /opt/algorand/.algorand/algod.token 
EOF

ENV ALGORAND_DATA="/opt/algorand/.algorand"

ENTRYPOINT ["entrypoint.sh"]
CMD ["algod"]
