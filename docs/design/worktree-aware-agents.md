# Worktree-Aware Agent Behavior - Design Document

## Problem Statement

Agents ignore git worktree setup and continue working in the main repository branch, violating GitFlow isolation principles.

### Root Cause Analysis

**Issue**: WS 00-052-00 created worktree at `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent` for isolated F052 implementation, but subsequent workstreams (00-052-01 through 00-052-07) executed in main dev branch instead.

**Root Cause**: No agent skills (@build, @oneshot) implement worktree detection or switching logic.

## Investigation Findings

### Current Behavior

1. **Worktree Creation** (WS 00-052-00):
   - Worktree created successfully: `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`
   - Location documented in: `docs/workstreams/WORKTREE.md`
   - Branch: `dev`

2. **Subsequent Execution** (WS 00-052-01 through 00-052-07):
   - Executed in main repository: `/Users/fall_out_bug/projects/vibe_coding/sdp`
   - 7 commits landed in `dev` instead of worktree
   - No worktree detection or switching occurred

3. **Skill Analysis**:
   - **@build skill**: No worktree detection logic
   - **@oneshot skill**: No worktree switching logic
   - **@feature skill**: No worktree creation step
   - **@design skill**: No worktree awareness

### Why This Happened

1. **Session Context Loss**: Worktree location documented in WS file, but not passed to agent session
2. **No Discovery Logic**: Agents don't check for `docs/workstreams/WORKTREE.md`
3. **No Git Worktree Command**: Agents don't run `git worktree list` to detect worktrees
4. **Working Directory Assumption**: Agents assume single working directory (main repo)

## Design Requirements

### Functional Requirements

**FR1: Worktree Discovery**
- Agents MUST check for worktree existence before executing workstreams
- Agents MUST read `docs/workstreams/WORKTREE.md` if present
- Agents MUST run `git worktree list` to detect active worktrees

**FR2: Worktree Switching**
- Agents MUST change to worktree directory if detected
- Agents MUST execute all commands in worktree context
- Agents MUST verify worktree branch matches expected branch

**FR3: Worktree Documentation**
- WS that creates worktree MUST document location in standard location
- WS that creates worktree MUST include cleanup instructions
- WS documentation MUST be machine-readable for agent consumption

**FR4: Fallback Behavior**
- If no worktree found, continue in main repository
- If worktree corrupted, warn and fall back to main repository
- Document fallback decision in execution report

### Non-Functional Requirements

**NFR1: Backward Compatibility**
- Existing repositories without worktree MUST continue working
- No breaking changes to existing workflows

**NFR2: Performance**
- Worktree detection MUST add <100ms overhead
- Worktree switching MUST be idempotent (safe to call multiple times)

**NFR3: Reliability**
- Worktree corruption MUST NOT block agent execution
- Worktree switching failures MUST be logged and reported

**NFR4: GitFlow Compliance**
- Feature branches MUST use worktrees when available
- Dev branch MUST remain stable during feature development

## Proposed Solution

### Phase 1: Worktree Discovery (Immediate)

**Add to @build skill:**

```go
// DetectWorktree checks for worktree existence
func DetectWorktree(workDir string) (string, error) {
    // 1. Check for WORKTREE.md documentation
    worktreeDoc := filepath.Join(workDir, "docs/workstreams/WORKTREE.md")
    if info, err := os.Stat(worktreeDoc); err == nil && !info.IsDir() {
        // Parse worktree location from documentation
        location, err := parseWorktreeDoc(worktreeDoc)
        if err == nil {
            return location, nil
        }
    }

    // 2. Check git worktree list
    cmd := exec.Command("git", "worktree", "list", "--porcelain")
    cmd.Dir = workDir
    output, err := cmd.Output()
    if err != nil {
        return "", nil // No worktree, not an error
    }

    // Parse output for worktree paths
    // Return first worktree (if multiple, need disambiguation)
    return parseWorktreeList(output)
}
```

**Add to @oneshot skill:**

```markdown
## Pre-Execution Check

Before spawning @build agents:
1. Detect worktree location
2. If worktree exists, set WORKTREE_DIR environment variable
3. Pass worktree path to all agent spawns
```

