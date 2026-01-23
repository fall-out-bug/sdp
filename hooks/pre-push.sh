#!/bin/bash
# sdp/hooks/pre-push.sh
# Run regression tests before pushing

set -e

echo "🔍 Running pre-push checks..."
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

# Run regression tests
echo "1. Running regression tests..."
cd tools/hw_checker

if poetry run pytest tests/unit/ -m fast -q --tb=no 2>&1; then
    echo "✓ Regression tests passed"
else
    echo "⚠️  Regression tests failed"
    echo "   Run: cd tools/hw_checker && poetry run pytest tests/unit/ -m fast -v"
    # Don't block push, just warn
fi

# Check coverage if coverage report exists
if [ -f ".coverage" ]; then
    echo ""
    echo "2. Checking coverage..."
    COVERAGE=$(poetry run python -c "import coverage; cov = coverage.Coverage(); cov.load(); cov.report(file='/dev/stdout', show_missing=False)" 2>/dev/null | grep -oP '\d+%' | head -1 || echo "0%")

    COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')

    if [ "$COVERAGE_NUM" -lt 80 ]; then
        echo "⚠️  Coverage is below 80% (currently: ${COVERAGE})"
        echo "   Consider adding more tests"
    else
        echo "✓ Coverage is ${COVERAGE} (≥ 80%)"
    fi
fi

cd - > /dev/null

echo ""
echo "✅ Pre-push checks complete"
