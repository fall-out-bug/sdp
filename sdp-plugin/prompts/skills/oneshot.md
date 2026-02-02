---
name: oneshot
description: Autonomous multi-agent execution with checkpoints, resume, and PR-less modes
tools: Read, Write, Edit, Bash, AskUserQuestion
version: 3.0.0
---

# @oneshot - Multi-Agent Execution

Execute all workstreams for a feature using multiple agents in parallel with checkpoint/resume.

## When to Use

- Feature has multiple workstreams
- Want autonomous execution with progress tracking
- Need checkpoint/resume for long-running features

## Invocation

```bash
@oneshot bd-0001                    # Standard (creates PR)
@oneshot bd-0001 --agents 5         # Use 5 agents
@oneshot bd-0001 --auto-approve     # Skip PR, deploy directly
@oneshot bd-0001 --sandbox          # Deploy to sandbox only
@oneshot bd-0001 --dry-run          # Preview changes
@oneshot bd-0001 --resume <id>      # Resume from checkpoint
@oneshot bd-0001 --background       # Run in background
```

## Execution Modes

| Mode | PR | Production | Use Case |
|------|-------|------------|----------|
| **Standard** | Yes | Yes | Production releases |
| **--auto-approve** | No | Yes | Trusted features, rapid iteration |
| **--sandbox** | No | No | Testing, staging |
| **--dry-run** | N/A | N/A | Preview changes |

All modes enforce quality gates: coverage ≥80%, LOC <200, type hints, no `except: pass`.

## Workflow

### Step 1: Load Execution Graph

```python
from sdp.beads import create_beads_client
from sdp.design.graph import DependencyGraph

client = create_beads_client()
graph = DependencyGraph()

for task in client.list_tasks(parent_id=feature_id):
    graph.add(WorkstreamNode(
        ws_id=task.id,
        depends_on=[d.task_id for d in task.dependencies],
        oneshot_ready=task.sdp_metadata.get("oneshot_ready", True),
    ))

execution_order = graph.topological_sort()
```

### Step 2: Initialize Checkpoint

```python
checkpoint = {
    "feature": feature_id,
    "agent_id": f"agent-{datetime.utcnow().strftime('%Y%m%d-%H%M%S')}",
    "status": "in_progress",
    "completed_ws": [],
    "execution_order": execution_order,
}

# Save to .oneshot/{feature_id}-checkpoint.json
```

### Step 3: Execute with Multi-Agent

```python
from sdp.beads import MultiAgentExecutor

executor = MultiAgentExecutor(client, num_agents=args.get("agents", 3))
result = executor.execute_feature(feature_id, checkpoint=checkpoint)

if result.success:
    checkpoint["status"] = "completed"
else:
    checkpoint["status"] = "failed"
    print(f"Resume with: @oneshot {feature_id} --resume {checkpoint['agent_id']}")
```

### Step 4: Two-Stage Review

After all WS complete:
1. **Automated:** `@review {feature_id}`
2. **Human UAT:** Manual testing (5-10 min)

## Output

**Success:**
```
✅ Feature complete! Executed 4 workstreams
   Agents: 3, Rounds: 2, Duration: ~15 min
   Checkpoint: .oneshot/bd-0001-checkpoint.json

Next: @review bd-0001 → Manual UAT → @deploy bd-0001
```

**Failure:**
```
❌ Execution failed: bd-0001.3

Resume: @oneshot bd-0001 --resume agent-20260126-120000
```

## Checkpoint Format

```json
{
  "feature": "bd-0001",
  "agent_id": "agent-20260126-120000",
  "status": "in_progress|completed|failed",
  "completed_ws": ["bd-0001.1", "bd-0001.2"],
  "execution_order": ["bd-0001.1", "bd-0001.2", "bd-0001.3"],
  "started_at": "2026-01-26T12:00:00Z"
}
```

## Key Features

| Feature | Description |
|---------|-------------|
| **Auto dependencies** | Beads DAG tracks dependencies |
| **Parallel execution** | Independent tasks run in parallel |
| **Checkpoint/resume** | Continue from interruption |
| **Background mode** | Long-running features |
| **Audit logging** | All --auto-approve logged to .sdp/audit.log |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| No tasks executing | `bd list --parent {id}` to check workstreams |
| Agents not utilized | Increase with `--agents 5` |
| Tasks failing | `bd show {id}` for details, then resume |
| Wrong order | Check `graph.topological_sort()` |

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@oneshot bd-0001` | Execute with 3 agents |
| `@oneshot bd-0001 --resume <id>` | Resume from checkpoint |
| `bd ready` | List ready tasks |
| `bd graph {id}` | Show dependency graph |
| `@review {feature}` | Automated review |

---

**Version:** 3.0.0  
**See Also:** `@idea`, `@design`, `@build`, `@review`
