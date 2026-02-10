#!/bin/bash
# sdp/scripts/check_complexity.sh
# Go: run go vet for basic checks. Python/radon removed.

set -e

TARGET_PATH=${1:-"sdp-plugin"}

echo "🔍 Complexity / quality check (Go)"
echo "==================================="
echo "Target: $TARGET_PATH"
echo ""

if [ ! -d "$TARGET_PATH" ]; then
    echo "⚠️ Directory not found: $TARGET_PATH (skipping)"
    exit 0
fi

if [ -f "$TARGET_PATH/go.mod" ]; then
    (cd "$TARGET_PATH" && go vet ./...)
    echo ""
    echo "✓ go vet passed"
else
    echo "⚠️ No go.mod in $TARGET_PATH (skipping)"
fi
exit 0
