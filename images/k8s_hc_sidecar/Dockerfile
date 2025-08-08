# syntax=docker/dockerfile:1.3-labs
FROM alpine:3@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

RUN apk add --no-cache bash netcat-openbsd jq curl ts expect

ENTRYPOINT ["bash", "/opt/entrypoint.sh"]
