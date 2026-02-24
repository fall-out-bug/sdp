---
name: design
description: System design with progressive disclosure, produces workstream files
---

# @design

Multi-agent design (Arch + Security + SRE) with progressive discovery blocks.

## When to Use

After @idea, or directly from a feature description. Creates workstream files with AC and scope.

## Workflow

### 1. Load requirements

- `docs/intent/{task_id}.json` or `docs/drafts/idea-*.md` if available
- Or: use the feature description directly

### 2. Progressive discovery — unless --quiet

3-5 discovery blocks, 2-3 questions each:
- **Architecture**: What components change? What's the data model?
- **Security**: Any auth, crypto, or boundary concerns?
- **Operations**: Any monitoring, logging, or CI concerns?

After each block: Continue / Skip / Done

### 3. Generate workstream files

Create `docs/workstreams/backlog/00-FFF-SS.md` for each deliverable.

**Required sections:**

```markdown
# 00-FFF-SS: Feature Name — Step Description

Feature: FFFF (sdp_dev-XXXX)
Phase: N
Status: Backlog

## Goal

One paragraph: what and why.

## Scope Files

List exact file paths or directory prefixes this workstream touches.
Used by sdp-guard for boundary checking and CI scope-compliance.

- internal/evidence/
- cmd/sdp-evidence/main.go

## Dependencies

- 00-FFF-01: prerequisite workstream (if any)

## Acceptance Criteria

Specific, testable, binary (pass/fail):

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] go build ./... passes
- [ ] go test ./internal/evidence/... passes
```

### 4. Create Beads issues

```bash
bd create --title="WS FFF-SS: Short title" --type=task
```

Append to `.beads-sdp-mapping.jsonl`:
```json
{"sdp_id":"00-FFF-SS","beads_id":"sdp_dev-XXXX","updated_at":"2026-..."}
```

**ALWAYS verify counts match:**
```bash
echo "Mapping: $(wc -l < .beads-sdp-mapping.jsonl)"
echo "Backlog:  $(ls docs/workstreams/backlog/*.md | wc -l)"
```

### 5. Update INDEX.md

Add new workstreams to the appropriate phase table in `docs/workstreams/INDEX.md`.

## Modes

| Mode | Blocks |
|------|--------|
| Default | 3-5 discovery blocks |
| --quiet | 2 blocks (Architecture + Data only) |

## Output

- Workstream files in `docs/workstreams/backlog/`
- `docs/drafts/{task_id}-design.md` (architecture notes)
- Updated `docs/workstreams/INDEX.md`
- Updated `.beads-sdp-mapping.jsonl`

## See Also

- @idea — Requirements
- @feature — Full planning orchestrator
- @build — Execute single workstream
- @oneshot — Execute all workstreams
