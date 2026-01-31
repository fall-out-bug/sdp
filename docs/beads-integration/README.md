# SDP + Beads Integration

> **Status:** PoC Phase
> **Worktree:** `feature/beads-integration`
> **Goal:** Replace F012 orchestrator with Beads for multi-agent task coordination

---

## What is This?

Experimental integration between **SDP** (Spec-Driven Protocol) and **Beads** (git-backed issue tracker for AI agents).

**Key Idea:** Use Beads as the backbone for task graph management, multi-agent coordination, and conflict-free concurrent execution. SDP focuses on what it does best: TDD execution, quality gates, and platform adapters.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Agentic Workflow                     │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐         ┌──────────────┐             │
│  │    Beads     │────────▶│     SDP      │             │
│  │              │         │              │             │
│  │ • Task graph │         │ • TDD        │             │
│  │ • Dependencies       │ • Quality    │             │
│  │ • State      │         │ • Skills     │             │
│  │ • Multi-agent│         │ • Execution  │             │
│  └──────────────┘         └──────────────┘             │
│       ▲                         ▲                       │
│       │                         │                       │
│       └───────────┬─────────────┘                       │
│                   │                                     │
│            Git as Storage                                │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

**Responsibilities:**

| Component | Responsibilities |
|-----------|-----------------|
| **Beads** | Task graph, dependencies, state persistence, multi-agent coordination, conflict-free IDs |
| **SDP** | TDD workflow, quality gates, platform adapters, skills (@idea, @design, @build) |

---

## Quick Start

### 1. Install Dependencies (Optional - for real Beads)

```bash
# Install Go 1.24+
# macOS
brew install go

# Install Beads
go install github.com/steveyegge/beads/cmd/bd@latest

# Initialize Beads in project
cd /path/to/project
bd init
```

### 2. Use Mock Client (Development)

```python
from sdp.beads import create_beads_client, BeadsSyncService

# Create mock client (no Beads required)
client = create_beads_client(use_mock=True)

# Create a task
from sdp.beads import BeadsTaskCreate, BeadsPriority

task = client.create_task(BeadsTaskCreate(
    title="Add user authentication",
    description="Implement email/password auth",
    priority=BeadsPriority.HIGH,
))

print(f"Created task: {task.id}")  # bd-0001
```

### 3. Use Real Beads (Production)

```python
from sdp.beads import create_beads_client, BeadsSyncService
from pathlib import Path

# Create real client (requires Beads installed)
client = create_beads_client(project_dir=Path.cwd())

# Create a task via Beads CLI
task = client.create_task(BeadsTaskCreate(
    title="Add user authentication",
    description="Implement email/password auth",
    priority=BeadsPriority.HIGH,
))

# Get ready tasks
ready = client.get_ready_tasks()
print(f"Ready to work on: {ready}")
```

---

## Multi-Agent Workflow Example

```python
from sdp.beads import create_beads_client, BeadsTaskCreate, BeadsDependency, BeadsDependencyType

# Initialize client
client = create_beads_client(use_mock=True)

# Create a feature task
feature = client.create_task(BeadsTaskCreate(
    title="User Authentication",
    priority=BeadsPriority.HIGH,
))

# Decompose into workstreams (sub-tasks)
ws1 = client.create_task(BeadsTaskCreate(
    title="Domain entities",
    parent_id=feature.id,
))

ws2 = client.create_task(BeadsTaskCreate(
    title="Repository layer",
    parent_id=feature.id,
    dependencies=[
        BeadsDependency(task_id=ws1.id, type=BeadsDependencyType.BLOCKS)
    ],
))

ws3 = client.create_task(BeadsTaskCreate(
    title="Service layer",
    parent_id=feature.id,
    dependencies=[
        BeadsDependency(task_id=ws2.id, type=BeadsDependencyType.BLOCKS)
    ],
))

# Get ready tasks (ws1 is ready, ws2/ws3 blocked by ws1)
ready = client.get_ready_tasks()
print(f"Ready: {ready}")  # [ws1.id]

# Agent 1 executes ws1
# After completion, ws2 becomes ready automatically!
ready = client.get_ready_tasks()
print(f"Ready: {ready}")  # [ws2.id]

# Agent 2 can now execute ws2 in parallel with Agent 3 working on something else
```

---

## Sync Service

Convert between SDP workstreams and Beads tasks:

```python
from sdp.beads import create_beads_client, BeadsSyncService
from pathlib import Path

client = create_beads_client(use_mock=True)
sync = BeadsSyncService(client)

# Sync workstream → Beads
result = sync.sync_workstream_to_beads(
    ws_file=Path("docs/workstreams/backlog/00-001-01.md"),
    ws_data={
        "ws_id": "00-001-01",
        "title": "Domain entities",
        "goal": "Create user domain model",
        "status": "backlog",
        "size": "MEDIUM",
        "dependencies": [],
        "acceptance_criteria": [
            {"id": "ac1", "text": "User entity exists", "checked": False},
        ],
    },
)

if result.success:
    print(f"Created Beads task: {result.beads_id}")
else:
    print(f"Error: {result.error}")
```

