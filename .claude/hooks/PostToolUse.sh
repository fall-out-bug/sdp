#!/bin/bash
# PostToolUse hook - runs after any tool use.
#
# Claude Code passes tool-use metadata via argv; we detect git-commit calls
# and sync the beads transport so the dolt mirror stays fresh.
set -euo pipefail

# Resolve repo root (two levels up from .claude/hooks/) so the hook works
# regardless of the caller's cwd.
HOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${HOOK_DIR}/../.." && pwd)"

LOG_FILE="${REPO_ROOT}/.claude/hooks/PostToolUse.log"
mkdir -p "$(dirname "${LOG_FILE}")"

{
  echo "🔔 PostToolUse hook triggered"
  date
} >> "${LOG_FILE}"

# Sync beads after git commit (use $* so word-array joins cleanly for glob match).
if [[ "$*" == *"git commit"* ]]; then
  echo "📦 Syncing Beads after commit..." >> "${LOG_FILE}"
  "${REPO_ROOT}/scripts/beads_transport.sh" export >> "${LOG_FILE}" 2>&1 || true
fi
