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

### Phase 3: Technical Interview (5-8 questions)

AskUserQuestion about:
- Technical approach (architecture, storage, failure mode)
- Tradeoffs (security vs performance, complexity vs speed)
- Integration points

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

## Power User Flags

- `--vision-only` -- Only create vision, skip planning
- `--no-interview` -- Skip questions, use defaults
- `--update-vision` -- Update existing PRODUCT_VISION.md
- `--spec PATH` -- Start from existing spec

## Output

- `PRODUCT_VISION.md` -- Project manifesto
- `docs/drafts/idea-{slug}.md` -- Full spec
- `docs/intent/{slug}.json` -- Machine-readable intent
