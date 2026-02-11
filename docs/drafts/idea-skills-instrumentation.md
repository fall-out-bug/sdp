# F056: Full Skills Instrumentation

> Beads: sdp-2mkm | Priority: P1

---

## Problem

F054 instruments only `@build`. The remaining 18 skills produce no evidence. When `@review` finds a security issue, when `@design` makes a decomposition decision, when `@idea` gathers requirements — none of this is recorded.

Full pipeline coverage means every AI action has a trace.

## Solution

Instrument all remaining skills to emit evidence events via the evidence writer from F054.

## Skill → Event Mapping

| Skill | Event Type | Key Data |
|-------|-----------|----------|
| `@review` | `verification` | findings, severity, actual review output, model used |
| `@deploy` | `approval` | merge approval, gates passed, who approved, target branch |
| `@design` | `plan` | decomposition decisions, dependency graph, WS count |
| `@idea` | `plan` | questions asked, answers received (drive mode), requirements captured |
| `@vision` | `plan` | strategic decisions, expert agent outputs |
| `@reality` | `verification` | codebase analysis results, health score, gaps found |
| `@oneshot` | `plan` + `approval` | execution plan, checkpoints, completion status |
| `@prototype` | `generation` | code generated, ship mode decisions |
| `@hotfix` | `generation` + `approval` | emergency fix, fast-track approval |
| `@bugfix` | `generation` + `verification` | fix applied, tests run |
| `@issue` | `plan` | severity classification, routing decision |
| `@debug` | `verification` | hypotheses tested, root cause identified |

## Constraints

- Each skill = 1 workstream (instrument independently)
- Non-blocking: evidence failure never breaks the skill
- Use `evidence.Emitter` from F054 (no new infrastructure)
- Skills are Claude skills (SKILL.md), not Go code — instrumentation happens in CLI commands they call

## Users

- Same as F054 — developers, enterprise, ops
- Added: reviewers who want to know what `@review` actually checked

## Success Metrics

- All 19 skills emit at least 1 event type
- `sdp log show` shows events from every skill type
- Full chain: `@idea` → `@design` → `@build` → `@review` → `@deploy` = complete evidence trail

## Dependencies

- F054 (evidence writer, schema, emitter)

## Notes

- Highest leverage: `@review` and `@deploy` (complete the generation → verification → approval chain)
- Lowest effort: skills that already call CLI commands (just add evidence emit to the command)
- Highest effort: skills that are pure Claude skills with no CLI commands
