#!/bin/bash
# sdp/hooks/validators/session-quality-check.sh
# Run at end of agent turn to check overall quality

# CD to project directory
cd "$(dirname "${BASH_SOURCE[0]}")/../.." || exit 0

echo "Running session quality checks..."

# Quick regression check (only if tests exist and poetry is available)
if [ -d "tests/unit" ] && command -v poetry &> /dev/null; then
    echo ""
    echo "=== Quick Regression Check ==="
    poetry run pytest tests/unit/ -m fast -q --tb=no 2>/dev/null && {
        echo "Fast tests: PASSED"
    } || {
        echo "WARNING: Some fast tests may be failing"
    }
fi

# Check for uncommitted TODO/FIXME in staged files
if git rev-parse --git-dir > /dev/null 2>&1; then
    echo ""
    echo "=== Staged Files Check ==="
    STAGED_PY=$(git diff --cached --name-only --diff-filter=ACM | grep "\.py$" || true)
    if [ -n "$STAGED_PY" ]; then
        TODO_IN_STAGED=$(echo "$STAGED_PY" | xargs grep -l "TODO\|FIXME" 2>/dev/null || true)
        if [ -n "$TODO_IN_STAGED" ]; then
            echo "WARNING: TODO/FIXME found in staged files:"
            echo "$TODO_IN_STAGED"
        else
            echo "No TODO/FIXME in staged Python files"
        fi
    fi
fi

echo ""
echo "Session quality check complete"
exit 0
