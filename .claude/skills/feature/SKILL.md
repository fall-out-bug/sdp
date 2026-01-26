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

### Phase 5: Transition to @design

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
