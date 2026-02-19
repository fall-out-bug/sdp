#!/bin/sh
# Pre-push hook: session validation + build/test before pushing
# Part of F065 - Agent Git Safety Protocol

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Get push target
REMOTE="$1"
URL="$2"
CURRENT_BRANCH=$(git branch --show-current)

# =========================================
# Session Validation (F065 Git Safety)
# =========================================

# Check session
if [ -f ".sdp/session.json" ]; then
    # Check if jq is available
    if command -v jq >/dev/null 2>&1; then
        EXPECTED_REMOTE=$(jq -r '.expected_remote' .sdp/session.json 2>/dev/null)
        EXPECTED_PUSH_TARGET="origin/$CURRENT_BRANCH"

        # Validate push target
        if [ -n "$EXPECTED_REMOTE" ] && [ "$EXPECTED_REMOTE" != "null" ] && [ "$EXPECTED_REMOTE" != "$EXPECTED_PUSH_TARGET" ]; then
            echo ""
            echo "ERROR: Push target mismatch!"
            echo "  Expected: $EXPECTED_REMOTE"
            echo "  Would push to: $EXPECTED_PUSH_TARGET"
            echo ""
            echo "Fix: git branch --set-upstream-to=$EXPECTED_REMOTE"
            echo ""
            exit 1
        fi
    fi
fi

# =========================================
# Protected Branch Check (F065 Feature Branch Enforcement)
# =========================================

# Prevent pushing to protected branches
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "dev" ]; then
    echo ""
    echo "ERROR: Direct push to $CURRENT_BRANCH is not allowed!"
    echo ""
    echo "Create a feature branch and use PR workflow:"
    echo "  git checkout -b feature/F###"
    echo "  git push origin feature/F###"
    echo "  # Then create a PR"
    echo ""
    exit 1
fi

# =========================================
# Build and Test
# =========================================

if [ -d "sdp-plugin" ]; then
    cd sdp-plugin
    go build ./...
    go test ./... -count=1 -short
fi

# Run staged guard checks before push
if command -v sdp >/dev/null 2>&1; then
    if sdp guard check --help 2>/dev/null | grep -q -- "--staged"; then
        sdp guard check --staged
    fi
fi
