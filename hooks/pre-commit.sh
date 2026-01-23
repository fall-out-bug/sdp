#!/bin/bash
# sdp/hooks/pre-commit.sh
# Git pre-commit hook for quality checks
# Install: ln -sf ../../sdp/hooks/pre-commit.sh .git/hooks/pre-commit

set -e

echo "🔍 Pre-commit checks"
echo "===================="

# Get list of staged files
STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM)

if [ -z "$STAGED_FILES" ]; then
    echo "No staged files, skipping checks"
    exit 0
fi

# Check 0: Not committing to main directly
echo ""
echo "Check 0: Branch check"
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
    echo "⚠️ Warning: Committing directly to $CURRENT_BRANCH"
    echo "  Consider using a feature branch: git checkout -b feature/{slug}"
fi
echo "✓ Branch: $CURRENT_BRANCH"

# Check 1: No time estimates in WS files
echo ""
echo "Check 1: No time estimates in WS files"
WS_FILES=$(echo "$STAGED_FILES" | grep "workstreams/.*\.md$" || true)

    if [ -n "$WS_FILES" ]; then
        TIME_ESTIMATES=$(git diff --cached -- $WS_FILES | grep -E "дн[яей]|час[ов]|недел|day|hour|week" | grep -vEi "^-|elapsed|duration|sla|telemetry" || true)
        if [ -n "$TIME_ESTIMATES" ]; then
            echo "❌ Time estimates found in WS files:"
            echo "$TIME_ESTIMATES"
            echo ""
            echo "Remove time-based estimates (days/hours/weeks)."
            echo "Or label as telemetry: 'Elapsed (telemetry): ...' or 'SLA target: ...'"
            echo "Use scope metrics instead (LOC, tokens)."
            exit 1
        fi
        echo "✓ No time estimates"
    else
    echo "  No WS files staged"
fi

# Check 1b: Workstreams layout (WS-only)
echo ""
echo "Check 1b: Workstreams layout"
if echo "$STAGED_FILES" | grep -q "workstreams/"; then
    if ! python3 sdp/scripts/check_workstreams_layout.py; then
        exit 1
    fi
else
    echo "  No workstreams files staged"
fi

# Check 2: No tech debt markers
echo ""
echo "Check 2: No tech debt markers"
# Check only code files, exclude shell scripts (which may have "tech debt" in error messages)
CODE_FILES=$(echo "$STAGED_FILES" | grep -E "\.(py|md|yml|yaml|json)$" | grep -v "\.sh$" || true)
if [ -n "$CODE_FILES" ]; then
    # Exclude lines that mention "No Tech Debt" as a rule (headers, documentation)
    # Also exclude: tables (|), assertions, zero/0 indicators, rule explanations
    TECH_DEBT=$(git diff --cached -- $CODE_FILES | \
        grep -iE "tech.?debt|сделаем.?потом|временн.*решени|later.*fix" | \
        grep -v "^-" | \
        grep -viE "no.?tech.?debt|⛔|запрещено|forbidden|не.*допуск|не.*потом|grep.*потом|потому.что" | \
        grep -viE "zero.tech|tech.debt.*0|assert.*tech|blocks.completion|cannot.be.marked|\|.*tech.*debt.*\||# (no|found)" || true)
else
    TECH_DEBT=""
fi

if [ -n "$TECH_DEBT" ]; then
    echo "❌ Tech debt markers found:"
    echo "$TECH_DEBT"
    echo ""
    echo "Fix the issue now, don't defer it."
    exit 1
fi
echo "✓ No tech debt markers"

# Check 3: Python files - basic checks
echo ""
echo "Check 3: Python code quality"
PY_FILES=$(echo "$STAGED_FILES" | grep "\.py$" || true)

if [ -n "$PY_FILES" ]; then
    # Check for bare except
    BARE_EXCEPT=$(git diff --cached -- $PY_FILES | grep -E "^\+.*except:" | grep -v "except.*:" || true)
    if [ -n "$BARE_EXCEPT" ]; then
        echo "❌ Bare except found:"
        echo "$BARE_EXCEPT"
        echo ""
        echo "Use specific exception types."
        exit 1
    fi
    
    # Check for pass in except
    EXCEPT_PASS=$(git diff --cached -- $PY_FILES | grep -A1 "^\+.*except" | grep "^\+.*pass$" || true)
    if [ -n "$EXCEPT_PASS" ]; then
        echo "⚠️ Warning: except: pass found"
        echo "Consider logging the exception."
    fi
    
    echo "✓ Python checks passed"
else
    echo "  No Python files staged"
fi

# Check 4: Clean Architecture (domain imports)
echo ""
echo "Check 4: Clean Architecture"
DOMAIN_FILES=$(echo "$STAGED_FILES" | grep "domain/.*\.py$" || true)

