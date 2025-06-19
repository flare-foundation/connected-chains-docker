# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:22.04 as build

ARG ALGORAND_VERSION=v4.1.2-stable
ARG GO_VERSION=1.23.9

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
    gosu \
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
    
COPY scripts/entrypoint.go /tmp/entrypoint.go
RUN go build -o /tmp/entrypoint /tmp/entrypoint.go
    
RUN cp /usr/bin/jq /tmp/jq && \
    cp /usr/bin/curl /tmp/curl && \
    cp /usr/sbin/gosu /tmp/gosu 

ARG USER_ID=10200
ARG GROUP_ID=10200

RUN groupadd -g ${GROUP_ID} algo && \
    useradd -u ${USER_ID} -g algo -d /opt/algorand algo && \ 
    cp /etc/passwd /tmp/passwd && \
    cp /etc/group /tmp/group

RUN mkdir -p /opt/algorand/.algorand && \
    chown -R algo:algo /opt/algorand && \
    touch /opt/algorand/.algorand/algod.token

FROM gcr.io/distroless/base-debian12:latest as final

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

WORKDIR /opt/algorand
ENV PATH=/opt/algorand/go/bin:/opt/algorand/scripts:/usr/local/go/bin:$PATH
ENV GOPATH=/opt/algorand/go
ENV ALGORAND_DATA=/opt/algorand/.algorand

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /tmp/group /etc/group

COPY --from=build /tmp/jq /usr/bin/jq
COPY --from=build /tmp/curl /usr/bin/curl
COPY --from=build /tmp/gosu /usr/bin/gosu

COPY --from=build /opt/algorand/go/ /opt/algorand/go/
COPY --from=build /opt/algorand/algorand-defaults/ /opt/algorand/algorand-defaults/
COPY --from=build /opt/algorand/.algorand/ /opt/algorand/.algorand/

COPY --from=build /tmp/entrypoint /opt/algorand/scripts/entrypoint
COPY default-config/ /opt/algorand/default-config/

ENTRYPOINT ["/opt/algorand/scripts/entrypoint"]
CMD ["algod"]
