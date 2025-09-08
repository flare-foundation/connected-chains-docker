# syntax=docker/dockerfile:1.3-labs
FROM debian:13@sha256:6d87375016340817ac2391e670971725a9981cfc24e221c47734681ed0f6c0f5 as final

ARG VERSION=v1.0.0

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