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

## See Also

- @debug — Root cause analysis
- @hotfix — P0
- @bugfix — P1/P2
