#!/bin/bash
# Post-oneshot hook: build and test after /oneshot completion

set -e

FEATURE_ID="$1"

if [[ -z "$FEATURE_ID" ]]; then
  echo "Usage: post-oneshot.sh F{XX}"
  exit 1
fi

echo "🧪 Running post-oneshot checks for $FEATURE_ID..."

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -d "sdp-plugin" ]; then
  echo ""
  echo "=== Build & Test (sdp-plugin) ==="
  cd sdp-plugin
  go build ./...
  go test ./... -count=1 -short
  cd "$REPO_ROOT"
fi

echo ""
echo "✅ All post-oneshot checks passed for $FEATURE_ID"
