# ⚠️ DEPRECATED: Moved to Skills

All command prompts have been migrated to skills in `.claude/skills/`.

## Migration

| Old Location | New Location |
|--------------|--------------|
| `prompts/commands/build.md` | `.claude/skills/build/SKILL.md` |
| `prompts/commands/review.md` | `.claude/skills/review/SKILL.md` |
| `prompts/commands/design.md` | `.claude/skills/design/SKILL.md` |
| `prompts/commands/idea.md` | `.claude/skills/idea/SKILL.md` |
| `prompts/commands/deploy.md` | `.claude/skills/deploy/SKILL.md` |
| `prompts/commands/oneshot.md` | `.claude/skills/oneshot/SKILL.md` |
| `prompts/commands/issue.md` | `.claude/skills/issue/SKILL.md` |
| `prompts/commands/hotfix.md` | `.claude/skills/hotfix/SKILL.md` |
| `prompts/commands/bugfix.md` | `.claude/skills/bugfix/SKILL.md` |

## Why?

1. **Single source of truth:** Skills only, no duplication
2. **Shorter prompts:** ≤100 lines for better agent compliance
3. **Better structure:** Clear workflow with quality gates
4. **Reference docs:** Detailed specs in `docs/reference/`

## Full Documentation

Detailed specifications moved to `docs/reference/`:
- `docs/reference/build-spec.md`
- `docs/reference/review-spec.md`
- `docs/reference/design-spec.md`

See [Migration Guide](../docs/migration/prompts-to-skills.md) for details.

---

**Deprecated in:** v0.6.0  
**Migration Path:** Use skills directly (`@build`, `@review`, `@design`, etc.)  
**See Also:** [Breaking Changes](../docs/migrations/breaking-changes.md)
