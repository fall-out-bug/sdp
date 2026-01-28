---
name: oneshot
description: Autonomous multi-agent execution using Beads ready detection with checkpoints, resume capability, and PR-less execution modes. Executes all feature workstreams in parallel with dependency tracking.
tools: Read, Write, Edit, Bash, AskUserQuestion
version: 2.2.0-workflow-efficiency
---

# @oneshot - Multi-Agent Execution (Workflow Efficiency Integration)

Execute all workstreams for a feature using multiple agents in parallel, with Beads automatically tracking dependencies, unblocking tasks, checkpoint/resume capability, and optional PR-less execution modes for rapid iteration.

## When to Use

- Feature has multiple workstreams that can run in parallel
- Want to execute entire feature autonomously
- After `@design` has created workstreams with execution graphs
- For hands-off execution with progress tracking
- For background execution with resume capability
- For PR-less execution when speed is critical (NEW)

## Beads vs Markdown Workflow

**This skill uses Beads for task discovery and multi-agent coordination with checkpoints.**

For traditional markdown workflow, use `prompts/commands/oneshot.md` instead.

## Invocation

```bash
@oneshot bd-0001
# or with custom agent count
@oneshot bd-0001 --agents 5
# or with resume
@oneshot bd-0001 --resume abc123xyz
# or background execution
@oneshot bd-0001 --background

# NEW: Workflow Efficiency modes (F014)
@oneshot bd-0001 --auto-approve    # Skip PR, deploy directly (trusted features)
@oneshot bd-0001 --sandbox         # Skip PR, deploy to sandbox only
@oneshot bd-0001 --dry-run         # Preview changes without execution
@oneshot bd-0001 --auto-approve --dry-run  # Preview + auto-approve
```

**Environment Variables:**
- `BEADS_USE_MOCK=true` - Use mock Beads (default for dev)
- `BEADS_USE_MOCK=false` - Use real Beads CLI (requires Go + bd installed)

## Execution Modes (NEW from F014)

@oneshot supports three execution modes for different use cases:

### Standard Mode (default)
```bash
@oneshot bd-0001
```
- **Behavior:** Creates PR, requires approval before deployment
- **Use Case:** Production deployments requiring oversight
- **Workflow:** Execute workstreams → Create PR → Wait for approval → Deploy
- **Duration:** ~3h 45m (includes PR wait time)

### Auto-Approve Mode
```bash
@oneshot bd-0001 --auto-approve
```
- **Behavior:** Skips PR, deploys directly after execution
- **Use Case:** Trusted features, rapid iteration, sandbox environments
- **Workflow:** Execute workstreams → Deploy immediately
- **Duration:** ~45 min (5x faster than standard)
- **Requirements:**
  - Quality gates still enforced (coverage ≥80%, LOC <200, type hints)
  - Destructive operations require manual confirmation
  - Audit logging enabled automatically

### Sandbox Mode
```bash
@oneshot bd-0001 --sandbox
```
- **Behavior:** Skips PR, deploys to sandbox environment only
- **Use Case:** Testing without production risk
- **Workflow:** Execute workstreams → Deploy to sandbox
- **Duration:** ~45 min
- **Requirements:**
  - Target environment must be configured as "sandbox"
  - Same quality gates as auto-approve
  - No production deployment possible

### Dry-Run Mode (preview)
```bash
@oneshot bd-0001 --dry-run
@oneshot bd-0001 --auto-approve --dry-run  # Combine modes
```
- **Behavior:** Preview changes without executing
- **Shows:** Workstreams to execute, files to create/modify, destructive operations
- **Use Case:** Validate before actual execution
- **Output:** Summary of planned changes + confirmation prompt

### Mode Comparison

| Mode | PR Required | Production | Duration | Use When |
|------|-------------|------------|----------|----------|
| **Standard** | ✅ Yes | ✅ Yes | ~3h 45m | Production releases |
| **Auto-approve** | ❌ No | ✅ Yes | ~45 min | Trusted features, rapid iteration |
| **Sandbox** | ❌ No | ❌ No | ~45 min | Testing, staging |
| **Dry-run** | N/A | N/A | <1 min | Preview changes |

## Audit Logging (NEW from F014)

All `--auto-approve` executions are logged to `.sdp/audit.log`:

```json
{
  "timestamp": "2026-01-28T10:00:00Z",
  "user": "developer@example.com",
  "feature": "bd-0001",
  "mode": "auto-approve",
  "workstreams_executed": 4,
  "result": "success",
  "deployment_target": "production"
}
```

View audit logs:
```bash
sdp audit --last 10
```

## Quality Gates (Enforced in ALL modes)

