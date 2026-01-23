# /debug Runbook

**Purpose:** Systematic debugging workflow for bug analysis and root cause isolation.

**Prerequisites:**
- Bug symptoms observed
- `/debug` command available
- Access to logs and codebase

## When to Use

- Unknown root cause
- Bug symptoms inconsistent
- Need systematic investigation
- Before `/hotfix` or `/bugfix`

## Workflow

### Phase 1: Gather Context

```bash
# 1. Document symptoms clearly
# What: What exactly is happening?
# When: When does it occur?
# Consistency: Always/sometimes/rarely?
# Expected: What should happen instead?

# 2. Collect evidence
tail -100 logs/app.log | grep ERROR
grep -A 10 "Traceback" logs/*.log
git log --oneline --since="24 hours ago"

# 3. Run /debug command
/debug "API returns 500 on /submissions endpoint"
```

### Phase 2: Systematic Debugging (5 Phases)

The `/debug` command implements 5-phase workflow:

**Phase 1: Symptom Documentation**
- Document observed behavior precisely
- Note timing and consistency
- State expected vs actual behavior
- Collect evidence (logs, traces, changes)

**Phase 2: Hypothesis Formation**
- List 3+ possible root causes
- Rank by probability (HIGH/MEDIUM/LOW)
- Provide supporting evidence
- Suggest quick tests

**Phase 3: Systematic Elimination**
- Test each hypothesis
- Use specific verification commands
- Document results (CONFIRMED/REJECTED)
- Collect evidence

**Phase 4: Root Cause Isolation**
- Document precisely:
  - What: root cause description
  - Where: file, line, function
  - Why: step-by-step failure chain
  - Why Not Caught: missing tests, race condition, edge case

**Phase 5: Impact Chain Analysis**
- Analyze affected components
- Determine severity (P0 CRITICAL, P1 HIGH, P2 MEDIUM, P3 LOW)
- Assess business impact

### Phase 3: Failsafe Rule (3 Strikes)

Track debugging attempts:

| Attempt | Outcome | Notes |
|---------|---------|-------|
| 1 | Fix Attempt #1 | {description} |
| 2 | Fix Attempt #2 | {description} |
| 3 | Fix Attempt #3 | {description} |

**After 3 failed attempts:**
- Create bug report in `tools/hw_checker/docs/issues/issue-{ID}.md`
- Document all attempts and outcomes
- Escalate to human developer
- Route to `/hotfix` (if P0) or `/bugfix` (if P1/P2)

## Example Session

```
/debug "API returns 500 on /submissions endpoint"
```

**Expected Output:**

```markdown
## Phase 1: Symptom Documentation

**Observed Behavior:**
- API returns 500 on POST /submissions
- Occurs when submission size > 10MB
- Consistent: always fails for large files

**Expected Behavior:**
- API should accept submissions up to 100MB

**Evidence:**
- Error logs show "MemoryError: Unable to allocate"
- Recent commit: "b1a2c3f" increased payload validation

## Phase 2: Hypothesis Formation

### Hypothesis #1: Memory Limit Too Low (HIGH - 70%)
**Supporting Evidence:**
- Error shows "Unable to allocate memory"
- Recent commit changed payload validation

**Quick Test:**
```bash
docker stats hw_checker_api
# Check memory usage
```

### Hypothesis #2: Database Connection Timeout (MEDIUM - 20%)
**Supporting Evidence:**
- Large submissions take longer
- Could be hitting connection timeout

**Quick Test:**
```bash
grep "connection timeout" logs/app.log
```

## Phase 3: Systematic Elimination

### Testing Hypothesis #1

**Method:**
```bash
docker stats hw_checker_api
```

**Result:** ✅ CONFIRMED

**Evidence Found:**
- Memory usage hits 512MB limit on > 10MB files
- Container OOM kill event in logs

**Conclusion:**
This is the root cause. Memory limit is too low for large file uploads.

## Phase 4: Root Cause Identified

**What:**
Memory limit set to 512MB, insufficient for > 10MB file uploads

**Where:**
- File: `docker-compose.yml`
- Line: 15
- Function: N/A

**Why:**
1. Large files require ~5x memory for processing
2. Memory limit not updated when file size increased
3. Recent commit changed validation without updating memory limit

**Why Not Caught:**
- [x] Missing test case: large file upload
- [ ] Race condition
- [ ] Edge case not considered
- [ ] Configuration issue
- [ ] External dependency failure

## Phase 5: Impact Chain Analysis

**Affected Components:**
- API server (memory crash)
- User submissions (all failing for large files)
- Downstream grading (no submissions to grade)

**Severity:** P1 HIGH

**Business Impact:**
- Users cannot submit large assignments
- Grading system not receiving submissions
- Production degradation (partial outage)
```

## Routing After Debugging

### P0 CRITICAL

```bash
/hotfix "increase memory limit to 2GB" --issue-id=001
```

**When:**
- System down
- Data loss
- Security breach
- Critical production issue

### P1 HIGH

```bash
/bugfix "handle large submissions gracefully" --feature=F23 --issue-id=002
```

**When:**
- Feature broken
- Production degradation
- Data integrity issue

### P2 MEDIUM

```bash
/bugfix "improve error message for memory limit" --feature=F23
```

**When:**
- Non-critical bug
- UX issue
- Performance problem

### P3 LOW

Create GitHub issue for future work

## Integration with Other Commands

```
/issue "API returns 500 on submissions"
    ↓
/debug "investigate root cause"
    ↓
[Root cause found]
    ↓
/hotfix (if P0) or /bugfix (if P1/P2)
    ↓
/codereview F{XX}
    ↓
/deploy F{XX}
```

## Troubleshooting

### /debug not found

**Problem:** Command not available

**Solution:**
- Claude Code: Check `.claude/skills/debug/SKILL.md` exists
- Cursor: Check `.cursor/commands/debug.md` exists
- OpenCode: Check `.opencode/commands/debug.md` exists

### Debugging hangs on one hypothesis

**Problem:** Stuck testing same hypothesis repeatedly

**Solution:**
1. Check if fix attempts are being tracked
2. Verify 3-strikes rule is enforced
3. Manually escalate to human if stuck
4. Create bug report and move to `/bugfix`

### Can't determine severity

**Problem:** Unsure whether P0, P1, P2, or P3

**Solution:**

| Severity | Criteria | Example |
|----------|-----------|---------|
| **P0 CRITICAL** | System down, data loss, security | API returning 500 for all users |
| **P1 HIGH** | Feature broken, degradation | Large file uploads failing |
| **P2 MEDIUM** | Non-critical bug, UX issue | Poor error message |
| **P3 LOW** | Enhancement, minor issue | Spelling error in docs |

## Best Practices

1. **Document everything:** Don't skip symptom documentation
2. **Test hypotheses:** Don't guess, verify with evidence
3. **Use failsafe:** Stop after 3 attempts, don't infinite loop
4. **Escalate appropriately:** P0 → /hotfix, P1/P2 → /bugfix
5. **Create bug reports:** Document unresolved issues for future

## References

- Master prompt: `sdp/prompts/commands/issue.md` → Section 4.0
- Issue command: `sdp/prompts/commands/issue.md`
- Hotfix command: `sdp/prompts/commands/hotfix.md`
- Bugfix command: `sdp/prompts/commands/bugfix.md`
