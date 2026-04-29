#!/usr/bin/env bash
# F141-04: pre-commit drift gate.
# Optional install via `sdp init` (F141-03) or manually:
#   ln -sf ../../scripts/hooks/sdp-doctor-precommit.sh .git/hooks/pre-commit
#
# Checks that .sdp/generated/ is in sync with sdp.manifest.yaml before every
# commit.  Exits 1 if drift is detected, so the commit is aborted.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

if ! go run ./cmd/sdp doctor adapters; then
  echo ""
  echo "❌  sdp doctor: adapter drift detected."
  echo "    Fix: run \`sdp generate-adapters --write --out .sdp/generated\`"
  echo "    then \`git add .sdp/generated\` and retry the commit."
  exit 1
fi
