#!/usr/bin/env bash
# Validate working-directory input
# Usage: validate-input.sh <working-directory>
# Exit codes: 0 = valid, 1 = invalid

set -euo pipefail

WORKING_DIR="${1:-.}"

# Validate working directory input to prevent path traversal
# Check for absolute paths (not allowed)
if [[ "$WORKING_DIR" = /* ]]; then
  echo "❌ Error: Absolute paths not allowed in working-directory" >&2
  echo "   Received: $WORKING_DIR" >&2
  echo "   Use relative path from repository root" >&2
  exit 1
fi

# Check for path traversal attempts
if [[ "$WORKING_DIR" == *".."* ]]; then
  echo "❌ Error: Path traversal not allowed in working-directory" >&2
  echo "   Received: $WORKING_DIR" >&2
  echo "   Use relative path without '..'" >&2
  exit 1
fi

echo "✅ Working directory validated: $WORKING_DIR"
exit 0
