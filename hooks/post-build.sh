#!/bin/sh
# Post-build hook for /build workflow.
# Usage: ./post-build.sh {WS-ID} [beads-status]
# On success: bd update {beads_id} --status completed
WS_ID="${1:-}"
STATUS="${2:-completed}"
if [ -n "$WS_ID" ] && command -v bd >/dev/null 2>&1; then
  # Resolve beads_id from .beads-sdp-mapping.jsonl if present
  if [ -f .beads-sdp-mapping.jsonl ]; then
    beads_id=$(grep "\"sdp_id\": \"$WS_ID\"" .beads-sdp-mapping.jsonl 2>/dev/null | head -1 | sed 's/.*"beads_id": "\([^"]*\)".*/\1/')
    if [ -n "$beads_id" ]; then
      bd update "$beads_id" --status "$STATUS" 2>/dev/null || true
    fi
  fi
fi
exit 0
