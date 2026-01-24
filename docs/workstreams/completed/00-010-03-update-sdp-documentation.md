---
ws_id: 00-500-03
feature: F010
status: completed
size: MEDIUM
project_id: 00
---

## WS-00-500-03: Update SDP Documentation

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- SDP PROTOCOL.md documents PP-FFF-SS naming convention
- Project ID registry is documented
- Cross-project dependencies are explained
- Migration guide for existing WS is provided

**Acceptance Criteria:**
- [x] AC1: PROTOCOL.md has PP-FFF-SS section with examples
- [x] AC2: Project ID registry table documented (00-05)
- [x] AC3: Cross-project dependency pattern explained
- [x] AC4: Migration guide for existing WS provided
- [x] AC5: All SDP command prompts reference new format

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

The new WS naming convention `PP-FFF-SS` needs to be documented across SDP. This includes the core protocol document, command prompts, and migration guides for existing projects.

### Зависимость

00-500-02 (PP-FFF-SS format must be implemented first)

### Входные файлы

- `sdp/docs/PROTOCOL.md` — core protocol document
- `sdp/prompts/commands/design.md` — design command
- `sdp/prompts/commands/build.md` — build command
- `sdp/docs/workstreams/TEMPLATE.md` — WS template

### Шаги

1. **Update PROTOCOL.md** with PP-FFF-SS section

   Add to `sdp/docs/PROTOCOL.md`:

   ```markdown
   ## Workstream Naming Convention (PP-FFF-SS)

   ### Format

   ```
   PP-FFF-SS
   ├─ PP: Project ID (2 digits, 00-99)
   ├─ FFF: Feature ID (3 digits, 000-999)
   └─ SS: Workstream Sequence (2 digits, 00-99)
   ```

   ### Project ID Registry

   | ID | Project | Description |
   |----|---------|-------------|
   | 00 | **SDP Protocol** | Universal meta-protocol (uses itself) |
   | 02 | hw_checker | Homework validation system |
   | 03 | mlsd | ML System Design course |
   | 04 | bdde | Big Data course |
   | 05 | msu_ai_masters | Meta-repo configuration |
   | 01 | *Reserved* | Available for future use |

   **Principle:** PP = who owns the workstream. All projects (02-05) use SDP (00) as their tool.

   ### Examples

   | WS ID | Project | Feature | Description |
   |-------|---------|---------|-------------|
   | 00-500-01 | SDP | F500 | Sync SDP content |
   | 00-410-01 | SDP | F410 | Contract-driven WS spec |
   | 02-150-01 | hw_checker | F150 | Config fixes |
   | 02-201-01 | hw_checker | F201 | Multi-IDE parity |
   | 03-100-01 | mlsd | F100 | Question domain |
   | 04-050-01 | bdde | F050 | Data pipeline |

   ### Cross-Project Dependencies

   Projects can depend on SDP workstreams:

   ```yaml
   # In hw_checker (02-150-03.md):
   ---
   depends_on:
     - 00-100-05  # SDP Protocol WS-100-05
   ---
   ```

   **Rule:** Projects (02-05) may depend on SDP (00), but SDP does not depend on specific projects.
   ```

2. **Create migration guide** at `sdp/docs/migration/ws-naming-migration.md`:

   ```markdown
   # WS Naming Migration Guide

   ## Old Format → New Format

   | Old Format | New Format | Example |
   |------------|------------|---------|
   | `WS-FFF-SS` | `PP-FFF-SS` | WS-193-01 → 00-193-01 |
   | `WS-FFF-SS` | `PP-FFF-SS` | WS-150-01 → 02-150-01 |

   ## Per-Project Migration

   ### SDP (Project 00)

   ```bash
   # All SDP workstreams get 00- prefix
   WS-190-01 → 00-190-01
   WS-191-01 → 00-191-01
   # etc.
   ```

   ### hw_checker (Project 02)

   ```bash
   # All hw_checker workstreams get 02- prefix
   WS-150-01 → 02-150-01
   WS-201-01 → 02-201-01
   # etc.
   ```

   ### mlsd (Project 03)

   ```bash
   # mlsd workstreams get 03- prefix
   WS-100-01 → 03-100-01
   ```

   ### bdde (Project 04)

   ```bash
   # bdde workstreams get 04- prefix
   ```

   ## Automated Migration Script

   ```bash
   cd sdp
   poetry run python scripts/migrate_ws_format.py --project-id 00
   ```

   ## Manual Updates Required

   1. Update WS file frontmatter
   2. Rename WS files
   3. Update INDEX.md references
   4. Update cross-WS dependencies
   ```

