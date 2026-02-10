#!/bin/bash
# Post-WS-Complete Hook - Reminder to close workstream/beads after completion

set -e

WS_ID=${1:-}

if [ -z "$WS_ID" ]; then
    echo "Usage: hooks/post-ws-complete.sh <WS_ID> [bypass] [reason]"
    exit 1
fi

echo "WS $WS_ID complete. Run: bd close <beads_id> --reason 'WS completed'; bd sync"
exit 0
