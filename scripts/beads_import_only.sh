#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

if bd sync --help >/dev/null 2>&1; then
  exec bd sync --import-only
fi

bd init --from-jsonl --force --quiet --database beads --prefix sdp
