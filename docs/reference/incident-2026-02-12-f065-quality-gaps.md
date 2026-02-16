# Incident Analysis: F065 Quality Gaps (2026-02-12)

## Summary

F065 (Agent Git Safety Protocol) was executed by orchestrator but failed review due to quality gaps that were not caught during execution.

## Problems Found

| Issue | Severity | Detection Point | Root Cause |
|-------|----------|-----------------|------------|
| LOC violations (5 files) | P1 | Review | No LOC gate in @build |
| Code duplication | P2 | Review | No duplication check |
| Missing logging | P1 | Review | Not in AC/quality gates |
| Partial worktree support | P2 | Review | No worktree tests |

## Root Cause Analysis

### 1. LOC Violations Not Caught

**What happened:**
- wrapper.go: 311 LOC (limit: 200)
- recovery.go: 342 LOC (limit: 200)
- cmd/session.go: 280 LOC
- internal/session.go: 237 LOC
- guard_context.go: 247 LOC

**Why @build didn't catch it:**
```
@build quality gates:
1. Tests pass ✅
2. Coverage ≥ 80% ✅
3. Linters clean ✅
4. ❌ NO LOC CHECK
```

**Fix:** Add LOC check to @build quality gates

### 2. Code Duplication Not Detected

**What happened:**
- `getCurrentBranch()` exists in 3 files
- Each copy slightly different

**Why @build didn't catch it:**
- No duplication detection in quality gates
- Agents implement locally without checking for existing utilities

**Fix:** Add `dupl` linter or manual duplication check

### 3. Missing Requirements Not in AC

**What happened:**
- Logging was expected by SRE review but not in any workstream AC
- Worktree support was mentioned but not tested

**Why:**
- ACs focused on functionality, not operational concerns
- No "non-functional requirements" checklist

**Fix:** Add operational requirements to all workstreams (logging, error handling, testing)

### 4. Orchestrator Continued Past Quality Issues

**What happened:**
- Orchestrator completed all 6 workstreams
- Review found FAIL on multiple dimensions
- Had to create fix-up issues

**Why:**
- Orchestrator only stops for "CRITICAL blockers"
- LOC violations are not considered CRITICAL
- Quality is verified at END, not during execution

**Fix:** Pre-commit quality gates that block commits

## Prevention Mechanisms

### Level 1: @build Quality Gates (CRITICAL)

Add to `.claude/skills/build/SKILL.md`:

```markdown
## Quality Gates (MANDATORY)

Before commit, ALL gates must pass:

1. **Tests**: `go test ./... -short` → 0 failures
2. **Coverage**: `go test -cover` → ≥ 80%
3. **Linters**: `go vet ./...` → 0 issues
4. **LOC Check**: `wc -l <file>` → ≤ 200 lines
5. **Duplication**: `dupl` or manual check → ≤ 5%

### LOC Gate Script

```bash
# For each scope_file
for file in $SCOPE_FILES; do
  loc=$(wc -l < "$file")
  if [ "$loc" -gt 200 ]; then
    echo "ERROR: $file is $loc LOC (max: 200)"
    echo "Split into smaller files before committing"
    exit 1
  fi
done
```

### Failing Gate = BLOCK

If any gate fails:
1. DO NOT commit
2. Fix the issue
3. Re-run all gates
4. Only then commit
```

### Level 2: Guard Enforcement

Add to `sdp guard check`:

```go
func checkLOC(filePath string) error {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    lines := strings.Count(string(content), "\n")
    if lines > 200 {
        return fmt.Errorf("%s: %d LOC exceeds 200 limit", filePath, lines)
    }
    return nil
}
```

Hook into pre-commit:
```bash
# .git/hooks/pre-commit
sdp guard check --loc
```

### Level 3: Orchestrator Verification

Add to `.claude/agents/orchestrator.md`:

```markdown
## Per-Workstream Quality Verification

After each @build completes:

1. Run quality gates locally:
   ```bash
   # Verify LOC
   for f in $(cat $WS_FILE | grep scope_files -A10 | grep -oE 'sdp-plugin/[^ ]+'); do
     loc=$(wc -l < "$f" 2>/dev/null || echo 0)
     if [ "$loc" -gt 200 ]; then
       echo "ERROR: $f exceeds 200 LOC"
       exit 1
     fi
   done
   ```

2. If gate fails:
   - Log: "Quality gate failed: LOC violation in $f"
   - DO NOT proceed to next WS
   - Fix the issue or escalate

3. Only proceed when ALL gates pass
```

### Level 4: Workstream Template Requirements

Add to all workstream templates:

```markdown
## Non-Functional Requirements

- [ ] Logging: Key operations logged
- [ ] Error handling: All errors wrapped with context
- [ ] Testing: Edge cases covered
- [ ] Documentation: Public APIs documented
- [ ] LOC: All files < 200 lines
```

## Immediate Actions

| Action | Owner | Status |
|--------|-------|--------|
| Add LOC check to @build | TBD | Pending |
| Create sdp guard --loc command | TBD | Pending |
| Update orchestrator.md | TBD | Pending |
| Update workstream template | TBD | Pending |
| Fix remaining LOC violations | TBD | In Progress |

## Lessons Learned

1. **Quality gates must be automated** - Manual review is too late
2. **LOC is a proxy for complexity** - Large files = hard to maintain
3. **Non-functional requirements matter** - Add to every WS
4. **Verify during execution, not after** - Block early, fix early

## Related

- Incident: docs/reference/incident-2026-02-12-agent-branch-confusion.md
- Design: docs/plans/2026-02-12-agent-git-safety-protocol.md
- Beads issues: sdp-12fj, sdp-49xl, sdp-talt, sdp-7ov9, sdp-sfgy, sdp-ilx3, sdp-bxw6, sdp-t7vx
