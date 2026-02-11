# Multi-stage build for SDP
FROM golang:1.21-alpine AS builder

# Build version can be overridden at build time
ARG VERSION=dev
ENV VERSION=${VERSION}

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -trimpath \
    -o /usr/local/bin/sdp \
    ./cmd/sdp

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git

# Copy binary from builder
COPY --from=builder /usr/local/bin/sdp /usr/local/bin/sdp

# Copy documentation
COPY --from=builder /src/README.md /README.md

ENTRYPOINT ["/usr/local/bin/sdp"]
CMD ["--version"]