Regardless of execution mode, all quality gates are enforced:

1. **Test Coverage:** ≥80% (pytest --cov-fail-under=80)
2. **File Size:** <200 LOC per file
3. **Type Hints:** mypy --strict
4. **Error Handling:** No `except: pass`
5. **Clean Architecture:** No layer violations

## Workflow

### Step 1: Read Execution Graph (NEW from ai-comm)

```python
from sdp.beads import create_beads_client, MultiAgentExecutor
from sdp.design.graph import DependencyGraph
import os

use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
client = create_beads_client(use_mock=use_mock)

# Load workstreams and build graph
graph = DependencyGraph()
feature = client.get_task(feature_id)

# Get all sub-tasks (workstreams) for this feature
all_tasks = client.list_tasks(parent_id=feature_id)

for task in all_tasks:
    metadata = task.sdp_metadata or {}
    graph.add(WorkstreamNode(
        ws_id=task.id,
        title=task.title,
        depends_on=[d.task_id for d in task.dependencies],
        oneshot_ready=metadata.get("oneshot_ready", True),
        estimated_loc=metadata.get("estimated_loc"),
        estimated_duration=metadata.get("estimated_duration"),
    ))

# Get correct execution order
execution_order = graph.topological_sort()
print(f"Execution order: {' → '.join(execution_order)}")
```

### Step 2: Initialize or Resume Checkpoint (NEW from ai-comm)

**New execution:**
```python
import json
from datetime import datetime, timezone

agent_id = f"agent-{datetime.utcnow().strftime('%Y%m%d-%H%M%S')}"
checkpoint_path = f".oneshot/{feature_id}-checkpoint.json"

checkpoint = {
    "feature": feature_id,
    "agent_id": agent_id,
    "status": "in_progress",
    "completed_ws": [],
    "current_ws": None,
    "execution_order": execution_order,
    "started_at": datetime.now(timezone.utc).isoformat(),
    "metrics": {
        "ws_total": len(execution_order),
        "ws_completed": 0,
    }
}

# Save checkpoint
os.makedirs(".oneshot", exist_ok=True)
with open(checkpoint_path, "w") as f:
    json.dump(checkpoint, f, indent=2)

print(f"📋 Checkpoint created: {checkpoint_path}")
print(f"   Agent ID: {agent_id}")
```

**Resume execution (NEW):**
```python
# Read existing checkpoint
if args.get("resume"):
    with open(checkpoint_path) as f:
        checkpoint = json.load(f)

    agent_id = checkpoint["agent_id"]
    completed = checkpoint["completed_ws"]

    print(f"🔄 Resuming from checkpoint")
    print(f"   Agent ID: {agent_id}")
    print(f"   Completed: {len(completed)}/{checkpoint['metrics']['ws_total']} WS")

    # Skip already completed workstreams
    execution_order = [ws for ws in execution_order if ws not in completed]
```

### Step 3: Initialize Multi-Agent Executor

```python
# Get agent count from argument or default
num_agents = args.get("agents", 3)

executor = MultiAgentExecutor(client, num_agents=num_agents)
```

### Step 4: Execute Feature (MERGED)

```python
# Execute all workstreams for feature
result = executor.execute_feature(
    feature_id,
    checkpoint=checkpoint,  # NEW
    mock_success=not args.get("debug", False)
)

if result.success:
    print(f"✅ Feature complete! Executed {result.total_executed} workstreams")

    # Update checkpoint (NEW)
    checkpoint["status"] = "completed"
    checkpoint["current_ws"] = None
    checkpoint["metrics"]["ws_completed"] = result.total_executed
    checkpoint["completed_at"] = datetime.now(timezone.utc).isoformat()

    with open(checkpoint_path, "w") as f:
        json.dump(checkpoint, f, indent=2)

    print(f"📋 Checkpoint updated: {checkpoint_path}")
else:
    print(f"❌ Execution failed: {result.error}")
    print(f"   Failed tasks: {result.failed_tasks}")

    # Update checkpoint with failure state (NEW)
    checkpoint["status"] = "failed"
    checkpoint["failed_tasks"] = result.failed_tasks
    checkpoint["error"] = result.error

    with open(checkpoint_path, "w") as f:
        json.dump(checkpoint, f, indent=2)

    print(f"📋 Checkpoint saved. Resume with: @oneshot {feature_id} --resume {agent_id}")
```

### Step 5: Monitor Progress

The executor automatically:
1. Discovers ready tasks via `get_ready_tasks()`
2. Executes ready tasks in parallel (up to `num_agents`)
3. Updates checkpoint after each WS completion (NEW)
4. Waits for completion
5. Repeats until no tasks remain
6. Reports summary

