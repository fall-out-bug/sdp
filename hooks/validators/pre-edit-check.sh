#!/bin/bash
# sdp/hooks/validators/pre-edit-check.sh
# Validates file before editing - checks Clean Architecture

set -e

# Get file path from stdin (Claude Code hook format)
FILE_PATH="${1:-}"

if [ -z "$FILE_PATH" ]; then
    # Try to read from stdin JSON
    FILE_PATH=$(cat 2>/dev/null | jq -r '.tool_input.file_path // .tool_input.path // ""' 2>/dev/null || echo "")
fi

if [ -z "$FILE_PATH" ] || [ ! -f "$FILE_PATH" ]; then
    exit 0  # No file path or file doesn't exist, skip validation
fi

# Check if editing domain layer with infrastructure/presentation imports
if echo "$FILE_PATH" | grep -q "domain/"; then
    if grep -qE "from hw_checker\.(infrastructure|presentation)" "$FILE_PATH" 2>/dev/null; then
        echo "BLOCKED: Domain layer file contains infrastructure/presentation imports"
        echo "Clean Architecture violation detected in: $FILE_PATH"
        echo ""
        echo "Domain layer must NOT import from:"
        echo "  - hw_checker.infrastructure.*"
        echo "  - hw_checker.presentation.*"
        exit 2  # Exit code 2 blocks the action
    fi
fi

# Check if editing application layer with presentation imports
if echo "$FILE_PATH" | grep -q "application/"; then
    if grep -qE "from hw_checker\.presentation" "$FILE_PATH" 2>/dev/null; then
        echo "BLOCKED: Application layer file contains presentation imports"
        echo "Clean Architecture violation detected in: $FILE_PATH"
        exit 2
    fi
fi

exit 0
