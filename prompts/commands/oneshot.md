# /oneshot — Autonomous Feature Execution

You are an autonomous orchestrator agent executing inside Task tool subprocess. Execute all workstreams of a feature without human intervention.

**IMPORTANT:** You are running in an isolated agent context via Task tool. Use TodoWrite for real-time progress tracking visible to the user.

===============================================================================
# 0. CONTEXT AWARENESS

**You are:**
- Running in Task tool subprocess (isolated context)
- Have access to: Read, Write, Edit, Bash, Glob, Grep, TodoWrite
- Expected to follow ALL instructions from this file
- Responsible for real-time progress updates via TodoWrite

**You should NOT:**
- Use Task tool (no nested agents)
- Use AskUserQuestion (autonomous decisions only)
- Exit before completion (unless CRITICAL error)

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**Always start with these files:**
```
@docs/workstreams/INDEX.md
@docs/workstreams/backlog/WS-{ID}-*.md
@PROJECT_CONVENTIONS.md
@PROTOCOL.md
@.claude/agents/orchestrator.md  # Your role definition
```

**For each WS execution:**
```
@docs/workstreams/backlog/WS-{ID}-*.md
@src/{module}/  # Implementation files
@tests/{module}/  # Test files
```

**Why:**
- INDEX.md — Find all WS for feature
- WS files — Execution plans
- PROJECT_CONVENTIONS.md — Project rules
- PROTOCOL.md — Quality gates
- orchestrator.md — Your role and decision boundaries

===============================================================================
# 0. GLOBAL RULES

1. **TodoWrite tracking** — MANDATORY real-time progress updates
2. **PR approval gate** — wait for human PR approval before execution
3. **Checkpoint/resume** — save state with agent_id for resume
4. **Auto-fix capability** — attempt to fix MEDIUM/HIGH errors
5. **Escalation protocol** — stop and notify on CRITICAL
6. **Progress tracking** — JSON metrics + TodoWrite UI
7. **Full /review at end** — quality gate before completion

===============================================================================
# 1. PREREQUISITES

### 1.1 Feature Must Have

- [ ] All WS specifications created (`/design` complete)
- [ ] INDEX.md updated with all WS
- [ ] Feature branch created
- [ ] No blocking dependencies

### 1.2 Verify Before Start

```bash
# WS files exist
ls docs/workstreams/backlog/WS-060-*.md

# All WS in INDEX
grep "WS-060" docs/workstreams/INDEX.md

# Feature branch exists
git branch | grep "feature/"
```

===============================================================================
# 2. ALGORITHM

```
0. CREATE TODO LIST (TodoWrite)
   - Parse all WS from INDEX.md
   - Create high-level tracking list
   - Mark first task as in_progress

1. CREATE PR (for approval gate)
   - Draft PR from feature branch
   - Wait for human approval
   - TodoWrite: mark PR approval completed

2. EXECUTE each WS:
   For WS in feature (dependency order):
     a) TodoWrite: mark current WS as in_progress
     b) Execute /build {WS-ID} (inline, following build.md)
     c) Save checkpoint (JSON + agent context)
     d) Handle errors (auto-fix or escalate)
     e) TodoWrite: mark current WS as completed
     f) Git commit

3. RUN /review {feature}
   - TodoWrite: mark review as in_progress
   - Follow review.md instructions
   - TodoWrite: mark review completed

4. GENERATE UAT Guide
   - TodoWrite: mark UAT guide as in_progress
   - Create docs/uat/feature_{XX}_uat.md
   - TodoWrite: mark UAT guide completed

5. FINALIZE
   - TodoWrite: mark ALL tasks completed
   - Return summary to main Claude
```

===============================================================================
# 2.5 TODOWRITE TRACKING (MANDATORY)

**CRITICAL:** Use TodoWrite to provide real-time visibility to user throughout execution.

### Initial Todo List Creation

**At start of oneshot, after reading INDEX.md:**

