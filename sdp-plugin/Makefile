VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
GO_BUILD := CGO_ENABLED=0 GOFLAGS=-buildvcs=false go build

.PHONY: build build-all clean test install deps fmt lint version

# Build for current platform
build:
	@echo "Building sdp v$(VERSION)..."
	$(GO_BUILD) $(LDFLAGS) -o bin/sdp ./cmd/sdp
	@echo "✓ Built bin/sdp"
	@ls -lh bin/sdp

# Build for all platforms
build-all:
	@echo "Building sdp v$(VERSION) for all platforms..."
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) $(LDFLAGS) -o bin/sdp-darwin-arm64 ./cmd/sdp
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o bin/sdp-darwin-amd64 ./cmd/sdp
	GOOS=linux GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o bin/sdp-linux-amd64 ./cmd/sdp
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(LDFLAGS) -o bin/sdp-windows-amd64.exe ./cmd/sdp
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

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./... 2>&1 | grep -v "coverage: 100.0%" | grep "coverage:"
	@echo ""
	@echo "Checking packages below 80%..."
	@$(eval LOW_COV := $(shell go test -cover ./... 2>&1 | grep -E "coverage: (7[0-9]|[0-6][0-9])" | wc -l | tr -d ' '))
	@if [ "$(LOW_COV)" -gt "0" ]; then \
		echo "Packages below 80%:"; \
		go test -cover ./... 2>&1 | grep -E "coverage: (7[0-9]|[0-6][0-9])"; \
	fi

# Coverage for internal packages only (excludes cmd/ and internal/quality)
test-coverage-internal:
	@echo "Running coverage for internal packages..."
	@go test -cover ./internal/... 2>&1 | grep "coverage:"
	@echo ""
	@echo "Summary (excluding internal/quality):"
	@go test -cover $$(go list ./internal/... | grep -v "internal/quality$$") 2>&1 | grep -E "coverage: (7[0-9]|[0-6][0-9])" > /dev/null && echo "⚠️  Some packages below 80%" || echo "✅ All internal packages >= 80%"

# CI coverage check (excludes packages with external dependencies)
test-coverage-ci:
	@echo "CI Coverage Check"
	@echo "=================="
	@go test -cover $$(go list ./... | grep -v -E "cmd/sdp$$|internal/quality$$") 2>&1 | grep "coverage:"
	@echo ""
	@LOW=$$(go test -cover $$(go list ./... | grep -v -E "cmd/sdp$$|internal/quality$$") 2>&1 | grep -E "coverage: (7[0-9]|[0-6][0-9])" | wc -l | tr -d ' '); \
	if [ "$$LOW" -gt "0" ]; then \
		echo "❌ Coverage check FAILED - $$LOW packages below 80%"; \
		exit 1; \
	else \
		echo "✅ Coverage check PASSED - all packages >= 80%"; \
	fi
