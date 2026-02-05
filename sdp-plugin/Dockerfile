# Multi-stage build for SDP
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version={{ .Version }}" \
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
COPY --from=builder /src/LICENSE /LICENSE

ENTRYPOINT ["/usr/local/bin/sdp"]
CMD ["--version"]
