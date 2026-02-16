# /idea — Requirements Gathering

When calling `/idea {description}`:

1. Load full prompt: `@.claude/skills/idea.md`
2. Execute Mandatory Initial Dialogue
3. Create draft in `docs/drafts/idea-{slug}.md`
4. Output summary for user

## Quick Reference

**Input:** Feature description from user
**Output:** `docs/drafts/idea-{slug}.md`
**Next:** `/design idea-{slug}`
