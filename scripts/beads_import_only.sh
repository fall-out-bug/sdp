#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

if bd sync --help >/dev/null 2>&1; then
  exec bd sync --import-only
fi

# bd >= 0.59 removed `bd sync`; rebuild the Dolt-backed database from the
# tracked JSONL snapshot so repo-based workflows still have deterministic state.
bd init --from-jsonl --force --quiet --database beads --prefix sdp
