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

## Write Plan (F101)

Before modifying any file, emit a write plan covering PRODUCT_VISION.md, PRD.md, ROADMAP.md, and feature drafts:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason. Covers PRODUCT_VISION.md, PRD.md, ROADMAP.md, and feature drafts.
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"vision"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `sdp/schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @vision <target>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

## PRD Mode

Detect project type (service/library/cli) from structure. Scaffold PRD with type-appropriate sections. Generate diagrams from @prd annotations in code. Validate section limits. `@vision "name" --update` regenerates diagrams from annotations.

## When to Use

Initial setup, quarterly review, major pivot, new market.

## Output

PRODUCT_VISION.md, docs/prd/PRD.md, docs/roadmap/ROADMAP.md, docs/drafts/feature-*.md

## See Also

- @idea — Feature-level requirements
- @feature — Planning orchestrator
