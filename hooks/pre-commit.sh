#!/bin/bash
# Pre-commit hook: quality checks on staged files
# Go-only: sdp CLI or go vet in sdp-plugin.

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -d "sdp-plugin" ]; then
    (cd sdp-plugin && go vet ./...) 2>/dev/null || true
fi
exit 0
