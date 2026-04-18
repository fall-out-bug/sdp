---
name: feature
description: Feature planning orchestrator (discovery -> idea -> ux -> design -> workstream tree)
version: 8.0.0
depends_on: "@discovery v1"
changes:
  - v8: Full product discovery flow with @discovery, @ux, impact analysis
  - Added --quick (skip @discovery), --infra (skip @ux)
  - Step 3.5: Impact analysis after @design
examples:
  - "@feature 'Add user authentication' --default # Full interactive pipeline"
  - "@feature 'Add user authentication' --quick # Skip to @design only, 0 questions"
  - "@feature 'Add payment processing' --auto # Non-interactive, from roadmap/plan docs"
---

# @feature

Orchestrate product discovery, requirements, UX research, and workstream design.

**Phase 0:** This skill targets Go projects (e.g. `go build`/`go test` in acceptance criteria). Language-agnostic expansion is planned.

## Three Explicit Modes (No Auto-Detection)

**Key Principle:** User must explicitly choose the mode. No hidden context detection. Same input + same mode = identical output.

| Mode | When to use | Steps | Questions |
|------|-------------|-------|-----------|
| `--default` (or no flag) | New/exploratory feature. Full interactive discovery. | 0, 1, 2, 2.5, 3, 3.5, 4 | Interactive (3-5 questions) |
| `--quick` | User knows what they want, just needs workstreams. | 3 only | **0 questions** - goes directly to @design |
| `--auto` | Feature already described in roadmap/plan. Non-interactive. | 0, 3, 4 only | **0 questions** - reads from docs/ROADMAP.md |

### Mode Behavior Guarantee

**Deterministic:** Each mode produces identical behavior given identical input. No context sniffing, no heuristics, no "smart defaults."

- `--default`: Always asks the same questions in the same order
- `--quick`: Always skips to @design with zero questions
- `--auto`: Always reads from roadmap/plan docs, never asks questions

---

## --auto Mode (Recommended for Roadmap Features)

For features already defined in `docs/ROADMAP.md`:

### Step A: Extract from Roadmap (Non-Interactive)

1. Find the feature in the roadmap: `rg "F0\d\d" docs/ROADMAP.md -A 10`
2. Extract: feature ID, description, success criteria, listed deliverables
3. Identify scope: what files/packages this touches (from deliverables and codebase)

**No questions asked.** Read from docs only.

### Step B: Auto-Generate Workstreams

For each deliverable, create workstream files using this format:

```markdown
---
ws_id: 00-FFF-SS
feature_id: FFFF
status: open
priority: P1
size: M
depends_on: []
ws_kind: leaf|aggregate
parent_ws_id: null|00-FFF-SS
dispatch_lifecycle: active
---

# 00-FFF-SS: Feature Name — Step Description

## Goal
One paragraph: what this workstream does and why.

## Scope Files
- path/to/file/or/dir (exact files or directory prefixes)
- ...

## Beads
- primary: sdplab-XXXX      # leaf only
- finding: sdplab-YYYY      # optional

## Acceptance Criteria
- [ ] Specific, testable criterion 1
- [ ] Specific, testable criterion 2
- [ ] go build ./... passes
- [ ] go test ./... passes
```

**Rules:**
- `aggregate` = container, no `primary` issue
- `leaf` = executable, may have one `primary` issue
- `parent_ws_id` links leaf to aggregate (max 1 nesting layer)

### Step C: Create Beads Issues

For each leaf workstream:
```bash
bd create --title="WS FFF-SS: Short title" --type=task
```

Update `.beads-sdp-mapping.jsonl`:
```json
{"sdp_id":"00-FFF-SS","beads_id":"sdp_dev-XXXX","updated_at":"2026-..."}
```

### Step D: Report

Output: feature ID, workstream count, file names, Beads IDs, ready-to-run command (`@build 00-FFF-01` or `@oneshot F0FF`).

---

## --quick Mode (@design Only, Zero Questions)

For users who know what they want and just need workstreams. Zero questions, deterministic.

**Steps:** Run @design directly, produce workstream files, skip roadmap/UX/impact analysis.

**When to use:** Clear feature description, no product research needed, immediate workstreams required.

---

## --default Mode (Full Interactive Pipeline)

Standard mode for new/exploratory features. Full discovery with interactive questions.

### Step 0: Roadmap Pre-Check

Use `@discovery "feature description"` or manually: extract keywords, `rg` docs for overlap, present Overlap Report (HIGH/MEDIUM), gate on user resolution.

### Step 1: Quick Interview (3-5 questions)

Problem, Users, Success. Gate: if vague (<200 words), ask clarification. Use @discovery output if available.

### Step 2: @idea

`@idea "..." --spec docs/drafts/discovery-{slug}.md` (use @discovery output if Step 0 ran)

### Step 2.5: @ux — unless --infra

Auto-trigger when @idea output has user-facing keywords (ui, user, interface, dashboard, form) and lacks infra (K8s, CRD, CLI-only).

### Step 3: @design

`@design {task_id}` — workstream files in docs/workstreams/backlog/ using **Workstream file format** above.

### Step 3.5: Impact Analysis

Read scope files, grep/rg for conflicts. Categorize: FILE CONFLICT, DATA BOUNDARY, DEPENDENCY CHAIN, PRIORITY SHIFT. Present report, user acknowledges.

### Step 4: Verify Outputs

Check discovery brief, idea spec, ux output, workstreams exist, direct execution targets are leaf workstreams.

---

## Key Principle: Protocol is Invisible

User sees: feature description → workstreams created → ready to build. Workstream files, scope declarations, beads IDs are plumbing.

## Completion

When all workstreams are created and verified, output:

```
@feature complete. Feature {ID}: {count} workstreams created.
  Aggregate: 00-{FFF}-00
  Leaves: 00-{FFF}-01 .. 00-{FFF}-{NN}

Next: @build 00-{FFF}-01  or  @oneshot F{XX}
```

## See Also

@discovery — Product discovery gate | @idea — Requirements | @ux — UX research | @design — Workstream planning | @build — Execute leaf workstream | @oneshot — Execute all ready leaf workstreams