### Phase 2: Worktree Switching (Short-term)

**Modify @build skill workflow:**

```markdown
### Step 0: Worktree Detection (NEW)

```bash
# Detect worktree
WORKTREE=$(detect-worktree)
if [ -n "$WORKTREE" ]; then
    echo "✓ Worktree detected: $WORKTREE"
    cd "$WORKTREE" || {
        echo "⚠️ Cannot access worktree, falling back to main repo"
        WORKTREE=""
    }
else
    echo "ℹ️ No worktree, using main repository"
fi

# Verify we're on correct branch
CURRENT_BRANCH=$(git branch --show-current)
if [ -n "$WORKTREE" ]; then
    echo "Current branch: $CURRENT_BRANCH"
fi
```

**Add to workstream template:**

```markdown
## Worktree Configuration

**Worktree Path:** `/Users/fall_out_bug/projects/vibe_coding/sdp-f052-impl`
**Expected Branch:** `feature/f052-implementation`
**Fallback Branch:** `dev`

**If worktree not found:** Create worktree first
```

### Phase 3: Worktree Creation (Long-term)

**Add to @design skill:**

```markdown
## Worktree Strategy

For features requiring isolation:
1. Create WS: `{project}-000-worktree-setup`
2. Worktree location: `../{project}-{feature}-impl`
3. Branch: `feature/{feature}-implementation`
4. Document in: `docs/workstreams/WORKTREE.md`

Subsequent WS automatically use worktree.
```

**Add to @feature skill:**

```markdown
### Question: Isolation Required?

- [ ] Feature requires isolated branch (GitFlow)
- [ ] Feature modifies core infrastructure
- [ ] Feature has high risk of breaking changes

If YES, add worktree setup WS as first workstream.
```

## Implementation Plan

### Workstream 1: Worktree Discovery (sdp-abmu)

**Goal:** Add worktree detection to @build and @oneshot skills

**Acceptance Criteria:**
- AC1: @build checks for `docs/workstreams/WORKTREE.md`
- AC2: @build runs `git worktree list` as fallback
- AC3: @oneshot passes worktree context to spawned agents
- AC4: Tests verify worktree detection works

**Files:**
- Modify: `.claude/skills/build/SKILL.md` (add worktree detection step)
- Modify: `.claude/skills/oneshot/SKILL.md` (add worktree context passing)
- Create: `src/sdp/git/worktree.go` (worktree detection utilities)
- Create: `src/sdp/git/worktree_test.go` (tests)

**Priority:** P1 (blocks proper GitFlow compliance)

### Workstream 2: Worktree Switching (future)

**Goal:** Add worktree switching to @build skill

**Acceptance Criteria:**
- AC1: @build changes to worktree directory before executing
- AC2: @build verifies worktree branch
- AC3: @build falls back to main repo on error
- AC4: Tests verify switching works

**Files:**
- Modify: `.claude/skills/build/SKILL.md` (add worktree switching logic)
- Modify: `src/sdp/git/worktree.go` (add switching function)
- Tests for switching and fallback

**Priority:** P2 (important but not blocking)

### Workstream 3: Worktree Documentation (future)

**Goal:** Standardize worktree documentation format

**Acceptance Criteria:**
- AC1: WS template includes worktree section
- AC2: WORKTREE.md template standardized
- AC3: @design skill generates worktree WS
- AC4: Tests verify documentation parsing

**Files:**
- Create: `docs/templates/WORKTREE.md.template`
- Modify: `.claude/skills/design/SKILL.md` (add worktree generation)
- Create: `src/sdp/git/worktree_parser.go` (parse WORKTREE.md)
- Tests for parsing

**Priority:** P3 (nice to have)

## Testing Strategy

### Unit Tests

**Worktree Detection:**
```go
func TestDetectWorktree_FromDocumentation(t *testing.T) {
    // Test reading docs/workstreams/WORKTREE.md
    workDir := setupTestRepo(t)
    createWorktreeDoc(t, workDir, "/path/to/worktree")

    detected, err := DetectWorktree(workDir)
    assert.NoError(t, err)
    assert.Equal(t, "/path/to/worktree", detected)
}

