---
name: prototype
description: Rapid prototyping shortcut for experienced vibecoders
tools:
  - Read
  - Write
  - Bash
  - Glob
  - Grep
  - AskUserQuestion
version: 1.1.0
---

# @prototype - Rapid Prototyping Shortcut

Ultra-fast feature planning: 15-min interview → 1-3 workstreams → immediate execution with relaxed gates.

> **Speed over discipline.** Tech debt tracked for later cleanup.

## When to Use

- Experienced developers who know the codebase
- Need working prototype FAST (same day)
- Technical debt acceptable initially

**Don't use for:** Production features, team projects, security-critical code.

## Workflow

### Step 1: Ultra-Fast Interview (5 Questions)

```
AskUserQuestion:
1. Problem: User pain point | New capability | Technical debt
2. Scope: Backend only | Frontend only | Full stack
3. Dependencies: None | APIs | Database
4. Risks: None known | Technical uncertainty | Dependencies
5. Success: User can do X | Performance gain | Bug fix
```

### Step 2: Generate Workstreams

| Scope | Workstreams |
|-------|-------------|
| Backend only | 1 WS: Backend Implementation |
| Frontend only | 1 WS: Frontend Implementation |
| Full stack | 3 WS: Backend, Frontend, Integration |

### Step 3: Launch @oneshot

```python
Skill(skill="oneshot", args={"feature_id": feature_id, "mode": "prototype"})
```

## Quality Gate Overrides

| Gate | Normal | Prototype |
|------|--------|-----------|
| TDD | Required | Optional |
| Coverage | ≥80% | None |
| File Size | <200 LOC | No limit |
| Architecture | Clean | Monolithic OK |

**Non-Negotiable:**
- Code MUST compile and run
- No crashes
- Feature works end-to-end
- Basic security (no XSS, SQL injection)

## Tech Debt Tracking

All violations auto-tracked as Beads issues (priority 3):
- Files > 200 LOC
- Missing test coverage
- Architecture violations

## Command Flags

```bash
@prototype <description> [--feature=FFF] [--workstreams=N] [--skip-interview] [--immediate]
```

## Output

```
docs/drafts/prototype-{feature_id}.md     # Interview summary
docs/workstreams/backlog/00-FFF-*.md      # 1-3 workstreams
```

## Follow-up Paths

1. **Fix Tech Debt** → `@review F{feature_id}`
2. **Refactor Properly** → `@feature "{description}" --based-on=F{feature_id}`
3. **Discard** → Start over

## See Also

- `@feature` — Full feature planning with strict gates
- `@oneshot` — Autonomous execution
- `@build` — Single workstream execution
