#!/bin/bash
# sdp/scripts/check_complexity.sh
# Check cyclomatic complexity using radon
# Reads thresholds from quality-gate.toml

set -e

# Default thresholds (can be overridden by quality-gate.toml)
MAX_CC=10
MAX_AVG_CC=5

# Load configuration from quality-gate.toml if it exists
if [ -f "quality-gate.toml" ]; then
    # Extract max_cc (complexity max_cc)
    CONFIG_MAX_CC=$(grep -A2 '\[complexity\]' quality-gate.toml | grep 'max_cc' | grep -oE '[0-9]+' || echo "")
    if [ -n "$CONFIG_MAX_CC" ]; then
        MAX_CC=$CONFIG_MAX_CC
    fi

    # Extract max_average_cc
    CONFIG_AVG_CC=$(grep -A2 '\[complexity\]' quality-gate.toml | grep 'max_average_cc' | grep -oE '[0-9]+' || echo "")
    if [ -n "$CONFIG_AVG_CC" ]; then
        MAX_AVG_CC=$CONFIG_AVG_CC
    fi
fi

TARGET_PATH=${1:-"."}
MIN_GRADE=${2:-"C"}

echo "🔍 Complexity Check (Radon)"
echo "============================"
echo "Target: $TARGET_PATH"
echo "Max CC: $MAX_CC"
echo "Max Average CC: $MAX_AVG_CC"
echo ""

# Check if radon is available (directly or via poetry)
if command -v radon &> /dev/null; then
    RADON_CMD="radon"
elif command -v poetry &> /dev/null; then
    RADON_CMD="poetry run radon"
else
    echo "❌ radon not found"
    echo ""
    echo "Install: pip install radon"
    echo "Or: poetry add radon --group dev"
    exit 1
fi

# Run radon cc - show complexity for each function
echo "Cyclomatic Complexity by function:"
echo "-----------------------------------"
$RADON_CMD cc "$TARGET_PATH" -a -s --min "A"

# Check if any function exceeds max_cc
VIOLATIONS=$($RADON_CMD cc "$TARGET_PATH" -s --min "A" 2>/dev/null | grep -E "[F|E|D]" | grep -oE "\([0-9]+\)" | grep -oE "[0-9]+" | awk -v max="$MAX_CC" '$1 > max' || echo "")

if [ -n "$VIOLATIONS" ]; then
    echo ""
    echo "❌ Complexity violations found (max: $MAX_CC)"
    echo "$VIOLATIONS"
    echo ""
    echo "Fix: Refactor complex functions into smaller functions"
    echo "  - Extract logic into separate helper functions"
    echo "  - Use strategy pattern for complex conditionals"
    echo "  - Follow Single Responsibility Principle"
    exit 1
fi

# Check average complexity
AVG_CC=$($RADON_CMD cc "$TARGET_PATH" -a --min "A" 2>/dev/null | grep "Average" | grep -oE "[0-9.]+" || echo "0")

if [ -n "$AVG_CC" ] && [ $(echo "$AVG_CC > $MAX_AVG_CC" | bc -l 2>/dev/null || echo "0") -eq 1 ]; then
    echo ""
    echo "❌ Average complexity too high: $AVG_CC (max: $MAX_AVG_CC)"
    echo ""
    echo "Fix: Overall codebase needs refactoring"
    echo "  - Break down complex modules"
    echo "  - Simplify control flow"
    exit 1
fi

echo ""
echo "✓ Complexity check passed"
echo "  Max CC: $MAX_CC"
echo "  Average CC: $AVG_CC (≤ $MAX_AVG_CC)"
