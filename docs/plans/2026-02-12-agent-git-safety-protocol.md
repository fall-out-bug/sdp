# Agent Git Safety Protocol

> **Status:** Research complete
> **Date:** 2026-02-12
> **Goal:** Prevent agents from committing to wrong worktrees/branches

---

## Table of Contents

1. [Overview](#overview)
2. [Worktree State File](#1-worktree-state-file)
3. [Git Operation Wrapper](#2-git-operation-wrapper)
4. [CWD Recovery Protocol](#3-cwd-recovery-protocol)
5. [Agent Prompt Constraints](#4-agent-prompt-constraints)
6. [Git Hooks Protection](#5-git-hooks-protection)
7. [Worktree Locking](#6-worktree-locking)
8. [Validation Commands](#7-validation-commands)
9. [Implementation Plan](#implementation-plan)
10. [Success Metrics](#success-metrics)

---

## Overview

### Problem Statement

Agents working with git worktrees experience branch confusion:
1. **Branch tracking corruption** — feature branches track wrong remotes
2. **Cross-feature commits** — commits appear in wrong branch history
3. **Direct pushes to dev** — bypassing PR workflow
4. **CWD reset** — shell returns to main repo after Bash commands

### Goals

1. **Agent cannot commit to wrong worktree/branch** — enforced at multiple levels
2. **Automatic verification before each git operation** — no manual discipline required
3. **Feature isolation** — one feature = one worktree = one branch
4. **Recovery after CWD reset** — state persists across shell resets

### Key Decisions

| Aspect | Decision |
|--------|----------|
| State Storage | Per-worktree `.sdp/session.json` |
| Git Commands | Wrapped via `sdp git` with validation |
| CWD Recovery | Hybrid: git worktree list + WS metadata |
| Agent Guidance | GIT_SAFETY.md + inline checklists |
| Hooks | Pre-commit/pre-push with session validation |
| Locking | Per-worktree lock file with status |
| CLI Commands | `sdp guard context check/show/clean` |

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     DEFENSE IN DEPTH                            │
├─────────────────────────────────────────────────────────────────┤
│  Layer 1: Agent Prompts                                         │
│  - GIT_SAFETY.md guidelines                                     │
│  - Inline checklists in skills                                  │
│  - Context verification requirements                            │
├─────────────────────────────────────────────────────────────────┤
│  Layer 2: CLI Wrapper                                           │
│  - sdp git <command> validates before execution                 │
│  - sdp guard context check                                      │
│  - Automatic CWD recovery                                       │
├─────────────────────────────────────────────────────────────────┤
│  Layer 3: Git Hooks                                             │
│  - pre-commit: verify session matches branch                    │
│  - pre-push: verify branch tracking correct                     │
│  - post-checkout: update session file                           │
├─────────────────────────────────────────────────────────────────┤
│  Layer 4: State Files                                           │
│  - .sdp/session.json (worktree context)                         │
│  - .sdp/worktree-lock.json (active status)                      │
│  - Hash verification for tamper detection                       │
└─────────────────────────────────────────────────────────────────┘
```

---

## 1. Worktree State File

> **Expert:** Kelsey Hightower (DevOps/K8s)
> **Principles:** declarative config, immutable infrastructure, GitOps

### Solution: Per-Worktree Session State

Each worktree maintains a session file that defines its identity:

**File:** `.sdp/session.json`

```json
{
  "version": "1.0",
  "worktree_path": "/Users/user/projects/sdp-F063",
  "feature_id": "F063",
  "expected_branch": "feature/F063",
  "expected_remote": "origin/feature/F063",
  "created_at": "2026-02-12T10:00:00Z",
  "created_by": "sdp worktree create",
  "hash": "sha256:abc123..."
}
```

### Fields

| Field | Purpose | Validation |
|-------|---------|------------|
| `worktree_path` | Absolute path to this worktree | Must match actual location |
| `feature_id` | Feature this worktree serves | Must exist in workstream files |
| `expected_branch` | Branch name for this worktree | Must match `git branch --show-current` |
| `expected_remote` | Remote tracking branch | Must match `git rev-parse --abbrev-ref @{u}` |
| `hash` | Tamper detection | SHA256 of file content |

### Lifecycle

1. **Creation:** `sdp worktree create F063` generates session file
2. **Validation:** Every git operation validates session matches reality
3. **Update:** Manual switch requires explicit `sdp session update`
4. **Deletion:** `sdp worktree delete F063` removes session file

### Error Handling

```bash
# Session file missing
ERROR: No session file found. Run: sdp session init

# Branch mismatch
ERROR: Session expects branch 'feature/F063' but on 'dev'
FIX: sdp session sync  # or git checkout feature/F063

# Hash mismatch (tampering detected)
ERROR: Session file corrupted or tampered
FIX: sdp session repair --force
```

---

## 2. Git Operation Wrapper

> **Expert:** Troy Hunt (Security)
> **Principles:** defense in depth, least privilege, validate all inputs

### Solution: State-File Based Git Commands

All git operations go through `sdp git` wrapper:

```bash
# Instead of: git commit -m "message"
sdp git commit -m "message"

# Instead of: git push
sdp git push
```

### Validation Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    sdp git <command>                         │
├──────────────────────────────────────────────────────────────┤
│  1. Load .sdp/session.json                                   │
│  2. Verify worktree_path matches pwd                         │
│  3. Verify expected_branch matches current branch            │
│  4. Verify expected_remote matches tracking                  │
│  5. Verify hash is valid                                     │
│  6. Execute: git <command>                                   │
│  7. Verify still on expected_branch after command            │
└──────────────────────────────────────────────────────────────┘
```

### Command Categories

| Category | Commands | Validation |
|----------|----------|------------|
| **Safe** | status, log, diff, show | Session check only |
| **Write** | add, commit, reset | Full validation + post-check |
| **Remote** | push, fetch, pull | Remote tracking verification |
| **Branch** | checkout, branch, merge | Session update required |

### Implementation

```go
func (g *GitWrapper) Execute(cmd string, args []string) error {
    // 1. Load session
    session, err := g.loadSession()
    if err != nil {
        return fmt.Errorf("session required: %w", err)
    }

    // 2. Validate CWD
    cwd, _ := os.Getwd()
    if cwd != session.WorktreePath {
        return fmt.Errorf("wrong worktree: expected %s, in %s",
            session.WorktreePath, cwd)
    }

    // 3. Validate branch
    currentBranch := g.getCurrentBranch()
    if currentBranch != session.ExpectedBranch {
        return fmt.Errorf("wrong branch: expected %s, on %s",
            session.ExpectedBranch, currentBranch)
    }

    // 4. Execute command
    if err := g.runGitCommand(cmd, args...); err != nil {
        return err
    }

    // 5. Post-check for write commands
    if g.isWriteCommand(cmd) {
        newBranch := g.getCurrentBranch()
        if newBranch != session.ExpectedBranch {
            return fmt.Errorf("CRITICAL: branch changed during command!")
        }
    }

    return nil
}
```

---

## 3. CWD Recovery Protocol

> **Expert:** Martin Kleppmann (Distributed Systems)
> **Principles:** eventual consistency, idempotency, partition tolerance

### Solution: Hybrid Recovery Strategy

When CWD resets to main repo, recover worktree context from multiple sources:

```
┌──────────────────────────────────────────────────────────────┐
│                  CWD Recovery Protocol                       │
├──────────────────────────────────────────────────────────────┤
│  Source 1: .sdp/session.json (primary)                       │
│  - Fastest, most reliable                                    │
│  - Contains full context                                     │
│                                                              │
│  Source 2: git worktree list (fallback)                      │
│  - Lists all worktrees with paths                            │
│  - Match by feature ID in path                               │
│                                                              │
│  Source 3: Workstream metadata (last resort)                 │
│  - docs/workstreams/backlog/00-FEATURE-*.md                  │
│  - Extract feature_id from frontmatter                        │
└──────────────────────────────────────────────────────────────┘
```

### Recovery Commands

```bash
# Check current context
sdp guard context check
# Output: ✅ In worktree /path/to/sdp-F063, branch feature/F063

# Find worktree for feature
sdp guard context find F063
# Output: /path/to/sdp-F063

# Change to correct worktree
sdp guard context go F063
# Equivalent to: cd $(sdp guard context find F063)
```

### Agent Integration

Agents should call recovery before git operations:

```bash
# Before any git command
sdp guard context check || sdp guard context go $EXPECTED_FEATURE
git status
```

### Idempotency

Recovery is idempotent - can run multiple times safely:
- If already in correct worktree: no-op
- If in wrong worktree: error with guidance
- If CWD reset: recover to expected worktree

---

## 4. Agent Prompt Constraints

> **Expert:** Dan Abramov (React/State)
> **Principles:** single responsibility, lift state only when needed, colocation

### Solution: GIT_SAFETY.md + Inline Checklists

**File:** `.claude/GIT_SAFETY.md`

```markdown
# Git Safety Guidelines for Agents

## CRITICAL RULES

1. **NEVER run git commands without verifying context**
2. **ALWAYS check pwd before git operations**
3. **ALWAYS verify branch before commits**
4. **NEVER push to branches you weren't assigned to**

## Pre-Flight Checklist

Before ANY git operation, run:

```bash
pwd                    # Verify worktree path
git branch --show-current   # Verify branch name
sdp guard context check     # Validate session
```

## Required Pattern

```bash
# CORRECT
sdp guard context check && git commit -m "message"

# WRONG - no verification
git commit -m "message"

# WRONG - assumes CWD is correct
cd /path && git commit  # CWD resets after!
```

## Error Recovery

If context check fails:
1. DO NOT proceed with git command
2. Run: sdp guard context find $EXPECTED_FEATURE
3. Run: sdp guard context go $EXPECTED_FEATURE
4. Retry git command

## Forbidden Commands

- `git push --force` — NEVER
- `git checkout main` without explicit user request
- `git merge` across feature branches
- `git rebase` without user approval
```

### Skill Integration

Add inline checklists to relevant skills:

```markdown
## Git Operations (from @build, @bugfix, @deploy)

**MANDATORY before any git command:**

```bash
# Step 1: Verify context
pwd
git branch --show-current
sdp guard context check

# Step 2: If check fails, recover
sdp guard context go $FEATURE_ID

# Step 3: Only then proceed
git add .
git commit -m "..."
```
```

### Verification in Agent Prompts

Orchestrator and builder agents must include:

```
You are working in worktree: /path/to/sdp-F063
Feature: F063
Expected branch: feature/F063

BEFORE any git operation:
1. Run: sdp guard context check
2. If fails: Run: sdp guard context go F063
3. Then proceed with git command

NEVER skip these steps. Your CWD may reset after tool calls.
```

---

## 5. Git Hooks Protection

> **Expert:** Troy Hunt (Security)
> **Principles:** defense in depth, least privilege, validate all inputs

### Solution: Pre-Commit/Pre-Push Validation

**Hook:** `.git/hooks/pre-commit`

```bash
#!/bin/bash
set -e

# Get current context
CURRENT_BRANCH=$(git branch --show-current)
CURRENT_DIR=$(pwd)

# Check for session file
if [ -f ".sdp/session.json" ]; then
    EXPECTED_BRANCH=$(jq -r '.expected_branch' .sdp/session.json)
    EXPECTED_DIR=$(jq -r '.worktree_path' .sdp/session.json)

    # Validate branch
    if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "ERROR: Branch mismatch!"
        echo "  Expected: $EXPECTED_BRANCH"
        echo "  Current:  $CURRENT_BRANCH"
        echo ""
        echo "Run: sdp session sync"
        exit 1
    fi

    # Validate directory
    if [ "$CURRENT_DIR" != "$EXPECTED_DIR" ]; then
        echo "ERROR: Directory mismatch!"
        echo "  Expected: $EXPECTED_DIR"
        echo "  Current:  $CURRENT_DIR"
        exit 1
    fi
fi

# Check for cross-feature commits
PARENT_BRANCH=$(git rev-parse --abbrev-ref HEAD@{upstream} 2>/dev/null || echo "")
if [[ "$PARENT_BRANCH" == "origin/dev" && "$CURRENT_BRANCH" != "dev" ]]; then
    # Feature branch should track origin/feature/X, not origin/dev
    echo "WARNING: Feature branch tracking origin/dev instead of origin/$CURRENT_BRANCH"
    echo "Run: git branch --set-upstream-to=origin/$CURRENT_BRANCH"
fi

exit 0
```

**Hook:** `.git/hooks/pre-push`

```bash
#!/bin/bash
set -e

# Get push target
REMOTE="$1"
URL="$2"

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)

# Check session
if [ -f ".sdp/session.json" ]; then
    EXPECTED_REMOTE=$(jq -r '.expected_remote' .sdp/session.json)

    # Build expected remote ref
    EXPECTED_PUSH_TARGET="origin/$CURRENT_BRANCH"

    # Validate push target
    if [ "$EXPECTED_REMOTE" != "$EXPECTED_PUSH_TARGET" ]; then
        echo "ERROR: Push target mismatch!"
        echo "  Expected: $EXPECTED_REMOTE"
        echo "  Current would push to: $EXPECTED_PUSH_TARGET"
        echo ""
        echo "Fix branch tracking:"
        echo "  git branch --set-upstream-to=$EXPECTED_REMOTE"
        exit 1
    fi
fi

# Prevent pushing to protected branches
if [[ "$CURRENT_BRANCH" == "main" || "$CURRENT_BRANCH" == "dev" ]]; then
    echo "ERROR: Direct push to $CURRENT_BRANCH is not allowed!"
    echo "Create a feature branch and use PR workflow."
    exit 1
fi

exit 0
```

### Hook Installation

```bash
# Install hooks (run once per worktree)
sdp hooks install

# Or manually
cp hooks/pre-commit .git/hooks/
cp hooks/pre-push .git/hooks/
chmod +x .git/hooks/pre-commit .git/hooks/pre-push
```

---

## 6. Worktree Locking

> **Expert:** Sam Newman (Architecture)
> **Principles:** bounded context, single responsibility, loose coupling

### Solution: Per-Worktree Lock File

**File:** `.sdp/worktree-lock.json`

```json
{
  "version": "1.0",
  "status": "active",
  "feature_id": "F063",
  "locked_at": "2026-02-12T10:00:00Z",
  "locked_by": "agent-20260212-100000",
  "expires_at": "2026-02-12T18:00:00Z",
  "agent_id": "agent-20260212-100000",
  "session_id": "sess-abc123"
}
```

### Lock States

| State | Meaning | Actions Allowed |
|-------|---------|-----------------|
| `active` | Agent working on feature | All git operations |
| `paused` | Agent taking break | Read-only operations |
| `completed` | Feature done, ready for cleanup | None (archive) |
| `expired` | Lock timed out | Can be claimed by new agent |

### Lock Commands

```bash
# Claim worktree for feature
sdp lock claim F063 --agent=$AGENT_ID

# Check lock status
sdp lock status
# Output: F063: active (agent-123, expires 18:00)

# Release lock
sdp lock release F063

# Extend lock (for long operations)
sdp lock extend F063 --duration=4h

# Force unlock (admin only)
sdp lock force-unlock F063 --reason="stale lock"
```

### Conflict Resolution

When two agents try to claim same worktree:

```
Agent A: sdp lock claim F063
→ Creates lock with agent-A

Agent B: sdp lock claim F063
→ ERROR: F063 locked by agent-A (expires 18:00)
→ Options:
   1. Wait for expiry
   2. Ask agent-A to release
   3. Use force-unlock (requires admin)
```

---

## 7. Validation Commands

> **Expert:** Kelsey Hightower (DevOps/K8s)
> **Principles:** declarative config, immutable infrastructure, GitOps

### Solution: `sdp guard context` Command Group

```bash
# Check current context (quick validation)
sdp guard context check
# Output:
#   ✅ Worktree: /path/to/sdp-F063
#   ✅ Branch: feature/F063
#   ✅ Tracking: origin/feature/F063
#   ✅ Session: valid

# Show full context details
sdp guard context show
# Output:
#   Worktree Path: /path/to/sdp-F063
#   Feature ID: F063
#   Current Branch: feature/F063
#   Expected Branch: feature/F063
#   Remote Tracking: origin/feature/F063
#   Session Valid: yes
#   Lock Status: active (agent-123)

# Clean up stale sessions
sdp guard context clean
# Removes expired locks and invalid sessions

# Repair corrupted session
sdp guard context repair
# Rebuilds session.json from git state
```

### Integration Points

| Command | When to Use |
|---------|-------------|
| `sdp guard context check` | Before every git operation |
| `sdp guard context show` | Debugging context issues |
| `sdp guard context find F063` | Locate worktree for feature |
| `sdp guard context go F063` | Change to feature worktree |
| `sdp guard context clean` | Periodic cleanup |
| `sdp guard context repair` | After manual git operations |

### Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| 0 | All checks pass | Proceed |
| 1 | Context mismatch | Run recovery |
| 2 | No session file | Initialize session |
| 3 | Lock conflict | Wait or resolve |
| 4 | Hash mismatch | Repair session |

---

## Implementation Plan

### Phase 1: Foundation (MVP)

**Goal:** Basic session tracking and validation

- [ ] Implement `.sdp/session.json` format
- [ ] Create `sdp session init` command
- [ ] Create `sdp session sync` command
- [ ] Add `sdp guard context check` command
- [ ] Update `sdp worktree create` to initialize session
- [ ] Write unit tests for session validation

**Files to create/modify:**
```
sdp-plugin/
├── cmd/sdp/
│   ├── session.go          # NEW: session commands
│   └── guard_context.go    # NEW: context validation
├── internal/
│   └── session/
│       ├── session.go      # NEW: session management
│       └── validator.go    # NEW: validation logic
```

### Phase 2: Git Wrapper

**Goal:** Safe git operations through wrapper

- [ ] Implement `sdp git` wrapper command
- [ ] Add pre/post validation for write commands
- [ ] Add remote tracking validation
- [ ] Integrate with session management
- [ ] Write integration tests

**Files to create/modify:**
```
sdp-plugin/
├── cmd/sdp/
│   └── git_wrapper.go      # NEW: git command wrapper
├── internal/
│   └── git/
│       ├── wrapper.go      # NEW: git wrapper logic
│       └── validator.go    # NEW: git validation
```

### Phase 3: Hooks

**Goal:** Git hooks for automatic protection

- [ ] Create pre-commit hook script
- [ ] Create pre-push hook script
- [ ] Add `sdp hooks install` command
- [ ] Add post-checkout hook for session update
- [ ] Test hooks with various scenarios

**Files to create/modify:**
```
hooks/
├── pre-commit              # NEW: commit validation
├── pre-push                # NEW: push validation
└── post-checkout           # NEW: session update
```

### Phase 4: Agent Integration

**Goal:** Update all skills and agents with safety checks

- [ ] Create `.claude/GIT_SAFETY.md`
- [ ] Update `@build` skill with context checks
- [ ] Update `@bugfix` skill with context checks
- [ ] Update `@deploy` skill with context checks
- [ ] Update `orchestrator.md` agent spec
- [ ] Update `builder.md` agent spec

**Files to create/modify:**
```
.claude/
├── GIT_SAFETY.md           # NEW: safety guidelines
├── skills/
│   ├── build/SKILL.md      # UPDATE: add checks
│   ├── bugfix/SKILL.md     # UPDATE: add checks
│   └── deploy/SKILL.md     # UPDATE: add checks
└── agents/
    ├── orchestrator.md     # UPDATE: add context awareness
    └── builder.md          # UPDATE: add context awareness
```

### Phase 5: Locking (Optional)

**Goal:** Prevent concurrent agent conflicts

- [ ] Implement `.sdp/worktree-lock.json` format
- [ ] Create `sdp lock claim/release/status` commands
- [ ] Add lock expiry mechanism
- [ ] Integrate with session management

---

## Success Metrics

### Quantitative

| Metric | Baseline | Target |
|--------|----------|--------|
| Wrong-branch commits | 3/week | 0/week |
| Branch tracking errors | 5/week | 0/week |
| Agent context confusion | 2/session | 0/session |
| Hook rejection rate | N/A | Track for tuning |

### Qualitative

| Metric | How to Measure |
|--------|----------------|
| Agent confidence | No "what branch am I on?" questions |
| User trust | No manual intervention needed |
| Safety feel | Pre-commit/push hooks work seamlessly |
| Recovery ease | Single command fixes context issues |

### Verification

After implementation:

```bash
# Test 1: Session validation
cd /path/to/sdp-F063
sdp guard context check
# Expected: ✅ All checks pass

# Test 2: Branch mismatch detection
git checkout dev
sdp guard context check
# Expected: ERROR - branch mismatch

# Test 3: Wrong worktree detection
cd /path/to/sdp  # main repo
sdp guard context check --expect F063
# Expected: ERROR - not in F063 worktree

# Test 4: Pre-commit hook
echo "test" >> test.txt && git add test.txt
git commit -m "test"
# Expected: Hook validates session

# Test 5: Pre-push hook
git push
# Expected: Hook validates remote tracking
```

---

## Related Documents

- [Incident Report: Agent Branch Confusion](../reference/incident-2026-02-12-agent-branch-confusion.md)
- [PROTOCOL.md](../PROTOCOL.md) - Core SDP specification
- [PRINCIPLES.md](../reference/PRINCIPLES.md) - Defense in depth principle

---

**Version:** 1.0.0
**Authors:** 7 Expert Agents (Kelsey Hightower, Troy Hunt, Martin Kleppmann, Dan Abramov, Sam Newman)
**Review Status:** Ready for implementation approval
