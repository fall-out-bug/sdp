---
name: design
description: System design with progressive disclosure, produces workstream files
version: 2.1.0
---

# @design

Multi-agent design (Arch + Security + SRE) with progressive discovery blocks.

## When to Use

After @idea, or directly from a feature description. Creates workstream files with AC and scope.

## Pre-flight

**Before creating draft:** `ls docs/drafts/idea-*` — do not duplicate. If an idea draft already covers this topic, reuse or extend it instead of creating a new one.

**Beads from review:** In scope by default. Mark OOS only with explicit justification (e.g. "duplicate", "superseded by prior work").

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

**When source is beads (review findings):** For each bead, run `bd show <id>` and grep the codebase for the fix. If already fixed, run `bd close <id>` and remove from scope. Do not create WS for beads that are already addressed.

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

## Write Plan (F101)

Before modifying any file, emit a write plan covering workstream files, design docs, and INDEX.md:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason. Covers workstream files, design docs, and INDEX.md updates.
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"design"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @design <target>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

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

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @idea — Requirements
- @feature — Full planning orchestrator
- @build — Execute single workstream
- @oneshot — Execute all workstreams
