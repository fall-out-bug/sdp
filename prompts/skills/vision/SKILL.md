---
name: vision
description: Strategic product planning - vision, PRD, roadmap from expert analysis
---

# @vision

Transform project ideas into vision, PRD, and roadmap.

## Workflow

1. **Interview** (3-5 questions) — Problem, users, success, MVP, competitors
2. **7 Expert Analysis** — Product, Market, Technical, UX, Business, Growth, Risk. **Output format (per expert):** `## {Expert}\n**Summary:** 2-3 sentences.\n**Key finding:** One actionable insight.\n**Risk/opportunity:** If any.` Avoid vague essays; be concrete.
3. **Generate artifacts** — PRODUCT_VISION.md, docs/prd/PRD.md (or docs/PROJECT_MAP.md), docs/roadmap/ROADMAP.md
4. **Extract features** — docs/drafts/feature-{slug}.md for P0/P1

## PRD Mode

Detect project type (service/library/cli) from structure. Scaffold PRD with type-appropriate sections. Generate diagrams from @prd annotations in code. Validate section limits. `@vision "name" --update` regenerates diagrams from annotations.

## When to Use

Initial setup, quarterly review, major pivot, new market.

## Output

PRODUCT_VISION.md, docs/prd/PRD.md, docs/roadmap/ROADMAP.md, docs/drafts/feature-*.md

## See Also

- @idea — Feature-level requirements
- @feature — Planning orchestrator
