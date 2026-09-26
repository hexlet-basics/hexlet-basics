# syntax=docker/dockerfile:1

# The Go image: one binary per process (ADR-0010) plus the atlas CLI for the
# pre-upgrade migration job, so a release ships exactly one Go artifact. The
# chart picks the process with `command:` — /server, /worker or atlas.
#
# The builder's Go and the copied atlas must match mise.toml (go, and the
# atlas-community pin from ADR-0014); move them together.

FROM golang:1.27.1 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
# Static binaries: every dependency is pure Go (go-git for course clones, the
# Docker SDK for the runner), so the runtime image needs no libc. -trimpath
# keeps build paths out of the binary, which is what makes rebuilds of the
# same commit byte-identical.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/server && \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/worker ./cmd/worker

# Community (Apache-2.0) build, not the MSA-licensed default image (ADR-0014).
FROM arigaio/atlas:1.3.0-community AS atlas

# distroless/static ships CA certificates and /tmp and runs as uid 65532. The
# api-check pod reaches the host Docker socket through the node's docker group
# (securityContext.supplementalGroups in the chart), so nothing about Docker is
# baked in here.
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=atlas /atlas /usr/local/bin/atlas
COPY --from=build /out/server /server
COPY --from=build /out/worker /worker
COPY migrations /migrations
USER nonroot:nonroot
EXPOSE 3001
ENTRYPOINT []
CMD ["/server"]