### Step 6: Two-Stage Review (NEW from ai-comm)

After all WS complete:

**Stage 1: Automated Review**
```bash
@review {feature_id}
```

**Stage 2: Human UAT**
- Manual testing (5-10 min)
- Approval for deploy

## Checkpoint Format (NEW)

```json
{
  "feature": "bd-0001",
  "agent_id": "agent-20260126-120000",
  "status": "in_progress",
  "completed_ws": ["bd-0001.1", "bd-0001.2"],
  "current_ws": "bd-0001.3",
  "execution_order": ["bd-0001.1", "bd-0001.2", "bd-0001.3"],
  "started_at": "2026-01-26T12:00:00Z",
  "metrics": {
    "ws_total": 3,
    "ws_completed": 2
  },
  "completed_at": null,
  "failed_tasks": [],
  "error": null
}
```

## Execution Flow

### Example: Feature with 5 Workstreams (MERGED)

**Initial state:**
```
bd-0001 (Feature)
├── bd-0001.1 (Domain) [READY, oneshot_ready=true]
├── bd-0001.2 (Repository) [READY, oneshot_ready=true]
├── bd-0001.3 (Service) [BLOCKED by bd-0001.2, oneshot_ready=true]
├── bd-0001.4 (API) [READY, oneshot_ready=true]
└── bd-0001.5 (Tests) [BLOCKED by bd-0001.4, oneshot_ready=false]  # Manual UAT required
```

**Checkpoint initialized:**
```json
{
  "feature": "bd-0001",
  "agent_id": "agent-20260126-120000",
  "status": "in_progress",
  "completed_ws": [],
  "execution_order": ["bd-0001.1", "bd-0001.2", "bd-0001.3", "bd-0001.4"],
  "metrics": {"ws_total": 4, "ws_completed": 0}
}
```

**Round 1 (3 agents):**
- Agent 1: Executes bd-0001.1
- Agent 2: Executes bd-0001.2
- Agent 3: Executes bd-0001.4

**After Round 1 (checkpoint updated):**
```json
{
  "completed_ws": ["bd-0001.1", "bd-0001.2", "bd-0001.4"],
  "current_ws": "bd-0001.3",
  "metrics": {"ws_total": 4, "ws_completed": 3}
}
```

```
bd-0001.1 [CLOSED] ✅
bd-0001.2 [CLOSED] ✅
bd-0001.3 [READY - unblocked by bd-0001.2]
bd-0001.4 [CLOSED] ✅
bd-0001.5 [READY - unblocked by bd-0001.4, oneshot_ready=false]  # Stop for manual UAT
```

**Round 2 (3 agents):**
- Agent 1: Executes bd-0001.3
- Agent 2: Skips bd-0001.5 (oneshot_ready=false)
- Agent 3: (idle - no ready tasks)

**After Round 2 (final checkpoint):**
```json
{
  "status": "completed",
  "completed_ws": ["bd-0001.1", "bd-0001.2", "bd-0001.3", "bd-0001.4"],
  "metrics": {"ws_total": 4, "ws_completed": 4},
  "completed_at": "2026-01-26T12:15:00Z"
}
```

```
All autonomous workstreams [CLOSED] ✅
bd-0001.5 [READY] - Awaiting manual UAT
Feature autonomous execution complete!
```

## Output

**Success:**
```
✅ Feature complete! Executed 4 workstreams

Execution summary:
- Workstreams executed: 4
- Agents used: 3
- Rounds: 2
- Duration: ~15 min
- Checkpoint: .oneshot/bd-0001-checkpoint.json

Remaining (manual UAT):
- bd-0001.5: Tests
```

**Failure (with resume):**
```
❌ Execution failed: 2 tasks failed

Failed tasks:
- bd-0001.3: Service layer
- bd-0001.4: API endpoints

Checkpoint saved: .oneshot/bd-0001-checkpoint.json

Resume with:
  @oneshot bd-0001 --resume agent-20260126-120000
```

## Background Execution (NEW)

```bash
# Start in background
@oneshot bd-0001 --background

# Output:
⏳ Starting background execution...
   Task ID: xyz789
   Output: /tmp/agent_xyz789.log
   Checkpoint: .oneshot/bd-0001-checkpoint.json

# Continue working...

# Check progress
Read(/tmp/agent_xyz789.log)

# Resume if interrupted
@oneshot bd-0001 --resume agent-20260126-120000
```

## Example Session

