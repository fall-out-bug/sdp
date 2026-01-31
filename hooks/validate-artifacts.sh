#!/bin/bash
# SDP Artifact Validation Hook
# Checks for required SDP artifacts (draft, intent, WS) before commits

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Checking SDP artifacts..."

# Track errors
errors=0

# Function to check if a feature has required artifacts
check_feature_artifacts() {
    local feature_id="$1"
    local feature_num="$2"  # Just the number part (e.g., "013" from F013)

    # Check for draft (try multiple patterns)
    local draft_file=""
    for pattern in "docs/drafts/idea-${feature_num}.md" "docs/drafts/idea-f${feature_num}-*.md"; do
        if ls $pattern 2>/dev/null | grep -q .; then
            draft_file=$(ls $pattern 2>/dev/null | head -1)
            break
        fi
    done
    if [ -f "$draft_file" ]; then
        echo -e "${GREEN}✓${NC} Draft: $draft_file"
    else
        echo -e "${YELLOW}⚠${NC} Draft missing: $draft_file"
    fi

    # Check for intent (try multiple patterns)
    local intent_file=""
    for pattern in "docs/intent/${feature_num}.json" "docs/intent/f${feature_num}-*.json"; do
        if ls $pattern 2>/dev/null | grep -q .; then
            intent_file=$(ls $pattern 2>/dev/null | head -1)
            break
        fi
    done
    if [ -f "$intent_file" ]; then
        # Validate intent schema
        if poetry run python -c "
from sdp.schema.validator import IntentValidator
import json
try:
    with open('$intent_file') as f:
        IntentValidator().validate(json.load(f))
    exit(0)
except Exception as e:
        print(f'Validation failed: {e}')
    exit(1)
" 2>/dev/null; then
            echo -e "${GREEN}✓${NC} Intent (valid): $intent_file"
        else
            echo -e "${RED}✗${NC} Intent invalid: $intent_file"
            ((errors++))
        fi
    else
        echo -e "${YELLOW}⚠${NC} Intent missing: $intent_file"
    fi

    # Check for workstreams
    local ws_pattern="docs/workstreams/backlog/00-${feature_num}-*.md"
    local alt_pattern="docs/workstreams/backlog/*-${feature_num}-*.md"
    local ws_count=$(ls $ws_pattern 2>/dev/null | wc -l | tr -d ' ')
    if [ "$ws_count" -eq 0 ]; then
        ws_count=$(ls $alt_pattern 2>/dev/null | wc -l | tr -d ' ')
        ws_pattern="$alt_pattern"
    fi

    if [ "$ws_count" -gt 0 ]; then
        echo -e "${GREEN}✓${NC} Workstreams: $ws_count files found"
        ls -1 $ws_pattern 2>/dev/null | head -5
        if [ "$ws_count" -gt 5 ]; then
            echo -e "  ... and $((ws_count - 5)) more"
        fi
    else
        echo -e "${YELLOW}⚠${NC} No workstreams found for $feature_id"
    fi
}

# Detect feature from current changes
# Look for changes in docs/ or skills/ that indicate feature work
changed_files=$(git diff --cached --name-only --diff-filter=ACM)

# Check if any feature-related files changed
if echo "$changed_files" | grep -q "docs/intent/"; then
    echo "📝 Intent files detected - validating associated feature"

    # Extract feature from intent files
    for intent_file in $(echo "$changed_files" | grep "docs/intent/" || true); do
        # Match patterns like 013.json or f013-ai-comm.json
        feature_num=$(basename "$intent_file" | sed -E 's/^f?0*([0-9]+).*/\1/')
        # Pad to 3 digits if needed
        feature_num=$(printf "%03d" "$feature_num" 2>/dev/null)
        if [[ -n "$feature_num" ]]; then
            check_feature_artifacts "F${feature_num}" "$feature_num"
        fi
    done
fi

if echo "$changed_files" | grep -q ".claude/skills/"; then
    echo "🛠️  Skills changed - verify artifacts exist"
    # Check for recent features by scanning draft files
    for draft_file in docs/drafts/idea-*.md docs/drafts/idea-f*.md; do
        if [ -f "$draft_file" ]; then
            # Extract feature number from patterns like idea-013.md or idea-f013-ai-comm.md
            feature_num=$(basename "$draft_file" | sed -E 's/idea-?f?0*([0-9]+).*/\1/')
            if [ -n "$feature_num" ]; then
                check_feature_artifacts "F${feature_num}" "$feature_num"
            fi
        fi
    done
fi

# Summary
if [ $errors -eq 0 ]; then
    echo -e "\n${GREEN}✓ Artifact validation passed${NC}"
    exit 0
else
    echo -e "\n${RED}✗ Artifact validation failed with $errors error(s)${NC}"
    exit 1
fi