3. **Update command prompts** to reference PP-FFF-SS

   In `sdp/prompts/commands/design.md`, update WS format section:

   ```markdown
   ### 3.2 Substream Format (STRICT)

   ```
   PP-FFF-SS
   ```

   PP = Project ID (2 digits, 00-99)
   FFF = Feature ID (3 digits, 000-999)
   SS = Sequence (2 digits, 00-99)
   ```

   **✅ Правильно:** `00-060-01`, `02-060-01`, `00-060-10`
   **❌ Неправильно:** `WS-060-1`, `60-01`, `WS-060-A`
   ```

4. **Update TEMPLATE.md** with comprehensive frontmatter example

   In `sdp/docs/workstreams/TEMPLATE.md`:

   ```markdown
   ---
   ws_id: PP-FFF-SS
   feature: FFFF
   title: Short descriptive title
   status: backlog|active|completed
   size: SMALL|MEDIUM|LARGE
   project_id: PP
   capability_tier: T0|T1|T2|T3
   execution_profile: contract-driven|exploratory
   depends_on:
     - PP-FFF-SS
   github_issue: NNN
   created: YYYY-MM-DD
   updated: YYYY-MM-DD
   ---

   ## WS-PP-FFF-SS: Title

   ### Goal

   **What should work after this WS:**
   - [ ] Specific functionality
   - [ ] Measurable outcome

   **Acceptance Criteria:**
   - [ ] AC1: Testable condition
   - [ ] AC2: Testable condition

   **⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

   ---
   ```

### Ожидаемый результат

- Complete documentation of PP-FFF-SS convention
- Migration guide for existing WS
- All prompts updated with new format

### Scope Estimate

- Файлов: ~8 (PROTOCOL.md, migration guide, prompts)
- Строк: ~1200 (MEDIUM)
- Токенов: ~4000

### Критерий завершения

```bash
cd sdp

# Verify documentation updated
grep -q "PP-FFF-SS" docs/PROTOCOL.md
grep -q "Project ID Registry" docs/PROTOCOL.md
test -f docs/migration/ws-naming-migration.md

# Verify prompts updated
grep -q "PP-FFF-SS" prompts/commands/design.md
grep -q "PP-FFF-SS" prompts/commands/build.md

# Verify template updated
grep -q "project_id:" docs/workstreams/TEMPLATE.md

echo "✅ SDP documentation updated for PP-FFF-SS"
```

### Ограничения

- НЕ менять сам format PP-FFF-SS
- НЕ переименовывать существующие WS в этом WS (только документация)
- НЕ добавлять новые поля в frontmatter без обсуждения

---

## Execution Report

**Executed by:** Claude
**Date:** 2026-01-24

### Goal Status
- [x] AC1: PROTOCOL.md has PP-FFF-SS section with examples — ✅
- [x] AC2: Project ID registry table documented (00-05) — ✅
- [x] AC3: Cross-project dependency pattern explained — ✅
- [x] AC4: Migration guide for existing WS provided — ✅
- [x] AC5: All SDP command prompts reference new format — ✅

**Goal Achieved:** ✅ YES

### Files Changed

| File | Action | LOC |
|------|--------|-----|
| `sdp/PROTOCOL.md` | modified | +65 (PP-FFF-SS section) |
| `sdp/docs/migration/ws-naming-migration.md` | created | 170 |
| `sdp/prompts/commands/design.md` | modified | +20 (updated format section) |
| `sdp/prompts/commands/build.md` | modified | +5 (updated examples) |

### Self-Check Results
```bash
$ grep -q "PP-FFF-SS" sdp/PROTOCOL.md && echo "✓ PROTOCOL.md updated"
✓ PROTOCOL.md updated

$ grep -q "Project ID Registry" sdp/PROTOCOL.md && echo "✓ Registry documented"
✓ Registry documented

$ test -f sdp/docs/migration/ws-naming-migration.md && echo "✓ Migration guide exists"
✓ Migration guide exists

$ grep -q "PP-FFF-SS" sdp/prompts/commands/design.md && echo "✓ design.md updated"
✓ design.md updated

$ grep -q "PP-FFF-SS" sdp/prompts/commands/build.md && echo "✓ build.md updated"
✓ build.md updated
```

### SDP Repo Commits
- `1bf5c85` - docs(sdp): WS-00-500-03 - Update SDP documentation for PP-FFF-SS

### Verification
- PP-FFF-SS naming convention documented in PROTOCOL.md ✅
- Project ID registry (00-05) with descriptions ✅
- Cross-project dependency pattern with examples ✅
- Migration guide with per-project instructions ✅
- Command prompts updated with new format ✅
- Backward compatibility noted (WS-FFF-SS still supported) ✅

### Next Steps
- Continue with WS-00-500-04: Configure SDP Submodule
- After F500 completion: `/codereview F500`
