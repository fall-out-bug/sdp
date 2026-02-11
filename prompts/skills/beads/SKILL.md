---
name: beads
description: Beads task tracker integration for SDP workflows.
---

# Beads Integration

Unified interface for working with Beads task tracker from SDP workflows.

## When to Use

Use this skill when:
- Creating new workstreams (auto-register with Beads)
- Checking task status and dependencies
- Updating task status after completion
- Syncing workstream state with Beads

## Quick Reference

| Action | Command | Beads Command |
|--------|---------|---------------|
| Create task | `bd create` | `bd create -w <ws-id>` |
| Check ready | `bd ready` | `bd ready --json` |
| Add dependency | `bd dep add` | `bd dep add <task> <depends-on>` |
| Update status | `bd update` | `bd update <ws-id> --status completed` |
| Show task | `bd show` | `bd show <task-id>` |
| Sync mapping | Auto | Updates `.beads-sdp-mapping.jsonl` |

## SDP ↔ Beads Mapping

SDP maintains a mapping file `.beads-sdp-mapping.jsonl`:

```json
{"sdp_id": "00-050-01", "beads_id": "sdp-x8p", "updated_at": "2026-02-05T12:04:08.943705"}
```

**Format:**
- `sdp_id`: Workstream ID (PP-FFF-SS format)
- `beads_id`: Beads task ID (auto-generated)
- `updated_at`: Last sync timestamp

## Integration Points

### 1. Workstream Creation (@design, @feature)

When creating new workstreams, automatically register with Beads:

```bash
# After creating WS file
bd create \
  --title "WS-00-050-01: Go Project Setup" \
  --description "Foundation for Go migration" \
  --status backlog \
  --metadata '{"ws_id": "00-050-01", "feature": "F050", "size": "MEDIUM"}'

# Update mapping
echo '{"sdp_id": "00-050-01", "beads_id": "'$(bd last --id)'", "updated_at": "'$(date -u +%Y-%m-%dT%H:%M:%S.%N)'"}' \
  >> .beads-sdp-mapping.jsonl
```

### 2. Dependency Setup (@design)

After creating all workstreams, set up dependencies:

```bash
# For each dependency in workstream
bd dep add sdp-gtw sdp-x8p  # 00-050-02 depends on 00-050-01
bd dep add sdp-o8h sdp-x8p  # 00-050-03 depends on 00-050-01
bd dep add sdp-645 "sdp-x8p,sdp-gtw,sdp-o8h"  # 00-050-04 depends on 1,2,3
```

**Workflow:**
1. Parse `depends_on` from workstream frontmatter
2. Map WS IDs to Beads IDs using `.beads-sdp-mapping.jsonl`
3. Execute `bd dep add` for each dependency

### 3. Ready Check (@build, @oneshot)

Before executing workstream, check if dependencies are satisfied:

```bash
# Get ready tasks
bd ready --json

# Parse output
ready_tasks=$(bd ready --json | jq -r '.[].id')

# Check if current WS is ready
if echo "$ready_tasks" | grep -q "sdp-x8p"; then
  echo "✅ WS-00-050-01 is ready to execute"
else
  echo "❌ WS-00-050-01 is blocked by dependencies"
  exit 1
fi
```

### 4. Status Update (@build completion)

After workstream completes, update Beads status:

```bash
# Success
bd update sdp-x8p --status completed

# Failure
bd update sdp-x8p --status failed --notes "Coverage too low: 75% < 80%"

# In Progress (for @oneshot)
bd update sdp-x8p --status in-progress
```

**Auto-trigger:** Add to post-build hook in `hooks/post-build.sh`

### 5. Sync State (continuous)

Keep Beads and SDP in sync:

```bash
# Sync mapping file
bd sync

# Validate mapping
python scripts/validate_beads_mapping.py
```

## Skill Integration

### @build Integration

```markdown
User: @build 00-050-01

Claude:
→ Mapping WS ID to Beads ID: 00-050-01 → sdp-x8p
→ Checking dependencies: ✅ All satisfied
→ Updating status: sdp-x8p → in-progress
→ [Execute workstream...]
→ Workstream complete
→ Updating status: sdp-x8p → completed
```

### @oneshot Integration

```markdown
User: @oneshot F050

Claude:
→ Loading feature workstreams...
→ Mapping WS IDs to Beads IDs: 13 workstreams
→ Querying bd ready: [sdp-x8p, sdp-gtw, sdp-o8h]
→ Wave 1: Executing 3 ready tasks in parallel
→ [Execute...]
→ Updating Beads status: sdp-x8p → completed
→ [Continue with next wave...]
```

### @design Integration

```markdown
User: @design idea-f050

Claude:
→ Creating 13 workstreams...
→ Registering with Beads: bd create --title "WS-00-050-01"...
→ Mapping: 00-050-01 → sdp-x8p
→ Setting up dependencies: bd dep add sdp-gtw sdp-x8p...
→ Verifying dependency graph: ✅ No cycles
→ Migration command: poetry run sdp beads migrate docs/workstreams/backlog/ --real
```

## Helper Commands

### Check WS Status

