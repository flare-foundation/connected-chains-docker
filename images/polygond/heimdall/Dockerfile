FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed

ARG HEIMDALL_DIR=/heimdall
ENV HEIMDALL_DIR=$HEIMDALL_DIR

RUN apk add --no-cache \
       ca-certificates \
       tini && \
       mkdir -p ${HEIMDALL_DIR}

WORKDIR ${HEIMDALL_DIR}
COPY heimdalld /usr/local/bin/
COPY heimdallcli /usr/local/bin/
COPY builder/files/genesis-mainnet-v1.json ${HEIMDALL_DIR}/
COPY builder/files/genesis-testnet-v4.json ${HEIMDALL_DIR}/

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh

EXPOSE 1317 26656 26657

ENTRYPOINT ["entrypoint.sh"]