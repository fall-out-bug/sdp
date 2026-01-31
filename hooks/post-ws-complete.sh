#!/bin/bash
# Post-WS-Complete Hook - Verifies completion before status change

set -e

WS_ID=${1:-}
BYPASS=${2:-false}
REASON=${3:-""}

if [ -z "$WS_ID" ]; then
    echo "Usage: hooks/post-ws-complete.sh <WS_ID> [bypass] [reason]"
    exit 1
fi

# Run Python verification
if [ "$BYPASS" = "true" ]; then
    if [ -z "$REASON" ]; then
        echo "❌ Bypass requires reason argument"
        echo "Usage: hooks/post-ws-complete.sh <WS_ID> true \"<reason>\""
        exit 1
    fi
    
    python3 -m sdp.hooks.ws_complete "$WS_ID" --bypass --reason "$REASON"
else
    python3 -m sdp.hooks.ws_complete "$WS_ID"
fi

exit $?
