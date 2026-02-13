#!/bin/bash
# Drift Detection Script for Prompt Trees
# WS-067-02: AC5
#
# Usage: ./hooks/check-prompt-drift.sh
#
# Checks for:
# 1. Duplicate prompt files outside canonical prompts/ directory
# 2. Symlink integrity for .claude/agents and .claude/skills

set -e

CANONICAL="prompts/"
ERRORS=0

echo "=== Prompt Tree Drift Detection ==="

# Check 1: No duplicate prompt files outside canonical
echo "Checking for duplicate prompt trees..."
DUPLICATES=$(find . -path ./prompts -prune -o -path ./.git -prune -o -name "*.md" -path "*prompts/agents*" -print 2>/dev/null || true)
if [ -n "$DUPLICATES" ]; then
    echo "ERROR: Duplicate prompt files found outside canonical path:"
    echo "$DUPLICATES"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ No duplicate prompt trees"
fi

# Check 2: Symlink integrity
echo ""
echo "Checking symlinks..."

for link in .claude/agents .claude/skills; do
    if [ -L "$link" ]; then
        target=$(readlink "$link")
        if [ -d "$link" ]; then
            echo "✅ $link -> $target (valid)"
        else
            echo "ERROR: $link points to non-existent target: $target"
            ERRORS=$((ERRORS + 1))
        fi
    else
        echo "ERROR: $link is not a symlink"
        ERRORS=$((ERRORS + 1))
    fi
done

# Summary
echo ""
if [ $ERRORS -eq 0 ]; then
    echo "✅ All checks passed"
    exit 0
else
    echo "❌ $ERRORS error(s) found"
    exit 1
fi
