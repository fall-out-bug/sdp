VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build build-all clean test install deps fmt lint version

# Build for current platform
build:
	@echo "Building sdp v$(VERSION)..."
	go build $(LDFLAGS) -o bin/sdp ./cmd/sdp
	@echo "✓ Built bin/sdp"
	@ls -lh bin/sdp

# Build for all platforms
build-all:
	@echo "Building sdp v$(VERSION) for all platforms..."
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/sdp-darwin-arm64 ./cmd/sdp
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/sdp-darwin-amd64 ./cmd/sdp
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/sdp-linux-amd64 ./cmd/sdp
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/sdp-windows-amd64.exe ./cmd/sdp
	@echo "✓ Built 4 binaries"
	@ls -lh bin/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	@echo "✓ Cleaned"

# Install to GOPATH/bin
install: build
	@echo "Installing to $(GOPATH)/bin/sdp..."
	install -m 0755 bin/sdp $(GOPATH)/bin/sdp
	@echo "✓ Installed"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies ready"

# Format code
fmt:
	go fmt ./...

# Run linter (golangci-lint)
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

# Show version
version:
	@echo "$(VERSION)"
