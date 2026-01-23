# /oneshot — One-Shot Feature Implementation

Ты — orchestrator agent. Выполняешь всю фичу за один проход (one-shot).

===============================================================================
# 0. MISSION

**Выполнить ВСЕ workstreams фичи автономно, соблюдая зависимости и качество.**

Ты НЕ спрашиваешь разрешения между WS. Ты САМ:
- Выбираешь следующий WS
- Выполняешь его
- Проверяешь результат
- Решаешь что дальше

Человек вмешивается ТОЛЬКО если CRITICAL блокер.

===============================================================================
# 1. INPUT

```bash
/oneshot F60
```

Где `F60` — feature ID.

===============================================================================
# 2. INITIALIZATION

### 2.1 Check Git Branch (GitFlow)

```bash
# Проверь что ты в feature branch
CURRENT_BRANCH=$(git branch --show-current)

if [[ "$CURRENT_BRANCH" != "feature/"* ]]; then
  echo "⛔ ERROR: Not on feature branch"
  echo "Current: $CURRENT_BRANCH"
  echo "Expected: feature/{slug}"
  echo ""
  echo "Run /design first to create feature branch"
  exit 1
fi

echo "✓ Branch: $CURRENT_BRANCH"
```

### 2.2 Check/Resume from Checkpoint

```bash
# Проверь наличие checkpoint
FEATURE_ID="F60"
CHECKPOINT_FILE=".oneshot/${FEATURE_ID}-checkpoint.json"

if [[ -f "$CHECKPOINT_FILE" ]]; then
  echo "📍 Found checkpoint: $CHECKPOINT_FILE"
  
  # Прочитай состояние
  COMPLETED_WS=$(jq -r '.completed_ws[]' "$CHECKPOINT_FILE")
  CURRENT_WS=$(jq -r '.current_ws' "$CHECKPOINT_FILE")
  STATUS=$(jq -r '.status' "$CHECKPOINT_FILE")
  
  echo "Status: $STATUS"
  echo "Completed: $COMPLETED_WS"
  echo "Current: $CURRENT_WS"
  
  if [[ "$STATUS" == "blocked" ]]; then
    echo "⚠️ Previous execution was BLOCKED"
    echo "Reason: $(jq -r '.blocked_reason' "$CHECKPOINT_FILE")"
    echo ""
    echo "Options:"
    echo "1. Resume: /oneshot $FEATURE_ID --resume"
    echo "2. Restart: /oneshot $FEATURE_ID --restart"
    exit 1
  fi
  
  # Auto-resume
  echo "Resuming from: $CURRENT_WS"
else
  echo "🆕 Starting fresh execution"
  mkdir -p .oneshot
fi
```

### 2.3 Create PR for Approval (GitFlow)

```bash
# Create PR: feature/{slug} → develop
FEATURE_ID="F60"
FEATURE_SLUG="lms-integration"

# Push feature branch if not already pushed
git push origin feature/${FEATURE_SLUG}

# Create PR via GitHub CLI (if available)
if command -v gh &> /dev/null; then
  gh pr create \
    --base develop \
    --head feature/${FEATURE_SLUG} \
    --title "Feature F${FEATURE_ID}: ${FEATURE_TITLE}" \
    --body "## Workstreams

$(ls tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md | \
   xargs -I {} basename {} | sed 's/^/- /')

## Scope

**Total WS:** $(ls tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md | wc -l)
**Estimated LOC:** ~{total_loc}

## Execution Plan

This PR will be auto-executed by /oneshot F${FEATURE_ID}

**Approval required before execution.**

## Checklist

- [ ] WS specifications reviewed
- [ ] Architecture aligned with PROJECT_MAP
- [ ] No duplicate WS in INDEX
- [ ] All dependencies clear

/oneshot will start after approval." \
    --label "oneshot,F${FEATURE_ID}" \
    --reviewer @human
  
  PR_URL=$(gh pr view --json url -q .url)
  echo "✓ PR created: $PR_URL"
else
  echo "⚠️ GitHub CLI not available"
  echo "Create PR manually: feature/${FEATURE_SLUG} → develop"
fi
```

### 2.4 Wait for PR Approval

```markdown
⏳ Waiting for PR approval...

PR: {url}
Status: PENDING REVIEW

Options:
1. Wait for human approval (recommended)
2. Skip approval with: /oneshot F{XX} --no-approval (dangerous!)

Approval required from:
- @human (maintainer)

Once approved, /oneshot will automatically start execution.
```

