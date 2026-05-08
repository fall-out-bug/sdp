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
| Sync | `scripts/beads_transport.sh export` |

## Integration Points

- **@build/@oneshot** — Check `bd ready` before WS, `bd update` after
- **@design** — `bd create` for new WS, `bd dep add` for dependencies
- **Mapping** — `.beads-sdp-mapping.jsonl` links WS ID to beads ID

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @build — Uses beads for dependency check
- @oneshot — Wave execution
- AGENTS.md — `bd ready`, `bd show`, `bd update`, `bd close`, `scripts/beads_transport.sh export`