```bash
# Check if workstream is ready
bd-ready() {
  local ws_id=$1
  local beads_id=$(grep "\"sdp_id\": \"$ws_id\"" .beads-sdp-mapping.jsonl | jq -r '.beads_id')
  bd ready --json | jq -r ".[] | select(.id == \"$beads_id\")"
}

# Usage
bd-ready "00-050-01"
# Output: {"id": "sdp-x8p", "title": "WS-00-050-01", ...}
```

### List Blocked Workstreams

```bash
# Show which workstreams are blocked
bd-blocked() {
  local ws_id=$1
  local beads_id=$(grep "\"sdp_id\": \"$ws_id\"" .beads-sdp-mapping.jsonl | jq -r '.beads_id')
  bd show "$beads_id" --json | jq -r '.blocking[]'
}

# Usage
bd-blocked "00-050-13"
# Output: ["sdp-x8p", "sdp-gtw", "sdp-o8h", ...]
```

### Update Multiple Tasks

```bash
# Batch update workstreams
bd-batch-update() {
  local status=$1
  shift
  local ws_ids=("$@")

  for ws_id in "${ws_ids[@]}"; do
    local beads_id=$(grep "\"sdp_id\": \"$ws_id\"" .beads-sdp-mapping.jsonl | jq -r '.beads_id')
    echo "Updating $ws_id ($beads_id) → $status"
    bd update "$beads_id" --status "$status"
  done
}

# Usage
bd-batch-update completed "00-050-01" "00-050-02" "00-050-03"
```

## Error Handling

### Missing Mapping

```bash
# If WS ID not found in mapping
if ! grep -q "\"sdp_id\": \"$ws_id\"" .beads-sdp-mapping.jsonl; then
  echo "❌ Error: Workstream $ws_id not registered with Beads"
  echo "Run: bd create --title \"WS-$ws_id\""
  exit 1
fi
```

### Beads Command Failed

```bash
# Wrap Beads commands with error handling
bd-safe() {
  if ! output=$(bd "$@" 2>&1); then
    echo "❌ Beads command failed: bd $*"
    echo "Error: $output"
    return 1
  fi
  echo "$output"
}

# Usage
bd-safe show "sdp-x8p"
```

### Sync Conflicts

```bash
# Resolve conflicts between SDP and Beads
bd-resolve() {
  local ws_id=$1
  echo "Resolving conflict for $ws_id..."

  # Check SDP status (from frontmatter)
  local sdp_status=$(grep "^status:" "docs/workstreams/backlog/$ws_id.md" | awk '{print $2}')

  # Check Beads status
  local beads_id=$(grep "\"sdp_id\": \"$ws_id\"" .beads-sdp-mapping.jsonl | jq -r '.beads_id')
  local beads_status=$(bd show "$beads_id" --json | jq -r '.status')

  if [ "$sdp_status" != "$beads_status" ]; then
    echo "⚠️ Status mismatch: SDP=$sdp_status, Beads=$beads_status"
    echo "Updating Beads to match SDP..."
    bd update "$beads_id" --status "$sdp_status"
  fi
}
```

## Best Practices

1. **Always update mapping** after creating workstreams
2. **Check dependencies** before execution (@build, @oneshot)
3. **Update status** after completion (post-build hook)
4. **Sync regularly** to prevent drift
5. **Validate mapping** before major operations

## Migration Guide

### Existing Workstreams → Beads

```bash
# Migrate all backlog workstreams
poetry run sdp beads migrate docs/workstreams/backlog/ --real

# Validate migration
cat .beads-sdp-mapping.jsonl | wc -l  # Should match workstream count
```

### Manual Registration

```bash
# For individual workstreams
bd-register() {
  local ws_id=$1
  local ws_file="docs/workstreams/backlog/$ws_id.md"

  # Extract metadata
  local title=$(grep "^## WS-$ws_id" "$ws_file" | sed 's/## //')
  local status=$(grep "^status:" "$ws_file" | awk '{print $2}')
  local size=$(grep "^size:" "$ws_file" | awk '{print $2}')

  # Create Beads task
  local beads_id=$(bd create \
    --title "$title" \
    --status "$status" \
    --metadata "{\"ws_id\": \"$ws_id\", \"size\": \"$size\"}" \
    --output-id)

  # Update mapping
  echo "{\"sdp_id\": \"$ws_id\", \"beads_id\": \"$beads_id\", \"updated_at\": \"$(date -u +%Y-%m-%dT%H:%M:%S.%N)\"}" \
    >> .beads-sdp-mapping.jsonl

  echo "✅ Registered $ws_id → $beads_id"
}
```

## Related Skills

- `/build` - Uses Beads for dependency checking
- `/oneshot` - Uses Beads for wave execution
- `/design` - Registers workstreams with Beads
- `/verify-workstream` - Checks Beads status before execution

## Quick Reference Card

```
┌──────────────────────────────────────────┐
│     BEADS INTEGRATION: KEY COMMANDS      │
├──────────────────────────────────────────┤
│  bd ready → Check available tasks        │
│  bd create → Register new workstream     │
│  bd dep add → Setup dependencies         │
│  bd update → Update task status          │
│  Mapping → .beads-sdp-mapping.jsonl      │
└──────────────────────────────────────────┘
```

**Remember:** Beads is the source of truth for task status. SDP syncs with Beads, not vice versa.
