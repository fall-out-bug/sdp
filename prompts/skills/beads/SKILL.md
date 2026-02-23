---
name: beads
description: Beads task tracker integration for SDP workflows.
---

# @beads

Beads integration for SDP. Mapping: `.beads-sdp-mapping.jsonl` (sdp_id → beads_id).

## Quick Reference

| Action | Command |
|--------|---------|
| Ready tasks | `bd ready` |
| Show task | `bd show <id>` |
| Update status | `bd update <id> --status completed` |
| Create | `bd create --title="..." --type=task` |
| Dependencies | `bd dep add <task> <depends-on>` |
| Sync | `bd sync` |

## Integration Points

- **@build/@oneshot** — Check `bd ready` before WS, `bd update` after
- **@design** — `bd create` for new WS, `bd dep add` for dependencies
- **Mapping** — `.beads-sdp-mapping.jsonl` links WS ID to beads ID

## See Also

- @build — Uses beads for dependency check
- @oneshot — Wave execution
- AGENTS.md — `bd ready`, `bd show`, `bd update`, `bd close`, `bd sync`
