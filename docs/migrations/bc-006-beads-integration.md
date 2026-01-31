# Breaking Change:     content="Use REST, not GraphQL",

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0

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
