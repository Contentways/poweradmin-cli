# Copyright (c) 2026 Contentways
# SPDX-License-Identifier: MIT

# ---- Builder ----
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Download dependencies first — cached unless go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build statically linked binaries for all target platforms.
COPY . .
RUN mkdir -p linux/amd64 linux/arm64 && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-s -w" \
        -o linux/amd64/poweradmin . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
        -ldflags="-s -w" \
        -o linux/arm64/poweradmin .

# ---- Final image ----
# Use distroless/static for a minimal, secure runtime image.
# No shell, no package manager, no unnecessary binaries.
FROM gcr.io/distroless/static:nonroot

# OCI image annotations.
LABEL org.opencontainers.image.title="poweradmin-cli" \
      org.opencontainers.image.description="CLI for managing Poweradmin DNS" \
      org.opencontainers.image.source="https://github.com/Contentways/poweradmin-cli" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.vendor="Contentways"

# GoReleaser sets TARGETPLATFORM (e.g. linux/amd64) via docker buildx.
# The binary for the target platform is copied from the builder stage.
ARG TARGETPLATFORM
COPY --from=builder /build/${TARGETPLATFORM}/poweradmin /usr/local/bin/poweradmin

# Run as nonroot user (uid 65532) provided by distroless.
USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/poweradmin"]
