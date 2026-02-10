#!/usr/bin/env bash
# Download and install SDP CLI binary
# Usage: download-sdp.sh <version> <output-dir> [cache-dir]
# Environment: RUNNER_TEMP must be set
# Output: Only the binary directory path (all other output to stderr)

set -euo pipefail

VERSION="${1:-latest}"
OUTPUT_DIR="${2:-$RUNNER_TEMP/.bin}"
CACHE_DIR="${3:-$RUNNER_TEMP/sdp-cache}"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture names
case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "❌ Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

# Download from GitHub releases
BINARY_NAME="sdp-${OS}-${ARCH}"
CACHE_KEY="${BINARY_NAME}-${VERSION}"

# Check cache first
CACHED_BINARY="$CACHE_DIR/$CACHE_KEY"
if [ -f "$CACHED_BINARY" ]; then
  echo "♻️  Using cached binary from: $CACHED_BINARY" >&2
  mkdir -p "$OUTPUT_DIR"
  cp "$CACHED_BINARY" "$OUTPUT_DIR/sdp"
  chmod +x "$OUTPUT_DIR/sdp"

  # Verify cached binary works
  if "$OUTPUT_DIR/sdp" --help >/dev/null 2>&1; then
    echo "✅ SDP CLI loaded from cache: $OUTPUT_DIR" >&2
    echo "$OUTPUT_DIR"
    exit 0
  else
    echo "⚠️  Cached binary is corrupted, re-downloading..." >&2
    rm -f "$CACHED_BINARY"
  fi
fi

if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/fall-out-bug/sdp/releases/latest/download/${BINARY_NAME}"
else
  DOWNLOAD_URL="https://github.com/fall-out-bug/sdp/releases/download/${VERSION}/${BINARY_NAME}"
fi

echo "Downloading SDP CLI version: $VERSION" >&2
echo "From: $DOWNLOAD_URL" >&2

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Download with retry logic (up to 3 attempts)
MAX_RETRIES=3
RETRY_COUNT=0
DOWNLOAD_SUCCESS=false
USED_FALLBACK=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ] && [ "$DOWNLOAD_SUCCESS" = false ]; do
  if [ $RETRY_COUNT -gt 0 ]; then
    echo "Retry $RETRY_COUNT/$MAX_RETRIES..."
    sleep 2  # Backoff before retry
  fi

  if curl -fsSL --retry 3 --retry-delay 2 "$DOWNLOAD_URL" -o "$OUTPUT_DIR/sdp" 2>/dev/null; then
    DOWNLOAD_SUCCESS=true
    echo "✅ Download successful" >&2
  else
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
      echo "❌ Failed to download after $MAX_RETRIES attempts" >&2

      # Fallback for CI/testing: try to use local sdp if available
      if command -v sdp >/dev/null 2>&1; then
        echo "⚠️  Using local sdp from PATH as fallback" >&2
        LOCAL_SDP_PATH=$(command -v sdp)
        mkdir -p "$OUTPUT_DIR"
        cp "$LOCAL_SDP_PATH" "$OUTPUT_DIR/sdp"
        chmod +x "$OUTPUT_DIR/sdp"
        DOWNLOAD_SUCCESS=true
        USED_FALLBACK=true
      else
        exit 1
      fi
    fi
  fi
done

# Download and verify SHA256 checksum (skip for local fallback)
if [ "$USED_FALLBACK" = true ]; then
  echo "⚠️  Using local fallback, skipping checksum verification" >&2
else
  CHECKSUM_URL="${DOWNLOAD_URL}.sha256"
  echo "Downloading checksum: $CHECKSUM_URL"

  if curl -fsSL --retry 2 "$CHECKSUM_URL" -o "$OUTPUT_DIR/sdp.sha256" 2>/dev/null; then
    echo "✅ Checksum downloaded" >&2

    # Verify checksum
    echo "Verifying SHA256 checksum..." >&2
    if cd "$OUTPUT_DIR" && sha256sum -c sdp.sha256 2>/dev/null; then
      echo "✅ SHA256 checksum verified" >&2
    else
      echo "❌ SHA256 checksum verification failed" >&2
      echo "   The downloaded binary may be corrupted or tampered with" >&2
      rm -f "$OUTPUT_DIR/sdp" "$OUTPUT_DIR/sdp.sha256"
      exit 1
    fi
  else
    echo "⚠️  Warning: Checksum not available, skipping verification" >&2
    echo "   This is less secure but continuing anyway" >&2
  fi
fi

chmod +x "$OUTPUT_DIR/sdp"

# Verify binary works
if ! "$OUTPUT_DIR/sdp" --help >/dev/null 2>&1; then
  echo "❌ Downloaded binary is not executable or corrupted" >&2
  exit 1
fi

# Cache the binary for future runs (skip for local fallback)
if [ "$USED_FALLBACK" = false ]; then
  mkdir -p "$CACHE_DIR"
  cp "$OUTPUT_DIR/sdp" "$CACHED_BINARY"
  echo "💾 Cached binary to: $CACHED_BINARY" >&2
fi

echo "✅ SDP CLI installed to: $OUTPUT_DIR" >&2

# Output path for sourcing (this is the only stdout output)
echo "$OUTPUT_DIR"
