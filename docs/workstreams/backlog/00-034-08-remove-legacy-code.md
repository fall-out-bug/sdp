---
ws_id: 00-034-08
feature: F034
status: READY
complexity: SMALL
project_id: "00"
---

# Workstream: Remove Legacy/Deprecated Code

**ID:** 00-034-08  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** SMALL (~500 LOC removal)

---

## Goal

Remove all deprecated/legacy code with 0-40% coverage to improve overall coverage and reduce maintenance burden.

---

## Context

Analysis identified legacy code consuming ~1,200 lines with minimal/zero coverage:

**Categories:**
1. **Deprecated re-export wrappers** (0% coverage) — Already replaced
2. **Migration scripts** (0% coverage) — One-time use
3. **Unused extensions** (14-35% coverage) — Superseded
4. **Legacy GitHub CLI** (0% coverage) — Unused
5. **Health checks** (0% coverage) — Integration-only

**Impact:** Removing ~1,200 lines → Coverage: 79% → **81%+**

---

## Scope

### In Scope — Delete These Files
✅ **Deprecated re-exports (0% coverage, ~35 lines):**
- `src/sdp/cli.py` (2 lines) — re-exports from cli/main
- `src/sdp/core/workstream.py` (4 lines) — re-exports from core/workstream/
- `src/sdp/core/feature.py` (4 lines) — re-exports from core/feature/
- `src/sdp/beads/sync.py` (3 lines) — re-exports from beads/sync/
- `src/sdp/validators/ws_completion.py` (5 lines) — re-exports from ws_completion/
- `src/sdp/core/workstream/models.py` (4 lines) — deprecated, use domain/

✅ **Migration scripts (0% coverage, ~159 lines):**
- `src/sdp/scripts/migrate_models.py` (83 lines) — one-time migration
- `src/sdp/scripts/migrate_workstream_ids.py` (73 lines) — one-time migration
- `src/sdp/scripts/__init__.py` (3 lines)

✅ **Unused GitHub CLI (0% coverage, ~74 lines):**
- `src/sdp/github/cli/__init__.py` (74 lines) — superseded by main CLI

✅ **Legacy health checks (0% coverage, ~42 lines):**
- `src/sdp/health_checks/beads.py` (42 lines) — integration-only, unused

✅ **Deprecated extensions (14% coverage, ~101 lines):**
- `src/sdp/cli_extension.py` (101 lines) — superseded by extensions/

✅ **Unused tier metrics (0% coverage, ~87 lines):**
- `src/sdp/core/tier_metrics.py` (58 lines) — unused
- `src/sdp/core/tier_promoter.py` (29 lines) — unused

✅ **Unused contract validator (0% coverage, ~33 lines):**
- `src/sdp/core/contract_validator.py` (33 lines) — unused

✅ **Unused GitHub integrations (0% coverage, ~64 lines):**
- `src/sdp/github/deploy_integration.py` (37 lines) — unused
- `src/sdp/github/design_integration.py` (27 lines) — unused

**Total removal:** ~601 lines of 0-14% coverage code

### Out of Scope — Keep But May Improve Later
❌ **Low coverage but potentially useful:**
- `cli/workstream.py` (29%) — active CLI commands
- `adapters/codex.py` (18%) — adapter for Codex platform
- `prd/*` modules (11-35%) — PRD generation features
- `github/*` modules (7-35%) — GitHub integration (partially used)

---

## Dependencies

**Depends On:**
- None (can start immediately)

**Blocks:**
- None

---

## Acceptance Criteria

- [ ] All deprecated re-export files deleted
- [ ] Migration scripts moved to `scripts/archive/` or deleted
- [ ] Unused GitHub CLI deleted
- [ ] Legacy health checks deleted
- [ ] Unused tier/contract modules deleted
- [ ] All imports updated (no broken imports)
- [ ] All tests pass after cleanup
- [ ] Coverage increased: 79% → 81%+

---

## Implementation Plan

### Task 1: Remove Deprecated Re-exports

