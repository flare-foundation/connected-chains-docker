FROM golang:1.24.6@sha256:2c89c41fb9efc3807029b59af69645867cfe978d2b877d475be0d72f6c6ce6f6

ARG BOR_DIR=/bor
ENV BOR_DIR=$BOR_DIR

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y \
    && mkdir -p /bor

WORKDIR ${BOR_DIR}
COPY . .
RUN make bor

RUN cp build/bin/bor /usr/local/bin/

ENV SHELL /bin/bash
EXPOSE 8545 8546 8547 30303 30303/udp

ENTRYPOINT ["bor"]