---
ws_id: 00-500-02
feature: F010
status: completed
size: MEDIUM
project_id: 00
---

## WS-00-500-02: Update SDP to Use PP-FFF-SS Naming

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- SDP workstream parser supports `PP-FFF-SS` format (2-digit project, 3-digit feature, 2-digit sequence)
- SDP WS template updated with new frontmatter fields
- Project ID 00 reserved for SDP Protocol itself
- Existing SDP WS files updated to new format

**Acceptance Criteria:**
- [x] AC1: Workstream parser accepts `PP-FFF-SS` format
- [x] AC2: WS template includes `project_id` field in frontmatter
- [x] AC3: SDP workstreams use new WS format (00-XXX-YY)
- [x] AC4: All existing SDP WS files renamed/updated to new format
- [x] AC5: Tests validate new format (00-500-01, 00-500-02, etc.)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

The unified WS naming convention uses `PP-FFF-SS` where:
- **PP** = Project ID (00=SDP, 02=hw_checker, 03=mlsd, 04=bdde, 05=meta)
- **FFF** = Feature ID (3 digits)
- **SS** = Workstream sequence (2 digits)

SDP itself uses Project ID 00, so its workstreams will be `00-XXX-YY`.

This WS updates SDP's own workstream parsing and templates to support the new format.

### Зависимость

00-500-01 (SDP repo must have content first)

### Входные файлы

- `sdp/src/sdp/core/workstream.py` — current WS parser
- `sdp/docs/workstreams/TEMPLATE.md` — current WS template
- `sdp/docs/workstreams/INDEX.md` — current WS index
- `sdp/tests/unit/test_workstream.py` — existing tests

### Шаги

1. **Update workstream parser** to support PP-FFF-SS format

   In `sdp/src/sdp/core/workstream.py`:

   ```python
   import re
   from dataclasses import dataclass
   from pathlib import Path
   from typing import Optional

   @dataclass
   class WorkstreamID:
       """Parsed workstream ID in PP-FFF-SS format."""

       project_id: int  # 00-99
       feature_id: int  # 000-999
       sequence: int  # 00-99

       def __str__(self) -> str:
           return f"{self.project_id:02d}-{self.feature_id:03d}-{self.sequence:02d}"

       @classmethod
       def parse(cls, ws_id: str) -> "WorkstreamID":
           """Parse WS ID string like '00-500-01'."""
           # Support both formats: PP-FFF-SS and WS-FFF-SS (legacy)
           pattern = r"^(\d{2})-(\d{3})-(\d{2})$"
           match = re.match(pattern, ws_id)
           if not match:
               raise ValueError(f"Invalid WS ID format: {ws_id}. Expected PP-FFF-SS")
           project_id, feature_id, sequence = match.groups()
           return cls(
               project_id=int(project_id),
               feature_id=int(feature_id),
               sequence=int(sequence)
           )

       @property
       def is_sdp(self) -> bool:
           """Check if this is an SDP Protocol workstream (Project 00)."""
           return self.project_id == 0
   ```

2. **Update WS template** with new frontmatter

   In `sdp/docs/workstreams/TEMPLATE.md`:

   ```markdown
   ---
   ws_id: PP-FFF-SS
   feature: FFFF
   status: backlog|active|completed
   size: SMALL|MEDIUM|LARGE
   project_id: PP
   depends_on:
     - PP-FFF-SS  # optional
   ---

   ## WS-PP-FFF-SS: Title

   ### Goal
   ...
   ```

3. **Update SDP workstreams** to use Project ID 00

   Existing SDP workstreams need prefix `00-`:
   - WS-190-01 → 00-190-01
   - WS-191-01 → 00-191-01
   - WS-192-01 → 00-192-01
   - WS-193-01 → 00-193-01
   - WS-194-01 → 00-194-01

4. **Create migration script** for renaming WS files:

   ```python
   # sdp/scripts/migrate_ws_format.py
   import re
   from pathlib import Path

   def migrate_ws_file(ws_file: Path) -> Path:
       """Rename WS file from old to new format."""
       # Read frontmatter
       content = ws_file.read_text()

       # Extract ws_id
       match = re.search(r'ws_id:\s*(WS-[\d-]+)', content)
       if not match:
           print(f"⚠️  No ws_id found in {ws_file}")
           return ws_file

       old_id = match.group(1)

       # Convert to new format (add 00- prefix for SDP)
       if old_id.startswith("WS-"):
           new_id = f"00-{old_id[3:]}"  # WS-190-01 → 00-190-01
       else:
           new_id = old_id

       # Update frontmatter
       content = content.replace(f"ws_id: {old_id}", f"ws_id: {new_id}")
       content = content.replace(f"## {old_id}:", f"## {new_id}:")

       # Add project_id if missing
       if "project_id:" not in content:
           content = re.sub(
               r"(ws_id:.*\n)",
               rf"\1project_id: 00\n",
               content
           )

       # Write updated content
       ws_file.write_text(content)

       # Rename file
       new_filename = ws_file.name.replace(old_id.replace("WS-", ""), new_id.replace("-", "-")[:8])
       new_path = ws_file.parent / new_filename
       ws_file.rename(new_path)

       print(f"✓ {ws_file.name} → {new_filename}")
       return new_path
   ```

