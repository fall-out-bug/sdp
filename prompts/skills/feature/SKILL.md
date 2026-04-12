---
name: feature
description: Feature planning orchestrator (discovery -> idea -> ux -> design -> workstream tree)
version: 8.0.0
depends_on: "@discovery v1"
changes:
  - v8: Full product discovery flow with @discovery, @ux, impact analysis
  - Added --quick (skip @discovery), --infra (skip @ux)
  - Step 3.5: Impact analysis after @design
---

# @feature

Orchestrate product discovery, requirements, UX research, and workstream design.

**Phase 0:** This skill targets Go projects (e.g. `go build`/`go test` in acceptance criteria). Language-agnostic expansion is planned.

## Modes

| Mode | When to use | Steps |
|------|-------------|-------|
| `--auto` | Feature already described in roadmap/plan. Generate workstreams directly. | 0, 3, 4 only |
| `--quick` | User knows what they want. Skip roadmap pre-check. | 1, 2, 3, 4 |
| Default | New/exploratory feature. Full discovery. | 0, 1, 2, 2.5, 3, 3.5, 4 |

---

## --auto Mode (Recommended for Roadmap Features)

For features already defined in `docs/roadmap/ROADMAP.md` or `docs/workstreams/INDEX.md`:

### Step A: Extract from Roadmap

1. Find the feature in the roadmap: `rg "F0\d\d" docs/roadmap/ROADMAP.md -A 10`
2. Extract: feature ID, description, success criteria, listed deliverables
3. Identify scope: what files/packages this touches (from deliverables and codebase)

### Step B: Auto-Generate Workstreams

For each deliverable in the feature, create a workstream file:

```
docs/workstreams/backlog/00-FFF-SS.md
```

Use one of two shapes:

- **Leaf workstream** — directly executable contract slice
- **Aggregate workstream** — non-executable container or roll-up over `2+` leaf workstreams

Only leaf workstreams are direct `@build` targets.

**Workstream file format:**

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

- path/to/file/or/dir (exact files or directory prefixes this WS touches)
- ...

## Beads

- primary: sdplab-XXXX      # leaf only
- finding: sdplab-YYYY      # optional on leaf or aggregate

## Acceptance Criteria

- [ ] Specific, testable criterion 1
- [ ] Specific, testable criterion 2
- [ ] go build ./... passes (Phase 0: Go; other languages later)
- [ ] go test ./... passes
```

Rules:

- `aggregate` must not have a `primary` Beads issue
- `leaf` may have one open `primary`
- use `parent_ws_id` only when a leaf belongs to an aggregate
- maximum nesting depth is one aggregate layer

### Step C: Create Beads Issues

For each executable leaf workstream created:
```bash
bd create --title="WS FFF-SS: Short title" --type=task
```

Update `.beads-sdp-mapping.jsonl`:
```json
{"sdp_id":"00-FFF-SS","beads_id":"sdp_dev-XXXX","updated_at":"2026-..."}
```

Aggregate workstreams do not get a `primary` execution issue. If an aggregate needs
tracking for roll-up risk, use a `finding` issue instead.

### Step D: Validate Shapes

```bash
echo "Leafs with primary: $(rg -l \"^- primary:\" docs/workstreams/backlog/00-FFF-*.md | wc -l)"
echo "Mappings:           $(rg -c '\"sdp_id\":\"00-FFF-' .beads-sdp-mapping.jsonl)"
# Primary mappings must match executable leaf workstreams, not total backlog files
```

### Step E: Report

Output:
- Feature ID + number of workstreams created
- Workstream file names
- Beads issue IDs
- Ready-to-run command: first leaf `@build 00-FFF-01` or `@oneshot F0FF`

---

## Default/Interactive Mode

### Step 0: Roadmap Pre-Check — unless --quick

Use `@discovery "feature description"` for roadmap pre-check and product research. Or manually:
1. Extract 3-5 keywords from feature description
2. `rg "<kw1>|<kw2>|<kw3>" docs/ -t md -l`
3. Analyze: ROADMAP overlap, workstream scope overlap, docs/drafts/idea-*.md
4. Present Overlap Report (HIGH/MEDIUM). User resolves: different / extend / supersede / more detail
5. Gate: proceed only after user resolves

### Step 1: Quick Interview (3-5 questions)

Problem, Users, Success. Gate: if vague (<200 words), ask clarification. If @discovery ran: use its output.

### Step 2: @idea

`@idea "..." --spec docs/drafts/discovery-{slug}.md` (if Step 0 ran) or `@idea "..."` (if --quick)

### Step 2.5: @ux — unless --infra

Auto-trigger when @idea output has user-facing keywords (ui, user, interface, dashboard, form) and lacks infra (K8s, CRD, CLI-only).

### Step 3: @design

`@design {task_id}` — workstream files in docs/workstreams/backlog/

Produces workstream files using the **Workstream file format** above.

### Step 3.5: Impact Analysis

Read scope files. grep/rg for conflicts. Categorize: FILE CONFLICT, DATA BOUNDARY, DEPENDENCY CHAIN, PRIORITY SHIFT. Present report. User acknowledges.

### Step 4: Verify Outputs

Check discovery brief, idea spec, ux output, workstreams exist, and that direct
execution targets are leaf workstreams rather than aggregates.

---

## Key Principle: Protocol is Invisible

The user sees:
- Feature description → workstreams created → ready to build

The workstream files, scope declarations, and beads IDs are plumbing.
The user is only asked to annotate if they want to (not required).

## See Also

- @discovery — Product discovery gate (roadmap pre-check)
- @idea — Requirements
- @ux — UX research
- @design — Workstream planning
- @build — Execute single executable leaf workstream
- @oneshot — Execute all ready leaf workstreams for a feature
