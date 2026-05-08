#!/bin/sh
# Pre-build hook for /build workflow.
# Usage: ./pre-build.sh {WS-ID}
WS_ID="${1:-}"
if [ -n "$WS_ID" ]; then
  if command -v sdp >/dev/null 2>&1; then
    sdp guard activate "$WS_ID" 2>/dev/null || true
  fi
fi
exit 0
