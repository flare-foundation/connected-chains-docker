# syntax=docker/dockerfile:1.3-labs
FROM ubuntu:20.04

ARG INDEXER_VERSION=2.10.0

ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/London"

WORKDIR /opt/algorand-indexer
ENV PATH=/opt/algorand-indexer/scripts:/opt/algorand-indexer:$PATH

ARG USER_ID=10201
ARG GROUP_ID=10201

RUN groupadd -g ${GROUP_ID} algoindexer && \
    useradd -u ${USER_ID} -g algoindexer -d /opt/algorand-indexer algoindexer

RUN <<-EOF
    apt-get -y update && \
    apt-get -y install jq curl gosu netcat-openbsd bzip2 && \
    curl -L -o /tmp/algorand-indexer_linux_amd64_${INDEXER_VERSION}.tar.bz2 https://github.com/algorand/indexer/releases/download/${INDEXER_VERSION}/algorand-indexer_linux_amd64_${INDEXER_VERSION}.tar.bz2 && \
    tar -xf /tmp/algorand-indexer_linux_amd64_${INDEXER_VERSION}.tar.bz2 -C /opt && \
    cp -r /opt/algorand-indexer_linux_amd64_${INDEXER_VERSION}/* /opt/algorand-indexer && \
    ls -la /opt/algorand-indexer && \
    rm -rf /opt/algorand-indexer_linux_amd64_${INDEXER_VERSION}
    rm /tmp/algorand-indexer_linux_amd64_${INDEXER_VERSION}.tar.bz2
EOF

COPY scripts/ /opt/algorand-indexer/scripts/

ENTRYPOINT ["entrypoint.sh"]
CMD ["run-indexer.sh"]
