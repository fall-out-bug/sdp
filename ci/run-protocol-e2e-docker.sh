#!/usr/bin/env bash
# Run protocol E2E in Docker (install + init + CLI + optional LLM).
# Usage: ./ci/run-protocol-e2e-docker.sh
# Set GLM_API_KEY for Phase 5 (opencode LLM code generation).

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "Building protocol-e2e image..."
docker build -f ci/Dockerfile.protocol-e2e \
  --build-arg SDP_PLUGIN_PATH=sdp-plugin \
  -t sdp-protocol-e2e:latest .

echo "Running protocol-e2e (GLM_API_KEY=${GLM_API_KEY:+set})..."
docker run --rm -e GLM_API_KEY="${GLM_API_KEY}" sdp-protocol-e2e:latest
