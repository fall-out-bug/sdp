---
name: feature
description: Unified entry point for feature development with progressive disclosure
tools: Read, Write, Edit, Bash, AskUserQuestion, Skill
---

# /feature - Unified Feature Development

Progressive disclosure workflow: vision -> requirements -> planning -> execution.

## When to Use

- Starting new feature (recommended for all)
- Exploring feature idea
- Creating MVP from scratch
- Power users can skip to @idea or @design directly

## Workflow

### Phase 1: Vision Interview (3-5 questions)

AskUserQuestion about:
- **Mission**: What problem do we solve?
- **Users**: Who are we building for?
- **Success Metrics**: How do we measure success?

### Phase 2: Generate PRODUCT_VISION.md

Create or update `PRODUCT_VISION.md` at project root.

Format:
```markdown
# PRODUCT_VISION.md

> **Last Updated:** YYYY-MM-DD
> **Version:** 1.0

## Mission

[Product mission statement]

## Users

1. **[User type]**
2. **[User type]**

## Success Metrics

- [ ] [Metric 1]
- [ ] [Metric 2]

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| [Aspect] | [Decision] | [Why] |

## Non-Goals

- [Not doing X]
- [Not doing Y]
```

**After creating Strategic Tradeoffs table:** Log each tradeoff as a decision:

```bash
# Find project root (look for .git)
root_dir=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

# For each row in Strategic Tradeoffs:
sdp decisions log --type="tradeoff" \
  --question="Aspect: [Aspect]" \
  --decision="[Decision]" \
  --rationale="[Rationale]" \
  --feature-id="{FXXX}" \
  --maker="user"
```

### Phase 3: Technical Interview (5-8 questions)

AskUserQuestion about:
- Technical approach (architecture, storage, failure mode)
- Tradeoffs (security vs performance, complexity vs speed)
- Integration points

**After each technical decision:** Log the decision:

```bash
sdp decisions log \
  --type=technical \
  --question="[What was the question?]" \
  --decision="[What was decided?]" \
  --rationale="[Why this choice?]" \
  --alternatives="[Option 1],[Option 2]" \
  --feature-id="{FXXX}" \
  --maker=user
```

### Phase 4: Generate intent.json

Validate against `docs/schema/intent.schema.json`, save to `docs/intent/{slug}.json`

```bash
# Create intent directory
mkdir -p docs/intent

# Create and validate
from sdp.schema.validator import IntentValidator
from sdp.schema.models import Intent

intent = Intent.from_dict({
    "problem": "...",
    "users": [...],
    "success_criteria": [...]
})

# Validate
validator = IntentValidator()
validator.validate(intent.to_dict())

# Save
import json
with open("docs/intent/{slug}.json", "w") as f:
    json.dump(intent.to_dict(), f, indent=2)
```

### Phase 5: Create Requirements Draft (REQUIRED)

Create `docs/drafts/idea-{slug}.md` with full specification:

```markdown
# {Feature Name}

> **Feature ID:** FXXX
> **Status:** Draft
> **Created:** YYYY-MM-DD

## Problem
[From interview]

## Users
[From interview]

## Success Criteria
[From interview]

## Goals
[Primary goals]

## Non-Goals
[Out of scope]

## Technical Approach
[From technical interview]
```

### Phase 6: Transition to @design

Call `/design` with full context (vision + intent).

### Phase 7: Orchestrator Execution (OPTIONAL)

After @design completes workstream planning, optionally execute workstreams autonomously:

```bash
# Option 1: Interactive execution
@feature "Add user authentication"
# ... phases 1-6 complete ...
# Claude asks: "Execute workstreams now? (y/n)"

# Option 2: Immediate execution
@feature "Add user authentication" --execute

# Option 3: Resume from checkpoint
@feature "Add user authentication" --resume F001
```

**Orchestrator capabilities:**
- Executes workstreams in dependency order (topological sort)
- Saves checkpoints to `.oneshot/{feature}-checkpoint.json`
- Auto-retry on failures (max 2 retries)
- Real-time progress updates: `[HH:MM] Executing WS-XXX...`
- Resumable from interruptions

**Checkpoint format:**
```json
{
  "id": "F001",
  "feature_id": "F001",
  "status": "in_progress",
  "completed_workstreams": ["WS-001", "WS-002"],
  "current_workstream": "WS-003",
  "created_at": "2026-02-06T15:23:00Z",
  "updated_at": "2026-02-06T15:46:00Z"
}
```

**Progress tracking example:**
```
[15:23] Starting feature execution: F001
[15:23] Loading workstreams...
[15:23] Building dependency graph...
[15:23] Execution order: [WS-001 WS-002 WS-003]
[15:24] Executing WS-001 (1/3)...
[15:46] WS-001 complete (22m)
[15:46] Executing WS-002 (2/3)...
[16:05] WS-002 complete (19m)
[16:05] Executing WS-003 (3/3)...
[16:20] WS-003 complete (15m)
[16:20] Feature execution complete: 3/3 workstreams, 57m total
```

**Error handling:**
- HIGH/MEDIUM issues: Auto-retry (max 2), log context
- CRITICAL errors: Pause and escalate to human
- Workstream blocked: Mark checkpoint as failed, save state

**Integration with @build skill:**
- Each workstream executed via `@build {ws_id}`
- Beads integration: `bd update {ws_id} --status in_progress`
- On success: `bd close {ws_id} --reason "WS completed"`
- On failure: `bd update {ws_id} --status blocked`

## Power User Flags

- `--vision-only` -- Only create vision, skip planning
- `--no-interview` -- Skip questions, use defaults
- `--update-vision` -- Update existing PRODUCT_VISION.md
- `--spec PATH` -- Start from existing spec

## Output

- `PRODUCT_VISION.md` -- Project manifesto
- `docs/drafts/idea-{slug}.md` -- Full spec
- `docs/intent/{slug}.json` -- Machine-readable intent
