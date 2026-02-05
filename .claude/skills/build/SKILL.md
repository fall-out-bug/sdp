---
name: build
description: Execute workstream with TDD and guard enforcement
tools: Read, Write, StrReplace, Shell, Skill
---

# @build - Execute Workstream

Execute a single workstream following TDD discipline with automatic guard.

## Invocation (BEADS-001)

Accepts **both** formats:

- `@build 00-001-01` — WS-ID (PP-FFF-SS), resolve beads_id from `.beads-sdp-mapping.jsonl`
- `@build sdp-xxx` — Beads task ID directly

## Beads Integration (optional)

**When Beads is enabled** (bd installed, `.beads/` exists):

1. **Resolve ID:** ws_id → beads_id via mapping (if ws_id given)
2. **Before work:** `bd update {beads_id} --status in_progress`
3. **Execute:** TDD cycle
4. **On success:** `bd close {beads_id} --reason "WS completed" --suggest-next`
5. **On failure:** `bd update {beads_id} --status blocked`
6. **Before commit:** `bd sync`

**When Beads NOT enabled:** Skip all Beads steps. Use ws_id only.

**Detection:** Check if `bd --version` works and `.beads/` exists.

## Quick Reference

| Step | Action | Gate | Beads? |
|------|--------|------|--------|
| 0 | Detect Beads | Check `bd --version` + `.beads/` | Detection |
| 0a | Resolve beads_id | ws_id → mapping (if Beads) | Optional |
| 0b | Beads IN_PROGRESS | `bd update --status in_progress` (if Beads) | Optional |
| 1 | Activate guard | `sdp guard activate {ws_id}` | Always |
| 2 | Read WS spec | AC present and clear | Always |
| 3 | TDD cycle | `@tdd` for each AC | Always |
| 4 | Quality check | `sdp quality check` passes | Always |
| 5 | Beads CLOSED/blocked | `bd close` or `bd update --status blocked` (if Beads) | Optional |
| 6 | Beads sync + commit | `bd sync` then commit (if Beads) | Optional |
| 7 | Commit only | `git commit` (if no Beads) | Fallback |

## Workflow

### Step 0: Resolve Task ID

```bash
# Input: ws_id (00-001-01) OR beads_id (sdp-xxx)
# If ws_id: beads_id = grep mapping for sdp_id
# If beads_id: ws_id = grep mapping for beads_id (reverse lookup)
# Guard needs ws_id; Beads needs beads_id
beads_id=$(grep -m1 "\"sdp_id\": \"{WS-ID}\"" .beads-sdp-mapping.jsonl 2>/dev/null | grep -o '"beads_id": "[^"]*"' | cut -d'"' -f4)
ws_id=$(grep -m1 "\"beads_id\": \"{beads_id}\"" .beads-sdp-mapping.jsonl 2>/dev/null | grep -o '"sdp_id": "[^"]*"' | cut -d'"' -f4)
```

### Step 1: Beads IN_PROGRESS (when Beads enabled)

```bash
[ -n "$beads_id" ] && bd update "$beads_id" --status in_progress
```

### Step 2: Activate Guard

```bash
sdp guard activate {WS-ID}
```

**Gate:** Must succeed. If fails, WS not ready.

### Step 3: Read Workstream

```bash
Read("docs/workstreams/backlog/{WS-ID}-*.md")
```

Extract:
- Goal and Acceptance Criteria
- Input/Output files
- Steps to execute

### Step 4: TDD Cycle

For each AC, call internal TDD skill:

```
@tdd "AC1: {description}"
```

Cycle: Red → Green → Refactor

### Step 5: Quality Check

```bash
sdp quality check --module {module}
```

Must pass:
- Coverage ≥80%
- mypy --strict
- ruff (no errors)
- Files <200 LOC

### Step 6: Beads CLOSED or blocked

**On success:**
```bash
[ -n "$beads_id" ] && bd close "$beads_id" --reason "WS completed" --suggest-next
```

**On failure (quality check fails, TDD fails):**
```bash
[ -n "$beads_id" ] && bd update "$beads_id" --status blocked
```

### Step 7: Complete

```bash
# When Beads enabled: sync before commit
[ -d .beads ] && bd sync

sdp guard complete {WS-ID}
git add .
git commit -m "feat({scope}): {WS-ID} - {title}"
```

## Quality Gates

See [Quality Gates Reference](../../docs/reference/quality-gates.md)

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| No active WS | Guard not activated | `sdp guard activate` |
| File not in scope | Editing wrong file | Check WS scope |
| Coverage <80% | Missing tests | Add tests |

## See Also

- [BEADS-001 Phase 2.3](../../docs/workstreams/backlog/BEADS-001-skills-integration.md) — Beads @build spec
- [WorkstreamExecutor](../../src/sdp/beads/skills_build.py) — Python implementation
- [Full Build Spec](../../docs/reference/build-spec.md)
- [TDD Skill](../tdd/SKILL.md)
- [Guard Skill](../guard/SKILL.md)
