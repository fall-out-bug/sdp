#!/usr/bin/env bash
# Pre-tool constraint enforcement for SDP agent sessions.
# Blocks destructive commands and evaluates agent-constraints.yaml rules.

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

# Hard-blocked commands (always, regardless of constraints file)
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

# SDP agent-constraints.yaml enforcement.
# Reads the current phase from .sdp/checkpoints/ if available.
if command -v sdp-guard >/dev/null 2>&1 && [ -f ".sdp/agent-constraints.yaml" ]; then
  # Determine current phase from checkpoint (default: build)
  CURRENT_PHASE="build"
  for cp_file in .sdp/checkpoints/*.json; do
    if [ -f "$cp_file" ]; then
      PHASE_FROM_CP=$(jq -r '.phase // ""' "$cp_file" 2>/dev/null || true)
      if [ -n "$PHASE_FROM_CP" ] && [ "$PHASE_FROM_CP" != "done" ]; then
        CURRENT_PHASE="$PHASE_FROM_CP"
        break
      fi
    fi
  done

  # Check the command against constraint rules
  RESULT=$(sdp-guard --check-constraints --phase="$CURRENT_PHASE" --command="$COMMAND" 2>&1 || true)
  EXIT_CODE=$?

  if [ $EXIT_CODE -eq 2 ]; then
    # halt/escalate: stop agent session
    echo "$RESULT"
    echo "HALT: SDP constraint violation requires agent session to stop."
    exit 2
  elif [ $EXIT_CODE -eq 1 ]; then
    # block: reject this specific action
    echo "$RESULT"
    exit 2
  fi
  # warn (exit 0): log and continue
  if [ -n "$RESULT" ] && echo "$RESULT" | grep -q "WARN"; then
    echo "$RESULT" >&2
  fi
fi

exit 0
