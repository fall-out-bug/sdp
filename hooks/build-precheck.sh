#!/usr/bin/env bash
# F142-07: pre-build precheck — refuse to /build a workstream without a ws file.
#
# Usage: scripts/hooks/build-precheck.sh <WS-ID>
#   WS-ID example: 00-141-02 or 00-082-01
#
# Exit codes:
#   0  ws file exists and is not status=design-pending → proceed
#   1  ws file missing OR design-pending → block /build with a clear message
#   2  bad usage
#
# Rule: no workstream → no execution. Aligns with:
#   - scripts/deliver-pick.sh (refuses to pick leafless features)
#   - sdp-guard --ws <id> (errors when ws scope cannot be parsed)
#   - sdp doctor backlog (CI gate)

set -uo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: build-precheck.sh <WS-ID>" >&2
  exit 2
fi

WS_ID="$1"
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
WS_FILE="${REPO_ROOT}/docs/workstreams/backlog/${WS_ID}.md"

if [[ ! -f "$WS_FILE" ]]; then
  cat >&2 <<EOF
build-precheck: REFUSE — no ws file at ${WS_FILE}

Rule (F142-07): /build requires an executable workstream file. Either
  1. The WS-ID is wrong — check 'docs/workstreams/INDEX.md' or 'sdp doctor backlog'.
  2. The feature was created via /design (not /feature) and skipped ws scaffold autogen.
     Run: bd show <bead-id> to find the right scope, then back-fill the scaffold
     under docs/workstreams/backlog/${WS_ID}.md before re-running /build.

Without a ws file, the guard cannot enforce scope and the verdict cannot be validated.
EOF
  exit 1
fi

# Read frontmatter status. Treat missing/empty as "open" (proceed).
status="$(awk '
  BEGIN { in_fm=0 }
  /^---$/ { in_fm=!in_fm; next }
  in_fm && /^status:/ {
    sub(/^status:[[:space:]]*/, "")
    print
    exit
  }
' "$WS_FILE")"

if [[ "$status" == "design-pending" ]]; then
  cat >&2 <<EOF
build-precheck: REFUSE — ${WS_FILE} declares status=design-pending.

Rule (F142-07): scaffolds with status=design-pending exist to occupy the
ws path so picker can skip them, but they are not /build-ready by author
intent. Run /design <feature-id> to produce a real workstream first.
EOF
  exit 1
fi

exit 0
