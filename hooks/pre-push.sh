#!/bin/bash
# Pre-push hook: build and test before pushing

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -d "sdp-plugin" ]; then
    cd sdp-plugin
    go build ./...
    go test ./... -count=1 -short
fi
