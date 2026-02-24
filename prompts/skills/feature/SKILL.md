---
name: feature
description: Feature planning orchestrator (discovery -> idea -> ux -> design -> workstreams)
---

# @feature

Orchestrate product discovery, requirements, UX research, and workstream design.

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

**Workstream file format:**

```markdown
# 00-FFF-SS: Feature Name — Step Description

Feature: FFFF (sdp_dev-XXXX)
Phase: N
Status: Backlog

## Goal

One paragraph: what this workstream does and why.

## Scope Files

- path/to/file/or/dir (exact files or directory prefixes this WS touches)
- ...

## Dependencies

- 00-FFF-S1: prior workstream (if any)

## Acceptance Criteria

- [ ] Specific, testable criterion 1
- [ ] Specific, testable criterion 2
- [ ] go build ./... passes
- [ ] go test ./... passes
```

### Step C: Create Beads Issues

For each workstream created:
```bash
bd create --title="WS FFF-SS: Short title" --type=task
```

Update `.beads-sdp-mapping.jsonl`:
```json
{"sdp_id":"00-FFF-SS","beads_id":"sdp_dev-XXXX","updated_at":"2026-..."}
```

### Step D: Validate Counts

```bash
echo "Mapping: $(wc -l < .beads-sdp-mapping.jsonl)"
echo "Backlog:  $(ls docs/workstreams/backlog/*.md | wc -l)"
# Must be equal
```

### Step E: Report

Output:
- Feature ID + number of workstreams created
- Workstream file names
- Beads issue IDs
- Ready-to-run command: `@build 00-FFF-01` or `@oneshot F0FF`

---

## Default/Interactive Mode

### Step 0: Roadmap Pre-Check — unless --quick

1. Extract 3-5 keywords from feature description
2. `rg "<kw1>|<kw2>|<kw3>" docs/ -t md -l`
3. Analyze: ROADMAP overlap, workstream scope overlap, docs/drafts/idea-*.md
4. Present Overlap Report (HIGH/MEDIUM). User resolves: different / extend / supersede / more detail
5. Gate: proceed only after user resolves

### Step 1: Quick Interview (3-5 questions)

Problem, Users, Success. Gate: if vague (<200 words), ask clarification.

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

Check discovery brief, idea spec, ux output, workstreams exist.

---

## Key Principle: Protocol is Invisible

The user sees:
- Feature description → workstreams created → ready to build

The workstream files, scope declarations, and beads IDs are plumbing.
The user is only asked to annotate if they want to (not required).

## See Also

- @idea — Requirements
- @ux — UX research
- @design — Workstream planning
- @build — Execute single workstream
- @oneshot — Execute all workstreams for a feature
