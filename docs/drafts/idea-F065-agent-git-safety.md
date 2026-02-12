# F065: Agent Git Safety Protocol

> **Status:** Ready for implementation
> **Date:** 2026-02-12
> **Beads:** sdp-sge2
> **Design:** docs/plans/2026-02-12-agent-git-safety-protocol.md

## Problem Statement

Agents working with git worktrees experience branch confusion:
1. **Branch tracking corruption** — feature branches track wrong remotes
2. **Cross-feature commits** — commits appear in wrong branch history
3. **Direct pushes to dev** — bypassing PR workflow
4. **CWD reset** — shell returns to main repo after Bash commands

## Goals

1. **Agent cannot commit to wrong worktree/branch** — enforced at multiple levels
2. **Automatic verification before each git operation** — no manual discipline required
3. **Feature isolation** — one feature = one worktree = one branch
4. **Recovery after CWD reset** — state persists across shell resets
5. **Feature branch enforcement** — features MUST be implemented in feature branches, not dev/main

## Non-Goals

- **Concurrent agent locking** — orchestrators work on features with multiple agents, so per-worktree locking is not appropriate

## Solution Architecture

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
│  - Hash verification for tamper detection                       │
└─────────────────────────────────────────────────────────────────┘
```

## Components

### 1. Session State Management
- Per-worktree `.sdp/session.json` file
- Tracks: worktree_path, feature_id, expected_branch, expected_remote
- Hash verification for tamper detection
- Commands: `sdp session init/sync/repair`

### 2. Git Operation Wrapper
- `sdp git <command>` wrapper
- Validates session before execution
- Post-check for write commands
- Remote tracking verification

### 3. CWD Recovery & Validation
- `sdp guard context check/show/find/go/clean/repair`
- Hybrid recovery from session + git worktree list
- Automatic CWD correction

### 4. Git Hooks
- pre-commit: session validation
- pre-push: branch tracking validation
- post-checkout: session update

### 5. Agent Prompt Constraints
- `.claude/GIT_SAFETY.md` guidelines
- Inline checklists in @build, @bugfix, @deploy skills
- Orchestrator/builder agent updates

### 6. Feature Branch Enforcement
- Protocol update: features MUST use feature branches
- Guard enforcement: block commits to dev/main for feature work
- Workstream scope validation

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Wrong-branch commits | 3/week | 0/week |
| Branch tracking errors | 5/week | 0/week |
| Agent context confusion | 2/session | 0/session |

## Workstreams

| WS | Title | Priority | Dependencies |
|----|-------|----------|--------------|
| 00-065-01 | Session State Management | P0 | None |
| 00-065-02 | Git Operation Wrapper | P1 | 00-065-01 |
| 00-065-03 | CWD Recovery & Validation Commands | P0 | 00-065-01 |
| 00-065-04 | Git Hooks Protection | P1 | 00-065-01 |
| 00-065-05 | Agent Prompt Constraints | P1 | None |
| 00-065-06 | Feature Branch Enforcement | P0 | 00-065-01 |

## Out of Scope

- Concurrent agent locking (orchestrators work with multiple agents)
- UI/UX changes
- Remote repository management
