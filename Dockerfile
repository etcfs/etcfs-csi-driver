# EtcFS CSI driver
#
# Build:
#   docker build -t etcfs-csi --build-arg VERSION=$(git describe --tags --always) .

FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /usr/local/bin/etcfs-csi ./cmd/etcfs-csi

FROM alpine:3.20

# The node plugin issues bind mounts itself, so it needs no mount helper; only
# TLS roots for the controller's etcd client.
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/etcfs-csi /usr/local/bin/etcfs-csi

ENTRYPOINT ["/usr/local/bin/etcfs-csi"]
