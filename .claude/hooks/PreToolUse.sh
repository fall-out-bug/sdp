#!/usr/bin/env bash
# Block destructive git commands before execution.

set -euo pipefail

PAYLOAD="$(cat 2>/dev/null || true)"

if [ -z "$PAYLOAD" ]; then
  exit 0
fi

TOOL_NAME="$(echo "$PAYLOAD" | jq -r '.tool_name // .tool // ""' 2>/dev/null || true)"
COMMAND="$(echo "$PAYLOAD" | jq -r '.tool_input.command // .tool_input.cmd // ""' 2>/dev/null || true)"

if [ "$TOOL_NAME" != "Bash" ] || [ -z "$COMMAND" ]; then
  exit 0
fi

if echo "$COMMAND" | grep -Eiq '(^|[[:space:]])git[[:space:]]+reset[[:space:]]+--hard([[:space:]]|$)'; then
  echo "BLOCKED: destructive git command is not allowed: git reset --hard"
  exit 2
fi

if echo "$COMMAND" | grep -Eiq '(^|[[:space:]])git[[:space:]]+clean([[:space:]]|$)'; then
  echo "BLOCKED: destructive git command is not allowed: git clean"
  exit 2
fi

if echo "$COMMAND" | grep -Eiq '(^|[[:space:]])git[[:space:]]+checkout[[:space:]]+--([[:space:]]|$)'; then
  echo "BLOCKED: destructive git command is not allowed: git checkout --"
  exit 2
fi

if echo "$COMMAND" | grep -Eiq '(^|[[:space:]])git[[:space:]]+restore[[:space:]].*--source=.*[[:space:]]+--([[:space:]]|$)'; then
  echo "BLOCKED: destructive git command is not allowed: git restore --source ... --"
  exit 2
fi

exit 0
