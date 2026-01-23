#!/bin/bash
# sdp/hooks/post-build.sh
# Post-build checks for /build command
# Usage: ./post-build.sh WS-060-01 [module_path]

set -e

WS_ID=$1
MODULE=${2:-""}

if [ -z "$WS_ID" ]; then
    echo "❌ Usage: ./post-build.sh WS-ID [module_path]"
    exit 1
fi

echo "🔍 Post-build checks for $WS_ID"
echo "================================"

cd tools/hw_checker

# Check 1: Fast tests (regression)
echo ""
echo "Check 1: Regression tests"
if poetry run pytest tests/unit/ -m fast -q --tb=no 2>/dev/null; then
    FAST_COUNT=$(poetry run pytest tests/unit/ -m fast --collect-only -q 2>/dev/null | tail -1 | grep -oE "[0-9]+" | head -1 || echo "0")
    echo "✓ Regression tests passed ($FAST_COUNT tests)"
else
    echo "❌ Regression tests failed"
    echo "   Run: cd tools/hw_checker && poetry run pytest tests/unit/ -m fast -v"
    exit 1
fi

# Check 2: Linters
echo ""
echo "Check 2: Linters"

if [ -n "$MODULE" ]; then
    LINT_PATH="src/hw_checker/$MODULE"
else
    LINT_PATH="src/hw_checker/"
fi

# Ruff
if poetry run ruff check "$LINT_PATH" --quiet 2>/dev/null; then
    echo "✓ Ruff: no issues"
else
    echo "⚠️ Ruff found issues (run: ruff check $LINT_PATH)"
fi

# Mypy (optional, soft fail)
if poetry run mypy "$LINT_PATH" --ignore-missing-imports --no-error-summary 2>/dev/null; then
    echo "✓ Mypy: no issues"
else
    echo "⚠️ Mypy found issues (run: mypy $LINT_PATH)"
fi

# Check 3: TODO/FIXME
echo ""
echo "Check 3: No TODO/FIXME"
if [ -n "$MODULE" ]; then
    TODO_PATH="src/hw_checker/$MODULE"
else
    TODO_PATH="src/hw_checker/"
fi

TODO_COUNT=$(grep -rn "TODO\|FIXME\|HACK\|XXX" "$TODO_PATH" 2>/dev/null | wc -l || echo "0")
if [ "$TODO_COUNT" -eq 0 ]; then
    echo "✓ No TODO/FIXME markers"
else
    echo "❌ Found $TODO_COUNT TODO/FIXME markers"
    grep -rn "TODO\|FIXME\|HACK\|XXX" "$TODO_PATH" 2>/dev/null | head -5
    exit 1
fi

# Check 4: File sizes
echo ""
echo "Check 4: File sizes (< 200 LOC)"
LARGE_FILES=$(find "$TODO_PATH" -name "*.py" -exec wc -l {} \; 2>/dev/null | awk '$1 > 200 {print $2 " (" $1 " lines)"}')
if [ -z "$LARGE_FILES" ]; then
    echo "✓ All files < 200 LOC"
else
    echo "⚠️ Large files found:"
    echo "$LARGE_FILES"
fi

# Check 5: Import check
echo ""
echo "Check 5: Import check"
if [ -n "$MODULE" ]; then
    IMPORT_PATH="hw_checker.$MODULE"
    if python -c "import $IMPORT_PATH" 2>/dev/null; then
        echo "✓ Module imports successfully"
    else
        echo "⚠️ Module import failed (run: python -c 'import $IMPORT_PATH')"
    fi
else
    echo "  Skipped (no module specified)"
fi

# Check 6: Coverage (if tests exist for module)
echo ""
echo "Check 6: Coverage"
if [ -n "$MODULE" ]; then
    TEST_FILE="tests/unit/test_${MODULE}.py"
    if [ -f "$TEST_FILE" ]; then
        COV_RESULT=$(poetry run pytest "$TEST_FILE" --cov="hw_checker/$MODULE" --cov-report=term-missing --cov-fail-under=80 -q 2>&1)
        if echo "$COV_RESULT" | grep -q "PASSED\|passed"; then
            COV_PCT=$(echo "$COV_RESULT" | grep -oE "[0-9]+%" | head -1 || echo "N/A")
            echo "✓ Coverage: $COV_PCT (≥80%)"
        else
            echo "⚠️ Coverage check failed"
            echo "   Run: pytest $TEST_FILE --cov=hw_checker/$MODULE --cov-fail-under=80"
        fi
    else
        echo "  Skipped (no test file: $TEST_FILE)"
    fi
