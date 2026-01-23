#!/bin/bash
# sdp/hooks/verification-completion.sh
# Verification hook for /build completion
# Enforces evidence-based claims and detects red flag phrases
# Usage: ./verification-completion.sh <ws-file-path>

set -e

WS_FILE="$1"

if [ -z "$WS_FILE" ]; then
    echo "❌ Usage: ./verification-completion.sh <ws-file-path>" >&2
    exit 1
fi

if [ ! -f "$WS_FILE" ]; then
    echo "❌ WS file not found: $WS_FILE" >&2
    exit 1
fi

echo "🔍 Verification check for completion claims"
echo "============================================"

# Check if Execution Report section exists
if ! grep -q "### Execution Report\|## Execution Report" "$WS_FILE"; then
    echo "❌ Execution Report section not found" >&2
    echo "   Add '### Execution Report' section to WS file" >&2
    exit 1
fi

# Extract Execution Report section (from "### Execution Report" to next "##" or "---" or end of file)
EXEC_REPORT=$(awk '
    BEGIN { in_section=0 }
    /^### Execution Report|^## Execution Report/ { in_section=1; print; next }
    in_section && /^---$/ { exit }
    in_section && /^## [^E]/ { exit }
    in_section { print }
' "$WS_FILE")

if [ -z "$EXEC_REPORT" ]; then
    echo "❌ Execution Report section is empty" >&2
    exit 1
fi

# Red flag phrases to detect
# Single words (use word boundaries)
SINGLE_WORD_FLAGS=("should" "probably" "might" "may")
# Multi-word phrases (use simple pattern matching)
MULTI_WORD_FLAGS=("seems to" "seems like" "appears to")

RED_FLAG_FOUND=0
RED_FLAG_PHRASES=""

# Check single-word flags with word boundaries
for flag in "${SINGLE_WORD_FLAGS[@]}"; do
    if echo "$EXEC_REPORT" | grep -qi "\b${flag}\b"; then
        RED_FLAG_FOUND=1
        if [ -z "$RED_FLAG_PHRASES" ]; then
            RED_FLAG_PHRASES="$flag"
        else
            RED_FLAG_PHRASES="$RED_FLAG_PHRASES, $flag"
        fi
    fi
done

# Check multi-word flags (no word boundaries needed)
for flag in "${MULTI_WORD_FLAGS[@]}"; do
    if echo "$EXEC_REPORT" | grep -qi "${flag}"; then
        RED_FLAG_FOUND=1
        if [ -z "$RED_FLAG_PHRASES" ]; then
            RED_FLAG_PHRASES="$flag"
        else
            RED_FLAG_PHRASES="$RED_FLAG_PHRASES, $flag"
        fi
    fi
done

if [ "$RED_FLAG_FOUND" -eq 1 ]; then
    echo "❌ Red flag phrases detected: $RED_FLAG_PHRASES" >&2
    echo "" >&2
    echo "   Iron Law: No completion without fresh verification" >&2
    echo "   Replace uncertain claims with command output evidence" >&2
    echo "" >&2
    echo "   Example:" >&2
    echo "   ❌ 'Tests should pass'" >&2
    echo "   ✅ '```bash" >&2
    echo "      $ pytest tests/unit/test_service.py -v" >&2
    echo "      ===== 15 passed in 0.5s =====" >&2
    echo "      ```'" >&2
    exit 1
fi

# Check for command output evidence (bash/shell code blocks with actual output)
# Look for code blocks that contain command output (not just commands)
EVIDENCE_FOUND=0

# Extract all bash/shell code blocks from Execution Report
CODE_BLOCKS=$(echo "$EXEC_REPORT" | awk '
    BEGIN { in_block=0; block="" }
    /^```(bash|sh|shell)$/ { in_block=1; block=""; next }
    /^```$/ { 
        if (in_block) {
            print block
            block=""
            in_block=0
        }
        next
    }
    in_block { 
        if (block != "") block=block "\n"
        block=block $0
    }
    END {
        if (in_block && block != "") {
            print block
        }
    }
')

if [ -n "$CODE_BLOCKS" ]; then
    # Check if any code block contains actual output (not just commands)
    # Evidence = lines that look like output (contain "passed", "PASSED", "✓", "=====", etc.)
    while IFS= read -r block; do
        # Skip empty blocks
        if [ -z "$(echo "$block" | tr -d '[:space:]')" ]; then
            continue
        fi
        
        # Check if block contains output indicators
        # Strategy: Filter out ALL command lines (starting with $), then check remaining lines
        # Use sed to delete command lines, then check if any output remains
        NON_COMMAND_LINES=$(echo "$block" | sed '/^[[:space:]]*\$/d')
        
        if [ -n "$NON_COMMAND_LINES" ]; then
            # Check for result keywords in non-command lines
            if echo "$NON_COMMAND_LINES" | grep -qiE "(passed|failed|error|✓|=====|coverage|%|all checks|tests)"; then
                EVIDENCE_FOUND=1
                break
            fi
            
            # Check for numbers in non-command lines
            if echo "$NON_COMMAND_LINES" | grep -qE "[0-9]+"; then
                EVIDENCE_FOUND=1
                break
            fi
        fi
    done <<< "$CODE_BLOCKS"
fi

if [ "$EVIDENCE_FOUND" -eq 0 ]; then
    echo "❌ No command output evidence found" >&2
    echo "" >&2
    echo "   Iron Law: No completion without fresh verification" >&2
    echo "   Execution Report must include command output as evidence" >&2
    echo "" >&2
    echo "   Required format:" >&2
    echo "   ```bash" >&2
    echo "   $ pytest tests/unit/test_service.py -v" >&2
    echo "   ===== 15 passed in 0.5s =====" >&2
    echo "   ```" >&2
    echo "" >&2
    echo "   Not just commands, but actual output showing results" >&2
    exit 1
fi

echo "✓ No red flag phrases detected"
echo "✓ Command output evidence found"
echo ""
echo "✅ Verification check PASSED"
exit 0