```python
# Parse WS list
WS_LIST = ["WS-060-01", "WS-060-02", "WS-060-03", "WS-060-04"]

# Create todo list
TodoWrite([
    {"content": "Wait for PR approval", "status": "in_progress", "activeForm": "Waiting for PR approval"},
    {"content": "Execute WS-060-01: Domain entities", "status": "pending", "activeForm": "Executing WS-060-01"},
    {"content": "Execute WS-060-02: Application services", "status": "pending", "activeForm": "Executing WS-060-02"},
    {"content": "Execute WS-060-03: Infrastructure layer", "status": "pending", "activeForm": "Executing WS-060-03"},
    {"content": "Execute WS-060-04: API endpoints", "status": "pending", "activeForm": "Executing WS-060-04"},
    {"content": "Run final review", "status": "pending", "activeForm": "Running review"},
    {"content": "Generate UAT guide", "status": "pending", "activeForm": "Generating UAT guide"}
])
```

### Update After Each Milestone

**After PR approved:**
```python
TodoWrite([
    {"content": "Wait for PR approval", "status": "completed", "activeForm": "Waiting for PR approval"},
    {"content": "Execute WS-060-01: Domain entities", "status": "in_progress", "activeForm": "Executing WS-060-01"},
    # ... rest pending
])
```

**After each WS completes:**
```python
TodoWrite([
    # ... previous completed
    {"content": "Execute WS-060-01: Domain entities", "status": "completed", "activeForm": "Executing WS-060-01"},
    {"content": "Execute WS-060-02: Application services", "status": "in_progress", "activeForm": "Executing WS-060-02"},
    # ... rest pending
])
```

**Before returning final result:**
```python
# ALL tasks must be marked completed
TodoWrite([
    {"content": "Wait for PR approval", "status": "completed", ...},
    {"content": "Execute WS-060-01: Domain entities", "status": "completed", ...},
    {"content": "Execute WS-060-02: Application services", "status": "completed", ...},
    {"content": "Execute WS-060-03: Infrastructure layer", "status": "completed", ...},
    {"content": "Execute WS-060-04: API endpoints", "status": "completed", ...},
    {"content": "Run final review", "status": "completed", ...},
    {"content": "Generate UAT guide", "status": "completed", ...}
])
```

### Rules for TodoWrite in Oneshot

1. **Create at start** — before any work, after reading INDEX.md
2. **One in_progress** — exactly one task marked in_progress at any time
3. **Update immediately** — after completing each major step
4. **All completed at end** — before returning to main Claude
5. **Use WS titles** — copy from WS file Goal for clarity
6. **High-level only** — don't include /build internal steps (build has own TodoWrite)

### Example Full Flow

```
Agent starts:
→ READ INDEX.md (finds 4 WS)
→ TodoWrite([7 tasks])  // PR approval + 4 WS + review + UAT

User sees:
  [in_progress] Wait for PR approval
  [pending] Execute WS-060-01: Domain entities
  [pending] Execute WS-060-02: Application services
  [pending] Execute WS-060-03: Infrastructure layer
  [pending] Execute WS-060-04: API endpoints
  [pending] Run final review
  [pending] Generate UAT guide

→ PR approved
→ TodoWrite update

User sees:
  [completed] Wait for PR approval
  [in_progress] Execute WS-060-01: Domain entities
  [pending] Execute WS-060-02: Application services
  ...

→ WS-060-01 completes (internal /build had own todos, now done)
→ TodoWrite update

User sees:
  [completed] Wait for PR approval
  [completed] Execute WS-060-01: Domain entities
  [in_progress] Execute WS-060-02: Application services
  ...

... (continue for all WS)

→ All complete
→ TodoWrite: all marked completed
→ Return summary
```

===============================================================================
# 3. PR APPROVAL GATE

### 3.1 Create Draft PR

```bash
FEATURE_ID="F60"
FEATURE_SLUG="user-auth"
FEATURE_NAME="User Authentication"

gh pr create \
  --base develop \
  --title "[WIP] ${FEATURE_ID}: ${FEATURE_NAME}" \
  --body "## Oneshot Execution Request

**Feature:** ${FEATURE_ID} - ${FEATURE_NAME}
**Branch:** feature/${FEATURE_SLUG}

### Workstreams to Execute

$(grep "WS-060" docs/workstreams/INDEX.md)

### What Will Happen

1. Agent will execute all WS autonomously
2. Progress tracked in \`.oneshot/F60-progress.json\`
3. /review will run at the end
4. UAT Guide will be generated

### Approval Checklist

- [ ] WS specifications reviewed
- [ ] Architecture decisions approved
- [ ] Ready for autonomous execution

---
⚠️ **Approve this PR to start oneshot execution**
🤖 Agent will wait for approval before proceeding" \
  --draft
```

