# Breaking Change: 

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0


### 2. WS-FFF-SS → PP-FFF-SS Format

#### What Changed

Workstream IDs changed from `WS-FFF-SS` to `PP-FFF-SS` format.

**Old Format:**
```
WS-193-01 (SDP workstream)
WS-150-01 (hw_checker workstream)
```

**New Format:**
```
00-193-01 (SDP = project 00)
02-150-01 (hw_checker = project 02)
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| No project context in ID | Prefix identifies project (PP) |
| Collisions across projects | Unique project IDs prevent conflicts |
| Manual tracking of which project a WS belongs to | Explicit in the ID |

#### Migration Steps

**Step 1: Determine Your Project ID**

Check `docs/PROJECT_ID_REGISTRY.md` (or create it):

```toml
# Project IDs
[projects]
sdp = "00"        # SDP itself
hw_checker = "02" # Homework checker
mlsd = "03"       # ML system
bdde = "04"       # BDDE
```

**Step 2: Run Migration Script**

```bash
# Dry run to see what will change
python scripts/migrate_workstream_ids.py --dry-run

# Migrate SDP workstreams (project 00)
python scripts/migrate_workstream_ids.py --project-id 00

# Migrate other projects
python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker
```

**Step 3: Manual Updates (if not using script)**

Update workstream frontmatter:

```yaml
---
# OLD
ws_id: WS-193-01
feature: F193

# NEW
ws_id: 00-193-01
project_id: 00
feature: F193
---
```

**Step 4: Rename Files**

```bash
# Old
WS-193-01-extension-interface.md

# New
00-193-01-extension-interface.md
```

**Step 5: Update Cross-WS Dependencies**

```yaml
---
# OLD
depends_on:
  - WS-100-05

# NEW
depends_on:
  - 00-100-05
---
```

**Step 6: Update INDEX.md References**

```markdown
<!-- OLD -->
- [WS-193-01](WS-193-01-extension-interface.md)

<!-- NEW -->
- [00-193-01](00-193-01-extension-interface.md)
```

#### Before/After Comparison

**OLD (WS-FFF-SS):**
```yaml
---
ws_id: WS-193-01
feature: F193
status: backlog
size: MEDIUM
depends_on:
  - WS-100-05
---
```

**NEW (PP-FFF-SS):**
```yaml
---
ws_id: 00-193-01
project_id: 00
feature: F193
status: backlog
size: MEDIUM
depends_on:
  - 00-100-05
---
```

#### Timeline

- **Deprecated:** 2025-11-01 (v0.2)
- **Removed:** 2025-12-01 (v0.3.0)
- **Migration Support:** Ongoing (backward compatible)

#### Validation

