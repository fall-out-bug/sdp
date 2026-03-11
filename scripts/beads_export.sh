#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

if bd sync --help >/dev/null 2>&1; then
  exec bd sync
fi

# bd >= 0.59 stores issue state in Dolt. Export a repo snapshot so existing
# git-backed workflows continue to publish .beads/issues.jsonl.
bd export -o .beads/issues.jsonl
