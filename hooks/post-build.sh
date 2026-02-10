#!/bin/bash
# Post-build hook: quality checks after workstream execution
# Usage: ./post-build.sh WS-ID [module_path]
# Go-only: build and test sdp-plugin.

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -d "sdp-plugin" ]; then
    cd sdp-plugin
    go build ./...
    go test ./... -count=1 -short
fi
