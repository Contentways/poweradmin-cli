# Copyright (c) 2026 Contentways
# SPDX-License-Identifier: MIT

# GoReleaser builds the binary and places it in $TARGETPLATFORM/poweradmin
# within the build context. We only need to copy it into the final image.
FROM gcr.io/distroless/static:nonroot

# OCI image annotations.
LABEL org.opencontainers.image.title="poweradmin-cli" \
      org.opencontainers.image.description="CLI for managing Poweradmin DNS" \
      org.opencontainers.image.source="https://github.com/Contentways/poweradmin-cli" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.vendor="Contentways"

# GoReleaser sets TARGETPLATFORM (e.g. linux/amd64) via docker buildx.
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/poweradmin /usr/local/bin/poweradmin

# Run as nonroot user (uid 65532) provided by distroless.
USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/poweradmin"]