**Polling for approval:**

```bash
# Check PR status
while true; do
  PR_STATUS=$(gh pr view --json reviewDecision -q .reviewDecision)
  
  if [[ "$PR_STATUS" == "APPROVED" ]]; then
    echo "✅ PR APPROVED - starting execution"
    break
  elif [[ "$PR_STATUS" == "CHANGES_REQUESTED" ]]; then
    echo "❌ PR CHANGES REQUESTED"
    echo "Fix issues and re-run /oneshot"
    exit 1
  else
    echo "⏳ Still waiting for approval... (status: $PR_STATUS)"
    sleep 60  # Check every minute
  fi
done
```

### 2.5 Send Start Notification

```bash
# Audit log
bash sdp/notifications/audit-log.sh command_started "/oneshot" "${FEATURE_ID}"

# Send Telegram notification (if configured)
WS_COUNT=$(ls tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md | wc -l)
bash sdp/notifications/telegram.sh oneshot_started "${FEATURE_ID}" "${WS_COUNT}"
```

### 2.6 Read Feature Context

```bash
# Feature spec
cat tools/hw_checker/docs/specs/feature_60/feature.md

# Workstreams map
grep "F60" tools/hw_checker/docs/workstreams/INDEX.md

# Project context
cat tools/hw_checker/docs/PROJECT_MAP.md
```

### 2.7 Build Execution Plan

Создай план выполнения:

```markdown
## Execution Plan: F60

**Feature:** {название}
**Total WS:** {count}

### Dependency Graph

```
WS-060-01 (no deps)
    ↓
WS-060-02 (depends on 060-01)
    ↓
WS-060-03 (depends on 060-02)
    ↓
WS-060-04 (depends on 060-03)
```

### Execution Order

1. WS-060-01 (ready)
2. WS-060-02 (after 060-01)
3. WS-060-03 (after 060-02)
4. WS-060-04 (after 060-03)

**Estimated scope:** {sum of all WS LOC}
```

### 2.3 Confirm Start

```markdown
## Ready to Execute

Feature: F60 - {название}
Workstreams: 4
Order: sequential (dependencies)

Starting autonomous execution...
```

===============================================================================
# 3. EXECUTION LOOP

```python
# Псевдокод
while True:
    # 1. Get next WS
    next_ws = find_ready_ws(feature_id)
    
    if next_ws is None:
        break  # All done
    
    # 2. Execute
    result = execute_ws(next_ws)
    
    # 3. Check result
    if result.failed:
        if result.severity == "CRITICAL":
            stop_and_notify_human()
        else:
            fix_and_retry()
    
    # 4. Update INDEX
    update_index(next_ws, "completed")
    
    # 5. Log progress
    log_progress(feature_id)

# Final review
review_result = review_feature(feature_id)
return review_result
```

### 3.1 Find Ready WS

```bash
# Найти WS фичи
grep "| WS-060" tools/hw_checker/docs/workstreams/INDEX.md

# Проверить зависимости
# Для каждого WS прочитать секцию "Зависимость"
```

**Правила:**
- WS готов к выполнению если:
  - Статус: `backlog`
  - Зависимости: все `completed` или "Независимый"
  
**Порядок приоритета:**
1. WS без зависимостей (параллельно если можно)
2. WS с выполненными зависимостями
3. Сначала меньшие (SMALL → MEDIUM → LARGE)

### 3.2 Execute WS

Для каждого WS выполни:

```bash
# 1. Pre-build checks
bash sdp/hooks/pre-build.sh WS-{ID}

# 2. Audit log
bash sdp/notifications/audit-log.sh ws_started "WS-{ID}"

# 3. Execute (Phase 3)
# Следуй @sdp/prompts/structured/phase-3-implement.md
# - Read WS file
# - Execute TDD
# - Write code
# - Run tests
# - Append Execution Report

# 4. Post-build checks
bash sdp/hooks/post-build.sh WS-{ID}

# 5. Audit log (on success)
bash sdp/notifications/audit-log.sh ws_completed "WS-{ID}" "{LOC}" "{coverage}"

# 6. Git commit
git add .
git commit -m "feat(scope): WS-{ID} - {title}

{one-line description}

Goal: {goal statement}
Files: {count} files, {LOC} lines
Tests: {count} tests, {coverage}%"
```

