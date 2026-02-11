#!/usr/bin/env bash
# Run lightweight checks after editing GitHub workflow files.

set -euo pipefail

PAYLOAD="$(cat 2>/dev/null || true)"
if [ -z "$PAYLOAD" ]; then
  exit 0
fi

TOOL_NAME="$(echo "$PAYLOAD" | jq -r '.tool_name // .tool // ""' 2>/dev/null || true)"
FILE_PATH_RAW="$(echo "$PAYLOAD" | jq -r '.tool_input.file_path // .tool_input.path // .tool_input.target_file // ""' 2>/dev/null || true)"

case "$TOOL_NAME" in
  Edit|Write|MultiEdit) ;;
  *) exit 0 ;;
esac

if [ -z "$FILE_PATH_RAW" ]; then
  exit 0
fi

FILE_PATH="$FILE_PATH_RAW"
if [ ! -f "$FILE_PATH" ]; then
  FILE_PATH="$(pwd)/$FILE_PATH_RAW"
fi

if [ ! -f "$FILE_PATH" ]; then
  exit 0
fi

if ! echo "$FILE_PATH" | grep -Eq '/\.github/workflows/.*\.ya?ml$'; then
  exit 0
fi

WARN=0

if grep -q $'\t' "$FILE_PATH"; then
  echo "WARNING: workflow contains tab characters: $FILE_PATH"
  WARN=1
fi

if ! grep -Eq '^[[:space:]]*name:' "$FILE_PATH"; then
  echo "WARNING: workflow missing top-level 'name:' in $FILE_PATH"
  WARN=1
fi

if ! grep -Eq '^[[:space:]]*on:' "$FILE_PATH"; then
  echo "WARNING: workflow missing top-level 'on:' in $FILE_PATH"
  WARN=1
fi

if ! grep -Eq '^[[:space:]]*jobs:' "$FILE_PATH"; then
  echo "WARNING: workflow missing top-level 'jobs:' in $FILE_PATH"
  WARN=1
fi

if command -v actionlint >/dev/null 2>&1; then
  if ! actionlint "$FILE_PATH" >/dev/null 2>&1; then
    echo "WARNING: actionlint reported issues in $FILE_PATH"
    WARN=1
  fi
fi

if [ "$WARN" -eq 0 ]; then
  echo "Workflow check: OK ($FILE_PATH_RAW)"
fi

exit 0
