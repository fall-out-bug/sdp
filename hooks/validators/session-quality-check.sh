#!/bin/bash
# sdp/hooks/validators/session-quality-check.sh
# Run at end of agent turn to check overall quality

# CD to project directory
cd "$(dirname "${BASH_SOURCE[0]}")/../.." || exit 0

echo "Running session quality checks..."

# Quick Go build and test if sdp-plugin exists
if [ -d "sdp-plugin" ]; then
    echo ""
    echo "=== Quick Build & Test ==="
    (cd sdp-plugin && go build ./... && go test ./... -short -count=1) 2>/dev/null && {
        echo "Go build and tests: PASSED"
    } || {
        echo "WARNING: Go build or tests may be failing"
    }
fi

# Check for TODO/FIXME in staged files
if git rev-parse --git-dir > /dev/null 2>&1; then
    echo ""
    echo "=== Staged Files Check ==="
    STAGED=$(git diff --cached --name-only --diff-filter=ACM | grep -E "\.(go|py|ts|tsx|js)$" || true)
    if [ -n "$STAGED" ]; then
        TODO_IN_STAGED=$(echo "$STAGED" | xargs grep -l "TODO\|FIXME" 2>/dev/null || true)
        if [ -n "$TODO_IN_STAGED" ]; then
            echo "WARNING: TODO/FIXME found in staged files:"
            echo "$TODO_IN_STAGED"
        else
            echo "No TODO/FIXME in staged files"
        fi
    fi
fi

echo ""
echo "Session quality check complete"
exit 0