### 3.3 Handle Failures

Если WS провалился:

```markdown
## WS-{ID} FAILED

**Error:** {error message}
**Severity:** CRITICAL / HIGH / MEDIUM

### Analysis

[Что пошло не так]

### Decision

**If CRITICAL (блокирует всю фичу):**
- Save checkpoint: `.oneshot/F{XX}-checkpoint.json`
- Audit log: `bash sdp/notifications/audit-log.sh ws_failed "WS-{ID}" "{reason}"`
- Send notification: `bash sdp/notifications/telegram.sh oneshot_blocked "F{XX}" "WS-{ID}" "{reason}"`
- EXIT with error

**If HIGH (можно попробовать автофикс):**
1. Analyze error
2. Fix automatically (если очевидно)
3. Retry WS
4. If still fails → CRITICAL

**If MEDIUM (можно отложить):**
- Mark WS as "needs_review"
- Continue with other WS
- Report в final review
```
→ STOP, create BLOCKED checkpoint, notify human:

```bash
# Create BLOCKED checkpoint
cat > ".oneshot/F${FEATURE_ID}-checkpoint.json" <<EOF
{
  "feature": "F${FEATURE_ID}",
  "status": "blocked",
  "completed_ws": ["WS-060-01"],
  "current_ws": "WS-060-02",
  "blocked_reason": "{error message}",
  "blocked_at": "$(date -Iseconds)",
  "severity": "CRITICAL"
}
EOF

git add ".oneshot/F${FEATURE_ID}-checkpoint.json"
git commit -m "chore(oneshot): F${FEATURE_ID} BLOCKED at WS-060-02 - CRITICAL error"
```

```
⛔ CRITICAL BLOCKER: WS-{ID}

Error: {message}
Impact: Cannot continue with F{XX}

Required action:
1. {что нужно исправить}
2. {альтернативный план}

Checkpoint saved: .oneshot/F{XX}-checkpoint.json
Status: BLOCKED

To resume after fix:
  /oneshot F{XX} --resume

Waiting for human decision...
```

**If HIGH/MEDIUM (можно исправить):**
→ Auto-fix:
1. Analyze root cause
2. Adjust approach
3. Retry (max 2 attempts)
4. If still fails → escalate to CRITICAL
```

### 3.4 Update Progress & Checkpoint

После каждого WS:

```bash
# Calculate metrics
START_TIME=$(date +%s)
ELAPSED=$(($(date +%s) - START_TIME))
LOC_TOTAL=$(git diff --stat $(git rev-list --max-parents=0 HEAD) | tail -1 | awk '{print $4}')
WS_COMPLETED=$(ls .oneshot/completed-*.marker 2>/dev/null | wc -l)
WS_TOTAL=$(ls tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md | wc -l)

# Update checkpoint with full metrics
FEATURE_ID="F60"
CHECKPOINT_FILE=".oneshot/${FEATURE_ID}-checkpoint.json"

cat > "$CHECKPOINT_FILE" <<EOF
{
  "feature": "$FEATURE_ID",
  "status": "in-progress",
  "completed_ws": ["WS-060-01", "WS-060-02"],
  "current_ws": "WS-060-03",
  "pending_ws": ["WS-060-04"],
  "started_at": "$(date -Iseconds)",
  "last_updated": "$(date -Iseconds)",
  "blocked_reason": null,
  "metrics": {
    "ws_total": $WS_TOTAL,
    "ws_completed": $WS_COMPLETED,
    "ws_completion_pct": $(($WS_COMPLETED * 100 / $WS_TOTAL)),
    "loc_total": $LOC_TOTAL,
    "elapsed_seconds": $ELAPSED,
    "coverage_avg": null,
    "complexity_avg": null
  }
}
EOF

# Create progress JSON for external tools
cat > ".oneshot/${FEATURE_ID}-progress.json" <<EOF
{
  "command": "/oneshot",
  "feature": "$FEATURE_ID",
  "status": "executing",
  "progress": {
    "ws_total": $WS_TOTAL,
    "ws_completed": $WS_COMPLETED,
    "ws_current": "WS-060-03",
    "ws_pending": 1,
    "completion_pct": $(($WS_COMPLETED * 100 / $WS_TOTAL)),
    "metrics": {
      "loc_written": $LOC_TOTAL,
      "coverage_avg": null,
      "complexity_avg": null
    },
    "timing": {
      "started_at": "$(date -u -Iseconds -d @$START_TIME)",
      "elapsed_seconds": $ELAPSED,
      "elapsed_human": "$(($ELAPSED / 3600))h $(($ELAPSED % 3600 / 60))m"
    }
  }
}
EOF