else
    echo "  Skipped (no module specified)"
fi

# Check 7: Verification (after Execution Report is appended)
echo ""
echo "Check 7: Verification (run after Execution Report is appended)"
WS_DIR="tools/hw_checker/docs/workstreams"
WS_FILE=$(find "$WS_DIR" -name "${WS_ID}-*.md" 2>/dev/null | head -1)

if [ -n "$WS_FILE" ] && grep -q "### Execution Report\|## Execution Report" "$WS_FILE"; then
    # Run verification hook if Execution Report exists
    HOOK_DIR="$(cd "$(dirname "$0")" && pwd)"
    if [ -f "$HOOK_DIR/verification-completion.sh" ]; then
        if "$HOOK_DIR/verification-completion.sh" "$WS_FILE" 2>/dev/null; then
            echo "✓ Verification check passed"
        else
            echo "⚠️ Verification check failed (run manually after appending Execution Report)"
            echo "   Run: $HOOK_DIR/verification-completion.sh $WS_FILE"
        fi
    else
        echo "  Skipped (verification hook not found)"
    fi
else
    echo "  Skipped (Execution Report not yet appended)"
    echo "  Run verification after appending Execution Report:"
    echo "    $HOOK_DIR/verification-completion.sh <ws-file>"
fi

# Check 8: Git commit with WS-ID
echo ""
echo "Check 8: Git commit"
if [ "$SKIP_COMMIT_CHECK" = "1" ]; then
    echo "⚠️ Skipped (SKIP_COMMIT_CHECK=1)"
else
    LAST_COMMIT=$(git log -1 --oneline 2>/dev/null || echo "")
    if echo "$LAST_COMMIT" | grep -q "$WS_ID"; then
        echo "✓ Commit found: $LAST_COMMIT"
    else
        echo "❌ No commit with WS-ID '$WS_ID' found"
        echo "   Last commit: $LAST_COMMIT"
        echo ""
        echo "   Run: git add . && git commit -m 'feat(...): $WS_ID - ...'"
        echo ""
        echo "   To skip this check: SKIP_COMMIT_CHECK=1 ./post-build.sh $WS_ID"
        exit 1
    fi
fi

# Check 9: Execution Report in WS file
echo ""
echo "Check 9: Execution Report"
if [ -z "$WS_FILE" ]; then
    WS_FILE=$(find tools/hw_checker/docs/workstreams -name "${WS_ID}*.md" 2>/dev/null | head -1)
fi
if [ -z "$WS_FILE" ]; then
    echo "⚠️ WS file not found for $WS_ID"
else
    if grep -q "### Execution Report\|## Execution Report" "$WS_FILE"; then
        echo "✓ Execution Report found in $WS_FILE"
    else
        echo "❌ Execution Report NOT found in $WS_FILE"
        echo ""
        echo "   Append Execution Report to WS file before completing."
        echo "   See template in sdp/prompts/commands/build.md Section 7"
        exit 1
    fi
fi

# Check 10: GitHub sync (optional, warning only)
echo ""
echo "Check 10: GitHub sync"
if [ -z "$GITHUB_TOKEN" ]; then
    echo "  Skipped (GITHUB_TOKEN not set)"
else
    ISSUE_NUM=$(grep -oP 'github_issue:\s*\K\d+' "$WS_FILE" 2>/dev/null || echo "")
    if [ -n "$ISSUE_NUM" ] && [ "$ISSUE_NUM" != "null" ]; then
        echo "✓ GitHub issue linked: #$ISSUE_NUM"
        echo "  Reminder: Update issue status if needed"
    else
        echo "⚠️ No GitHub issue linked (github_issue: null)"
        echo "  Consider syncing: sdp-github sync $WS_FILE"
    fi
fi

echo ""
echo "================================"
echo "✅ Post-build checks PASSED"
echo ""
echo "Next steps:"
echo "1. Append Execution Report to WS file"
echo "2. Run verification: sdp/hooks/verification-completion.sh <ws-file>"
echo "3. Run /build for next WS (if any)"
echo "4. After all WS: /codereview {feature}"
