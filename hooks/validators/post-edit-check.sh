#!/bin/bash
# scripts/hooks/validators/post-edit-check.sh
# Validates file after editing - checks TODO/FIXME and file size

FILE_PATH="${1:-}"

if [ -z "$FILE_PATH" ] || [ ! -f "$FILE_PATH" ]; then
    exit 0
fi

WARNINGS=0

# Check for TODO/FIXME in Python files
if echo "$FILE_PATH" | grep -qE "\.py$"; then
    TODO_COUNT=$(grep -cE "TODO|FIXME|HACK|XXX" "$FILE_PATH" 2>/dev/null || echo "0")
    if [ "$TODO_COUNT" -gt 0 ]; then
        echo "WARNING: Found $TODO_COUNT TODO/FIXME/HACK/XXX markers in $FILE_PATH"
        echo "         Remove these before completing the workstream"
        WARNINGS=$((WARNINGS + 1))
    fi
    
    # Check file size
    LINE_COUNT=$(wc -l < "$FILE_PATH")
    if [ "$LINE_COUNT" -gt 200 ]; then
        echo "WARNING: File $FILE_PATH has $LINE_COUNT lines (limit: 200)"
        echo "         Consider splitting this file"
        WARNINGS=$((WARNINGS + 1))
    elif [ "$LINE_COUNT" -gt 150 ]; then
        echo "INFO: File $FILE_PATH has $LINE_COUNT lines (approaching 200 limit)"
    fi
    
    # Check for bare except
    if grep -qE "except\s*:" "$FILE_PATH" 2>/dev/null; then
        BARE_EXCEPT=$(grep -cE "except\s*:" "$FILE_PATH" 2>/dev/null || echo "0")
        echo "WARNING: Found $BARE_EXCEPT bare 'except:' clause(s) in $FILE_PATH"
        echo "         Use specific exception types"
        WARNINGS=$((WARNINGS + 1))
    fi
    
    # Check for missing type hints (functions without ->)
    FUNCS_WITHOUT_RETURN=$(grep -cE "^\s*def\s+\w+\([^)]*\)\s*:" "$FILE_PATH" 2>/dev/null || echo "0")
    if [ "$FUNCS_WITHOUT_RETURN" -gt 0 ]; then
        echo "INFO: $FUNCS_WITHOUT_RETURN function(s) may be missing return type hints"
    fi
fi

if [ "$WARNINGS" -gt 0 ]; then
    echo ""
    echo "Total warnings: $WARNINGS"
fi

# Don't block, just warn
exit 0
