#!/bin/sh
# SDP One-liner Installer
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
#   curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
#
# Default: project assets (prompts/hooks/config)
# --binary-only: install global CLI binary only

set -e

DEFAULT_REPO="fall-out-bug/sdp"
SDP_REPO="${SDP_REPO:-$DEFAULT_REPO}"
SDP_REF="${SDP_REF:-main}"
SDP_CLI_VERSION="${SDP_CLI_VERSION:-latest}"

SCRIPT_DIR=""
if [ -d "./scripts" ]; then
    SCRIPT_DIR="./scripts"
elif [ -d "$(dirname "$0")/scripts" ]; then
    SCRIPT_DIR="$(dirname "$0")/scripts"
fi

BINARY_ONLY=0
for arg in "$@"; do
    case "$arg" in
        --binary-only)
            BINARY_ONLY=1
            ;;
    esac
done

run_remote_script() {
    name="$1"
    shift
    url="https://raw.githubusercontent.com/${SDP_REPO}/${SDP_REF}/scripts/${name}"
    curl -fsSL "$url" | SDP_REPO="$SDP_REPO" SDP_REF="$SDP_REF" SDP_IDE="${SDP_IDE:-auto}" sh -s -- "$@"
}

if [ "$BINARY_ONLY" = "1" ]; then
    echo "📦 Installing SDP CLI binary..."
    if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/install.sh" ]; then
        exec sh "$SCRIPT_DIR/install.sh" "$SDP_CLI_VERSION"
    fi
    run_remote_script "install.sh" "$SDP_CLI_VERSION"
    exit $?
fi

echo "🔗 Installing SDP project assets (prompts/hooks/config)..."
if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/install-project.sh" ]; then
    exec sh "$SCRIPT_DIR/install-project.sh" "$@"
fi
run_remote_script "install-project.sh" "$@"
exit $?
