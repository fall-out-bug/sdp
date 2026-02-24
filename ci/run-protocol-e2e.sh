#!/usr/bin/env bash
# Local wrapper: docker build + docker run for protocol E2E
# Usage: GLM_API_KEY=... ./ci/run-protocol-e2e.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Protocol E2E (Docker) ==="
docker build -f "$REPO_ROOT/ci/Dockerfile.protocol-e2e" \
  -t sdp-protocol-e2e:latest "$REPO_ROOT"

echo ""
echo "=== Running protocol E2E test ==="
docker run --rm \
  -e GLM_API_KEY="${GLM_API_KEY:-}" \
  sdp-protocol-e2e:latest

echo ""
echo "Protocol E2E passed"
