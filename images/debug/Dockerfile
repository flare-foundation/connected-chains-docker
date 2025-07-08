# syntax=docker/dockerfile:1.3-labs
FROM debian:12 as final

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