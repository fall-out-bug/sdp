# SDP Workflow Decision Guide

**Date:** 2026-01-30
**Status:** Active
**Decision:** Beads-first workflow is the recommended approach

---

## Decision Summary

**✅ RECOMMENDED:** Beads-first workflow using `@feature` skill
**⚠️ ALTERNATIVE:** Traditional markdown workflow using direct `@idea` + `@design`
**❌ DEPRECATED:** Pure markdown files without Beads integration

---

## Workflows Comparison

### Beads-First Workflow (Recommended)

**Entry Point:** `@feature "Feature description"`

**Characteristics:**
- Progressive disclosure: Vision → Requirements → Planning → Execution
- Automatic task creation in Beads (git-backed issue tracker)
- All artifacts tracked in one place
- Natural dependency tracking
- Better for team collaboration

**When to Use:**
- ✅ **Default choice** for all new features
- ✅ Multi-person projects (shared task tracking)
- ✅ Complex features (10+ workstreams)
- ✅ Projects needing dependency tracking
- ✅ Teams wanting traceability

**Flow:**
```
@feature "Add user authentication"
  ↓ (generates PRODUCT_VISION.md)
@design feature-user-auth
  ↓ (creates Beads tasks)
@oneshot F01  (or execute workstreams individually)
```

**Artifacts Created:**
- `PRODUCT_VISION.md` (generated)
- Beads tasks (auto-created)
- Workstream markdown files (if needed)

**Pros:**
- Single source of truth (Beads)
- Automatic dependency management
- PR-based execution (@oneshot)
- Progress tracking built-in
- Easier handoffs between team members

**Cons:**
- Requires Beads CLI installation
- Git-based (requires commits)
- Slightly more overhead for tiny features

---

### Traditional Markdown Workflow (Alternative)

**Entry Point:** `@idea "Feature description"`

**Characteristics:**
- Direct task creation via AskUserQuestion
- Optional Beads integration
- Markdown-first documentation
- Faster for solo developers

**When to Use:**
- ⚠️ Solo development (single person)
- ⚠️ Quick prototypes (1-3 workstreams)
- ⚠️ Beads CLI not available
- ⚠️ Exploratory work (may abandon)

**Flow:**
```
@idea "Add user authentication"
  ↓ (interactive requirements gathering)
@design idea-user-auth
  ↓ (creates workstream markdown files)
@build WS-001-01
@build WS-001-02
...
```

**Artifacts Created:**
- `docs/drafts/idea-*.md` (requirements)
- `docs/workstreams/backlog/WS-*.md` (workstreams)
- Beads tasks (optional)

**Pros:**
- Faster to start (no PRs)
- Works without Beads
- Simpler for small features
- Markdown files are readable

**Cons:**
- Manual dependency tracking
- No centralized progress tracking
- Harder to see status across features
- Requires manual integration with Beads

---

## Migration Path

### From Markdown to Beads-First

**If you're using traditional markdown workflow:**

1. **Existing features:** Continue using markdown → natural migration point
2. **New features:** Start using `@feature` instead of `@idea`
3. **Hybrid approach:** Use `@feature` for planning, execute workstreams manually

**Steps to migrate a feature:**
```bash
# 1. Create Beads tasks from existing workstreams
@design existing-feature-name

# 2. Execute using Beads workflow
@oneshot F01
```

---

## Decision Matrix

| Situation | Recommended Workflow |
|-----------|---------------------|
| New feature (default) | `@feature` (Beads-first) |
| Team project (2+ people) | `@feature` (Beads-first) |
| Complex feature (10+ WS) | `@feature` (Beads-first) |
| Solo developer | Either (preference) |
| Quick prototype (1-3 WS) | `@idea` (markdown) |
| Beads not installed | `@idea` (markdown) |
| Exploratory work | `@idea` (markdown) |
| Production feature | `@feature` (Beads-first) |

---

## Implementation Notes

### Beads-First Workflow Details

See `.claude/skills/feature/SKILL.md` for full implementation.

**Key Steps:**
1. Vision: Generate PRODUCT_VISION.md
2. Requirements: Interactive deep interview (AskUserQuestion)
3. Planning: EnterPlanMode for workstream decomposition
4. Execution: @oneshot for autonomous execution

**Configuration:**
- Beads CLI: `pip install beads-cli`
- Repository: `git init` (if not already)
- Environment: `BEADS_USE_MOCK=false` (default)

### Traditional Markdown Workflow Details

See `.claude/skills/` for markdown-based prompts.

**Key Steps:**
1. Requirements: `@idea` → docs/drafts/
2. Planning: `@design` → docs/workstreams/
3. Execution: `@build WS-FFF-SS` per workstream
4. Review: `@review FFF`

**Configuration:**
- Optional: Beads CLI
- Repository: Git (for workstream files)
- Environment: None required

---

## FAQ

**Q: Should I delete markdown files if using Beads-first?**
A: No. Keep workstream markdown files as documentation. Beads tracks execution status, markdown describes intent.

**Q: Can I mix workflows?**
A: Yes, but not recommended for a single feature. Choose one workflow per feature for consistency.

**Q: What if Beads CLI changes?**
A: The markdown workflow is a stable fallback. SDP maintains both interfaces.

**Q: Which workflow is used in SDP's own development?**
A: SDP uses Beads-first for all features (2025+).

---

## Future Evolution

**Phase 1 (Current):** Both workflows supported, Beads-first recommended
**Phase 2 (2026-Q2):** Deprecation warnings for pure markdown
**Phase 3 (2026-Q4):** Markdown workflow becomes opt-in

**Migration Timeline:**
- ✅ 2025-01: Beads-first workflow implemented
- ✅ 2026-01: Decision document created (this file)
- ⏳ 2026-06: Add deprecation notices to `@idea` skill
- ⏳ 2026-12: Make markdown workflow opt-in flag

---

**Related Documentation:**
- [CLAUDE.md](../CLAUDE.md) - Quick reference
- [PROTOCOL.md](../PROTOCOL.md) - Full specification
- [`.claude/skills/feature/SKILL.md`](../.claude/skills/feature/SKILL.md) - Beads-first implementation
- [`.claude/skills/`](../.claude/skills/) - Skill definitions

---

**Decision Record:**
- **Proposed:** 2026-01-30
- **Approved:** 2026-01-30 (implicit - documented as recommendation)
- **Reviewer:** SDP maintainer
- **Status:** Active
