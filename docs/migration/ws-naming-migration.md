# WS Naming Migration Guide

## Old Format → New Format

| Old Format | New Format | Example |
|------------|------------|---------|
| `WS-FFF-SS` | `PP-FFF-SS` | WS-193-01 → 00-193-01 |
| `WS-FFF-SS` | `PP-FFF-SS` | WS-150-01 → 02-150-01 |

## Per-Project Migration

### SDP (Project 00)

All SDP workstreams get `00-` prefix:

```bash
WS-190-01 → 00-190-01
WS-191-01 → 00-191-01
WS-192-01 → 00-192-01
WS-193-01 → 00-193-01
WS-194-01 → 00-194-01
WS-410-01 → 00-410-01
# etc.
```

**Automated migration:**
```bash
cd sdp
python scripts/migrate_sdp_ws.py
```

### hw_checker (Project 02)

All hw_checker workstreams get `02-` prefix:

```bash
WS-001-01 → 02-001-01
WS-150-01 → 02-150-01
WS-201-01 → 02-201-01
# etc.
```

**Script for hw_checker:**
```bash
cd tools/hw_checker
python ../sdp/scripts/migrate_ws_format.py --project-id 02
```

### mlsd (Project 03)

```bash
WS-100-01 → 03-100-01
WS-110-01 → 03-110-01
```

### bdde (Project 04)

```bash
WS-050-01 → 04-050-01
```

### msu_ai_masters Meta-Repo (Project 05)

```bash
WS-500-01 → 05-500-01
WS-501-01 → 05-501-01
```

## Frontmatter Updates

### Old Format

```yaml
---
ws_id: WS-193-01
feature: F193
status: backlog
size: MEDIUM
---
```

### New Format

```yaml
---
ws_id: 00-193-01
feature: F193
status: backlog
size: MEDIUM
project_id: 00
---
```

Note the addition of `project_id: PP` field.

## File Renaming

Files should be renamed to match the new format:

```bash
# Old
WS-193-01-extension-interface.md

# New
00-193-01-extension-interface.md
```

## Manual Updates Required

1. **Update WS file frontmatter**
   - Add `project_id: PP` field
   - Update `ws_id` format

2. **Rename WS files**
   - Match filename to new ws_id format

3. **Update INDEX.md references**
   - Change WS-XXX-YY → PP-XXX-YY

4. **Update cross-WS dependencies**
   - Add project prefix to dependency IDs

## Cross-WS Dependencies

When updating dependencies, include the project ID:

```yaml
# Old format
depends_on:
  - WS-100-05

# New format (if dependency is in SDP)
depends_on:
  - 00-100-05
```

## Validation

After migration, validate:

```bash
# Check for remaining legacy format
grep -r "ws_id: WS-" docs/workstreams/

# Should return empty (all migrated)

# Verify new format
grep -r "project_id:" docs/workstreams/

# Should show all files with project_id
```

## Backward Compatibility

The SDP parser supports both formats:

- `PP-FFF-SS` (new format) → Explicit project ID
- `WS-FFF-SS` (legacy format) → Implicitly Project 00 (SDP)

Legacy format is interpreted as SDP workstreams for compatibility.

## Tools

- **SDP Parser:** `sdp.src.sdp.core.WorkstreamID` - validates and parses both formats
- **Migration Script:** `sdp/scripts/migrate_sdp_ws.py` - automated migration for SDP
- **Generic Migration:** `sdp/scripts/migrate_ws_format.py --project-id NN` - for any project

## Questions?

- See [SDP PROTOCOL.md](../PROTOCOL.md) for full specification
- See [PROJECT_ID_REGISTRY.md](../../docs/PROJECT_ID_REGISTRY.md) for project registry
