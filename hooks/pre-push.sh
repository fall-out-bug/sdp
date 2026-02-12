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

# AC7: Run staged guard checks before push
# Ensures all commits being pushed comply with guard policies
if command -v sdp &> /dev/null; then
    sdp guard check --staged
fi
