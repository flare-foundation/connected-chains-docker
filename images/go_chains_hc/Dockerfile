FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /node-healthcheck

FROM gcr.io/distroless/cc-debian13:nonroot@sha256:4cf9e68a5cbd8c9623480b41d5ed6052f028c44cc29f91b21590613ab8bec824

COPY --from=builder /node-healthcheck /node-healthcheck

EXPOSE 8080
ENTRYPOINT ["/node-healthcheck"]