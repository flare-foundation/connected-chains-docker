# syntax=docker/dockerfile:1.3-labs
FROM debian:13@sha256:0d01188e8dd0ac63bf155900fad49279131a876a1ea7fac917c62e87ccb2732d as final

ARG VERSION=v2.0.0

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