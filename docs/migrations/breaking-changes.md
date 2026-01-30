# SDP Breaking Changes Migration Guide

**Version:** v0.5.0
**Last Updated:** 2026-01-30

## Table of Contents

- [Introduction](#introduction)
- [Breaking Changes Summary](#breaking-changes-summary)
- [Detailed Migration Guides](#detailed-migration-guides)
  - [1. Consensus → Slash Commands (v1.2 → v0.3.0)](#1-consensus--slash-commands-v12--v030)
  - [2. WS-FFF-SS → PP-FFF-SS Format](#2-ws-fff-ss--pp-fff-ss-format)
  - [3. 4-Phase → Slash Commands](#3-4-phase--slash-commands)
  - [4. State Machine → File-based](#4-state-machine--file-based)
  - [5. JSON → Message Router](#5-json--message-router)
  - [6. Beads Integration](#6-beads-integration)
  - [7. QualityGateValidator Removal](#7-qualitygatevalidator-removal)
- [Troubleshooting](#troubleshooting)

---

## Introduction

This document helps you migrate between major versions of SDP by documenting all breaking changes, their rationale, and step-by-step migration instructions.

### What Are Breaking Changes?

A **breaking change** is any modification that breaks backward compatibility, requiring manual updates to your code, configuration, or workflow. Breaking changes occur when:

- APIs are removed or renamed
- File formats change structure
- Commands are deprecated or replaced
- Workflow steps are reordered

### Why We Document Breaking Changes

- **Transparency**: Clear communication about what changed and why
- **Migration Path**: Step-by-step instructions to upgrade safely
- **Timeline**: Deprecation warnings before removal
- **Rationale**: Understanding the "why" behind changes

---

## Breaking Changes Summary

| Change | Deprecated | Removed | Migration Effort | Impact |
|--------|------------|---------|------------------|--------|
| **1. Consensus → Slash Commands** | v1.2 | v0.3.0 | High | Complete workflow redesign |
| **2. WS-FFF-SS → PP-FFF-SS** | v0.2 | v0.3.0 | Medium | All workstream IDs |
| **3. 4-Phase → Slash Commands** | v0.1 | v0.3.0 | High | Agent coordination model |
| **4. State Machine → File-based** | v1.2 | v0.3.0 | High | `status.json` removal |
| **5. JSON → Message Router** | v0.4 | v0.5.0 | Medium | Agent messaging API |
| **6. Beads Integration** | N/A | v0.5.0 | Low | Optional feature |
| **7. QualityGateValidator Removal** | v0.4.9 | v0.5.0 | Low | Code validation |

---

## Detailed Migration Guides

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

#### Timeline

- **Deprecated:** 2025-12-01 (v1.2)
- **Removed:** 2026-01-01 (v0.3.0)
- **Migration Support:** Ended 2026-03-01

---

### 5. JSON → Message Router

#### What Changed

The **JSON-based messaging** was replaced with **Message Router** system.

**OLD (JSON Messaging):**
```json
// consensus/messages/inbox/developer/message-001.json
{
  "from": "architect",
  "to": "developer",
  "subject": "API design clarification",
  "body": "Please use REST for now...",
  "timestamp": "2025-12-31T10:00:00Z"
}
```

**NEW (Message Router):**
```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="architect",
    content="Please use REST for now...",
    recipient="developer",
))
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Writing JSON files manually is error-prone | Python API with type hints |
| No message validation | Message schema enforced by code |
| Hard to send messages between agents | Router handles delivery |
| No notification system | Integrated with Telegram |

#### Migration Steps

**Step 1: Install Message Router**

```bash
# Already included in SDP v0.5.0
poetry install
```

**Step 2: Convert JSON Messages to Python**

```python
# OLD: consensus/messages/inbox/developer/message-001.json
# {
#   "from": "architect",
#   "to": "developer",
#   "subject": "API design",
#   "body": "Use REST...",
#   "timestamp": "2025-12-31T10:00:00Z"
# }

# NEW: Use Message Router
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="architect",
    content="Use REST for API endpoints. GraphQL will be considered later.",
    recipient="developer",
))
```

**Step 3: Update Agent Prompts**

Old agents read JSON messages:

```markdown
# OLD: consensus/prompts/developer.md
## Context Files
Read consensus/messages/inbox/developer/ for messages
```

New agents use SendMessageRouter:

```markdown
# NEW: .claude/skills/build/SKILL.md
## Teammate Communication
Use SendMessage tool to communicate with other agents.
```

**Step 4: Remove Message JSON Files**

```bash
# After migrating to Message Router
rm -rf consensus/messages/

# Messages now sent via Python API
```

**Step 5: Enable Telegram Notifications (Optional)**

```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

```python
from sdp.unified.notifications.telegram import TelegramNotifier

notifier = TelegramNotifier(config=config)
notifier.send(Notification(
    type=NotificationType.SUCCESS,
    message="Workstream completed successfully",
))
```

#### Before/After Comparison

**OLD (JSON Messaging):**
```bash
# Architect sends message to developer
cat > consensus/messages/inbox/developer/msg-001.json << EOF
{
  "from": "architect",
  "to": "developer",
  "body": "Use REST, not GraphQL",
  "timestamp": "$(date -Iseconds)"
}
EOF

# Developer reads messages
ls consensus/messages/inbox/developer/
# msg-001.json
```

**NEW (Message Router):**
```python
# Architect sends message
router.send_message(Message(
    sender="architect",
    content="Use REST, not GraphQL",
    recipient="developer",
))

# Message delivered automatically + Telegram notification
```

#### Timeline

- **Deprecated:** 2025-12-15 (v0.4)
- **Removed:** 2026-01-15 (v0.5.0)
- **Migration Support:** Ongoing

---

### 6. Beads Integration

#### What Changed

**Beads CLI** integration was added as an optional task tracking system.

**What is Beads?**
- Command-line task tracker
- Hash-based task IDs (bd-0001, bd-0001.1, etc.)
- Dependency DAG (task blocking)
- Ready task detection

#### Why It Changed

Beads provides:
- **Better task tracking** than manual to-do lists
- **Dependency management** between workstreams
- **Ready detection** (which tasks can be executed now)
- **Git-friendly** (JSONL storage)

#### Migration Steps

**Step 1: Install Beads CLI (Optional)**

```bash
# Beads is optional - SDP works without it
cargo install beads-cli  # or: pip install beads

# Initialize Beads repository
beads init
```

**Step 2: Create Feature in Beads**

```python
from sdp.beads import create_beads_client
from sdp.beads.models import BeadsTaskCreate, BeadsStatus

client = create_beads_client(use_mock=True)  # Set False for real Beads

# Create feature
feature = client.create_task(BeadsTaskCreate(
    title="User Authentication",
    description="Add OAuth2 login flow",
))

# Decompose into workstreams
ws1 = client.create_task(BeadsTaskCreate(
    title="Domain model",
    parent_id=feature.id,
))
ws2 = client.create_task(BeadsTaskCreate(
    title="Database schema",
    parent_id=feature.id,
))
```

**Step 3: Add Dependencies**

```python
# Add dependency: ws2 blocked by ws1
client.add_dependency(ws2.id, ws1.id, dep_type="blocks")

# Get ready tasks (tasks with no blockers)
ready = client.get_ready_tasks()
# Returns: [ws1.id] (ws2 is blocked by ws1)
```

**Step 4: Update Workstream Status**

```python
# Mark workstream as complete
client.update_task_status(ws1.id, BeadsStatus.CLOSED)

# Check ready tasks again
ready = client.get_ready_tasks()
# Returns: [ws2.id] (ws1 is done, ws2 is now ready)
```

**Step 5: Use with @oneshot (Optional)**

The `@oneshot` skill uses Beads for progress tracking:

```bash
@oneshot beads-auth
# Executes all workstreams in dependency order
# Updates Beads status after each workstream
```

#### Before/After Comparison

**WITHOUT Beads (Manual Task Tracking):**
```bash
# Track tasks manually in TODO.md
echo "- [ ] WS-01: Domain model" >> TODO.md
echo "- [ ] WS-02: Database schema" >> TODO.md

# Manually check dependencies
grep "WS-01" TODO.md
```

**WITH Beads (Automatic Task Tracking):**
```python
# Create tasks in Beads
ws1 = client.create_task(BeadsTaskCreate(title="Domain model"))
ws2 = client.create_task(BeadsTaskCreate(title="Database schema"))

# Add dependency
client.add_dependency(ws2.id, ws1.id, dep_type="blocks")

# Get ready tasks
ready = client.get_ready_tasks()
# Returns: [ws1.id] (only ws1 is ready)
```

#### Timeline

- **Introduced:** 2026-01-01 (v0.5.0)
- **Status:** Optional feature
- **Migration Support:** N/A (opt-in)

#### Configuration

Beads is **disabled by default**. Enable in `.env`:

```bash
# .env
BEADS_ENABLED=true
BEADS_DB_PATH=./beads/tasks.jsonl
```

Or use mock (for testing):

```python
client = create_beads_client(use_mock=True)
```

---

### 7. QualityGateValidator Removal

#### What Changed

The `QualityGateValidator` class was removed as dead code (P1-04).

**OLD (QualityGateValidator):**
```python
from sdp.quality import QualityGateValidator

validator = QualityGateValidator()
violations = validator.validate_file("src/module.py")
```

**NEW (Direct Validation):**
```bash
# Use hooks or CLI directly
python -m sdp.quality validate path/to/file.py

# Or use pre-commit hook
git commit  # Automatically runs hooks/pre-commit.sh
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Unused class in codebase | Remove dead code (YAGNI) |
| Validation logic duplicated | Single source of truth in hooks |
| Unclear when to use class vs CLI | Use CLI or hooks (clear interface) |

#### Migration Steps

**Step 1: Find Usages of QualityGateValidator**

```bash
grep -r "QualityGateValidator" --include="*.py" .
```

**Step 2: Replace with Direct Validation**

```python
# OLD
from sdp.quality import QualityGateValidator

validator = QualityGateValidator()
violations = validator.validate_file("src/module.py")
if violations:
    print(f"Found {len(violations)} violations")

# NEW
import subprocess

result = subprocess.run(
    ["python", "-m", "sdp.quality", "validate", "src/module.py"],
    capture_output=True,
    text=True
)
if result.returncode != 0:
    print(f"Validation failed:\n{result.stdout}")
```

**Step 3: Use Git Hooks (Recommended)**

```bash
# Copy hook to .git/hooks/
cp scripts/hooks/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# Now validation runs automatically on commit
git commit  # Validates all changed files
```

**Step 4: Update Tests**

```python
# OLD
def test_quality_gate():
    validator = QualityGateValidator()
    violations = validator.validate_file("test_file.py")
    assert len(violations) == 0

# NEW
def test_quality_gate():
    result = subprocess.run(
        ["python", "-m", "sdp.quality", "validate", "test_file.py"],
        capture_output=True,
    )
    assert result.returncode == 0
```

#### Before/After Comparison

**OLD (QualityGateValidator):**
```python
from sdp.quality import QualityGateValidator

# Initialize
validator = QualityGateValidator()

# Validate
violations = validator.validate_directory("src/")
if violations:
    validator.print_report()
    exit(1)
```

**NEW (CLI + Hooks):**
```bash
# Option 1: Direct CLI
python -m sdp.quality validate src/

# Option 2: Pre-commit hook (automatic)
git commit  # Runs validation automatically

# Option 3: Make target
make quality  # Defined in Makefile
```

#### Timeline

- **Deprecated:** 2026-01-15 (v0.4.9)
- **Removed:** 2026-01-30 (v0.5.0)
- **Migration Support:** Ended 2026-02-15

#### Impact Analysis

**Files Removed:**
- `src/sdp/quality/validator.py` (QualityGateValidator class)

**Files Kept:**
- `src/sdp/quality/models.py` (QualityGateConfig, ValidationIssue)
- `src/sdp/quality/config.py` (TOML configuration parser)
- `scripts/hooks/pre-commit.sh` (Git hook)
- `scripts/hooks/post-commit.sh` (Git hook)

**Migration Effort:** Low (replace with CLI or hooks)

---

## Troubleshooting

### Issue: Legacy Workstream IDs Still Present

**Symptom:**
```bash
grep -r "ws_id: WS-" docs/workstreams/
# Returns results (should be empty)
```

**Solution:**
```bash
# Run migration script
python scripts/migrate_workstream_ids.py --project-id 00

# Verify
grep -r "ws_id: WS-" docs/workstreams/
# Should return empty
```

---

### Issue: "consensus/status.json not found" Error

**Symptom:**
```
FileNotFoundError: consensus/status.json
```

**Solution:**
You're using old documentation. Update to slash commands:

```bash
# OLD workflow
cp templates/phase-analyze.md docs/specs/epic/analyze/

# NEW workflow
@idea "Epic description"
```

---

### Issue: QualityGateValidator Import Error

**Symptom:**
```
ImportError: cannot import name 'QualityGateValidator' from 'sdp.quality'
```

**Solution:**
Replace with CLI validation:

```python
# OLD
from sdp.quality import QualityGateValidator

# NEW
import subprocess
subprocess.run(["python", "-m", "sdp.quality", "validate", "src/"])
```

---

### Issue: Beads Integration Not Working

**Symptom:**
```
Error: Beads CLI not found
```

**Solution:**
Beads is optional. Either:

1. **Install Beads:**
```bash
cargo install beads-cli
```

2. **Use Mock (for testing):**
```python
client = create_beads_client(use_mock=True)
```

3. **Disable Beads:**
```bash
# .env
BEADS_ENABLED=false
```

---

### Issue: Message Router Not Sending Messages

**Symptom:**
```python
router.send_message(Message(...))
# No error, but message not received
```

**Solution:**

Check recipient name:

```python
# WRONG
router.send_message(Message(
    recipient="developer",  # ❌ Must be teammate name
))

# RIGHT
router.send_message(Message(
    recipient="team-lead",  # ✅ Use actual teammate name
))
```

Read team config to find correct names:

```bash
cat ~/.claude/teams/{team-name}/config.json
```

---

## Need Help?

### Documentation

- [PROTOCOL.md](../PROTOCOL.md) - Full specification
- [CLAUDE.md](../CLAUDE.md) - Claude Code integration
- [docs/SITEMAP.md](../docs/SITEMAP.md) - Documentation index

### Migration Scripts

- `scripts/migrate_workstream_ids.py` - Workstream ID migration
- `scripts/migrate_sdp_ws.py` - SDP-specific migration

### Community

- GitHub Issues: [fall-out-bug/sdp](https://github.com/fall-out-bug/sdp/issues)
- Discussions: [GitHub Discussions](https://github.com/fall-out-bug/sdp/discussions)

---

**Last Updated:** 2026-01-30
**SDP Version:** v0.5.0
