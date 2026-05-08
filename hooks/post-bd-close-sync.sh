#!/usr/bin/env bash
# post-bd-close-sync.sh — Auto-sync workstream status after bd close.
#
# Moves workstream files from backlog/ to done/ and updates INDEX.md status
# entries when a beads issue is closed.
#
# Usage:
#   scripts/hooks/post-bd-close-sync.sh <issue-id> [<issue-id> ...]
#   scripts/hooks/post-bd-close-sync.sh sdplab-nai
#   scripts/hooks/post-bd-close-sync.sh sdplab-nai sdplab-abc
#
# Environment:
#   BD_POST_CLOSE_DRY_RUN=1   List intended changes without applying
#   BD_POST_CLOSE_QUIET=1     Suppress non-error output
#
# Brownfield-safe: exits 0 silently if workstream files or directories don't exist.
set -euo pipefail

# --- Paths ---
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
WS_DIR="${PROJECT_ROOT}/docs/workstreams"
BACKLOG_DIR="${WS_DIR}/backlog"
DONE_DIR="${WS_DIR}/done"
INDEX_FILE="${WS_DIR}/INDEX.md"

# --- Helpers ---
log() {
  if [[ -z "${BD_POST_CLOSE_QUIET:-}" ]]; then
    printf '[post-bd-close-sync] %s\n' "$*"
  fi
}

log_dry() {
  printf '[post-bd-close-sync] DRY RUN: %s\n' "$*"
}

