#!/bin/bash
# SDP Install Script (WS-067-06: AC7)
# Usage: curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/scripts/install.sh | bash
# Or: ./install.sh [version]

set -euo pipefail

VERSION="${1:-latest}"
REPO="fall-out-bug/sdp"
BINARY_NAME="sdp"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo "ERROR: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Map OS names for archive
case "$OS" in
    darwin) ARCHIVE_OS="Darwin" ;;
    linux) ARCHIVE_OS="Linux" ;;
    *)
        echo "ERROR: Unsupported OS: $OS"
        exit 1
        ;;
esac

# Resolve version
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "ERROR: Could not determine latest version"
        exit 1
    fi
fi

echo "Installing SDP ${VERSION} for ${OS}/${ARCH}..."

# Construct archive name (matches goreleaser naming)
ARCHIVE_NAME="${BINARY_NAME}_${VERSION:1}_${ARCHIVE_OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download archive
echo "Downloading ${ARCHIVE_NAME}..."
if ! curl -sSLf "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE_NAME}"; then
    echo "ERROR: Failed to download ${DOWNLOAD_URL}"
    echo ""
    echo "Available releases: https://github.com/${REPO}/releases"
    exit 1
fi

# Download and verify checksum
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
echo "Verifying checksum..."
curl -sSLf "$CHECKSUM_URL" -o "${TMP_DIR}/checksums.txt" || echo "Warning: Could not download checksums"

if [ -f "${TMP_DIR}/checksums.txt" ]; then
    cd "$TMP_DIR"
    if grep -q "${ARCHIVE_NAME}" checksums.txt; then
        if command -v sha256sum &> /dev/null; then
            sha256sum -c --ignore-missing checksums.txt || echo "Warning: Checksum verification failed"
        elif command -v shasum &> /dev/null; then
            shasum -a 256 -c checksums.txt 2>/dev/null || echo "Warning: Checksum verification failed"
        fi
    fi
    cd - > /dev/null
fi

# Extract binary
echo "Extracting..."
tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

# Find the binary (might be in a subdirectory)
BINARY_PATH=$(find "${TMP_DIR}" -name "${BINARY_NAME}" -type f | head -1)
if [ -z "$BINARY_PATH" ]; then
    echo "ERROR: Binary not found in archive"
    exit 1
fi

# Install
mkdir -p "${INSTALL_DIR}"
chmod +x "${BINARY_PATH}"
mv "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"

echo ""
echo "✅ SDP ${VERSION} installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo ""

# Check if in PATH
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo "⚠️  ${INSTALL_DIR} is not in your PATH"
    echo ""
    echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo "    export PATH=\"\${HOME}/.local/bin:\${PATH}\""
    echo ""
    echo "Then reload: source ~/.bashrc  # or ~/.zshrc"
fi

# Verify installation
"${INSTALL_DIR}/${BINARY_NAME}" version 2>/dev/null || echo "Run '${BINARY_NAME} version' to verify installation"
