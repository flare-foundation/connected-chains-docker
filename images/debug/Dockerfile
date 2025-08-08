# syntax=docker/dockerfile:1.3-labs
FROM debian:12@sha256:b6507e340c43553136f5078284c8c68d86ec8262b1724dde73c325e8d3dcdeba as final

RUN apt-get update && apt-get install -y \
    curl \
    jq \
    procps \
    netcat-openbsd \
    tcpdump \
    strace \
    net-tools \
    iproute2 \
    vim \
    nano \
    less \
    tree

ENTRYPOINT ["/bin/bash"]