```bash
# Decompose feature first
@idea "Add user auth"
# → bd-0001 + docs/intent/bd-0001.json

@design bd-0001
# → bd-0001.1, bd-0001.2, bd-0001.3 (sequential deps)
# + execution graph

# Execute all workstreams
@oneshot bd-0001 --agents 2

# Output:
⏳ Executing feature bd-0001 with 2 agents...
📋 Checkpoint created: .oneshot/bd-0001-checkpoint.json
   Agent ID: agent-20260126-120000
   Execution order: bd-0001.1 → bd-0001.2 → bd-0001.3

Round 1:
  🤖 Agent 1: Executing bd-0001.1 (Domain entities)
  🤖 Agent 2: Executing bd-0001.2 (Repository layer)
📋 Checkpoint updated: 2/3 WS completed

✅ bd-0001.1 complete
✅ bd-0001.2 complete

Round 2:
  🤖 Agent 1: Executing bd-0001.3 (Service layer)
📋 Checkpoint updated: 3/3 WS completed

✅ bd-0001.3 complete

✅ Feature complete! Executed 3 workstreams
📋 Checkpoint updated: .oneshot/bd-0001-checkpoint.json

Next steps:
  1. Automated review: @review bd-0001
  2. Manual UAT (5-10 min)
  3. Deploy: @deploy bd-0001
```

## Key Features

**Automatic Dependency Management:**
- No manual task ordering needed
- Beads DAG automatically tracks dependencies
- Tasks become ready as dependencies complete
- Topological sort from execution graph ensures correct order (NEW)

**Parallel Execution:**
- Independent tasks run in parallel
- Configurable agent count (1-10)
- ThreadPoolExecutor for concurrency

**Progress Tracking:**
- Real-time status updates
- Ready task discovery between rounds
- Comprehensive error reporting
- Checkpoint after each WS (NEW)

**Fault Tolerance (NEW):**
- Checkpoint/resume capability
- Failed tasks don't block other tasks
- Continues until no tasks ready
- Reports all failures at end
- Resume from interruption

**Two-Stage Review (NEW):**
- Automated review via @review
- Human UAT for final validation
- Clear approval workflow

## Benefits vs Manual Execution

| Aspect | Manual @build | @oneshot (Beads + ai-comm) |
|--------|---------------|---------------------------|
| **Task discovery** | Manual | `bd ready` |
| **Parallelization** | Manual | Auto (3 agents) |
| **Status tracking** | Manual file moves | Auto status updates |
| **Error handling** | Manual | Auto reporting |
| **Checkpoint/Resume** | None | ✅ (NEW) |
| **Background mode** | None | ✅ (NEW) |
| **Execution ordering** | Manual | Graph-based (NEW) |
| **Time to execute 5 WS** | ~30 min | ~10 min |

## Troubleshooting

**No tasks executing:**
```bash
# Check if feature has workstreams
bd list --parent bd-0001

# Check ready tasks
bd ready

# Verify workstreams are OPEN
bd list --status open

# Check execution graph (NEW)
python -c "from sdp.design.graph import DependencyGraph; import json; graph = DependencyGraph(); print(graph.topological_sort())"
```

**Agents not utilized:**
```bash
# Check dependency graph
bd graph bd-0001

# Increase agent count if tasks are independent
@oneshot bd-0001 --agents 5
```

**Tasks failing repeatedly:**
```bash
# Check task details for errors
bd show bd-0001.3

# View task status
bd status bd-0001.3

# Reset blocked tasks to retry
bd update bd-0001.3 --status open

# Resume from checkpoint (NEW)
@oneshot bd-0001 --resume agent-20260126-120000
```

**Execution interrupted (NEW):**
```bash
# Check checkpoint
cat .oneshot/bd-0001-checkpoint.json

# Resume from checkpoint
@oneshot bd-0001 --resume agent-20260126-120000
```

**Wrong execution order (NEW):**
```python
# Verify topological sort
from sdp.design.graph import DependencyGraph

graph = DependencyGraph()
# ... load workstreams from Beads ...
print(graph.topological_sort())
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@oneshot bd-0001` | Execute feature with 3 agents |
| `@oneshot bd-0001 --agents 5` | Use 5 agents |
| `@oneshot bd-0001 --resume <id>` | Resume from checkpoint (NEW) |
| `@oneshot bd-0001 --background` | Background execution (NEW) |
| `bd ready` | List ready tasks |
| `bd graph` | Show dependency graph |
| `bd show {id}` | View task details |
| `@review {feature}` | Automated review (NEW) |

## File Structure (NEW)

```
.oneshot/
├── bd-0001-checkpoint.json    # Execution checkpoint
├── bd-0002-checkpoint.json    # Another feature
└── ...
```

---

**Version:** 2.1.0-beads-ai-comm
**Status:** Beads + AI-Comm Integration
**See Also:** `@idea`, `@design`, `@build`, `@review`