```bash
# These files just re-export, safe to delete
rm src/sdp/cli.py
rm src/sdp/core/workstream.py
rm src/sdp/core/feature.py
rm src/sdp/beads/sync.py
rm src/sdp/validators/ws_completion.py
rm src/sdp/core/workstream/models.py
```

Update any imports from these files to use direct imports.

### Task 2: Archive Migration Scripts

```bash
# Move to archive (may need for historical reference)
mkdir -p scripts/archive
mv src/sdp/scripts/migrate_models.py scripts/archive/
mv src/sdp/scripts/migrate_workstream_ids.py scripts/archive/
rm src/sdp/scripts/__init__.py
```

### Task 3: Remove Unused Modules

```bash
# Unused GitHub CLI
rm -rf src/sdp/github/cli/

# Unused health checks
rm src/sdp/health_checks/beads.py

# Deprecated extension CLI
rm src/sdp/cli_extension.py

# Unused tier metrics
rm src/sdp/core/tier_metrics.py
rm src/sdp/core/tier_promoter.py

# Unused contract validator
rm src/sdp/core/contract_validator.py

# Unused GitHub integrations
rm src/sdp/github/deploy_integration.py
rm src/sdp/github/design_integration.py
```

### Task 4: Update Imports

Search for imports from deleted files:
```bash
rg "from sdp\.cli import" src/
rg "from sdp\.core\.workstream import" src/ | grep -v "from sdp.core.workstream\."
rg "from sdp\.scripts\." src/
rg "from sdp\.github\.cli" src/
```

Update to use correct imports.

### Task 5: Verify No Broken Imports

```bash
# Check Python imports
python -c "import sdp"

# Run tests
pytest tests/ -x

# Check coverage
pytest --cov=src/sdp --cov-report=term
```

---

## DO / DON'T

### Code Removal

**✅ DO:**
- Search for imports from deleted files first
- Update imports before deleting files
- Run tests after each deletion
- Keep git history (don't use `git rm --cached`)

**❌ DON'T:**
- Delete files without checking imports
- Remove potentially useful code
- Delete without commit message explaining why

---

## Files to Delete

**Deprecated re-exports (22 lines):**
- `src/sdp/cli.py`
- `src/sdp/core/workstream.py`
- `src/sdp/core/feature.py`
- `src/sdp/beads/sync.py`
- `src/sdp/validators/ws_completion.py`
- `src/sdp/core/workstream/models.py`

**Migration scripts (159 lines):**
- `src/sdp/scripts/migrate_models.py`
- `src/sdp/scripts/migrate_workstream_ids.py`
- `src/sdp/scripts/__init__.py`

**Unused modules (420 lines):**
- `src/sdp/github/cli/__init__.py`
- `src/sdp/health_checks/beads.py`
- `src/sdp/cli_extension.py`
- `src/sdp/core/tier_metrics.py`
- `src/sdp/core/tier_promoter.py`
- `src/sdp/core/contract_validator.py`
- `src/sdp/github/deploy_integration.py`
- `src/sdp/github/design_integration.py`

**Total:** ~601 lines

---

## Expected Coverage Impact

**Before:** 79% (8,921 statements, 1,862 uncovered)  
**After:** ~81% (8,320 statements, ~1,580 uncovered)

**Calculation:**
- Remove 601 lines
- Assume 580 are uncovered (average 96% uncovered for these files)
- New total: 8,320 statements
- New uncovered: 1,282
- New coverage: (8,320 - 1,282) / 8,320 = **84.6%**

**Conservative estimate:** 81-82% (some imports may need coverage)

---

## Test Plan

### Verification
```bash
# 1. Check no broken imports
python -c "import sdp; from sdp.cli.main import app"

# 2. Run full test suite
pytest tests/ -x

# 3. Check coverage
pytest --cov=src/sdp --cov-report=term --cov-fail-under=80

# 4. Verify key imports work
python -c "
from sdp.core.workstream.parser import parse_workstream
from sdp.core.feature.loader import load_feature
from sdp.beads.sync.sync_service import sync_beads
from sdp.validators.ws_completion.verifier import WSCompletionVerifier
"
```

---

**Version:** 1.0  
**Created:** 2026-01-31
