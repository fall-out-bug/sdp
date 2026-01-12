# /design — Analyze + Plan

When calling `/design {slug}`:

1. Load full prompt: `@prompts/commands/design.md`
2. Read PROJECT_MAP.md and INDEX.md
3. Read draft: `docs/drafts/idea-{slug}.md`
4. Create all WS files in `workstreams/backlog/`
5. Update INDEX.md
6. Output summary

## Quick Reference

**Input:** `docs/drafts/idea-{slug}.md`
**Output:** `docs/workstreams/backlog/WS-XXX-*.md`
**Next:** `/build WS-XXX-01`