if [ -n "$DOMAIN_FILES" ]; then
    BAD_IMPORTS=$(git diff --cached -- $DOMAIN_FILES | grep -E "^\+.*from hw_checker\.(infrastructure|presentation)" || true)
    if [ -n "$BAD_IMPORTS" ]; then
        echo "❌ Domain layer importing from infrastructure/presentation:"
        echo "$BAD_IMPORTS"
        echo ""
        echo "Domain layer should not depend on outer layers."
        exit 1
    fi
    echo "✓ Clean Architecture respected"
else
    echo "  No domain files staged"
fi

# Check 5: WS file format (if creating new WS)
echo ""
echo "Check 5: WS file format"
NEW_WS_FILES=$(echo "$STAGED_FILES" | grep "workstreams/backlog/WS-.*\.md$" || true)

if [ -n "$NEW_WS_FILES" ]; then
    for WS_FILE in $NEW_WS_FILES; do
        # Check Goal section exists
        if ! git show ":$WS_FILE" | grep -q "### 🎯"; then
            echo "❌ Missing Goal section in $WS_FILE"
            echo "Add '### 🎯 Цель (Goal)' section."
            exit 1
        fi
        
        # Check Acceptance Criteria exists
        if ! git show ":$WS_FILE" | grep -q "Acceptance Criteria"; then
            echo "❌ Missing Acceptance Criteria in $WS_FILE"
            exit 1
        fi
        
        # Check substream format (if applicable)
        WS_ID=$(basename "$WS_FILE" | grep -oE "WS-[0-9]{3}-[0-9]{2}" || true)
        if [ -n "$WS_ID" ]; then
            # Valid substream format
            echo "✓ $WS_FILE (substream format OK)"
        else
            WS_ID=$(basename "$WS_FILE" | grep -oE "WS-[0-9]{3}" || true)
            if [ -n "$WS_ID" ]; then
                echo "✓ $WS_FILE (main WS format OK)"
            else
                echo "⚠️ Warning: Unusual WS ID format in $WS_FILE"
            fi
        fi
    done
    echo "✓ WS format checks passed"
else
    echo "  No new WS files staged"
fi

echo ""
echo "===================="
echo "✅ Pre-commit checks PASSED"

# Check 6: Breaking changes detection
echo ""
echo "Check 6: Breaking changes"
if [ -f "tools/hw_checker/scripts/detect_breaking_changes.py" ]; then
    cd tools/hw_checker
    if python3 scripts/detect_breaking_changes.py --staged; then
        echo "✓ No breaking changes"
        cd - > /dev/null
    else
        # Breaking changes detected
        CHANGES_COUNT=$(grep -c "category=" BREAKING_CHANGES.md 2>/dev/null || echo "unknown")
        cd - > /dev/null
        
        # Send notification
        bash sdp/notifications/telegram.sh breaking_changes "$CHANGES_COUNT"
        
        echo ""
        echo "⚠️ Breaking changes detected!"
        echo "Review and commit:"
        echo "  - BREAKING_CHANGES.md"
        echo "  - MIGRATION_GUIDE.md"
        exit 1
    fi
else
    echo "  Breaking changes detection script not found (skipping)"
fi

# Check 7: Test quality
echo ""
echo "Check 7: Test quality"
TEST_FILES=$(echo "$STAGED_FILES" | grep -E "tests/.*test_.*\.py$|test_.*\.py$" || true)

if [ -n "$TEST_FILES" ]; then
    cd tools/hw_checker
    # Convert to absolute paths for check_test_quality.py
    ABS_TEST_FILES=""
    for f in $TEST_FILES; do
        if [ -f "$f" ]; then
            ABS_TEST_FILES="$ABS_TEST_FILES $(realpath "$f" 2>/dev/null || echo "$f")"
        elif [ -f "tools/hw_checker/$f" ]; then
            ABS_TEST_FILES="$ABS_TEST_FILES $(realpath "tools/hw_checker/$f" 2>/dev/null || echo "tools/hw_checker/$f")"
        fi
    done
    
    if [ -n "$ABS_TEST_FILES" ]; then
        if ! poetry run python scripts/check_test_quality.py --strict $ABS_TEST_FILES 2>&1; then
            cd - > /dev/null
            echo ""
            echo "❌ Test quality checks failed"
            echo "Fix violations before committing."
            echo "See: docs/test-quality-audit-report.md for details."
            exit 1
        fi
        echo "✓ Test quality checks passed"
    else
        echo "  No valid test files found in staged files"
    fi
    cd - > /dev/null
else
    echo "  No test files staged"
fi

echo ""
echo "===================="
echo "✅ All pre-commit checks PASSED"
