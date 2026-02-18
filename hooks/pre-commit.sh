#!/bin/bash
# Pre-commit hook: session validation + quality checks on staged files
# Part of F065 - Agent Git Safety Protocol

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Get current context
CURRENT_BRANCH=$(git branch --show-current)
CURRENT_DIR=$(pwd)

# =========================================
# Session Validation (F065 Git Safety)
# =========================================

# Check for session file
if [ -f ".sdp/session.json" ]; then
    # Check if jq is available
    if command -v jq &> /dev/null; then
        EXPECTED_BRANCH=$(jq -r '.expected_branch' .sdp/session.json 2>/dev/null)
        EXPECTED_DIR=$(jq -r '.worktree_path' .sdp/session.json 2>/dev/null)

        # Validate branch
        if [ -n "$EXPECTED_BRANCH" ] && [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
            echo ""
            echo "ERROR: Branch mismatch!"
            echo "  Expected: $EXPECTED_BRANCH"
            echo "  Current:  $CURRENT_BRANCH"
            echo ""
            echo "Run: sdp session sync"
            echo "  or: git checkout $EXPECTED_BRANCH"
            echo ""
            exit 1
        fi

        # Validate directory
        if [ -n "$EXPECTED_DIR" ] && [ "$CURRENT_DIR" != "$EXPECTED_DIR" ]; then
            echo ""
            echo "ERROR: Directory mismatch!"
            echo "  Expected: $EXPECTED_DIR"
            echo "  Current:  $CURRENT_DIR"
            echo ""
            echo "Run: sdp guard context go <feature-id>"
            echo ""
            exit 1
        fi

        # Check for cross-feature commits (feature branch tracking origin/dev)
        PARENT_BRANCH=$(git rev-parse --abbrev-ref HEAD@{upstream} 2>/dev/null || echo "")
        FEATURE_ID=$(jq -r '.feature_id' .sdp/session.json 2>/dev/null)
        if [[ -n "$FEATURE_ID" && "$PARENT_BRANCH" == "origin/dev" && "$CURRENT_BRANCH" != "dev" ]]; then
            echo ""
            echo "WARNING: Feature branch tracking origin/dev instead of origin/$CURRENT_BRANCH"
            echo "  Feature: $FEATURE_ID"
            echo "  Current tracking: $PARENT_BRANCH"
            echo ""
            echo "Fix: git branch --set-upstream-to=origin/$CURRENT_BRANCH"
            echo ""
        fi
    fi
fi

# =========================================
# Protected Branch Check (F065 Feature Branch Enforcement)
# =========================================

# Block commits to main/dev for feature work
if [ -f ".sdp/session.json" ]; then
    if [[ "$CURRENT_BRANCH" == "main" || "$CURRENT_BRANCH" == "dev" ]]; then
        FEATURE_ID=$(jq -r '.feature_id' .sdp/session.json 2>/dev/null)
        if [ -n "$FEATURE_ID" ] && [ "$FEATURE_ID" != "null" ]; then
            echo ""
            echo "ERROR: Cannot commit to $CURRENT_BRANCH for feature $FEATURE_ID"
            echo ""
            echo "Create a feature branch:"
            echo "  git checkout -b feature/$FEATURE_ID"
            echo ""
            exit 1
        fi
    fi
fi

# =========================================
# Quality Checks
# =========================================

# Run go vet if sdp-plugin exists
if [ -d "sdp-plugin" ]; then
    (cd sdp-plugin && go vet ./...) 2>/dev/null || true
fi

# Run guard checks for staged files
if command -v sdp &> /dev/null; then
    STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM)
    if [ -n "$STAGED_FILES" ]; then
        while IFS= read -r file; do
            [ -n "$file" ] || continue
            if ! sdp guard check "$file" >/dev/null 2>&1; then
                echo ""
                echo "ERROR: Guard check failed for staged file: $file"
                sdp guard check "$file" || true
                echo ""
                exit 1
            fi
        done <<EOF
$STAGED_FILES
EOF
    fi
fi

exit 0
