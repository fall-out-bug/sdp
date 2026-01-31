# Breaking Change: Consensus → Slash Commands

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v1.2 → v0.3.0

### 1. Consensus → Slash Commands (v1.2 → v0.3.0)

#### What Changed

The **Consensus Protocol** (v1.2) was replaced with **Slash Commands** (v0.3.0).

**Old Workflow (Consensus v1.2):**
```
Analyze → Plan → Execute → Review
```

**New Workflow (Slash Commands v0.3.0):**
```
/idea → /design → /build → /review → /deploy
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Complex 4-phase workflow required understanding entire protocol | Progressive disclosure: commands scale from simple to complex |
| State scattered across multiple files (`status.json`, artifacts) | Single source of truth in workstream files |
| Required reading 200+ line docs to start | `@feature` provides 5-min interactive interview |
| Rigid agent chain (Analyst→Architect→TechLead→Developer) | Flexible skill-based system |

#### Migration Steps

**Step 1: Update Your Mental Model**

Old concepts → New concepts:
- `Analyze phase` → `/idea` skill (interactive requirements)
- `Plan phase` → `/design` skill (workstream planning)
- `Execute phase` → `/build` skill (single workstream)
- `Review phase` → `/review` skill (quality check)

**Step 2: Migrate Your Epics**

For each epic in `docs/specs/`:

```bash
# OLD (Consensus v1.2)
docs/specs/epic-auth/
├── epic.md
├── consensus/
│   ├── status.json          # ❌ Remove
│   ├── artifacts/           # ❌ Remove
│   └── messages/            # ❌ Remove
└── implementation.md        # ❌ Remove

# NEW (Slash Commands v0.3.0)
docs/
├── drafts/
│   └── idea-auth.md         # ✅ /idea output
└── workstreams/
    └── backlog/
        ├── 00-AUTH-01.md   # ✅ /design output
        ├── 00-AUTH-02.md
        └── 00-AUTH-03.md
```

**Step 3: Convert status.json to Workstream Files**

Extract state from `status.json`:

```python
# OLD: consensus/status.json
{
  "epic_id": "EP-AUTH",
  "phase": "implementation",
  "workstreams": [
    {"id": "WS-01", "title": "Domain model", "status": "done"},
    {"id": "WS-02", "title": "Use cases", "status": "in_progress"}
  ]
}

# NEW: docs/workstreams/completed/00-AUTH-01.md
---
ws_id: 00-AUTH-01
feature: F001
status: completed
size: MEDIUM
---

# Domain Model

## Description
Define user and role entities...

## Acceptance Criteria
- [x] User entity with email/password
- [x] Role entity with permissions

## Execution Report
Completed: 2026-01-15
Coverage: 85%
```

**Step 4: Update Agent Prompts**

Old agent prompts are now **skills**:

```bash
# OLD: consensus/prompts/analyst.md
# ❌ Removed

# NEW: .claude/skills/idea/SKILL.md
# ✅ Interactive requirements gathering
```

**Step 5: Update Documentation Links**

Search your codebase for:
- `consensus/status.json` → Remove or replace with workstream files
- `consensus/artifacts/` → Replace with `docs/workstreams/`
- `prompts/structured/` → Replace with `.claude/skills/`

#### Before/After Comparison

**OLD (Consensus v1.2):**
```bash
# 1. Create epic
mkdir -p docs/specs/epic-auth/consensus
echo "# User Authentication" > docs/specs/epic-auth/epic.md

# 2. Initialize status
cat > docs/specs/epic-auth/consensus/status.json << EOF
{
  "epic_id": "EP-AUTH",
  "phase": "requirements",
  "mode": "full"
}
EOF

# 3. Run analyst agent (manual process)
# 4. Run architect agent (manual process)
# 5. Run tech lead agent (manual process)
```

**NEW (Slash Commands v0.3.0):**
```bash
# 1. Interactive requirements (5-min interview)
@idea "Add user authentication"
# → Creates: docs/drafts/idea-auth.md

# 2. Plan workstreams
@design idea-auth
# → Creates: docs/workstreams/backlog/00-AUTH-*.md

# 3. Execute first workstream
@build 00-AUTH-01
# → Moves to: docs/workstreams/completed/
```

#### Timeline

- **Deprecated:** 2025-12-01 (v1.2)
- **Removed:** 2026-01-01 (v0.3.0)
- **Migration Support:** Ends 2026-06-01

---

