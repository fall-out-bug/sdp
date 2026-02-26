#!/bin/sh
# Pre-push hook: go test -short, evidence validation for feature branches when internal/ or cmd/ changed.
# CWD = repo root. Exit 1 on any failure.
set -e

# 1. go test -short ./...
go test -short ./... || { echo "pre-push: go test -short failed" >&2; exit 1; }

# 2. If feature branch + diff touches internal/ or cmd/: require .sdp/evidence/*.json, validate
BRANCH=$(git branch --show-current)
case "$BRANCH" in
  feature/*) ;;
  *) exit 0 ;;  # Not a feature branch, skip evidence check
esac

# Collect all changed files from refs being pushed (stdin: local_ref local_sha remote_ref remote_sha)
CHANGED=""
while read local_ref local_sha remote_ref remote_sha; do
  [ -z "$local_sha" ] && continue
  if [ "$remote_sha" = "0000000000000000000000000000000000000000" ]; then
    CHANGED="$CHANGED $(git diff --name-only 4b825dc642cb6eb9a060e54bf8d69288fbee4904 $local_sha 2>/dev/null || true)"
  else
    CHANGED="$CHANGED $(git diff --name-only $remote_sha $local_sha 2>/dev/null || true)"
  fi
done

FEATURE_PATHS=$(echo "$CHANGED" | tr ' ' '\n' | grep -E '^internal/|^cmd/' || true)
EVIDENCE_FILES=$(echo "$CHANGED" | tr ' ' '\n' | grep '^\.sdp/evidence/.*\.json$' || true)

if [ -n "$FEATURE_PATHS" ] && [ -z "$EVIDENCE_FILES" ]; then
  echo "pre-push: feature branch with internal/ or cmd/ changes requires .sdp/evidence/*.json in push" >&2
  echo "Add evidence (see @build skill step 3b)" >&2
  exit 1
fi

if [ -z "$EVIDENCE_FILES" ]; then
  exit 0
fi

# 3. Validate evidence; if sdp-evidence not built: skip, warn
EVIDENCE_CMD=""
if command -v sdp-evidence >/dev/null 2>&1; then
  EVIDENCE_CMD="sdp-evidence"
elif [ -f bin/sdp-evidence ] && [ -x bin/sdp-evidence ]; then
  EVIDENCE_CMD="./bin/sdp-evidence"
else
  echo "pre-push: sdp-evidence not found, skipping evidence validation (run: make build-sdp-evidence)" >&2
  exit 0
fi

for f in $EVIDENCE_FILES; do
  if [ -f "$f" ]; then
    $EVIDENCE_CMD validate --require-pr-url=false "$f" || { echo "pre-push: evidence validation failed: $f" >&2; exit 1; }
  fi
done

exit 0