5. **Add tests** for new format

   In `sdp/tests/unit/test_workstream.py`:

   ```python
   import pytest
   from sdp.core.workstream import WorkstreamID

   def test_parse_sdp_ws_id():
       """Test parsing SDP workstream IDs (Project 00)."""
       ws = WorkstreamID.parse("00-500-01")
       assert ws.project_id == 0
       assert ws.feature_id == 500
       assert ws.sequence == 1
       assert ws.is_sdp is True
       assert str(ws) == "00-500-01"

   def test_parse_hw_checker_ws_id():
       """Test parsing hw_checker workstream IDs (Project 02)."""
       ws = WorkstreamID.parse("02-150-01")
       assert ws.project_id == 2
       assert ws.feature_id == 150
       assert ws.sequence == 1
       assert ws.is_sdp is False

   def test_invalid_ws_id_format():
       """Test rejection of invalid WS ID formats."""
       with pytest.raises(ValueError):
           WorkstreamID.parse("WS-190-01")  # Old format

       with pytest.raises(ValueError):
           WorkstreamID.parse("0-500-1")  # Wrong digit counts

   def test_project_id_registry():
       """Test Project ID registry validation."""
       valid_projects = {0, 2, 3, 4, 5}  # 00, 02, 03, 04, 05
       # Project 01 is reserved but unused

       for project_id in valid_projects:
           ws = WorkstreamID(project_id, 100, 1)
           assert ws.project_id in valid_projects
   ```

### Ожидаемый результат

- SDP parser handles `PP-FFF-SS` format correctly
- All SDP WS files use format `00-XXX-YY`
- Template updated for new format
- Tests validate new format

### Scope Estimate

- Файлов: ~10 (parser, template, tests, migration script)
- Строк: ~800 (MEDIUM)
- Токенов: ~3000

### Критерий завершения

```bash
cd sdp

# Run tests
pytest tests/unit/test_workstream.py -v

# Verify no legacy WS-XXX-YY format remains
! grep -r "ws_id: WS-" docs/workstreams/

# Verify all WS use PP-FFF-SS
grep -r "ws_id: [0-9][0-9]-" docs/workstreams/ | wc -l  # Should be > 0

# Verify Project ID 00 for SDP
grep "project_id: 00" docs/workstreams/backlog/*.md | wc -l  # Should match SDP WS count
```

### Ограничения

- НЕ менять сам format PP-FFF-SS (уже утвержден)
- НЕ переименовывать WS других проектов (только SDP/00)
- НЕ нарушать backward compatibility для парсера (должен понимать оба формата)

---

## Execution Report

**Executed by:** Claude
**Date:** 2026-01-24

### Goal Status
- [x] AC1: Workstream parser accepts PP-FFF-SS format — ✅
- [x] AC2: WS template includes project_id field in frontmatter — ✅
- [x] AC3: SDP workstreams use new WS format (00-XXX-YY) — ✅
- [x] AC4: All existing SDP WS files renamed/updated to new format — ✅
- [x] AC5: Tests validate new format (00-500-01, 00-500-02, etc.) — ✅

**Goal Achieved:** ✅ YES

### Files Changed

| File | Action | LOC |
|------|--------|-----|
| `sdp/src/sdp/core/workstream.py` | modified | +100 (WorkstreamID class) |
| `sdp/src/sdp/core/__init__.py` | modified | +2 (export WorkstreamID) |
| `sdp/templates/workstream-frontmatter.md` | created | 120 |
| `sdp/scripts/migrate_sdp_ws.py` | created | 95 |
| `sdp/tests/unit/core/test_workstream.py` | modified | +120 (10 new tests) |
| `sdp/docs/workstreams/backlog/*.md` | renamed/updated | 4 files migrated |

### Self-Check Results
```bash
$ cd /tmp/sdp && python -m pytest tests/unit/core/test_workstream.py -v
===== 18 passed in 0.17s =====

$ cd /tmp/sdp && python -m pytest tests/unit/core/ -q
===== 69 passed in 0.26s =====

$ grep "ws_id: 00-" docs/workstreams/backlog/*.md | wc -l
===== 4 (all WS migrated) =====

$ grep "ws_id: WS-" docs/workstreams/backlog/*.md | wc -l
===== 0 (no legacy format) =====
```

### SDP Repo Commits
- `8e8fd2d` - feat(sdp): WS-00-500-02 - Add PP-FFF-SS workstream ID parsing
- `f3c246f` - feat(sdp): WS-00-500-02 - Migrate SDP workstreams to PP-FFF-SS format

### Verification
- WorkstreamID class parses PP-FFF-SS format ✅
- Legacy WS-FFF-SS format still supported (backward compatibility) ✅
- Project ID validation against registry {0, 2, 3, 4, 5} ✅
- Type checkers: is_sdp, is_hw_checker, is_mlsd, is_bdde, is_meta_repo ✅
- Migration script created and tested ✅
- All 4 SDP workstreams migrated to 00-XXX-YY format ✅
- Template with project_id field created ✅

### Next Steps
- Sync SDP submodule changes to msu_ai_masters
- Continue with WS-00-500-03: Update SDP Documentation
