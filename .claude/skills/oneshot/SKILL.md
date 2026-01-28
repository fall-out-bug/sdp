---
name: oneshot
description: Autonomous multi-agent execution using Beads ready detection. Executes all feature workstreams in parallel with dependency tracking.
tools: Read, Write, Edit, Bash, AskUserQuestion
version: 2.0.0-beads
---

# @oneshot - Multi-Agent Execution (Beads Integration)

Execute all workstreams for a feature using multiple agents in parallel, with Beads automatically tracking dependencies and unblocking tasks.

## When to Use

- Feature has multiple workstreams that can run in parallel
- Want to execute entire feature autonomously
- After `@design` has created workstreams
- For hands-off execution with progress tracking

## Beads vs Markdown Workflow

**This skill uses Beads for task discovery and multi-agent coordination.**

For traditional markdown workflow, use `prompts/commands/oneshot.md` instead.

## Invocation

```bash
@oneshot bd-0001
# or with custom agent count
@oneshot bd-0001 --agents 5
```

**Environment Variables:**
- `BEADS_USE_MOCK=true` - Use mock Beads (default for dev)
- `BEADS_USE_MOCK=false` - Use real Beads CLI (requires Go + bd installed)

## Workflow

### Step 1: Initialize Multi-Agent Executor

```python
from sdp.beads import create_beads_client, MultiAgentExecutor
import os

use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
client = create_beads_client(use_mock=use_mock)

# Get agent count from argument or default
num_agents = args.get("agents", 3)

executor = MultiAgentExecutor(client, num_agents=num_agents)
```

### Step 2: Execute Feature

```python
# Execute all workstreams for feature
result = executor.execute_feature(feature_id)

if result.success:
    print(f"✅ Feature complete! Executed {result.total_executed} workstreams")
else:
    print(f"❌ Execution failed: {result.error}")
    print(f"   Failed tasks: {result.failed_tasks}")
```

### Step 3: Monitor Progress

The executor automatically:
1. Discovers ready tasks via `get_ready_tasks()`
2. Executes ready tasks in parallel (up to `num_agents`)
3. Waits for completion
4. Repeats until no tasks remain
5. Reports summary

## Execution Flow

### Example: Feature with 5 Workstreams

**Initial state:**
```
bd-0001 (Feature)
├── bd-0001.1 (Domain) [READY]
├── bd-0001.2 (Repository) [READY]
├── bd-0001.3 (Service) [BLOCKED by bd-0001.2]
├── bd-0001.4 (API) [READY]
└── bd-0001.5 (Tests) [BLOCKED by bd-0001.4]
```

**Round 1 (3 agents):**
- Agent 1: Executes bd-0001.1
- Agent 2: Executes bd-0001.2
- Agent 3: Executes bd-0001.4

**After Round 1:**
```
bd-0001.1 [CLOSED] ✅
bd-0001.2 [CLOSED] ✅
bd-0001.3 [READY - unblocked by bd-0001.2]
bd-0001.4 [CLOSED] ✅
bd-0001.5 [READY - unblocked by bd-0001.4]
```

**Round 2 (3 agents):**
- Agent 1: Executes bd-0001.3
- Agent 2: Executes bd-0001.5
- Agent 3: (idle - no ready tasks)

**After Round 2:**
```
All workstreams [CLOSED] ✅
Feature complete!
```

## Output

**Success:**
```
✅ Feature complete! Executed 5 workstreams

Execution summary:
- Workstreams executed: 5
- Agents used: 3
- Rounds: 2
- Duration: ~5 min
```

**Failure:**
```
❌ Execution failed: 2 tasks failed

Failed tasks:
- bd-0001.3: Service layer
- bd-0001.5: Tests

Check logs for details:
bd show bd-0001.3
bd show bd-0001.5
```

## Example Session

```bash
# Decompose feature first
@idea "Add user auth"
# → bd-0001

@design bd-0001
# → bd-0001.1, bd-0001.2, bd-0001.3 (sequential deps)

# Execute all workstreams
@oneshot bd-0001 --agents 2

# Output:
⏳ Executing feature bd-0001 with 2 agents...

Round 1:
  🤖 Agent 1: Executing bd-0001.1 (Domain entities)
  🤖 Agent 2: Executing bd-0001.2 (Repository layer)

✅ bd-0001.1 complete
✅ bd-0001.2 complete

Round 2:
  🤖 Agent 1: Executing bd-0001.3 (Service layer)

✅ bd-0001.3 complete

✅ Feature complete! Executed 3 workstreams
```

## Key Features

**Automatic Dependency Management:**
- No manual task ordering needed
- Beads DAG automatically tracks dependencies
- Tasks become ready as dependencies complete

**Parallel Execution:**
- Independent tasks run in parallel
- Configurable agent count (1-10)
- ThreadPoolExecutor for concurrency

**Progress Tracking:**
- Real-time status updates
- Ready task discovery between rounds
- Comprehensive error reporting

**Fault Tolerance:**
- Failed tasks don't block other tasks
- Continues until no tasks ready
- Reports all failures at end

## Benefits vs Manual Execution

| Aspect | Manual @build | @oneshot |
|--------|---------------|----------|
| **Task discovery** | Manual | `bd ready` |
| **Parallelization** | Manual | Auto (3 agents) |
| **Status tracking** | Manual file moves | Auto status updates |
| **Error handling** | Manual | Auto reporting |
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
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@oneshot bd-0001` | Execute feature with 3 agents |
| `@oneshot bd-0001 --agents 5` | Use 5 agents |
| `bd ready` | List ready tasks |
| `bd graph` | Show dependency graph |
| `bd show {id}` | View task details |

---

**Version:** 2.0.0-beads
**Status:** Beads Integration
**See Also:** `@idea`, `@design`, `@build`
