---
name: protocol-consistency
description: Audit consistency across workstream docs, CLI capabilities, and CI workflows.
---

# @protocol-consistency

Detect drift between docs, CLI, and CI.

## Workflow

1. **Verify CLI** — `sdp --help`, `sdp <cmd> --help` — commands in docs exist
2. **Validate WS schema** — Read `docs/workstreams/backlog/<ws-id>.md`, run `sdp doctor adapters` to check for drift
3. **Validate CI** — `rg "sdp .*" .github/workflows hooks scripts` — paths valid
4. **Report** — Source file, observed vs expected, risk, suggested fix
5. **Track** — `bd create --title="Protocol drift: ..." --type=task --priority=2`

## Output

Report: scope, blocking/non-blocking mismatches, findings, recommended fixes.

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |
