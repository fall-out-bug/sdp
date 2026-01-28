# F014: Workflow Efficiency

> **Feature ID:** F014
> **Status:** Draft
> **Created:** 2026-01-28
> **Intent:** docs/intent/f014-workflow-efficiency.json

## Problem

SDP workflow has two major bottlenecks preventing optimal "product team in a box" experience:

1. **PR Approval Bottleneck**: @oneshot requires PR approval before deployment, blocking autonomous execution and causing 80% of cycle time delays
2. **@idea Interview Overhead**: 6-12 questions taking 15-20 minutes overwhelm users and slow down the initial phase

**Current State:**
- 4-WS feature takes ~3h 45m from @idea to @deploy
- PR approval creates blocking wait in autonomous workflow
- New users get overwhelmed by extensive @idea questioning

**Target State:**
- 4-WS feature takes <45 min from @idea to @deploy (**5x throughput increase**)
- @idea interview takes 5-8 min with optional deep dive
- PR-less execution for trusted features

## Users

### Primary: Fast-iterating developers
- **Needs:** Rapid feedback cycles without waiting for PR approval
- **Pain points:**
  - Blocked by PR approval during @oneshot
  - Overwhelmed by too many @idea questions
  - Slow iteration cycles kill momentum

### Secondary: AI agents using @oneshot
- **Needs:** Autonomous execution without human-in-the-loop for trusted features
- **Pain points:**
  - PR gate breaks autonomous workflow
  - Cannot deploy to sandbox for testing without PR

## Success Criteria

1. **PR-less Adoption**: Users successfully use --auto-approve for trusted deployments
2. **Cycle Time**: Reduce @idea → @deploy from 3h 45m to <45 min
3. **@idea Speed**: Reduce interview from 15-20 min to 5-8 min
4. **Throughput**: 5x improvement in features completed per day

## Goals

1. **Remove PR approval bottleneck** with multiple execution modes
2. **Streamline @idea interview** to 3-5 critical questions with optional deep dive
3. **Maintain quality gates** while enabling faster execution
4. **Provide safe modes** (--sandbox, --dry-run) for risk mitigation

## Non-Goals

- Removing quality gates (coverage, LOC, type hints remain enforced)
- Removing @design step (still required for planning)
- Removing @review step (still required for quality check)
- Removing human oversight entirely (just making it optional)

## Technical Approach

### Part 1: PR Gate Removal

**Implementation:** Add execution modes to @oneshot skill

| Mode | Flag | Behavior | Use Case |
|------|------|----------|----------|
| **Trusted** | `--auto-approve` | Skip PR, deploy directly | Trusted features, rapid iteration |
| **Sandbox** | `--sandbox` | Skip PR, deploy to sandbox only | Testing without production risk |
| **Standard** | (default) | PR required (existing) | Production deployments |

**CLI Examples:**
```bash
# Trusted feature, no PR needed
@oneshot F060 --auto-approve

# Test in sandbox without PR
@oneshot F060 --sandbox

# Standard PR-required workflow
@oneshot F060
```

### Part 2: @idea Streamlining

**Implementation:** Two-round interview with progressive disclosure

**Round 1: Critical Questions (Required, 5-8 min)**
- Mission: What problem do we solve?
- Users: Who are we building for?
- Technical approach: Architecture, storage, failure mode
- Risk level: Critical vs non-critical

**Round 2: Deep Dive (Optional, 5-10 min)**
- Triggered by:
  - Ambiguity detected in Round 1 answers
  - User explicitly requests deep dive
  - High-risk feature detected

**Example Flow:**
```bash
@idea "Add user auth"
# Round 1: 4 questions (5-8 min)
# → Output: bd-0001 + docs/intent/bd-0001.json

# If ambiguous:
@idea "Add user auth" --deep-dive
# Round 2: 3-5 additional questions (5-10 min)
# → Updates intent with detailed answers
```

### Part 3: Risk Mitigation

**Quality Gates (Enforced Regardless of Mode):**
- Test coverage ≥80%
- File size <200 LOC
- Type hints (mypy --strict)
- No `except: pass`

**Additional Safeguards:**
1. **Destructive Operations Confirmation**: Manual confirmation for database changes, deletions
2. **Audit Logging**: Log all --auto-approve executions with timestamp and user
3. **Dry-Run Mode**: `--dry-run` flag to preview changes without execution

**Example:**
```bash
# Preview changes before executing
@oneshot F060 --auto-approve --dry-run

# Output shows:
# - Workstreams to execute: 4
# - Files to create: 12
# - Files to modify: 8
# - Destructive operations: 1 (database migration)
#
# Confirm? [y/N]
```

## Workstream Breakdown

### F014.01: @oneshot Execution Modes
**Status:** Backlog
**Size:** MEDIUM
**Dependencies:** None

Add `--auto-approve` and `--sandbox` flags to @oneshot skill:
- Modify MultiAgentExecutor to accept mode parameter
- Add pre-execution validation (quality gates still enforced)
- Add audit logging for --auto-approve executions
- Update @oneshot SKILL.md with new modes

### F014.02: @idea Two-Round Interview
**Status:** Backlog
**Size:** MEDIUM
**Dependencies:** None

Refactor @idea skill for two-round interview:
- Define 3-5 critical questions for Round 1
- Implement ambiguity detection logic
- Add `--deep-dive` flag for optional Round 2
- Update @idea SKILL.md with new workflow

### F014.03: Destructive Operations Detection
**Status:** Backlog
**Size:** SMALL
**Dependencies:** F014.01

Add safeguards for destructive operations:
- Detect database migrations, file deletions
- Require manual confirmation before execution
- Add `--dry-run` flag to preview changes

### F014.04: Audit Logging
**Status:** Backlog
**Size:** SMALL
**Dependencies:** F014.01

Implement audit trail for --auto-approve:
- Log to `.sdp/audit.log`
- Include: timestamp, user, feature, mode, result
- Add `sdp audit` command to view logs

### F014.05: Documentation & Examples
**Status:** Backlog
**Size:** SMALL
**Dependencies:** F014.01, F014.02

Update documentation:
- Add workflow examples for each mode
- Update TUTORIAL.md with quick-start guide
- Add decision tree: "@build vs @oneshot vs --auto-approve"
- Update CHANGELOG.md

## Expected Impact

| Metric | Baseline | Target | Improvement |
|--------|----------|--------|-------------|
| @idea → @deploy time | 3h 45m | <45 min | **5x faster** |
| @idea interview duration | 15-20 min | 5-8 min | **3x faster** |
| Features per developer per day | 1-2 | 5-10 | **5x throughput** |
| PR-less adoption rate | 0% | >60% | **New capability** |

## Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Quality degradation | High | Low | Quality gates enforced regardless |
| Destructive operations in auto-approve | High | Low | Detection + confirmation required |
| Loss of oversight | Medium | Medium | Audit logging + --dry-run mode |
| User confusion about modes | Low | Medium | Clear documentation + examples |

## Open Questions

1. Should --auto-approve require additional flags for production deployments?
2. What qualifies as "ambiguous" in @idea Round 1? Need heuristics.
3. Should audit logs be committed to git or kept locally?
4. How to handle rollback from --auto-approve deployments?

## Next Steps

1. **Review intent schema** with stakeholders
2. **Create detailed design** for @oneshot modes
3. **Implement F014.01** (@oneshot execution modes)
4. **Test with pilot users** before full rollout

---

**References:**
- Analysis: docs/plans/2025-01-26-sdp-analysis-design.md (lines 192-222)
- Product Vision: PRODUCT_VISION.md
- Schema: docs/schema/intent.schema.json
