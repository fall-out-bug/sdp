#!/usr/bin/env bash
# Test SDP install in Docker (clean env, OpenCode layout)
# Usage: ./ci/test-install-docker.sh [SDP_REF]
#   SDP_REF defaults to schema/coding-workflow-predicate for PR testing

set -e

SDP_REF="${1:-schema/coding-workflow-predicate}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== SDP Install Docker Test ==="
echo "SDP_REF: $SDP_REF"
echo ""

docker build -f "$REPO_ROOT/ci/Dockerfile.install-test" \
  --build-arg SDP_REF="$SDP_REF" \
  -t sdp-install-test:latest \
  "$REPO_ROOT"

echo ""
echo "=== Running install verification ==="
docker run --rm sdp-install-test:latest

echo ""
echo "✅ Docker install test passed"
