# Migration: prompts/commands/ → skills/

## Overview

In SDP v0.6.0, command prompts moved from `prompts/commands/` to `.claude/skills/`.

## Why?

| Problem | Solution |
|---------|----------|
| Duplicate sources | Single source in skills |
| 400-500 line prompts | ≤100 line skills |
| Agents lose focus | Short, structured skills |

## Changes

### Before (v0.5)

```
prompts/commands/build.md (443 lines)
.claude/skills/build/SKILL.md (142 lines)
```

Agents could follow either, leading to inconsistency.

### After (v0.6)

```
.claude/skills/build/SKILL.md (~80 lines)
docs/reference/build-spec.md (full details)
```

Single source, short prompt, details in docs.

## How to Migrate

1. Update any scripts referencing `prompts/commands/`
2. Use skills directly: `@build`, `@review`, etc.
3. For full specs, see `docs/reference/`

## Mapping

| Old | New | Full Spec |
|-----|-----|-----------|
| `prompts/commands/build.md` | `@build` | `docs/reference/build-spec.md` |
| `prompts/commands/review.md` | `@review` | `docs/reference/review-spec.md` |
| `prompts/commands/design.md` | `@design` | `docs/reference/design-spec.md` |
| `prompts/commands/idea.md` | `@idea` | `.claude/skills/idea/SKILL.md` |
| `prompts/commands/deploy.md` | `@deploy` | `.claude/skills/deploy/SKILL.md` |
| `prompts/commands/oneshot.md` | `@oneshot` | `.claude/skills/oneshot/SKILL.md` |
| `prompts/commands/issue.md` | `@issue` | `.claude/skills/issue/SKILL.md` |
| `prompts/commands/hotfix.md` | `@hotfix` | `.claude/skills/hotfix/SKILL.md` |
| `prompts/commands/bugfix.md` | `@bugfix` | `.claude/skills/bugfix/SKILL.md` |

## Breaking Changes

### Removed Files

All files in `prompts/commands/` except `README.md` have been removed.

### References in Code

If you have scripts that reference `prompts/commands/`, update them:

**Before:**
```bash
cat prompts/commands/build.md
```

**After:**
```bash
cat .claude/skills/build/SKILL.md
# Or for full spec:
cat docs/reference/build-spec.md
```

### Custom Commands

If you added custom commands in `prompts/commands/`, migrate them to `.claude/skills/`:

```bash
# 1. Create skill directory
mkdir -p .claude/skills/my-command

# 2. Use template
cp templates/skill-template.md .claude/skills/my-command/SKILL.md

# 3. Fill in details (keep ≤100 lines)

# 4. Create full spec
cp docs/reference/build-spec.md docs/reference/my-command-spec.md

# 5. Validate
sdp skill validate .claude/skills/my-command/SKILL.md
```

## Benefits

1. **Shorter prompts:** Better agent compliance
2. **Consistent structure:** All skills follow same template
3. **Separation of concerns:** Quick reference vs detailed spec
4. **Easier maintenance:** Single source of truth
5. **Better validation:** Automated skill validation

## Timeline

- **v0.5.x:** Both `prompts/commands/` and skills exist (duplication)
- **v0.6.0:** `prompts/commands/` deprecated, skills only
- **Future:** `prompts/commands/` directory removed entirely

## See Also

- [ADR-007: Skill Length Limit](../adr/007-skill-length-limit.md)
- [Skill Template](../../templates/skill-template.md)
- [Breaking Changes Guide](../migrations/breaking-changes.md)
