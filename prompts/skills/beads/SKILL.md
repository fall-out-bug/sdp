---
name: beads
description: Beads task tracker integration for SDP workflows.
---

# @beads

Beads integration for SDP. Mapping: `.beads-sdp-mapping.jsonl` (sdp_id → beads_id).

## Quick Reference

| Action | Command |
|--------|---------|
| Ready tasks | `sdp beads ready` |
| Show task | `sdp beads show <id>` |
| Update status | `sdp beads update <id> --status completed` |
| Create | `sdp beads create --title="..." --type=task` |
| Close | `sdp beads close <id> --reason "..."` |
| Dependencies | `bd dep add <task> <depends-on>` |
| Sync repo snapshot | `sdp beads sync` |

## Integration Points

- **@build/@oneshot** — Check `sdp beads ready` before WS, `sdp beads update/close` after
- **@design** — `sdp beads create` for new WS, `bd dep add` for dependencies
- **Mapping** — `.beads-sdp-mapping.jsonl` links WS ID to beads ID

## See Also

- @build — Uses beads for dependency check
- @oneshot — Wave execution
- AGENTS.md — `bd ready`, `bd show`, `bd update`, `bd close`, `./scripts/beads_export.sh`

`sdp beads create`, `update`, and `close` persist the repo snapshot automatically. Prefer them over raw `bd` mutations in agent workflows.
