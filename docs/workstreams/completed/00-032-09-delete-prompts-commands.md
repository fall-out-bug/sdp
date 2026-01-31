---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-06
- 00-032-07
- 00-032-08
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: Все файлы из `prompts/commands/` удалены
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`prompts/commands/README.md` содержит redirect notice'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: CHANGELOG.md содержит breaking change entry
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: ''
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`docs/migration/prompts-to-skills.md` создан'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-09
---

## 00-032-09: Delete prompts/commands/

### 🎯 Goal

**What must WORK after completing this WS:**
- `prompts/commands/` удалён после миграции в skills
- Redirect notice для обратной совместимости
- Migration guide обновлён

**Acceptance Criteria:**
- [ ] AC1: Все файлы из `prompts/commands/` удалены
- [ ] AC2: `prompts/commands/README.md` содержит redirect notice
- [ ] AC3: CHANGELOG.md содержит breaking change entry
- [ ] AC4: `docs/migration/prompts-to-skills.md` создан

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: После миграции в skills, `prompts/commands/` дублирует информацию.

**Solution**: Удалить с redirect notice.

### Dependencies

- **00-032-06**: Merge build.md
- **00-032-07**: Merge review.md  
- **00-032-08**: Merge design.md

### Steps

1. **Create redirect notice**

   ```markdown
   # prompts/commands/README.md
   # ⚠️ DEPRECATED: Moved to Skills
   
   All command prompts have been migrated to skills in `.claude/skills/`.
   
   ## Migration
   
   | Old Location | New Location |
   |--------------|--------------|
   | `prompts/commands/build.md` | `.claude/skills/build/SKILL.md` |
   | `prompts/commands/review.md` | `.claude/skills/review/SKILL.md` |
   | `prompts/commands/design.md` | `.claude/skills/design/SKILL.md` |
   | `prompts/commands/idea.md` | `.claude/skills/idea/SKILL.md` |
   | ... | ... |
   
   ## Why?
   
   1. Single source of truth (skills only)
   2. Shorter prompts (≤100 lines)
   3. Better agent compliance
   
   ## Full Documentation
   
   Detailed specifications moved to `docs/reference/`:
   - `docs/reference/build-spec.md`
   - `docs/reference/review-spec.md`
   - `docs/reference/design-spec.md`
   
   See [Migration Guide](../docs/migration/prompts-to-skills.md)
   ```

2. **Delete command files**

   ```bash
   rm prompts/commands/build.md
   rm prompts/commands/review.md
   rm prompts/commands/design.md
   rm prompts/commands/idea.md
   rm prompts/commands/deploy.md
   rm prompts/commands/oneshot.md
   rm prompts/commands/issue.md
   rm prompts/commands/hotfix.md
   rm prompts/commands/bugfix.md
   rm prompts/commands/test.md
   rm prompts/commands/prd.md
   ```

3. **Update CHANGELOG**

   ```markdown
   # CHANGELOG.md
   
   ## [0.6.0] - 2026-XX-XX
   
   ### Breaking Changes
   
   - **prompts/commands/ deprecated** — All command prompts moved to `.claude/skills/`
     - Migration: Use skills directly (`@build`, `@review`, etc.)
     - See: `docs/migration/prompts-to-skills.md`
   
   ### Added
   
   - Guard skill for pre-edit enforcement
   - Skill template standard (≤100 lines)
   - Traceability check in review
   ```

4. **Create migration guide**

   ```markdown
   # docs/migration/prompts-to-skills.md
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
   ```

### Output Files

- `prompts/commands/README.md` (redirect notice)
- `CHANGELOG.md` (updated)
- `docs/migration/prompts-to-skills.md`

### Completion Criteria

```bash
# Only README remains
ls prompts/commands/
# Expected: README.md only

# Migration guide exists
test -f docs/migration/prompts-to-skills.md
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
