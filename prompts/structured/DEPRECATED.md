# ⚠️ DEPRECATED: 4-Phase Workflow

**Status:** Superseded by slash commands  
**Date:** 2026-01-11  
**Replacement:** Use `/design`, `/build`, `/codereview`, `/deploy` instead

---

## Why Deprecated?

The 4-phase workflow (phase-1-analyze.md, phase-2-design.md, etc.) has been superseded by **slash commands** which provide:

1. **Better UX** - Single command vs multi-phase process
2. **Less duplication** - One source of truth
3. **Easier maintenance** - Update one file, not two
4. **GitHub integration** - Commands sync with GitHub Projects
5. **Checkpoint/resume** - `/oneshot` supports recovery

---

## Migration Guide

| Old (Phases) | New (Commands) |
|--------------|----------------|
| Phase 1 (Analyze) | `/idea` + `/design` |
| Phase 2 (Plan) | `/design` |
| Phase 3 (Implement) | `/build` or `/oneshot` |
| Phase 4 (Review) | `/codereview` |
| - | `/deploy` |

### Example

**Before (4-phase):**
```
1. Run phase-1-analyze.md → create WS map
2. Run phase-2-design.md → create WS files
3. Run phase-3-implement.md → implement
4. Run phase-4-review.md → review
```

**After (slash commands):**
```bash
/design idea-lms-integration  # Steps 1+2 combined
/oneshot F60                  # Step 3 (autonomous)
/codereview F60                   # Step 4
/deploy F60                   # New: deployment
```

---

## Files in This Directory

- `phase-1-analyze.md` (325 lines) → Use `/design`
- `phase-2-design.md` (325 lines) → Use `/design`
- `phase-3-implement.md` (420 lines) → Use `/build`
- `phase-4-review.md` (390 lines) → Use `/codereview`

**Total:** ~1460 lines now replaced by 4 commands

---

## When to Still Use Phases?

**Only if:**
- You prefer manual step-by-step control
- Working without slash command support
- Legacy workflow required

**But recommended:** Switch to slash commands for better experience.

---

## Master Prompts Location

All slash commands delegate to:

📁 **sdp/prompts/commands/**
- `idea.md` - Requirements gathering
- `design.md` - WS planning (496 lines)
- `build.md` - TDD implementation (400 lines)
- `review.md` - Code review (460+ lines)
- `deploy.md` - Deployment (480+ lines)
- `oneshot.md` - Autonomous execution (750+ lines)
- `issue.md` - Bug analysis (640+ lines)
- `hotfix.md` - Emergency fixes (420+ lines)
- `bugfix.md` - Quality fixes (530+ lines)

---

## Deprecation Timeline

- **2026-01-11:** Phase files marked deprecated
- **2026-02-01:** Consider archiving phase files
- **2026-03-01:** Remove phase files (if no usage)

---

For questions or issues with migration, see:
- `CLAUDE.md` - Main documentation
- `sdp/README.md` - Protocol overview
