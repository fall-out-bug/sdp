# Workstream ID Format Deprecation Notice

**Status:** Legacy format deprecated as of SDP v0.5.0 (2026-01-29)

---

## Summary

The legacy workstream ID format `WS-FFF-SS` is **deprecated** in favor of the standardized `PP-FFF-SS` format.

| Aspect | Legacy Format ❌ | New Format ✅ |
|--------|-----------------|---------------|
| **Format** | `WS-FFF-SS` | `PP-FFF-SS` |
| **Example** | `WS-001-01` | `00-001-01` |
| **Project ID** | Implicit (always SDP) | Explicit (00-99) |
| **Frontmatter** | `ws_id: WS-001-01` | `ws_id: 00-001-01`<br>`project_id: 00` |

---

## Why This Change?

### 1. **Multi-Project Support**

SDP is now used across multiple repositories. Explicit project IDs prevent conflicts:

```bash
# SDP (Project 00)
00-001-01  # SDP feature 001, workstream 01

# hw_checker (Project 02)
02-001-01  # hw_checker feature 001, workstream 01

# Different projects, same feature/workstream numbers → no conflict!
```

### 2. **Consistency**

All workstreams follow the same format:
- `PP` = Project ID (01-99)
- `FFF` = Feature number (001-999)
- `SS` = Workstream sequence (01-99)

### 3. **Cross-Project Dependencies**

Explicit project IDs enable dependencies across projects:

```yaml
depends_on:
  - 00-100-05  # Depends on SDP workstream
  - 02-050-01  # Depends on hw_checker workstream
```

---

## Migration Guide

### For Existing Workstreams

#### Option 1: Automated Migration (Recommended)

```bash
# Preview changes (safe, no modifications)
python scripts/migrate_workstream_ids.py --dry-run

# Migrate SDP workstreams
python scripts/migrate_workstream_ids.py --project-id 00

# Migrate other projects
python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker
```

**Features:**
- ✅ `--dry-run` mode for safe preview
- ✅ Updates frontmatter (`ws_id` and `project_id`)
- ✅ Renames files to match new format
- ✅ Updates cross-WS dependencies
- ✅ Comprehensive validation and error reporting
- ✅ Full test coverage (≥80%)

#### Option 2: Manual Migration

1. **Update frontmatter:**

```yaml
# Old
---
ws_id: WS-001-01
feature: F001
status: backlog
---

# New
---
ws_id: 00-001-01
feature: F001
status: backlog
project_id: 00
---
```

2. **Rename file:**

```bash
# Old
WS-001-01-test-workstream.md

# New
00-001-01-test-workstream.md
```

3. **Update dependencies:**

```yaml
# Old
depends_on:
  - WS-100-05

# New (if dependency is in SDP)
depends_on:
  - 00-100-05
```

---

## For New Workstreams

**Always use `PP-FFF-SS` format:**

```markdown
---
ws_id: 00-001-01
feature: F001
status: backlog
size: MEDIUM
project_id: 00
---
```

---

## Timeline

| Date | Milestone |
|------|-----------|
| 2026-01-29 | Legacy format deprecated |
| 2026-02-28 | Migration scripts available |
| 2026-06-01 | Legacy format support removed (planned) |

**Note:** SDP parser currently supports both formats for backward compatibility.

---

## Validation

After migration, verify:

```bash
# Check for remaining legacy format
grep -r 'ws_id: WS-' docs/workstreams/
# Should return empty (all migrated)

# Verify new format
grep -r 'project_id:' docs/workstreams/
# Should show all files with project_id

# Count files by format
find docs/workstreams -name 'WS-*.md' | wc -l  # Old format (should be 0)
find docs/workstreams -name '00-*.md' | wc -l  # New format
```

---

## Impact Analysis

### Files Affected

- `docs/workstreams/**/*.md` - All workstream files
- `docs/workstreams/INDEX.md` - Workstream index
- Cross-references in other workstreams

### Tools Affected

- ✅ **SDP Parser** - Supports both formats (backward compatible)
- ✅ **Migration Script** - Automated conversion available
- ✅ **Validation Scripts** - Updated to accept both formats

---

## Frequently Asked Questions

### Q: Do I need to migrate immediately?

**A:** Not immediately, but migration is recommended. The parser supports both formats, but new workstreams must use `PP-FFF-SS`.

### Q: Will old workstreams break?

**A:** No. The SDP parser interprets `WS-FFF-SS` as `00-FFF-SS` (SDP project) for backward compatibility.

### Q: How do I handle cross-project dependencies?

**A:** Use the full `PP-FFF-SS` format for all dependencies:

```yaml
depends_on:
  - 00-100-05  # SDP workstream
  - 02-050-01  # hw_checker workstream
```

### Q: What about GitHub issues?

**A:** GitHub issue titles should use the new format:

```
[00-001-01] Implement authentication flow
```

### Q: Can I use the migration script in CI/CD?

**A:** Yes! Use `--dry-run` in CI to check for legacy format:

```yaml
# .github/workflows/validate.yml
- name: Check workstream format
  run: python scripts/migrate_workstream_ids.py --dry-run
```

---

## Project ID Registry

| Project ID | Project | Repository |
|------------|---------|------------|
| 00 | SDP | `sdp` |
| 02 | hw_checker | `tools/hw_checker` |
| 03 | mlsd | `mlsd` |
| 04 | bdde | `bdde` |
| 05 | msu_ai_masters | `msu_ai_masters` |

**To register a new project:** Update this table and `docs/PROJECT_ID_REGISTRY.md`.

---

## Resources

- **Migration Script:** `scripts/migrate_workstream_ids.py`
- **Migration Guide:** `docs/migration/ws-naming-migration.md`
- **Protocol Spec:** `PROTOCOL.md` (Workstream Naming Convention section)
- **Tests:** `tests/unit/test_migrate_workstream_ids.py`

---

## Need Help?

1. Check the migration guide: `docs/migration/ws-naming-migration.md`
2. Run dry-run to preview: `python scripts/migrate_workstream_ids.py --dry-run`
3. Review test cases: `tests/unit/test_migrate_workstream_ids.py`
4. Open an issue in the repository

---

**Version:** SDP 0.5.0
**Updated:** 2026-01-29
**Author:** SDP Protocol Team
