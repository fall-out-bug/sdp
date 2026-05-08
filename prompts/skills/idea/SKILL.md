---
name: idea
description: Interactive requirements gathering with progressive disclosure
version: 5.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @idea - Requirements Gathering with Progressive Disclosure

Deep interviewing to capture comprehensive feature requirements using progressive disclosure (3-question cycles). Creates markdown spec, optionally creates Beads task.

---

## EXECUTE THIS NOW

When user invokes `@idea "feature description"`:

### Step 1: Read Context

Read existing project files to understand context:
- `PRODUCT_VISION.md` - Align with project goals
- `docs/specs/**/*` - Similar features

### Step 2: Progressive Interview (3-Question Cycles)

**Question Target:**
- Minimum: 12 questions (bounded exploration)
- Maximum: 27 questions (deep analysis)
- Average: 18-20 questions per feature

**3-Question Cycles:**

1. Ask 3 focused questions
2. Offer trigger point after each cycle
3. User chooses: continue / deep design / skip to @design

**Cycle 1 - Vision (3 questions):**
- What is the core mission of this feature?
- How does this align with PRODUCT_VISION.md?
- Who are the primary users?

**TRIGGER POINT (after each cycle):**
- Continue (more questions)
- Deep design (jump to @design with architectural exploration)
- Skip to @design (move to workstream decomposition)

**Cycle 2 - Problem & Users (3 questions):**
- What problem does this solve?
- What are the user pain points?
- What happens if we don't build this?

**Cycle 3 - Technical Approach (3 questions):**
- Storage/data requirements?
- Failure modes to handle?
- Integration points?

**Cycle 4 - UI/UX & Quality (3 questions):**
- UI/UX requirements?
- Performance targets?
- Security considerations?

**Cycle 5 - Testing & Edge Cases (3 questions):**
- Testing strategy?
- Edge cases to handle?
- Success metrics?

### Step 3: TMI Detection

If user provides extensive detail upfront:
- "detailed spec", "full implementation", "complete architecture"
- User writes >500 characters in initial prompt

Offer shortcuts:
- Continue with targeted questions (recommended)
- Skip to @design with detailed spec
- Use --quiet mode for minimal questions

### Step 4: Create Beads Task

```bash
bd create --title="{feature_title}" --type=feature --priority=2
```

Include in description:
- Context & Problem
- Goals & Non-Goals
- Technical Approach
- Concerns & Tradeoffs

### Step 5: Create Spec File

Create `docs/intent/{task_id}.json` with machine-readable intent.

**Output must validate against** `schema/intent.schema.json` ([intent.schema.json](../../schema/intent.schema.json)). **Required fields:** `problem`, `users`, `success_criteria`. Optional: `context`, `non_goals`, `risks`, `question_count`.

---

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive spec with tradeoffs explored

---

## Modes

| Mode | Questions | Purpose |
|------|-----------|---------|
| Default | 12-27 | Full progressive interview |
| `--quiet` | 3-5 | Minimal questions (core only) |
| `--spec path` | Varies | Use existing spec as base |

---

## --quiet Mode

Minimal questions (3-5 core only):
1. Mission?
2. Users?
3. Core requirement?

Skip deep-dive cycles, move directly to @design.

---

## Output

**Primary:** Beads task ID (e.g., `sdp-xxx`)

**Secondary:**
- `docs/intent/{task_id}.json` - Machine-readable intent
- Question count included in metadata

---

## Key Principles

1. **Progressive disclosure** - 3 questions at a time
2. **User-controlled depth** - trigger points after each cycle
3. **Respect brevity** - --quiet mode for experienced users
4. **No obvious questions** - explore tradeoffs, not yes/no
5. **TMI detection** - offer shortcuts when user over-explains

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@idea "feature"` | Create task with progressive interview |
| `@idea "feature" --quiet` | Minimal questions (3-5 core only) |
| `bd show {id}` | View task details |
| `@design {id}` | Decompose into workstreams |

---

## Few-Shot Examples

**Good — productive 3-question cycle (answers chain):**
- Q1: What is the core mission of this feature?  
  A: "Let users reset password via email when they forget it."  
- Q2: How does this align with PRODUCT_VISION.md?  
  A: "Vision says 'self-service account recovery'; this is the main path."  
- Q3: Who are the primary users?  
  A: "End users who forgot password; support team (fewer tickets)."  
→ Trigger: Continue (next cycle: problem/pain) or Skip to @design.

**Bad — single yes/no question:**
- Q: Do you need authentication? A: "Yes."  
Reason: No exploration. Each answer should inform the next question.

**Bad — TMI upfront (offer shortcut):**
User writes 500+ chars: "I need a full auth system with OAuth, MFA, session management, rate limiting..."
→ Offer: Continue with targeted questions (recommended) | Skip to @design | --quiet mode.

**Good:** Each answer informs the next; 3 questions per cycle; explore tradeoffs, not yes/no.

---

## Write Plan (F101)

Before creating the intent file or beads task, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (intent file, beads task, event log).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"idea"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @idea <feature>:
  CREATE: docs/intent/{task_id}.json — Machine-readable intent spec
  CREATE: beads task — Feature tracking issue
  MODIFY: .sdp/log/events.jsonl — Write plan event log

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

---

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- `@design` - Workstream decomposition
- `@build` - Execute workstream
- `@feature` - Orchestrator that calls @idea + @design
