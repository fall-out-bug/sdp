---
name: oneshot
description: Autonomous multi-agent execution with checkpoints, resume, and PR-less modes
tools: Task, Read, Bash
version: 4.0.0
---

# @oneshot - Multi-Agent Feature Execution

Execute all workstreams for a feature using autonomous orchestrator agent with checkpoint/resume.

## When to Use

- Feature has multiple workstreams (5-30 WS)
- Want autonomous execution with progress tracking
- Need checkpoint/resume for long-running features
- Parallel execution of independent workstreams

## Invocation

```bash
@oneshot F050                       # Execute feature F050
@oneshot F050 --agents 5            # Use 5 builder agents (default: 3)
@oneshot F050 --resume abc123       # Resume from checkpoint
@oneshot F050 --background          # Run in background
```

## How It Works

### Step 1: Identify Feature Workstreams

**Detect Beads:**
```bash
# Check if Beads is available
if bd --version &>/dev/null && [ -d .beads ]; then
    BEADS_ENABLED=true
else
    BEADS_ENABLED=false
fi
```

**Find workstreams:**
```bash
# Method 1: Beads (if enabled)
if [ "$BEADS_ENABLED" = true ]; then
    bd list --parent F050 --json | jq -r '.[].id'
fi

# Method 2: Markdown (always works)
grep -l "^feature: F050" docs/workstreams/backlog/*.md
```

### Step 2: Launch Orchestrator Agent

**CRITICAL:** Use Task tool to spawn general-purpose agent with orchestrator instructions:

```python
# Gather workstreams for this feature
ws_files = Glob("docs/workstreams/backlog/00-050-*.md")  # Adjust pattern for feature
workstreams = []
for f in ws_files:
    ws_id = parse_ws_id(f)  # e.g., "00-050-01" from filename
    # Get beads_id from mapping
    beads_id = Bash(f'grep "{{sdp_id: \\"{ws_id}\\"}}" .beads-sdp-mapping.jsonl | grep -o \'beads_id": "[^"]*"\' | cut -d\'"\' -f4')
    workstreams.append({"ws_id": ws_id, "beads_id": beads_id})

# Launch orchestrator agent
Task(
    subagent_type="general-purpose",
    prompt=f"""
You are executing feature {feature_id} autonomously as an orchestrator.

**READ FIRST:** Read(".claude/agents/orchestrator.md") - This is your specification

**Workstreams to execute:**
{chr(10).join([f"- {w['ws_id']}: {get_title(w['ws_id'])} ({w['beads_id']})" for w in workstreams])}

**Your workflow:**
1. Read .claude/agents/orchestrator.md - this defines your behavior
2. Build dependency graph:
   - For each WS, check dependencies in WS file
   - Or use: `bd show {{beads_id}}` to see blocking relationships
3. Execute in topological order:
   - For each WS: @build {{ws_id}}
   - @build handles: Beads status + TDD + quality gates + commit
4. Update checkpoint after each WS
5. On all complete: @review {feature_id}

**Progress format (timestamps required):**
[HH:MM] Executing {{ws_id}}...
[HH:MM] ✅ COMPLETE (Xm, Y% coverage, commit: abc123)

**Checkpoint file (.oneshot/{feature_id}-checkpoint.json):**
{{
  "feature": "{feature_id}",
  "agent_id": "agent-{timestamp}",
  "status": "in_progress",
  "completed_ws": ["{{ws_id}}", ...],
  "execution_order": ["{{ws_id}}", ...],
  "started_at": "{datetime.now().isoformat()}"
}}

**Escalate to human if:**
- CRITICAL errors (blockers)
- Circular dependencies
- Quality gate fails after retry

**Auto-fix:**
- HIGH/MEDIUM issues (max 2 retries per WS)
- Implementation details within WS scope
""",
    run_in_background=background_flag
)
```

### Step 3: Monitor Progress

**Foreground mode** (default):
```
Agent provides real-time updates:
→ [15:23] Executing 00-050-01...
→ [15:45] ✅ COMPLETE (22m, 85% coverage, commit: a1b2c3d)
→ [15:46] Executing 00-050-02...
...
```

**Background mode** (`--background`):
```bash
# Agent runs in background, returns task_id
task_id = "xyz789"

# Check progress anytime
Read("/tmp/agent_xyz789.log")

# Or wait for completion
TaskOutput(task_id="xyz789", block=True)
```

