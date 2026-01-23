#!/bin/bash
# WS Sync Hook for Claude Code
# Syncs workstream files to GitHub after Write/Edit
#
# Usage: ws-sync-hook.sh <file_path>
# Environment: GITHUB_TOKEN, GITHUB_REPO

set -e

FILE_PATH="${1:-}"

# Skip if not a WS file
if [[ ! "$FILE_PATH" =~ workstreams/.*/WS-.*\.md$ ]]; then
    exit 0
fi

# Skip if no GitHub token (local development without sync)
if [[ -z "${GITHUB_TOKEN:-}" ]]; then
    echo "Skipping GitHub sync: GITHUB_TOKEN not set"
    exit 0
fi

# Skip if no GitHub repo
if [[ -z "${GITHUB_REPO:-}" ]]; then
    echo "Skipping GitHub sync: GITHUB_REPO not set"
    exit 0
fi

echo "Syncing to GitHub: $FILE_PATH"

# Find sdp directory (hook is in sdp/hooks/validators/)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SDP_DIR="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# Convert relative path to absolute if needed
if [[ ! "$FILE_PATH" =~ ^/ ]]; then
    # Relative path - make absolute from workspace root
    WORKSPACE_ROOT="$(cd "$SDP_DIR/../.." && pwd)"
    FILE_PATH="$WORKSPACE_ROOT/$FILE_PATH"
fi

# Run sync (non-blocking - errors logged but don't fail)
cd "$SDP_DIR"
if poetry run sdp-github sync-ws "$FILE_PATH" 2>&1; then
    echo "GitHub sync complete"
else
    echo "Warning: GitHub sync failed (non-blocking)"
fi

exit 0
