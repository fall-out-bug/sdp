# Breaking Change: #### Timeline

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0

#### Timeline

- **Deprecated:** 2025-10-01 (v0.1)
- **Removed:** 2026-01-01 (v0.3.0)
- **Migration Support:** Ended 2026-02-01

---

### 4. State Machine → File-based

#### What Changed

The **state machine** model (`consensus/status.json`) was replaced with **file-based state**.

**OLD (State Machine):**
```json
// consensus/status.json
{
  "epic_id": "EP-AUTH",
  "phase": "implementation",
  "mode": "full",
  "blockers": [],
  "approvals": ["analyst", "architect"]
}
```

**NEW (File-based):**
```yaml
---
# docs/workstreams/backlog/00-AUTH-01.md
ws_id: 00-AUTH-01
status: backlog
size: MEDIUM
---

# Domain Model

## Description
...
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Single point of failure (status.json corruption) | Distributed state across workstream files |
| Requires locking for concurrent access | Workstreams are independent |
| Implicit state (must read status.json) | Explicit state in file location |
| Extra step to update state | State = file location |

#### Migration Steps

**Step 1: Extract State from status.json**

```python
import json
from pathlib import Path

# Read old state
status = json.loads(Path("docs/specs/epic-auth/consensus/status.json").read_text())

# Map phase to workstream locations
phase_to_location = {
    "requirements": "drafts/",
    "planning": "backlog/",
    "implementation": "in_progress/",
    "testing": "completed/",
    "done": "completed/"
}

location = phase_to_location.get(status["phase"])
print(f"Workstreams should be in: docs/workstreams/{location}")
```

**Step 2: Move Workstreams by Status**

```bash
# OLD: consensus/status.json defines state
{
  "workstreams": [
    {"id": "WS-01", "status": "done"},
    {"id": "WS-02", "status": "in_progress"},
    {"id": "WS-03", "status": "todo"}
  ]
}

# NEW: File location defines state
docs/workstreams/
├── completed/
│   └── 00-AUTH-01.md   # status: done
├── in_progress/
│   └── 00-AUTH-02.md   # status: in_progress
└── backlog/
    └── 00-AUTH-03.md   # status: backlog
```

**Step 3: Remove consensus Directory**

```bash
# After migration, remove old state files
rm -rf docs/specs/*/consensus/

# Keep only docs/drafts/ and docs/workstreams/
```

**Step 4: Update Validation Scripts**

Old scripts checked `status.json`:

```python
# OLD
def validate_phase(status_file):
    status = json.load(open(status_file))
    if status["phase"] not in PHASES:
        raise ValueError(f"Invalid phase: {status['phase']}")

# NEW
def validate_workstream(ws_file):
    frontmatter = parse_frontmatter(ws_file)
    if frontmatter["status"] not in STATUSES:
        raise ValueError(f"Invalid status: {frontmatter['status']}")
```

#### Before/After Comparison

**OLD (State Machine):**
```bash
# Check current state
cat docs/specs/epic-auth/consensus/status.json
# Output: {"phase": "implementation", "workstreams": [...]}

# Update state (requires validation)
python scripts/update_status.py --phase testing
# Validates against state machine rules
# Updates status.json atomically
```

**NEW (File-based):**
```bash
# Check current state
find docs/workstreams/ -name "*.md" -type f
# Output shows:
# - docs/workstreams/completed/00-AUTH-01.md
# - docs/workstreams/in_progress/00-AUTH-02.md
# - docs/workstreams/backlog/00-AUTH-03.md

# Update state (move file)
mv docs/workstreams/in_progress/00-AUTH-02.md \
   docs/workstreams/completed/00-AUTH-02.md
# No validation needed (state = location)
```