### Step 4: Resume from Checkpoint

If execution interrupted:

```bash
@oneshot F050 --resume agent-20260205-152300
```

Agent reads checkpoint file `.oneshot/F050-checkpoint.json` and continues from last completed workstream.

## Output

**Success:**
```
✅ Feature F050 Execution Complete

Agent: agent-20260205-152300
Duration: 3h 45m
Workstreams: 13/13 completed
Avg Coverage: 84%

Checkpoint: .oneshot/F050-checkpoint.json

Next Steps:
1. Human UAT (5-10 min)
2. @review F050 (automated)
3. @deploy F050 (if UAT passes)
```

**Failure:**
```
❌ Execution Failed: 00-050-09

Error: Circular dependency detected
Checkpoint saved: .oneshot/F050-checkpoint.json

Resume: @oneshot F050 --resume agent-20260205-152300
Or fix manually: @build 00-050-09
```

## Checkpoint Format

```json
{
  "feature": "F050",
  "agent_id": "agent-20260205-152300",
  "status": "in_progress",
  "completed_ws": ["00-050-01", "00-050-02"],
  "failed_ws": [],
  "execution_order": ["00-050-01", "00-050-02", "00-050-03", ...],
  "started_at": "2026-02-05T15:23:00Z"
}
```

## Orchestrator Agent Capabilities

The orchestrator agent (`.claude/agents/orchestrator.md`) has:

**Autonomous Decisions:**
- Execution order based on dependencies
- Implementation details within WS scope
- Test strategy for each workstream
- Auto-fix HIGH/MEDIUM issues (max 2 retries)

**Human Escalation:**
- CRITICAL errors that block feature
- Circular dependencies
- Scope overflow (>1500 LOC)
- Quality gate failures after 2 retries
- Architectural decisions not in spec

**Quality Standards:**
- All AC met ✅
- Coverage ≥ 80%
- All fast tests pass
- Linters clean (ruff, mypy)
- Clean Architecture compliance

## Key Features

| Feature | Description |
|---------|-------------|
| **Auto dependencies** | Beads DAG tracks dependencies |
| **Parallel execution** | Independent WS run in parallel |
| **Checkpoint/resume** | Continue from interruption |
| **Background mode** | Long-running features |
| **Progress tracking** | Real-time timestamps |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| No tasks executing | `bd list --parent {feature}` to check workstreams |
| Agent not starting | Check `.claude/agents/orchestrator.md` exists |
| Wrong execution order | Check dependency graph: `bd graph {feature}` |
| Background agent silent | Read `/tmp/agent_{task_id}.log` |
| Resume not working | Check checkpoint file: `.oneshot/{feature}-checkpoint.json` |

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@oneshot F050` | Execute with 3 agents |
| `@oneshot F050 --background` | Run in background |
| `@oneshot F050 --resume <id>` | Resume from checkpoint |
| `bd ready` | List ready tasks |
| `bd graph F050` | Show dependency graph |
| `@review F050` | Automated review |

## Example Execution

```bash
User: @oneshot F050

Claude:
→ Launching orchestrator agent...
→ Agent ID: agent-20260205-152300
→ Task ID: xyz789

[Agent Output]
→ Reading feature spec: docs/drafts/F050.md
→ Found 13 workstreams
→ Building dependency graph from Beads...
→ Execution order: 01→02→03→05→06→07→09→08→10→11→12→13→14

→ [15:23] Executing 00-050-01: Workstream Parser...
→ [15:45] ✅ COMPLETE (22m, 85% coverage, commit: a1b2c3d)
→ [15:46] Executing 00-050-02: TDD Runner...
→ [16:12] ✅ COMPLETE (26m, 82% coverage, commit: d4e5f6g)
...

→ All workstreams complete!
→ Running @review F050...
→ Review verdict: APPROVED

→ Feature complete! Duration: 3h 45m
→ Checkpoint: .oneshot/F050-checkpoint.json
→ Ready for human UAT
```

---

**Version:** 4.0.0 (Task-based orchestration)
**See Also:** `@idea`, `@design`, `@build`, `@review`, `@deploy`
**Agent:** `.claude/agents/orchestrator.md`
