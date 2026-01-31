#!/bin/bash
# sdp/hooks/pre-build.sh
# Pre-build checks for /build command
# Usage: ./pre-build.sh WS-060-01

set -e

WS_ID=$1

if [ -z "$WS_ID" ]; then
    echo "❌ Usage: ./pre-build.sh WS-ID"
    exit 1
fi

echo "🔍 Pre-build checks for $WS_ID"
echo "================================"

# Find WS file (support both SDP and hw_checker layouts)
REPO_ROOT=$(git rev-parse --show-toplevel)
WS_FILE=""
WS_DIR=""
for CANDIDATE in "docs/workstreams" "tools/hw_checker/docs/workstreams"; do
    FULL_DIR="$REPO_ROOT/$CANDIDATE"
    if [ -d "$FULL_DIR" ]; then
        FOUND=$(find "$FULL_DIR" -name "${WS_ID}-*.md" 2>/dev/null | head -1)
        if [ -n "$FOUND" ]; then
            WS_FILE="$FOUND"
            WS_DIR="$CANDIDATE"
            break
        fi
    fi
done

if [ -z "$WS_FILE" ]; then
    echo "❌ WS file not found: ${WS_ID}-*.md"
    echo "   Searched in: docs/workstreams, tools/hw_checker/docs/workstreams"
    exit 1
fi

echo "✓ Found: $WS_FILE"

# Check 1: Goal defined
echo ""
echo "Check 1: Goal defined"
if grep -q "### 🎯 Цель\|### 🎯 Goal" "$WS_FILE"; then
    echo "✓ Goal section found"
else
    echo "❌ Goal section not found"
    echo "   Add '### 🎯 Цель (Goal)' section to WS file"
    exit 1
fi

# Check 2: Acceptance Criteria
echo ""
echo "Check 2: Acceptance Criteria"
if grep -q "Acceptance Criteria" "$WS_FILE"; then
    AC_COUNT=$(grep -c "\- \[ \]" "$WS_FILE" || echo "0")
    echo "✓ Acceptance Criteria found ($AC_COUNT items)"
else
    echo "❌ Acceptance Criteria not found"
    echo "   Add Acceptance Criteria checklist to WS file"
    exit 1
fi

# Check 3: Scope not LARGE
echo ""
echo "Check 3: Scope check"
if grep -q "🔴.*LARGE\|LARGE.*РАЗБИТЬ" "$WS_FILE"; then
    echo "❌ Scope is LARGE — split WS first"
    echo "   Run /design to break down into smaller WS"
    exit 1
else
    echo "✓ Scope is acceptable"
fi

# Check 4: Dependencies completed (if any)
echo ""
echo "Check 4: Dependencies"
DEP=$(grep -A1 "### Зависимость\|### Dependency" "$WS_FILE" | tail -1 | tr -d '[]' | xargs)

if [ -z "$DEP" ] || [ "$DEP" = "Независимый" ] || [ "$DEP" = "Independent" ] || [ "$DEP" = "-" ]; then
    echo "✓ No dependencies (independent WS)"
else
    echo "  Dependencies: $DEP"
    # Check if dependency is completed
    INDEX_FILE="$REPO_ROOT/$WS_DIR/INDEX.md"
    if grep -q "$DEP.*completed\|$DEP.*✅" "$INDEX_FILE" 2>/dev/null; then
        echo "✓ Dependency $DEP is completed"
    else
        echo "⚠️ Warning: Dependency $DEP may not be completed"
        echo "  Check INDEX.md to verify status"
    fi
fi

# Check 5: Input files exist
echo ""
echo "Check 5: Input files"
INPUT_FILES=$(grep -A10 "### Входные файлы\|### Input" "$WS_FILE" | grep "^\- \`" | sed 's/.*`\([^`]*\)`.*/\1/' | head -5)

if [ -n "$INPUT_FILES" ]; then
    for FILE in $INPUT_FILES; do
        if [ -f "$FILE" ]; then
            echo "✓ $FILE exists"
        else
            echo "⚠️ $FILE not found (may be created during build)"
        fi
    done
else
    echo "  No input files specified"
fi

echo ""
echo "================================"
echo "✅ Pre-build checks PASSED"
echo ""
echo "Ready to execute: /build $WS_ID"
