#!/bin/bash
# sdp/hooks/post-codereview.sh
# Post-codereview checks for /codereview command
# Usage: ./post-codereview.sh F{XX}

set -e

FEATURE=$1
SKIP_UAT_CHECK="${SKIP_UAT_CHECK:-0}"
SKIP_MOVE_CHECK="${SKIP_MOVE_CHECK:-0}"

if [ -z "$FEATURE" ]; then
    echo "❌ Usage: ./post-codereview.sh F{XX}"
    echo "   Example: ./post-codereview.sh F191"
    exit 1
fi

# Extract feature number (F191 -> 191)
FEATURE_NUM=$(echo "$FEATURE" | grep -oE "[0-9]+")

echo "🔍 Post-codereview checks for $FEATURE"
echo "========================================"

WS_DIR="tools/hw_checker/docs/workstreams"
UAT_DIR="tools/hw_checker/docs/uat"

# Find all WS files for this feature
WS_FILES=$(find "$WS_DIR" -name "WS-${FEATURE_NUM}*.md" -o -name "WS-${FEATURE_NUM}-*.md" 2>/dev/null)
WS_COUNT=$(echo "$WS_FILES" | grep -c "WS-" || echo "0")

if [ "$WS_COUNT" -eq 0 ]; then
    echo "❌ No WS files found for feature $FEATURE"
    exit 1
fi

echo "Found $WS_COUNT workstreams for $FEATURE"

# Check 1: Review Results in all WS files
echo ""
echo "Check 1: Review Results in all WS files"
MISSING_REVIEW=0
for WS_FILE in $WS_FILES; do
    if grep -q "### Review Results" "$WS_FILE"; then
        WS_NAME=$(basename "$WS_FILE")
        VERDICT=$(grep -A5 "### Review Results" "$WS_FILE" | grep -oE "APPROVED|CHANGES REQUESTED" | head -1 || echo "N/A")
        echo "  ✓ $WS_NAME: $VERDICT"
    else
        echo "  ❌ Missing Review Results: $WS_FILE"
        MISSING_REVIEW=1
    fi
done

if [ "$MISSING_REVIEW" -eq 1 ]; then
    echo ""
    echo "❌ Some WS files missing Review Results"
    echo "   Append Review Results to each WS file"
    exit 1
fi
echo "✓ All WS have Review Results"

# Check 2: Determine overall verdict
echo ""
echo "Check 2: Overall verdict"
CHANGES_REQUESTED=$(echo "$WS_FILES" | xargs grep -l "CHANGES REQUESTED" 2>/dev/null | wc -l)
if [ "$CHANGES_REQUESTED" -gt 0 ]; then
    echo "⚠️ Overall verdict: CHANGES REQUESTED ($CHANGES_REQUESTED WS need fixes)"
    echo ""
    echo "Next steps:"
    echo "  1. Fix issues in WS files marked CHANGES REQUESTED"
    echo "  2. Re-run /build for affected WS"
    echo "  3. Re-run /codereview $FEATURE"
    exit 0  # Not an error, just not approved yet
fi
echo "✓ Overall verdict: APPROVED"

# Check 3: UAT Guide exists (only if APPROVED)
echo ""
echo "Check 3: UAT Guide"
UAT_FILE="$UAT_DIR/${FEATURE}-uat-guide.md"
ALT_UAT_FILE="$UAT_DIR/F${FEATURE_NUM}-uat-guide.md"

if [ "$SKIP_UAT_CHECK" = "1" ]; then
    echo "⚠️ Skipped (SKIP_UAT_CHECK=1)"
elif [ -f "$UAT_FILE" ] || [ -f "$ALT_UAT_FILE" ]; then
    FOUND_UAT=$(ls "$UAT_FILE" "$ALT_UAT_FILE" 2>/dev/null | head -1)
    echo "✓ UAT Guide found: $FOUND_UAT"
else
    echo "❌ UAT Guide NOT found"
    echo "   Expected: $UAT_FILE or $ALT_UAT_FILE"
    echo ""
    echo "   Create UAT Guide using template: sdp/templates/uat-guide.md"
    echo ""
    echo "   To skip: SKIP_UAT_CHECK=1 ./post-codereview.sh $FEATURE"
    exit 1
fi

# Check 4: Completed WS should be in completed/ folder
echo ""
echo "Check 4: WS file locations"
if [ "$SKIP_MOVE_CHECK" = "1" ]; then
    echo "⚠️ Skipped (SKIP_MOVE_CHECK=1)"
else
    MISPLACED=0
    for WS_FILE in $WS_FILES; do
        # Check if WS is completed (has APPROVED verdict)
        if grep -q "APPROVED" "$WS_FILE" && grep -q "status: completed" "$WS_FILE"; then
            # Should be in completed/ folder
            if echo "$WS_FILE" | grep -q "/backlog/"; then
                echo "  ⚠️ Completed WS in backlog/: $(basename "$WS_FILE")"
                MISPLACED=1
            fi
        fi
    done
    
    if [ "$MISPLACED" -eq 1 ]; then
        echo ""
        echo "⚠️ Some completed WS files are still in backlog/"
        echo "   Move them to completed/ folder:"
        echo "   git mv backlog/WS-*.md completed/$(date +%Y-%m)/"
        # Warning only, not blocking
    else
        echo "✓ All WS files in correct locations"
    fi
fi

# Check 5: Feature Summary output reminder
echo ""
echo "Check 5: Feature Summary"
echo "  Reminder: Output Feature Summary to user with:"
echo "  - Total WS count"
echo "  - Verdict per WS"
echo "  - Overall verdict"
echo "  - Next steps"

echo ""
echo "========================================"
echo "✅ Post-codereview checks PASSED"
echo ""
echo "Next steps:"
echo "  1. Human UAT: Review $UAT_FILE"
echo "  2. If UAT passes: /deploy $FEATURE"