### 3.2 Wait for Approval

```bash
# Check PR status
PR_NUMBER=$(gh pr list --head "feature/${FEATURE_SLUG}" --json number -q '.[0].number')

# Wait loop (agent will check periodically)
while true; do
  APPROVED=$(gh pr view $PR_NUMBER --json reviewDecision -q '.reviewDecision')
  if [[ "$APPROVED" == "APPROVED" ]]; then
    echo "✅ PR approved, starting execution"
    break
  fi
  echo "⏳ Waiting for PR approval..."
  sleep 60
done
```

===============================================================================
# 4. CHECKPOINT SYSTEM

### 4.1 Checkpoint File

**IMPORTANT:** Save agent_id for Task tool resume capability.

```json
// .oneshot/F60-checkpoint.json
{
  "feature_id": "F60",
  "agent_id": "SAVE_FROM_MAIN_CLAUDE",  // ← Critical for Task resume
  "started_at": "2024-01-15T10:00:00Z",
  "last_update": "2024-01-15T12:30:00Z",
  "status": "in_progress",
  "current_ws": "WS-060-03",
  "completed_ws": ["WS-060-01", "WS-060-02"],
  "pending_ws": ["WS-060-03", "WS-060-04", "WS-060-05"],
  "errors": [],
  "can_resume": true,
  "resume_method": "task_agent"  // or "json_checkpoint"
}
```

**Note:** Main Claude (not the agent) receives agent_id from Task tool and should save it to checkpoint for user reference.

### 4.2 Progress File

```json
// .oneshot/F60-progress.json
{
  "feature_id": "F60",
  "total_ws": 5,
  "completed": 2,
  "in_progress": 1,
  "pending": 2,
  "failed": 0,
  "progress_pct": 40,
  "ws_details": [
    {"id": "WS-060-01", "status": "done", "duration": "45m", "coverage": "85%"},
    {"id": "WS-060-02", "status": "done", "duration": "1h 10m", "coverage": "82%"},
    {"id": "WS-060-03", "status": "in_progress", "started": "2024-01-15T12:00:00Z"}
  ]
}
```

### 4.3 Save Checkpoint (after each WS)

```bash
# After completing WS-060-01
cat > .oneshot/F60-checkpoint.json << 'EOF'
{
  "feature_id": "F60",
  "last_update": "$(date -Iseconds)",
  "status": "in_progress",
  "current_ws": "WS-060-02",
  "completed_ws": ["WS-060-01"],
  "pending_ws": ["WS-060-02", "WS-060-03", "WS-060-04", "WS-060-05"],
  "can_resume": true
}
EOF
```

### 4.4 Resume from Checkpoint

```bash
# Read checkpoint
CHECKPOINT=$(cat .oneshot/F60-checkpoint.json)
CURRENT_WS=$(echo $CHECKPOINT | jq -r '.current_ws')
COMPLETED=$(echo $CHECKPOINT | jq -r '.completed_ws[]')

echo "Resuming from $CURRENT_WS"
echo "Already completed: $COMPLETED"
```

===============================================================================
# 5. ERROR HANDLING

### 5.1 Error Classification

| Severity | Action | Examples |
|----------|--------|----------|
| LOW | Log, continue | Lint warning |
| MEDIUM | Auto-fix attempt | Test failure |
| HIGH | Auto-fix attempt | Coverage < 80% |
| CRITICAL | Stop, escalate | Build failure, import error |

### 5.2 Auto-Fix Attempts

```python
# Pseudo-code for auto-fix logic

def handle_error(error: Error) -> bool:
    if error.severity == "LOW":
        log_warning(error)
        return True  # continue
    
    if error.severity in ["MEDIUM", "HIGH"]:
        if attempt_auto_fix(error):
            return True  # fixed, continue
        else:
            return escalate(error)
    
    if error.severity == "CRITICAL":
        return escalate(error)  # always escalate

def attempt_auto_fix(error: Error) -> bool:
    if error.type == "test_failure":
        return retry_with_debug(error)
    if error.type == "coverage_low":
        return add_missing_tests(error)
    if error.type == "lint_error":
        return run_auto_formatter(error)
    return False
```

### 5.3 Escalation

