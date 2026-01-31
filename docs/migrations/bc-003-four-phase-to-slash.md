# Breaking Change: 

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0


After migration, verify:

```bash
# Check for remaining legacy format
grep -r "ws_id: WS-" docs/workstreams/

# Should return empty (all migrated)

# Verify new format
grep -r "project_id:" docs/workstreams/

# Should show all files with project_id
```

---

### 3. 4-Phase → Slash Commands

#### What Changed

The **4-Phase Workflow** (Analyze, Plan, Execute, Review) was replaced with **Slash Commands**.

**Old Agent Chain (4-Phase):**
```
Analyst → Architect → TechLead → Developer → QA → DevOps
```

**New Skill System (Slash Commands):**
```
@idea → @design → @build → @review → @deploy
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Fixed agent chain doesn't match real workflows | Skills are composable |
| Every epic requires full chain (even bug fixes) | Different commands for different tasks |
| No progressive disclosure | Start with @feature, expand as needed |
| Agents defined in separate ADRs | Skills defined in `.claude/skills/` |

#### Migration Steps

**Step 1: Map Old Phases to New Commands**

| Old Phase | New Command | Description |
|-----------|-------------|-------------|
| `Analyze` | `/idea` or `@idea` | Interactive requirements gathering |
| `Plan` | `/design` or `@design` | Workstream decomposition |
| `Execute` | `/build` or `@build` | Single workstream execution |
| `Review` | `/review` or `@review` | Quality check |

**Step 2: Remove Old Phase Directories**

```bash
# OLD: 4-phase structure
docs/specs/epic-auth/
├── analyze/      # ❌ Remove
├── plan/         # ❌ Remove
├── execute/      # ❌ Remove
└── review/       # ❌ Remove

# NEW: workstream-based structure
docs/
├── drafts/
│   └── idea-auth.md
└── workstreams/
    ├── backlog/
    └── completed/
```

**Step 3: Convert Phase Artifacts to Workstreams**

Extract information from phase directories:

```bash
# OLD: docs/specs/epic-auth/analyze/requirements.md
# → NEW: docs/drafts/idea-auth.md (created by @idea)

# OLD: docs/specs/epic-auth/plan/implementation.md
# → NEW: docs/workstreams/backlog/00-AUTH-*.md (created by @design)

# OLD: docs/specs/epic-auth/execute/WS-01.md
# → NEW: docs/workstreams/completed/00-AUTH-01.md (created by @build)
```

**Step 4: Update Agent Instructions**

Old agent prompts are now skills:

```bash
# OLD: prompts/commands/analyst.md
# NEW: .claude/skills/idea/SKILL.md

# OLD: prompts/commands/architect.md
# NEW: .claude/skills/design/SKILL.md

# OLD: prompts/commands/developer.md
# NEW: .claude/skills/build/SKILL.md
```

**Step 5: Update Git Hooks**

Old hooks checked phase transitions:

```bash
# OLD: hooks/pre-phase-transition.sh
# ❌ Removed

# NEW: hooks/pre-build.sh
# ✅ Validates before @build
```

#### Before/After Comparison

**OLD (4-Phase Workflow):**
```bash
# 1. Analyst (analyze phase)
cp templates/phase-analyze.md docs/specs/epic-auth/analyze/instructions.md
# Claude reads instructions, creates requirements.md

# 2. Architect (plan phase)
cp templates/phase-plan.md docs/specs/epic-auth/plan/instructions.md
# Claude reads requirements, creates architecture.md

# 3. Tech Lead (plan phase)
cp templates/phase-techlead.md docs/specs/epic-auth/plan/instructions-tl.md
# Claude creates implementation.md with workstreams

# 4. Developer (execute phase)
cp templates/phase-execute.md docs/specs/epic-auth/execute/instructions.md
# Claude implements workstreams
```

**NEW (Slash Commands):**
```bash
# 1. Interactive requirements
@idea "Add user authentication"
# Claude asks questions, creates docs/drafts/idea-auth.md

# 2. Plan workstreams
@design idea-auth
# Claude explores codebase, creates docs/workstreams/backlog/00-AUTH-*.md

# 3. Execute workstream
@build 00-AUTH-01
# Claude follows TDD cycle, moves to completed/
```

#### Timeline