# Commit both files
git add "$CHECKPOINT_FILE" ".oneshot/${FEATURE_ID}-progress.json"
git commit -m "chore(oneshot): checkpoint F${FEATURE_ID} - WS-060-02 complete"
```

**Progress report:**

```markdown
## Progress: F60

| WS | Status | LOC | Coverage |
|----|--------|-----|----------|
| WS-060-01 | ✅ DONE | 350 | 85% |
| WS-060-02 | ✅ DONE | 800 | 82% |
| WS-060-03 | 🔄 IN PROGRESS | - | - |
| WS-060-04 | ⏳ WAITING | - | - |

**Completed:** 2/4 (50%)
**Next:** WS-060-03
**Checkpoint:** `.oneshot/F60-checkpoint.json` ✅
**Progress JSON:** `.oneshot/F60-progress.json` ✅

### Live Metrics (JSON)

```json
{
  "feature": "F60",
  "status": "executing",
  "progress": {
    "completion_pct": 50,
    "ws_completed": 2,
    "ws_total": 4,
    "loc_written": 1150,
    "elapsed": "1h 23m"
  }
}
```
```

===============================================================================
# 4. FINAL REVIEW

После выполнения ВСЕХ WS:

```bash
# Run post-oneshot hooks
bash sdp/hooks/post-oneshot.sh F60

# Auto-review
/codereview F60
```

Следуй `@sdp/prompts/commands/codereview.md`:
- Check all WS
- Generate UAT Guide
- Report verdict

### 4.1 If APPROVED

```bash
# Calculate duration
DURATION=$(($(date +%s) - START_TIME))
DURATION_HUMAN="$(($DURATION / 3600))h $(($DURATION % 3600 / 60))m"

# Audit log
bash sdp/notifications/audit-log.sh command_completed "/oneshot" "F60" "success"

# Send completion notification
bash sdp/notifications/telegram.sh oneshot_completed "F60" "$DURATION_HUMAN"
```

```markdown
## ✅ Feature F60 COMPLETE

**Status:** APPROVED
**Workstreams:** 4/4 completed
**Coverage:** {avg}%
**Regression:** ✅ all passed

### Summary

| Metric | Value |
|--------|-------|
| Total LOC | {sum} |
| Total tests | {count} |
| Avg coverage | {%} |
| Critical issues | 0 |

### Next Steps

1. Human UAT: `tools/hw_checker/docs/uat/F60-uat-guide.md`
2. After sign-off: `/deploy F60`

**Feature ready for human verification.**
```

### 4.2 If CHANGES REQUESTED

```markdown
## ⚠️ Feature F60 NEEDS FIXES

**Status:** CHANGES REQUESTED

### Issues

| WS | Severity | Issue |
|----|----------|-------|
| WS-060-02 | HIGH | Coverage 75% < 80% |
| WS-060-03 | CRITICAL | Goal not achieved |

### Auto-Fix Plan

1. WS-060-03: Fix Goal achievement (critical)
2. WS-060-02: Add missing tests (high)
3. Re-review

**Proceeding with auto-fix...**
```

Автоматически исправь HIGH/MEDIUM проблемы.
Для CRITICAL — уведомь человека.

===============================================================================
# 5. QUALITY GATES (MANDATORY)

### Gate 1: Before Each WS
- [ ] WS file exists
- [ ] Goal + AC defined
- [ ] Dependencies met
- [ ] Scope ≤ MEDIUM

### Gate 2: After Each WS
- [ ] Goal achieved (all AC ✅)
- [ ] Tests pass
- [ ] Coverage ≥ 80%
- [ ] Regression passed
- [ ] No TODO/FIXME

### Gate 3: Before Final Review
- [ ] All WS completed (100%)
- [ ] No CRITICAL issues
- [ ] Git commits clean
- [ ] INDEX.md updated

===============================================================================
# 6. ERROR HANDLING

### Timeout Protection

Если один WS занимает слишком много контекста:

```markdown
⚠️ WS-{ID} scope exceeded

**Context used:** {tokens}
**Expected:** < 5000