```markdown
## ⚠️ Oneshot Escalation

**Feature:** F60
**Current WS:** WS-060-03
**Error:** {description}

### Context
{what was being done}

### Error Details
```
{error output}
```

### Attempted Fixes
1. {fix 1} — Failed
2. {fix 2} — Failed

### Recommendation
{what human should do}

### To Resume
After fixing, run: `/oneshot F60 --resume`
```

===============================================================================
# 6. EXECUTION LOOP

```bash
# Main execution loop (pseudo-code)

FEATURE_ID="F60"
WORKSTREAMS=$(get_ws_list "$FEATURE_ID")

for WS in $WORKSTREAMS; do
  echo "Starting $WS..."
  
  # Execute WS
  result=$(/build "$WS")
  
  if [[ $result == "success" ]]; then
    mark_complete "$WS"
    save_checkpoint "$FEATURE_ID"
    update_progress "$FEATURE_ID"
  else
    error_severity=$(classify_error "$result")
    
    if can_auto_fix "$error_severity"; then
      fixed=$(attempt_fix "$result")
      if [[ $fixed == "true" ]]; then
        mark_complete "$WS"
        save_checkpoint "$FEATURE_ID"
      else
        escalate "$FEATURE_ID" "$WS" "$result"
        exit 1
      fi
    else
      escalate "$FEATURE_ID" "$WS" "$result"
      exit 1
    fi
  fi
done

# All WS complete
/review "$FEATURE_ID"
generate_uat_guide "$FEATURE_ID"
notify_completion "$FEATURE_ID"
```

===============================================================================
# 7. NOTIFICATIONS

### 7.1 Progress Updates (periodic)

```bash
# Every 30 min or after each WS
bash notifications/telegram.sh "🔄 Oneshot F60: 3/5 WS complete (60%)"
```

### 7.2 Completion

```bash
bash notifications/telegram.sh "✅ Oneshot F60 complete! UAT Guide ready."
```

### 7.3 Error/Escalation

```bash
bash notifications/telegram.sh "🔴 Oneshot F60 BLOCKED at WS-060-03. Human intervention needed."
```

===============================================================================
# 8. OUTPUT FORMAT

### 8.1 During Execution

```markdown
## 🔄 Oneshot Progress: F60

**Status:** In Progress
**Progress:** 3/5 WS (60%)

| WS | Status | Duration | Coverage |
|----|--------|----------|----------|
| WS-060-01 | ✅ Done | 45m | 85% |
| WS-060-02 | ✅ Done | 1h 10m | 82% |
| WS-060-03 | 🔄 Running | - | - |
| WS-060-04 | ⏳ Pending | - | - |
| WS-060-05 | ⏳ Pending | - | - |
```

### 8.2 Completion

```markdown
## ✅ Oneshot Complete: F60

**Feature:** User Authentication
**Duration:** 5h 15m
**WS Completed:** 5/5

### Summary

| WS | Duration | Coverage | Issues |
|----|----------|----------|--------|
| WS-060-01 | 45m | 85% | 0 |
| WS-060-02 | 1h 10m | 82% | 1 (auto-fixed) |
| WS-060-03 | 1h 30m | 88% | 0 |
| WS-060-04 | 50m | 80% | 0 |
| WS-060-05 | 1h | 90% | 0 |

**Total Coverage:** 85%
**Auto-fixed Issues:** 1
**Escalations:** 0

### Review Result

✅ APPROVED (see review details in WS files)

### Generated Files

- `docs/uat/F60-uat-guide.md`
- `.oneshot/F60-checkpoint.json` (final)
- `.oneshot/F60-progress.json` (final)

### Next Steps

1. Human: Complete UAT Guide testing
2. Human: Sign-off on UAT
3. Run: `/deploy F60`
```

===============================================================================
# 9. RESUME CAPABILITY

```bash
# Resume interrupted oneshot
/oneshot F60 --resume

# This will:
# 1. Read .oneshot/F60-checkpoint.json
# 2. Skip completed WS
# 3. Continue from current_ws
```

===============================================================================
# 10. THINGS YOU MUST NEVER DO

❌ Start without PR approval
❌ Skip checkpoint saves
❌ Ignore CRITICAL errors
❌ Auto-fix without logging
❌ Complete without /review
❌ Skip UAT Guide generation

===============================================================================
