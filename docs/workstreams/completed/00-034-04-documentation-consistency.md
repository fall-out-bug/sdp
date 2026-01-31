---
ws_id: 00-034-04
feature: F034
status: completed
complexity: SMALL
project_id: "00"
---

# Workstream: Documentation Consistency

**ID:** 00-034-04  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** SMALL (~200 edits)

---

## Goal

Устранить все inconsistencies в документации: стандартизировать WS ID формат, исправить битые ссылки, выровнять версии.

---

## Context

Выявленные проблемы:

| Issue | Count | Impact |
|-------|-------|--------|
| WS ID формат (WS-XXX vs PP-FFF-SS) | ~50 occurrences | HIGH (confusion) |
| Битые ссылки | ~10 links | MEDIUM (broken navigation) |
| Версия несоответствие (v0.5.0 vs v0.6.0) | ~8 places | LOW (confusion) |
| Skill docs без примеров | ~5 skills | MEDIUM (usability) |

---

## Scope

### In Scope
- ✅ Replace all `WS-XXX-YY` with `PP-FFF-SS` format
- ✅ Fix all broken links
- ✅ Align version to v0.6.0 everywhere
- ✅ Add examples to skill documentation
- ✅ Update README badges

### Out of Scope
- ❌ Code changes
- ❌ New documentation sections
- ❌ Translation updates

---

## Dependencies

**Depends On:**
- None (can start immediately)

**Blocks:**
- None (independent quick win)

---

## Acceptance Criteria

- [ ] `grep -r "WS-[0-9]" docs/ README.md CLAUDE.md` returns 0 results
- [ ] All links in docs/ are valid (verified by script)
- [ ] `grep -r "v0.5.0" .` returns 0 results (except CHANGELOG)
- [ ] All skills in `.claude/skills/` have usage examples

---

## Implementation Plan

### Task 1: Standardize WS ID Format

**Find & Replace:**
```
WS-001-01 → 00-001-01
WS-XXX-YY → PP-FFF-SS (in templates/examples)
WS-COMMENTS-01 → 00-COMMENTS-01 (or use numeric)
```

**Files to update:**
- [ ] `README.md` (lines 73, 128, 179-182)
- [ ] `CLAUDE.md` (multiple occurrences)
- [ ] `PROTOCOL.md` (examples section)
- [ ] `START_HERE.md`
- [ ] `docs/beginner/*.md`
- [ ] `docs/reference/commands.md`
- [ ] `templates/workstream.md`

### Task 2: Fix Broken Links

**Known broken links:**

| File | Broken Link | Fix |
|------|-------------|-----|
| `START_HERE.md:103` | `docs/TUTORIAL.md` | `docs/beginner/TUTORIAL.md` |
| `README.md:218` | `docs/TUTORIAL.md` | `docs/beginner/TUTORIAL.md` |
| `CLAUDE.md` | `docs/GLOSSARY.md` | `docs/reference/GLOSSARY.md` |
| `docs/internals/README.md` | `docs/SITEMAP.md` | Verify path |

**Verification script:**
```bash
# Find all markdown links and verify
grep -rohE '\[.*?\]\(.*?\.md\)' docs/ | \
  sed 's/.*(\(.*\))/\1/' | \
  while read link; do
    [ ! -f "docs/$link" ] && echo "BROKEN: $link"
  done
```

### Task 3: Align Version Numbers

**Current state:**
- `README.md:264` → v0.5.0
- `START_HERE.md:282` → v0.5.0
- `.cursorrules` → v0.6.0
- `pyproject.toml` → (check actual)

**Target:** v0.6.0 everywhere

**Files to update:**
- [ ] `README.md`
- [ ] `START_HERE.md`
- [ ] `PROTOCOL.md`
- [ ] `CLAUDE.md`
- [ ] `pyproject.toml` (source of truth)

### Task 4: Add Skill Documentation Examples

**Skills missing examples:**

| Skill | Missing |
|-------|---------|
| `prd/SKILL.md` | Usage example |
| `guard/SKILL.md` | Output example |
| `tdd/SKILL.md` | Full cycle example |
| `think/SKILL.md` | Agent output example |
| `init/SKILL.md` | Interactive flow example |

**Template for examples:**
```markdown
## Example

### Input
```bash
@skill "parameters"
```

### Output
```
→ Step 1: ...
→ Step 2: ...
✅ Complete
```
```

### Task 5: Update README Badges

**Current:** Coverage badge shows 91%
**Fix:** Update to actual coverage or remove until CI reports it

```markdown
[![Coverage](https://img.shields.io/badge/coverage-68%25-yellow.svg)](tests/)
```

Or link to CI coverage report when available.

---

## DO / DON'T

### Documentation

**✅ DO:**
- Use consistent formatting (PP-FFF-SS)
- Verify all links before committing
- Keep examples up to date
- Use relative paths for internal links

**❌ DON'T:**
- Mix WS ID formats
- Leave dead links
- Use absolute URLs for internal docs
- Forget to update table of contents

---

## Files to Modify

**Documentation:**
- [ ] `README.md`
- [ ] `CLAUDE.md`
- [ ] `PROTOCOL.md`
- [ ] `START_HERE.md`
- [ ] `docs/beginner/*.md` (6 files)
- [ ] `docs/reference/*.md` (15 files)
- [ ] `templates/workstream.md`

**Skills:**
- [ ] `.claude/skills/prd/SKILL.md`
- [ ] `.claude/skills/guard/SKILL.md`
- [ ] `.claude/skills/tdd/SKILL.md`
- [ ] `.claude/skills/think/SKILL.md`
- [ ] `.claude/skills/init/SKILL.md`

**Config:**
- [ ] `pyproject.toml` (version)

---

## Test Plan

### Verification
```bash
# No legacy WS format
grep -r "WS-[0-9]" docs/ README.md CLAUDE.md PROTOCOL.md
# Expected: 0 results

# All links valid
./scripts/check_links.sh
# Expected: 0 broken links

# Version consistent
grep -r "v0\.[0-9]\." . --include="*.md" | grep -v "v0.6.0" | grep -v CHANGELOG
# Expected: 0 results (except historical references)
```

### Manual Verification
- [ ] Read through START_HERE.md as new user
- [ ] Follow all tutorial links
- [ ] Try all skill examples

---

## Estimated Effort

| Task | Changes | Time |
|------|---------|------|
| WS ID format | ~50 replacements | Quick |
| Broken links | ~10 fixes | Quick |
| Version alignment | ~8 files | Quick |
| Skill examples | ~5 skills | Medium |
| **Total** | **~75 edits** | **SMALL** |

---

**Version:** 1.0  
**Created:** 2026-01-31