---

## File Structure

```
src/sdp/beads/
├── __init__.py       # Package exports
├── models.py         # Data models (BeadsTask, BeadsStatus, etc.)
├── client.py         # BeadsClient interface + implementations
│   ├── BeadsClient          # Abstract interface
│   ├── MockBeadsClient      # In-memory mock (dev/testing)
│   └── CLIBeadsClient       # Real Beads via subprocess
└── sync.py           # Bidirectional sync service
```

---

## Testing

Run tests with mock client (no Beads required):

```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp-beads-integration

# Run all tests
pytest tests/unit/beads/

# Run specific test
pytest tests/unit/beads/test_client.py -v

# Run with coverage
pytest --cov=src/sdp/beads tests/unit/beads/
```

---

## Installation (Real Beads)

If you want to use real Beads (not mock):

### Prerequisites

1. **Go 1.24+**
   ```bash
   # macOS
   brew install go

   # Verify
   go version
   ```

2. **Beads CLI**
   ```bash
   go install github.com/steveyegge/beads/cmd/bd@latest

   # Verify
   bd --version
   ```

3. **Initialize Beads in project**
   ```bash
   cd /path/to/project
   bd init

   # Creates:
   # .beads/
   # ├── beads.db          # SQLite cache
   # └── issues.jsonl      # Git-tracked tasks
   ```

### Switch from Mock to Real

```python
# Development (mock)
client = create_beads_client(use_mock=True)

# Production (real)
client = create_beads_client()  # Auto-detects Beads
```

---

## Comparison: SDP F012 vs Beads

| Feature | SDP F012 (Planned) | Beads (Existing) |
|---------|-------------------|------------------|
| **Conflict-free IDs** | Manual PP-FFF-SS | ✅ Hash-based (`bd-a3f8`) |
| **Multi-agent** | Custom orchestrator | ✅ Built-in |
| **Dependency graph** | Manual implementation | ✅ Native DAG |
| **State persistence** | Custom JSON files | ✅ SQLite + JSONL |
| **Git integration** | File-based only | ✅ Git-backed storage |
| **Concurrent creation** | Race conditions | ✅ Content addressing |
| **Ready detection** | Manual | ✅ Auto |
| **Daemon mode** | Planned | ✅ Background sync |

**Conclusion:** Beads already solves what F012 is trying to build.

---

## Next Steps

### Phase 1: Foundation (Current)

- ✅ Create BeadsClient interface
- ✅ Implement MockBeadsClient
- ✅ Implement CLIBeadsClient
- ✅ Create sync service
- 🚧 Write tests
- 🚧 Documentation

### Phase 2: Skills Integration

- 📋 Update @idea to create Beads tasks
- 📋 Update @design to create sub-task graphs
- 📋 Update @build to work with Beads IDs
- 📋 Add multi-agent @oneshot using Beads ready detection

### Phase 3: Migration

- 📋 Migrate existing workstreams to Beads
- 📋 Update all skills to use Beads IDs
- 📋 Deprecate F012 workstreams (use Beads instead)
- 📋 Performance benchmarking

---

## FAQ

**Q: Do I need Go installed?**
A: No, not for development. Use `create_beads_client(use_mock=True)`. Go is only needed for production use with real Beads.

**Q: Can I use this alongside existing SDP workstreams?**
A: Yes! The sync service converts between formats. You can migrate gradually.

**Q: What happens to F012?**
A: We're parking F012 for now. If Beads works well, we'll remove the duplicate workstreams from F012.

**Q: Is Beads required?**
A: No. SDP works standalone. Beads is an optional enhancement for multi-agent workflows.

**Q: How do I get started?**
A:
1. Use mock client for development: `create_beads_client(use_mock=True)`
2. Write tests with mock
3. Install Go + Beads when ready for production
4. Switch to real client: `create_beads_client()`

---

## Resources

- [Beads GitHub](https://github.com/steveyegge/beads)
- [Beads Documentation](https://github.com/steveyegge/beads/blob/main/docs/ADVANCED.md)
- [SDP PROTOCOL](../PROTOCOL.md)
- [Integration Design](../docs/plans/2025-01-28-beads-sdp-integration-design.md)

---

**Version:** 0.1.0 (PoC)
**Branch:** `feature/beads-integration`
**Worktree:** `/Users/fall_out_bug/projects/vibe_coding/sdp-beads-integration`
