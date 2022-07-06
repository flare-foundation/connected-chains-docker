# syntax=docker/dockerfile:1.3-labs
FROM alpine:3

RUN apk add --no-cache bash netcat-openbsd jq curl ts expect

ENTRYPOINT ["bash", "/opt/entrypoint.sh"]