func TestDetectWorktree_FromGitCommand(t *testing.T) {
    // Test git worktree list parsing
    workDir := setupTestRepo(t)
    createWorktree(t, workDir, "/path/to/worktree")

    detected, err := DetectWorktree(workDir)
    assert.NoError(t, err)
    assert.Equal(t, "/path/to/worktree", detected)
}

func TestDetectWorktree_NotFound(t *testing.T) {
    // Test no worktree scenario
    workDir := setupTestRepo(t)

    detected, err := DetectWorktree(workDir)
    assert.NoError(t, err)
    assert.Equal(t, "", detected)
}
```

**Worktree Switching:**
```go
func TestSwitchToWorktree_Success(t *testing.T) {
    workDir := setupTestRepo(t)
    createWorktree(t, workDir, "/path/to/worktree")

    err := SwitchToWorktree(workDir, "/path/to/worktree")
    assert.NoError(t, err)
    // Verify current directory changed
}

func TestSwitchToWorktree_NotFound(t *testing.T) {
    workDir := setupTestRepo(t)

    err := SwitchToWorktree(workDir, "/path/to/nonexistent")
    assert.Error(t, err)
    // Verify fallback to main repo
}
```

### Integration Tests

**End-to-End Worktree Workflow:**
```go
func TestWorktreeWorkflow_Full(t *testing.T) {
    // 1. Create worktree
    workDir := setupTestRepo(t)
    worktreePath := createWorktree(t, workDir, "feature-branch")

    // 2. Document worktree
    createWorktreeDoc(t, workDir, worktreePath)

    // 3. Run @build in worktree context
    agent := NewBuildAgent()
    result := agent.Execute("00-001-01")

    // 4. Verify work in worktree, not main repo
    assert.Equal(t, worktreePath, result.WorkDir)
    assert.NoFileExists(t, filepath.Join(workDir, "newfile.go"))
    assert.FileExists(t, filepath.Join(worktreePath, "newfile.go"))
}
```

### Manual Testing

**Scenario 1: New Feature with Worktree**
1. Run @feature to create workstreams
2. Run WS 00-XXX-00 to create worktree
3. Run @oneshot to execute remaining WS
4. Verify: All commits in worktree, not main dev

**Scenario 2: Existing Feature without Worktree**
1. Run @build on existing WS
2. Verify: Works in main repo (backward compatibility)

**Scenario 3: Worktree Corruption**
1. Create worktree
2. Corrupt worktree (delete .git directory)
3. Run @build
4. Verify: Falls back to main repo with warning

## Rollout Plan

### Phase 1: Investigation (COMPLETE)
- [x] Root cause identified
- [x] Design document created
- [x] Implementation plan defined

### Phase 2: Worktree Discovery (CURRENT)
- [ ] Implement DetectWorktree() function
- [ ] Add detection to @build skill
- [ ] Add detection to @oneshot skill
- [ ] Write tests
- [ ] Merge to dev

### Phase 3: Worktree Switching (FUTURE)
- [ ] Implement SwitchToWorktree() function
- [ ] Add switching to @build skill
- [ ] Write tests
- [ ] Merge to dev

### Phase 4: Worktree Creation (FUTURE)
- [ ] Add worktree generation to @design
- [ ] Create WORKTREE.md template
- [ ] Update WS templates
- [ ] Merge to dev

## Success Metrics

- **Worktree Detection Rate**: 100% (when worktree exists)
- **Worktree Switching Success Rate**: >95%
- **Backward Compatibility**: 100% (no regressions)
- **GitFlow Compliance**: 100% (features use worktrees)

## References

- Git Worktree Documentation: https://git-scm.com/docs/git-worktree
- F052 Workstream: docs/workstreams/completed/00-052-00-backup-and-worktree.md
- GitFlow Best Practices: docs/reference/gitflow.md
- Bug Report: sdp-abmu

---

**Document Version:** 1.0
**Last Updated:** 2026-02-08
**Status:** Draft - Ready for Implementation
**Author:** Investigation Agent (sdp-abmu)
