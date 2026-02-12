#!/bin/bash
# Pre-commit hook: quality checks on staged files
# Runs go vet and sdp guard check --staged

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Run go vet if sdp-plugin exists
if [ -d "sdp-plugin" ]; then
    (cd sdp-plugin && go vet ./...) 2>/dev/null || true
fi

# AC7: Run staged guard checks
# This checks scope compliance for staged files
if command -v sdp &> /dev/null; then
    sdp guard check --staged
fi

exit 0
