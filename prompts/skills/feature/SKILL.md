---
name: feature
description: Feature planning orchestrator (discovery -> idea -> ux -> design -> workstreams)
---

# @feature

Orchestrate product discovery, requirements, UX research, and workstream design.

## Workflow

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

### Step 3.5: Impact Analysis

Read scope files. grep/rg for conflicts. Categorize: FILE CONFLICT, DATA BOUNDARY, DEPENDENCY CHAIN, PRIORITY SHIFT. Present report. User acknowledges.

### Step 4: Verify Outputs

Check discovery brief, idea spec, ux output, workstreams exist.

## Flags

| Flag | Effect |
|------|--------|
| --quick | Skip Step 0 |
| --infra | Skip @ux |

## Output

Discovery brief, idea spec, ux doc (if user-facing), workstream files. Ready for @oneshot or @build.

## See Also

- @idea — Requirements
- @ux — UX research
- @design — Workstream planning
- @oneshot — Execution
