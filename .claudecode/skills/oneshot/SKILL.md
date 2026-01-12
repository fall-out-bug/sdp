---
name: oneshot
description: Autonomous execution of all workstreams in a feature. Manages dependencies, checkpoints, and quality gates.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /oneshot - Autonomous Feature Execution

Execute all workstreams in a feature autonomously with checkpoint/resume support.

## When to Use

- After `/design` completes WS planning
- To execute feature hands-off
- For parallel development
- When you trust the agent

## Invocation

```bash
/oneshot F60                    # Start fresh
/oneshot F60 --resume           # Resume from checkpoint
/oneshot F60 --no-approval      # Skip PR approval (dangerous!)
```

## Workflow

**IMPORTANT:** This skill delegates to the master prompt.

### Load Master Prompt

```bash
cat sdp/prompts/commands/oneshot.md
```

**This file contains:**
- Full autonomous execution algorithm
- PR approval gate (GitHub)
- Dependency resolution
- Checkpoint system (resume capability)
- Error handling (auto-fix MEDIUM, escalate CRITICAL)
- Progress tracking (JSON metrics)
- Final review integration
- Telegram notifications

### Execute Instructions

Follow `sdp/prompts/commands/oneshot.md`:

**Phase 1: Initialization**
1. Check if resume needed
2. Create PR for approval
3. Wait for human approval
4. Read feature context
5. Build execution plan

**Phase 2: Execution Loop**
- For each WS (dependency order):
  1. Pre-build checks
  2. Execute `/build WS-ID`
  3. Post-build checks
  4. Update checkpoint
  5. Update progress JSON
  6. Commit

**Phase 3: Error Handling**
- CRITICAL: Save checkpoint, notify, STOP
- HIGH: Auto-fix, retry (max 2x)
- MEDIUM: Mark for review, continue

**Phase 4: Final Review**
- Run post-oneshot hooks
- Execute `/review F{XX}`
- Generate UAT guide

## Key Features

From master prompt:

1. **PR Approval Gate** - human approval before execution
2. **Checkpoint System** - resume from failures
3. **Progress Tracking** - real-time JSON metrics
4. **Dependency Resolution** - correct execution order
5. **Auto-Fix** - handles MEDIUM/HIGH failures
6. **Telegram Alerts** - notifies on CRITICAL blocks
7. **Audit Log** - tracks all events

## Checkpoint Format

`.oneshot/F{XX}-checkpoint.json`:

```json
{
  "feature": "F60",
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

## Progress JSON

`.oneshot/F{XX}-progress.json`:

```json
{
  "feature": "F60",
  "status": "executing",
  "progress": {
    "completion_pct": 50,
    "ws_completed": 2,
    "ws_current": "WS-060-03",
    "loc_written": 1150,
    "elapsed": "1h 23m"
  }
}
```

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