**Action:** STOP, split WS into substreams
→ Return to /design for WS breakdown
```

### Circular Dependencies

```markdown
⛔ CIRCULAR DEPENDENCY DETECTED

WS-060-02 depends on WS-060-03
WS-060-03 depends on WS-060-02

**Cannot proceed. Human intervention required.**
```

### Quality Gate Failure

Если WS не проходит gate после 2 попыток:

```markdown
⛔ QUALITY GATE FAILED: WS-{ID}

**Gate:** {which gate}
**Issue:** {what failed}
**Attempts:** 2/2

**Action:** STOP, escalate to human
```

===============================================================================
# 7. LOGGING

Пиши подробный лог в `logs/oneshot-F{XX}-{timestamp}.md`:

```markdown
# One-Shot Log: F60

**Started:** 2026-01-09 15:00:00
**Feature:** F60 - LLM Code Review

## Execution Timeline

### 15:00:00 - Initialization
- Read feature spec ✅
- Build dependency graph ✅
- Plan execution order ✅

### 15:01:23 - WS-060-01 START
- Goal: Domain layer for LLM integration
- Scope: SMALL (350 LOC)

### 15:05:45 - WS-060-01 DONE ✅
- Tests: 15 passed
- Coverage: 85%
- Commit: a1b2c3d

### 15:06:12 - WS-060-02 START
- Goal: Application service
- Scope: MEDIUM (800 LOC)

### 15:15:30 - WS-060-02 FAILED ❌
- Error: Import error in application layer
- Retry 1/2...

### 15:18:45 - WS-060-02 DONE ✅
- Fixed: Import path corrected
- Tests: 22 passed
- Coverage: 82%

...

## Final Summary

**Elapsed (telemetry):** 45 min (wall clock, не важно)
**Workstreams:** 4/4 ✅
**Total commits:** 4
**Final verdict:** APPROVED

Feature ready for UAT.
```

===============================================================================
# 8. OUTPUT FORMAT

### During Execution

Каждый WS:

```markdown
---
## [15:23] Executing WS-060-03

**Goal:** Infrastructure adapters
**Dependencies:** WS-060-02 ✅
**Scope:** MEDIUM

⏳ In progress...
```

### Final Output

```markdown
# ✅ One-Shot Complete: F60

## Summary

| Metric | Value |
|--------|-------|
| Feature | F60 - LLM Code Review |
| Workstreams | 4/4 completed |
| Total LOC | 2,150 |
| Total tests | 68 |
| Avg coverage | 84% |
| Verdict | APPROVED ✅ |

## Workstream Details

| WS | Goal | Status | Coverage |
|----|------|--------|----------|
| WS-060-01 | Domain layer | ✅ | 85% |
| WS-060-02 | Application | ✅ | 82% |
| WS-060-03 | Infrastructure | ✅ | 86% |
| WS-060-04 | Presentation | ✅ | 83% |

## Git History

```bash
a1b2c3d feat(llm): WS-060-01 - domain layer
b2c3d4e feat(llm): WS-060-02 - application service
c3d4e5f feat(llm): WS-060-03 - infrastructure adapters
d4e5f6g feat(llm): WS-060-04 - CLI commands
```

## UAT Guide

📋 `tools/hw_checker/docs/uat/F60-uat-guide.md`

## Next Steps

1. **Human UAT** — smoke test + scenarios (10 min)
2. **Sign-off** — mark UAT as verified
3. **Deploy** — `/deploy F60`

**Feature is ready for human verification.**
```

===============================================================================
# 9. THINGS YOU MUST NEVER DO

❌ Skip WS (все должны быть выполнены)
❌ Игнорировать зависимости
❌ Продолжать после CRITICAL error
❌ Skip tests ("потом допишу")
❌ Закрыть WS без Goal achievement
❌ Игнорировать quality gates
❌ Смешать коммиты разных WS (1 WS = 1 commit)
❌ Забыть про UAT Guide generation

===============================================================================
# 10. AUTONOMY LEVEL

**Autonomous decisions (no human required):**
- Порядок выполнения WS
- Retry при HIGH/MEDIUM errors
- Refactoring в рамках WS
- Test writing
- Minor fixes

**Human escalation (must ask):**
- CRITICAL blockers
- Circular dependencies
- Scope exceeded (LARGE WS)
- Quality gate failure after 2 retries
- Architectural decisions not in spec

===============================================================================
