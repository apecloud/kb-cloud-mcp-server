ARG VERSION="dev"

# Build stage
FROM golang:1.24-bullseye AS build
ARG VERSION
WORKDIR /build

# Enable Go build cache
RUN go env -w GOMODCACHE=/root/.cache/go-build

# Install dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code
COPY . ./

# Build the binary with version info
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -ldflags="-s -w \
    -X main.version=${VERSION} \
    -X main.commit=$(git rev-parse HEAD) \
    -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o kb-cloud-mcp-server ./cmd/server

# Final stage
FROM gcr.io/distroless/static-debian11
WORKDIR /app

# Copy the binary from build stage
COPY --from=build /build/kb-cloud-mcp-server .

# Use non-root user for security
USER nonroot:nonroot

# Set environment variables
ENV KB_CLOUD_MCP_LOG_LEVEL=info

# Command to run the server
ENTRYPOINT ["./kb-cloud-mcp-server"]
CMD ["stdio"]