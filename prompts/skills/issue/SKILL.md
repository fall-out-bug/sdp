---
name: issue
description: Analyze bugs, classify severity (P0-P3), route to appropriate fix command (@hotfix, @bugfix, or backlog).
---

# @issue

Classify bugs and route to fix command.

## Severity → Route

| Severity | Signals | Route |
|----------|---------|-------|
| P0 | "production down", "crash", "blocked" | @hotfix |
| P1 | "doesn't work", "failing", "broken" | @bugfix |
| P2 | "edge case", "sometimes" | backlog |
| P3 | "cosmetic", "typo" | defer |

## Workflow

1. Document symptom, reproduction, environment
2. Form hypotheses, rank by likelihood
3. Test systematically
4. Classify per table above
5. `bd create --title="Bug: {desc}" --type=bug --priority={0-3}`
6. Route: @hotfix (P0), @bugfix (P1), or schedule WS (P2/P3)

## Output

Issue file, routing recommendation.

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @debug — Root cause analysis
- @hotfix — P0
- @bugfix — P1/P2
