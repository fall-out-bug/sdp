#!/usr/bin/env bash
# Download and install SDP CLI binary
# Usage: download-sdp.sh <version> <output-dir>
# Environment: RUNNER_TEMP must be set

set -euo pipefail

VERSION="${1:-latest}"
OUTPUT_DIR="${2:-$RUNNER_TEMP/.bin}"

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

if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/fall-out-bug/sdp/releases/latest/download/${BINARY_NAME}"
else
  DOWNLOAD_URL="https://github.com/fall-out-bug/sdp/releases/download/${VERSION}/${BINARY_NAME}"
fi

echo "Downloading SDP CLI version: $VERSION"
echo "From: $DOWNLOAD_URL"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Download with retry logic (up to 3 attempts)
MAX_RETRIES=3
RETRY_COUNT=0
DOWNLOAD_SUCCESS=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ] && [ "$DOWNLOAD_SUCCESS" = false ]; do
  if [ $RETRY_COUNT -gt 0 ]; then
    echo "Retry $RETRY_COUNT/$MAX_RETRIES..."
    sleep 2  # Backoff before retry
  fi

  if curl -fsSL --retry 3 --retry-delay 2 "$DOWNLOAD_URL" -o "$OUTPUT_DIR/sdp" 2>/dev/null; then
    DOWNLOAD_SUCCESS=true
    echo "✅ Download successful"
  else
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
      echo "❌ Failed to download after $MAX_RETRIES attempts" >&2
      exit 1
    fi
  fi
done

chmod +x "$OUTPUT_DIR/sdp"

# Verify binary works
if ! "$OUTPUT_DIR/sdp" --help >/dev/null 2>&1; then
  echo "❌ Downloaded binary is not executable or corrupted" >&2
  exit 1
fi

echo "✅ SDP CLI installed to: $OUTPUT_DIR"

# Output path for sourcing
echo "$OUTPUT_DIR"
