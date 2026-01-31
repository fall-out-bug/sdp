#!/bin/bash
# sdp/hooks/pre-push.sh
# Run regression tests before pushing to remote

set -e

echo "🔍 Running pre-push checks..."
echo ""

# Check if strict mode is enabled (hard blocking)
# Set SDP_HARD_PUSH=1 to enable hard blocking
# Default: warn-only for backward compatibility
HARD_PUSH=${SDP_HARD_PUSH:-0}

if [ "$HARD_PUSH" = "1" ]; then
    echo "🔒 HARD PUSH MODE enabled (SDP_HARD_PUSH=1)"
    echo "   Test failures will BLOCK push"
else
    echo "⚠️  Soft mode (set SDP_HARD_PUSH=1 to enable hard blocking)"
    echo "   Test failures will WARN but not block push"
fi

echo ""

# Change to project root
cd "$(git rev-parse --show-toplevel)"

# Get list of files to be pushed
FILES_TO_PUSH=$(git diff --name-only HEAD @{u} 2>/dev/null || echo "")

if [ -z "$FILES_TO_PUSH" ]; then
    # No upstream branch, compare with HEAD~1
    FILES_TO_PUSH=$(git diff --name-only HEAD~1 HEAD 2>/dev/null || echo "")
fi

# Check if any Python files are being pushed
PY_FILES=$(echo "$FILES_TO_PUSH" | grep "\.py$" || true)

if [ -z "$PY_FILES" ]; then
    echo "No Python files to push, skipping tests"
    echo ""
    echo "✅ Pre-push checks complete"
    exit 0
fi

HAS_FAILURES=0

# Run regression tests
echo "1. Running regression tests..."
if [ -d "tests" ] && command -v pytest &> /dev/null; then
    if pytest tests/ -x -q --tb=no 2>&1; then
        echo "✓ Regression tests passed"
    else
        echo "❌ Regression tests failed"
        echo ""
        echo "To fix this issue:"
        echo "  1. Run: pytest tests/ -v"
        echo "  2. Fix failing tests"
        echo "  3. Commit the fixes"
        echo "  4. Push again"
        echo ""
        if [ "$HARD_PUSH" = "1" ]; then
            echo "🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)"
            echo "To bypass: git push --no-verify"
            exit 1
        else
            echo "⚠️  WARNING: Push not blocked (SDP_HARD_PUSH not set)"
            echo "   Set SDP_HARD_PUSH=1 to enforce blocking in future"
            HAS_FAILURES=1
        fi
    fi
else
    echo "⚠️  pytest or tests/ directory not found, skipping tests"
fi

# Check coverage if coverage report exists
if [ -f ".coverage" ] && command -v coverage &> /dev/null; then
    echo ""
    echo "2. Checking coverage..."
    COVERAGE=$(python -c "import coverage; cov = coverage.Coverage(); cov.load(); cov.report(file='/dev/stdout', show_missing=False)" 2>/dev/null | grep -oP '\d+%' | head -1 || echo "0%")

    COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')

    if [ "$COVERAGE_NUM" -lt 80 ]; then
        echo "❌ Coverage is below 80% (currently: ${COVERAGE})"
        echo ""
        echo "To fix this issue:"
        echo "  1. Run: pytest --cov=src/sdp --cov-report=term-missing"
        echo "  2. Add tests for uncovered lines"
        echo "  3. Commit the tests"
        echo "  4. Push again"
        echo ""
        if [ "$HARD_PUSH" = "1" ]; then
            echo "🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)"
            echo "To bypass: git push --no-verify"
            exit 1
        else
            echo "⚠️  WARNING: Push not blocked (SDP_HARD_PUSH not set)"
            echo "   Set SDP_HARD_PUSH=1 to enforce blocking in future"
            HAS_FAILURES=1
        fi
    else
        echo "✓ Coverage is ${COVERAGE} (≥ 80%)"
    fi
fi

echo ""
if [ "$HARD_PUSH" = "1" ] && [ "$HAS_FAILURES" -eq 0 ]; then
    echo "✅ Pre-push checks passed (hard blocking mode)"
elif [ "$HARD_PUSH" = "0" ]; then
    if [ "$HAS_FAILURES" -eq 1 ]; then
        echo "⚠️  Pre-push checks complete (WARNING mode - failures detected but push allowed)"
        echo "   To enable hard blocking: export SDP_HARD_PUSH=1"
    else
        echo "✅ Pre-push checks complete (WARNING mode - all checks passed)"
        echo "   To enable hard blocking: export SDP_HARD_PUSH=1"
    fi
fi
