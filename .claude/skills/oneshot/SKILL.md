---
name: oneshot
description: Autonomous execution of all workstreams in a feature. Manages dependencies, checkpoints, and quality gates.
tools: Read, Bash, Task
---

# /oneshot - Autonomous Feature Execution

Execute all workstreams in a feature autonomously using Task tool for isolated agent execution.

## When to Use

- After `/design` completes WS planning
- To execute feature hands-off
- For features with 3+ workstreams
- When you want background execution
- For parallel development

## Invocation

```bash
/oneshot F60                    # Synchronous execution
/oneshot F60 --background       # Background execution (for large features)
/oneshot F60 --resume {agent_id} # Resume from interrupted execution
```

## Workflow

**IMPORTANT:** This skill uses Task tool to spawn an orchestrator agent in isolated context.

### Step 1: Validate Feature Exists

```bash
# Check feature specification
ls docs/specs/feature_60/ || echo "Feature F60 not found"

# Check workstreams exist
ls docs/workstreams/backlog/WS-060-*.md | wc -l

# Verify INDEX
grep "WS-060" docs/workstreams/INDEX.md
```

### Step 2: Parse Arguments

```bash
# Extract feature ID
FEATURE_ID="$1"  # e.g., F60

# Check for --background flag
BACKGROUND=false
if [[ "$2" == "--background" ]]; then
  BACKGROUND=true
fi

# Check for --resume flag
RESUME_AGENT_ID=""
if [[ "$2" == "--resume" ]]; then
  RESUME_AGENT_ID="$3"
fi
```

### Step 3: Launch Orchestrator Agent via Task Tool

**For fresh execution (synchronous):**

```python
Task(
    subagent_type="general-purpose",
    prompt=f"""You are an autonomous orchestrator agent for feature {FEATURE_ID}.

STEP 1: Read your instructions
- READ prompts/commands/oneshot.md for full orchestration workflow
- READ .claude/agents/orchestrator.md for agent-specific guidance

STEP 2: Understand the feature
- READ docs/specs/feature_{feature_num}/feature.md
- READ docs/workstreams/INDEX.md to find all WS for this feature
- READ all WS files: docs/workstreams/backlog/WS-{feature_num}-*.md

STEP 3: Execute autonomously
- Follow ALL steps from prompts/commands/oneshot.md
- Use TodoWrite to track progress in real-time
- Create PR and wait for approval
- Execute each WS following /build logic
- Run final /review
- Generate UAT guide

STEP 4: Return results
Report final status: completed WS, coverage, issues, next steps.

Execute feature {FEATURE_ID} completely following oneshot.md instructions.""",
    description=f"Orchestrating {FEATURE_ID}",
    run_in_background=false
)
```

**For background execution:**

```python
Task(
    subagent_type="general-purpose",
    prompt="(same as above)",
    description=f"Orchestrating {FEATURE_ID} (background)",
    run_in_background=true
)
```

Returns: `{"task_id": "xyz", "output_file": "/path/to/log"}`

**For resume:**

```python
Task(
    resume=RESUME_AGENT_ID,
    prompt="Continue execution from last checkpoint"
)
```

### Step 4: Monitor Execution

**Synchronous mode:**
- Wait for agent to complete
- Agent outputs progress via TodoWrite (visible in UI)
- Receive final result when done

**Background mode:**
- User continues working
- Check progress: `tail -f {output_file}`
- Or: `Read({output_file})` to see current state
- Notification when agent completes

### Step 5: Display Results

When agent completes, show summary:

```markdown
## ✅ Feature {FEATURE_ID} Execution Complete

**Agent ID:** {agent_id} (use for resume if needed)
**Duration:** {elapsed_time}
**Workstreams:** {completed}/{total}
**Coverage:** avg {coverage}%

### Executed Workstreams
- WS-060-01: Domain entities ✅ (45m, 85% coverage)
- WS-060-02: Application services ✅ (1h 10m, 82% coverage)
- WS-060-03: Infrastructure ✅ (50m, 88% coverage)

### Review Status
{review_verdict}

### Next Steps
1. Human UAT (5-10 min)
2. `/deploy {FEATURE_ID}` if UAT passes

**Agent ID for resume:** {agent_id}
```

## Key Features

**Claude Code Integration:**

1. **Task Tool Orchestration** - isolated agent with clean context
2. **TodoWrite Progress** - real-time UI updates during execution
3. **Background Execution** - run long features async
4. **Resume Capability** - continue from agent_id checkpoint
5. **Parallel Tool Calls** - faster validation (pytest + ruff + mypy simultaneously)

**From Master Prompt (oneshot.md):**

6. **PR Approval Gate** - human approval before execution
7. **Checkpoint System** - JSON files for state persistence
8. **Dependency Resolution** - correct WS execution order
9. **Auto-Fix** - handles MEDIUM/HIGH failures autonomously
10. **Error Escalation** - CRITICAL issues stop and notify
11. **Final Review** - quality gate before UAT

## Resume Strategy

**Two mechanisms work together:**

### 1. Task Agent Resume (Primary)

```bash
# After interruption, agent returns agent_id
/oneshot F60 --resume {agent_id}

# Task tool resumes with full context
Task(resume="{agent_id}", prompt="Continue from where you stopped")
```

**Advantages:**
- Built into Claude Code
- Preserves full conversation context
- No manual checkpoint management

### 2. JSON Checkpoints (Backup)

`.oneshot/F{XX}-checkpoint.json`:

```json
{
  "feature": "F60",
  "agent_id": "abc123xyz",  // ← for Task resume
  "status": "in-progress",
  "completed_ws": ["WS-060-01", "WS-060-02"],
  "current_ws": "WS-060-03",
  "started_at": "2026-01-11T10:00:00Z",
  "metrics": {
    "ws_total": 4,
    "ws_completed": 2,
    "loc_total": 1150,
    "coverage_avg": 84
  }
}
```

**Used for:**
- Cross-session resume
- Manual recovery
- Metrics tracking

## Error Escalation

| Severity | Action |
|----------|--------|
| CRITICAL | Stop, checkpoint, notify human |
| HIGH | Auto-fix, retry max 2x, escalate if fails |
| MEDIUM | Mark needs_review, continue |

## Output

Final summary:

```markdown
## ✅ Feature F60 COMPLETE

**Status:** APPROVED
**Workstreams:** 4/4 completed
**Duration:** 3h 45m
**Coverage:** avg 86%

### Next Steps
1. Human UAT (5-10 min)
2. `/deploy F60` if passes
```

## Hooks Integration

- `pre-build.sh` - before each WS
- `post-build.sh` - after each WS
- `post-oneshot.sh` - after all WS (integration/e2e tests)

## Master Prompt Location

📄 **sdp/prompts/commands/oneshot.md** (750+ lines)

**Why reference?**
- Complex orchestration logic
- Checkpoint/resume mechanics
- Error handling strategies
- Too long to duplicate

## Quick Reference

**Input:** Feature ID  
**Output:** All WS executed + Review + UAT guide  
**Next:** Human UAT → `/deploy F{XX}`

## Safety

- **PR approval required** (default)
- **Checkpoints every WS** (resume capability)
- **Human intervention** on CRITICAL failures
- **Post-oneshot tests** before review