# Extract workstream ID from a bd issue.
# Strategy:
#   1. Try bd show --json to get the title (e.g. "F124-05: ...")
#   2. Parse title for pattern like F124-05 or F124-5 (feature-step)
#   3. Convert to ws_id format 00-124-05
#   4. Fallback: scan backlog/*.md frontmatter for a "## Beads" section mentioning the issue id
resolve_ws_id() {
  local issue_id="$1"
  local ws_id=""

  # Strategy 1: Extract from bd issue title via --json
  if command -v bd >/dev/null 2>&1; then
    local title=""
    title="$(bd show "$issue_id" --json 2>/dev/null | grep -o '"title":"[^"]*"' | head -1 | sed 's/"title":"//;s/"//')" || true
    if [[ -n "$title" ]]; then
      # Match patterns like F124-05 or F124-5 (with or without zero-padding)
      local feature_step=""
      feature_step="$(printf '%s' "$title" | grep -oE 'F[0-9]+-[0-9]+' | head -1)" || true
      if [[ -n "$feature_step" ]]; then
        # Parse F<feature>-<step> into 00-FFF-SS
        local feature_num step_num
        feature_num="$(printf '%s' "$feature_step" | sed 's/F\([0-9]*\)-.*/\1/')"
        step_num="$(printf '%s' "$feature_step" | sed 's/F[0-9]*-\([0-9]*\)/\1/')"
        # Remove leading zeros for printf, then re-pad
        feature_num="$((10#$feature_num))"
        step_num="$((10#$step_num))"
        ws_id="$(printf '00-%03d-%02d' "$feature_num" "$step_num")"
      fi
    fi
  fi

  # Strategy 2: Scan backlog frontmatter for the issue id in ## Beads section
  if [[ -z "$ws_id" && -d "$BACKLOG_DIR" ]]; then
    local candidate=""
    for candidate in "$BACKLOG_DIR"/*.md; do
      [[ -f "$candidate" ]] || continue
      # Look for the beads issue ID in the Beads section of the workstream
      # Match "- <issue_id>" or "- <issue_id>: ..."
      if grep -q "^- ${issue_id}\(:\|$\)" "$candidate" 2>/dev/null; then
        ws_id="$(basename "$candidate" .md)"
        break
      fi
    done
  fi

  printf '%s' "$ws_id"
}

# Move a workstream file from backlog/ to done/.
# Returns 0 on success or if already done, 1 on unexpected error.
move_workstream() {
  local ws_id="$1"
  local src="${BACKLOG_DIR}/${ws_id}.md"
  local dst="${DONE_DIR}/${ws_id}.md"

  # Already in done/ — idempotent success
  if [[ -f "$dst" ]]; then
    log "already in done/: ${ws_id}.md"
    return 0
  fi

  # Not in backlog — nothing to do (brownfield-safe skip)
  if [[ ! -f "$src" ]]; then
    log "no backlog file for ${ws_id} (skipping)"
    return 0
  fi

  if [[ -n "${BD_POST_CLOSE_DRY_RUN:-}" ]]; then
    log_dry "would move ${src} -> ${dst}"
    return 0
  fi

  # Ensure done/ directory exists
  mkdir -p "$DONE_DIR"

  # Update the status field in frontmatter before moving
  # Match: open, in_progress, in-progress, backlog → done
  if command -v sed >/dev/null 2>&1; then
    sed -i.bak -E 's/^status: (open|in_progress|in-progress|backlog)$/status: done/' "$src" 2>/dev/null && rm -f "${src}.bak" || true
  fi

  mv "$src" "$dst"
  log "moved ${ws_id}.md -> done/"
}

# Update the Workstream Status section in INDEX.md for a given ws_id.
# Changes the status column from Backlog/In Progress to Done.
update_index_status() {
  local ws_id="$1"

  if [[ ! -f "$INDEX_FILE" ]]; then
    log "INDEX.md not found (skipping status update)"
    return 0
  fi

  # Find the row with this ws_id in the status tables and change status to Done
  # The ws_id appears as the first column like | 00-124-05 |
  local pattern="^| ${ws_id} |"
  if grep -q "$pattern" "$INDEX_FILE"; then
    if [[ -n "${BD_POST_CLOSE_DRY_RUN:-}" ]]; then
      log_dry "would update INDEX.md status for ${ws_id} -> Done"
      return 0
    fi

    # Replace Backlog or In Progress with Done on the matching row
    # Using | as delimiter since rows contain pipes; escape carefully
    sed -i.bak "/${pattern}/s/| Backlog |/| Done |/;/${pattern}/s/| In Progress |/| Done |/" "$INDEX_FILE" 2>/dev/null && rm -f "${INDEX_FILE}.bak" || true
    log "updated INDEX.md: ${ws_id} -> Done"
  else
    log "${ws_id} not found in INDEX.md status table (skipping)"
  fi
}

# Update the Features table (top of INDEX.md) status column for the feature.
# When all workstreams for a feature are done, flip the feature status.
update_feature_status() {
  local ws_id="$1"

  if [[ ! -f "$INDEX_FILE" ]]; then
    return 0
  fi

  # Extract feature number from ws_id (00-124-05 -> F124)
  local feature_num
  feature_num="$(printf '%s' "$ws_id" | sed 's/00-\([0-9]*\)-.*/\1/')"
  local feature_id="F${feature_num}"

  # Check if all backlog workstreams for this feature are done
  local backlog_count=0
  if [[ -d "$BACKLOG_DIR" ]]; then
    backlog_count="$(find "$BACKLOG_DIR" -name "00-${feature_num}-*.md" -type f 2>/dev/null | wc -l | tr -d ' ')"
  fi

  if [[ "$backlog_count" -eq 0 ]]; then
    # All workstreams done — check if feature row exists in INDEX.md
    local feature_pattern="| \\*\\*${feature_id}\\*\\* |"
    if grep -q "$feature_pattern" "$INDEX_FILE"; then
      if [[ -n "${BD_POST_CLOSE_DRY_RUN:-}" ]]; then
        log_dry "would update INDEX.md feature ${feature_id} -> Done"
        return 0
      fi

      sed -i.bak "/${feature_pattern}/s/| Backlog |/| Done |/;/${feature_pattern}/s/| In Progress |/| Done |/" "$INDEX_FILE" 2>/dev/null && rm -f "${INDEX_FILE}.bak" || true
      log "updated INDEX.md feature ${feature_id} -> Done (all workstreams complete)"
    fi
  fi
}

# --- Main ---
main() {
  if [[ $# -eq 0 ]]; then
    echo 'Usage: post-bd-close-sync.sh <issue-id> [<issue-id> ...]' >&2
    exit 2
  fi

  # Brownfield guard: if workstream directory doesn't exist, exit cleanly
  if [[ ! -d "$WS_DIR" ]]; then
    log "docs/workstreams/ not found (skipping)"
    exit 0
  fi

  local moved=0
  local skipped=0

  for issue_id in "$@"; do
    log "processing issue: ${issue_id}"

    local ws_id
    ws_id="$(resolve_ws_id "$issue_id")"

    if [[ -z "$ws_id" ]]; then
      log "no workstream found for issue ${issue_id} (skipping)"
      ((skipped++)) || true
      continue
    fi

    log "resolved ${issue_id} -> ${ws_id}"

    move_workstream "$ws_id"
    update_index_status "$ws_id"
    update_feature_status "$ws_id"
    ((moved++)) || true
  done

  log "done: ${moved} moved, ${skipped} skipped"
}

main "$@"
