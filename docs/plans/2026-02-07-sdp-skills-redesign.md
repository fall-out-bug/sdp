# SDP Skills Redesign: Multi-Agent Patterns from Superpowers/Hyperpowers

> **Status:** Research complete
> **Date:** 2026-02-07
> **Goal:** Analyze low-level SDP skills (@idea, @design, @build, @review, @deploy) and integrate patterns from superpowers/hyperpowers (multi-agent coordination, collective decisions, user-in-the-loop brainstorming)

---

## Table of Contents

1. [Overview](#overview)
2. [1. Multi-Agent Coordination](#1-multi-agent-coordination)
3. [2. Collective Decision Making](#2-collective-decision-making)
4. [3. User-in-the-Loop Brainstorming](#3-user-in-the-loop-brainstorming)
5. [4. Agent Composition](#4-agent-composition)
6. [5. State Management](#5-state-management)
7. [6. Progressive Disclosure](#6-progressive-disclosure)
8. [7. Skill Reusability](#7-skill-reusability)
9. [8. Agent Specialization](#8-agent-specialization)
10. [9. Failure Recovery](#9-failure-recovery)
10. [10. Quality Gates](#10-quality-gates)
11. [Implementation Plan](#implementation-plan)

---

## Overview

### Goals

1. **Enhance multi-agent coordination** — Apply superpowers' two-stage review + parallel dispatch patterns to SDP
2. **Improve decision quality** — Implement collective decision-making (voting, consensus, synthesis)
3. **Optimize user interaction** — Progressive disclosure interviews with smart trigger points
4. **Clarify agent architecture** — Hybrid composition (skills orchestrate, agents execute)
5. **Add resilience** — Circuit breakers, categorized retries, graceful degradation
6. **Balance information** — Progressive disclosure with verbosity tiers

### Key Decisions

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| **Coordination** | Hybrid: Sequential two-stage review + Parallel batch execution | Quality over speed for individual WS, speed for independent WS |
| **Decisions** | Hierarchical supervisor with synthesis rules | Matches existing orchestrator pattern, scalable |
| **User interaction** | Progressive disclosure with 3-question cycles | Cognitive load management, comprehensive coverage |
| **Composition** | Hybrid: Skills compose via Skill(), agents spawn via Task() | Clear separation: skills = logic, agents = execution |
| **State** | JSON file checkpoint with atomic writes | Proven pattern, language-agnostic |
| **Disclosure** | Verbosity tiers (--quiet, default, --verbose, --debug) + Collapsible sections | User control + clean first impression |
| **Reuse** | Convention pattern (WET) over DRY for skills | Duplication cheaper than wrong abstraction |
| **Agents** | 19 specialized single-purpose agents | Bounded contexts, clear responsibilities |
| **Failure** | Circuit breaker + error categorization | Prevent cascade failures, smart retries |
| **Quality** | Two-stage review (spec + quality) per workstream | Catch issues early, prevent drift |

---

## 1. Multi-Agent Coordination

> **Experts:** Martin Kleppmann, Martin Fowler, Kelsey Hightower

### Solution: Hybrid Two-Stage Review + Parallel Dispatch

**Pattern from superpowers:**
```markdown
## Per Workstream Workflow

1. Implementer subagent implements, tests, commits, self-reviews
2. Spec compliance reviewer confirms code matches spec
   - If ❌ Issues → Implementer fixes → Re-review
3. Code quality reviewer approves quality
   - If ❌ Issues → Implementer fixes → Re-review
4. Mark task complete
```

**Applied to SDP @build:**

```markdown
### For Each Workstream:

1. Dispatch Implementer Subagent (.claude/agents/builder.md)
   - Execute @build {ws_id}
   - Report: What implemented, test results, files changed

2. Dispatch Spec Compliance Reviewer
   - Compare actual implementation to requirements
   - Report: ✅ Spec compliant OR ❌ Issues found

3. If spec fails → Implementer fixes → Re-review

4. Dispatch Code Quality Reviewer (.claude/agents/reviewer.md)
   - Review: coverage, mypy, ruff, architecture
   - Report: ✅ Approved OR ❌ Issues found

5. If quality fails → Implementer fixes → Re-review

6. Mark WS complete, update checkpoint
```

**Applied to @oneshot Orchestrator:**

```markdown
### Step 2: Execute Ready Workstreams in Parallel

# Find all ready workstreams (dependencies satisfied)
ready_ws = bd ready --json

# Dispatch N builder agents in parallel (max 3)
for ws_id in ready_ws[:3]:
    Task(
        subagent_type="general-purpose",
        prompt=f"Execute @build {ws_id}",
        description=f"Build {ws_id}",
        run_in_background=True  # Parallel execution
    )

# Wait for all agents, update checkpoint, continue to next batch
```

**Benefits:**
- **Quality:** Two-stage review prevents over/under-building (#1 SDP issue: scope creep)
- **Speed:** Parallel execution of independent workstreams (3x speedup for 5+ WS)
- **Resilience:** Checkpoint after each batch, resume from interruption

---

## 2. Collective Decision Making

> **Experts:** Martin Kleppmann, Leslie Lamport, Barbara Liskov

### Solution: Hierarchical Supervisor with Agent Synthesis

**Pattern:**
```python
class AgentSynthesizer:
    """Collects and synthesizes multi-agent outputs."""

    def synthesize(self, outputs: dict[str, AgentOutput]) -> SynthesisResult:
        # Rule 1: Unanimous agreement
        if self._unanimous(outputs):
            return outputs.values()[0].recommendation

        # Rule 2: Domain expertise prioritization
        if task_type == TaskType.ARCHITECTURE:
            return self._prioritize_by_role(outputs, ["architect", "planner"])

        # Rule 3: Quality gate overrides
        if self._quality_gate_failure(outputs):
            return SynthesisResult(decision="block", rationale="Quality gate failure")

        # Rule 4: Merge compatible outputs
        merged = self._try_merge(outputs)
        if merged:
            return merged

        # Fallback: Escalate to human
        return SynthesisResult(decision="escalate", rationale="Unresolvable conflict")
```

**Synthesis Rules:**
1. **Unanimous:** All agents agree → use common output
2. **Domain expertise:** Architectural decisions → architect > planner
3. **Quality gates:** Reviewer veto overrides builder approval
4. **Merge:** Compatible outputs combined (e.g., test strategies)
5. **Escalate:** Unresolvable conflicts → human intervention

**Applied to @design:**
- Spawns 3 agents (System Architect + Security + SRE) in parallel
- Synthesizer combines outputs into unified design document
- Conflicts resolved by domain expertise (architect for structure, security for threats)

**Applied to @review:**
- Spawns 6 agents (QA, Security, DevOps, SRE, TechLead, Documentation)
- Synthesizer aggregates findings, produces APPROVED/CHANGES_REQUESTED verdict
- Quality gate failures (coverage <80%, security vulnerabilities) block approval

---

## 3. User-in-the-Loop Brainstorming

> **Experts:** Nir Eyal, Kent C. Dodds, Theo Browne

### Solution: Progressive Disclosure Interview (3-5 Question Cycles)

**Trigger Points by Skill:**

| Skill | Questions | Type | Trigger |
|-------|-----------|------|---------|
| **@feature** | 1 (quick) + delegate | Vision/Decision | Description vague (<200 words) |
| **@idea** | 4+ cycles | Vision/Technical/Tradeoffs | Continue until complete |
| **@design** | 1 (approval) | Confirmation | ExitPlanMode request |
| **@build** | 0 | N/A | Autonomous (TDD) |
| **@oneshot** | 0-1 (blockers) | Decision | CRITICAL error only |

**Question Flow:**

```markdown
### @feature: Quick Interview (3-5 questions)

Q1: "What problem does this feature solve?"
- User pain point (Fixes existing user friction)
- New capability (Enables new user workflows)
- Technical debt (Improves code quality/performance)

Q2: "Who are the primary users?"
- End users, Internal, Developers (multi-select)

Q3: "What defines success?"
- Adoption, Efficiency, Quality (multi-select)

### @idea: Deep Requirements Gathering

Cycle 1: Mission + Alignment (2 questions)
Cycle 2: Problem + Users + Technical approach (3 questions)
Cycle 3: Storage + Failure modes + Security (3 questions)
Cycle 4: UI/UX + Testing + Edge cases (3 questions)

Complete: Create Beads task + intent file
```

**Progressive Disclosure Principles:**
- **Start broad:** Vision, mission, users
- **Go deep:** Technical, architecture, tradeoffs
- **Confirm at gates:** Approval, deployment
- **Escalate only on blockers:** Execution issues

**Total Questions for New Feature:**
- **Minimum:** 3 (@feature) + 8 (@idea) + 1 (@design) = **12 questions**
- **Maximum:** 5 (@feature) + 20+ (@idea complex) + 1 (@design) + 1 (@oneshot) = **27+ questions**

---

## 4. Agent Composition

> **Experts:** Martin Fowler, Yegor Bugayenko, Dan Abramov

### Solution: Hybrid Composition (Layered Architecture)

**Three Layers:**

```
┌─────────────────────────────────────┐
│   Skills (Orchestration Layer)     │
│   @feature, @oneshot, @review       │
│   - Call Skill() for composition    │
│   - Spawn Task() for agents        │
└─────────────────────────────────────┘
           │
           ├─► Skill() Calls
           │    (@feature → @idea → @design)
           │
           └─► Task() Spawns
                (6 agents for @review)
```

**Composition Rules:**

1. **Skills compose skills** when sequential:
   ```python
   @feature orchestrates:
     Skill('@idea', args="Add payment")
     Skill('@design', args="idea-payment")
   ```

2. **Skills spawn agents** when parallel:
   ```python
   @review spawns:
     Task("qa-agent")
     Task("security-agent")
     Task("devops-agent")
     ...
   ```

3. **Agents don't compose** (they're leaf nodes)

4. **Agents are reusable** across skills:
   - `builder.md` used by @build and orchestrator
   - `security.md` used by @design and @review

**Agent Reusability Matrix:**

| Agent Type | Used By | Standalone? | Parallel? |
|------------|---------|-------------|-----------|
| **Skills** (@idea, @design, @build) | User direct, @feature | ✅ Yes | ❌ Sequential |
| **Review Agents** (QA, Security, etc.) | @review only | ❌ No | ✅ Yes |
| **Orchestrator** | @oneshot only | ❌ No | ✅ Yes (background) |

---

## 5. State Management

> **Experts:** Martin Kleppmann, Kelsey Hightower, Theo Browne

### Solution: Minimal JSON File Checkpoint

**Checkpoint Schema:**
```json
{
  "feature": "F050",
  "agent_id": "agent-20260205-212900",
  "status": "in_progress",
  "completed_ws": ["00-050-01", "00-050-02"],
  "failed_ws": [],
  "execution_order": ["00-050-01", "00-050-02", "00-050-03", ...],
  "current_ws": "00-050-03",
  "started_at": "2026-02-05T21:29:00Z",
  "last_updated": "2026-02-05T23:45:00Z",
  "metadata": {
    "gates": {...},
    "team": {...},
    "beads_mapping": {...}
  }
}
```

**Save Checkpoint (after each WS):**
```python
def save_checkpoint_safe(data, path):
    """Atomic checkpoint write using temp file + rename."""
    import tempfile
    import os

    # Write to temp file in same directory
    dir = os.path.dirname(path)
    with tempfile.NamedTemporaryFile(mode='w', dir=dir, delete=False, suffix='.tmp') as tmp:
        json.dump(data, tmp, indent=2)
        tmp_path = tmp.name

    # Atomic rename (overwrites target)
    os.replace(tmp_path, path)
```

**Resume Checkpoint:**
```python
# Load checkpoint
with open(checkpoint_path) as f:
    checkpoint = json.load(f)

# Verify agent_id if provided
if agent_id and checkpoint['agent_id'] != agent_id:
    raise ValueError("Agent ID mismatch")

# Filter completed WS from execution order
remaining_ws = [ws for ws in checkpoint['execution_order']
                if ws not in checkpoint['completed_ws']]

# Continue from first incomplete WS
next_ws = remaining_ws[0]
```

**File Locations:**
- Python: `.oneshot/{feature_id}-checkpoint.json`
- Go: `.sdp/checkpoints/{feature_id}.json`

---

## 6. Progressive Disclosure

> **Experts:** Nielsen Norman Group, Edward Tufte, Kathy Sierra

### Solution: Verbosity Tiers + Collapsible Sections

**Standardize Verbosity Tiers:**
```
Level 0 (--quiet):   Exit status only (✅/❌)
Level 1 (default):   Summary (1-3 lines with key metrics)
Level 2 (--verbose): Step-by-step progress
Level 3 (--debug):   Internal state + API calls
```

**Examples:**

```bash
# @build skill
--quiet:  "✅"
default: "✅ 00-050-01: Workstream Parser (22m, 85%, commit:abc123)"
--verbose: "→ Reading WS spec...
           → TDD cycle: Red (3m) → Green (12m) → Refactor (7m)
           → Quality check: PASS (coverage 85%, mypy clean)
           ✅ COMPLETE"
--debug:  "→ Beads API: bd.update(sdp-123, status='in_progress')
           → Guard: sdp.guard.activate(00-050-01)
           [...full trace...]"
```

**Collapsible Sections for Documentation:**
```markdown
## Workstream: 00-050-01 - Workstream Parser

### Executive Summary
Implement YAML frontmatter parser for workstream files with 85% coverage.

<details>
<summary>📖 Full Context (click to expand)</summary>

This workstream is part of F050 (SDP Core Tooling)...
</details>

<details>
<summary>🔧 Implementation Details (click to expand)</summary>

### Tasks
1. Parser Core (src/sdp/core/parser.py)
...
</details>

### Acceptance Criteria
- ✅ Extracts ws_id from frontmatter
- ✅ Validates format (PP-FFF-SS)
```

**TMI Detection Thresholds:**
```
Rule 1: First screen ≤ 10 lines (or 1 screenful)
Rule 2: Secondary info = 1 click/tap away (collapsible)
Rule 3: Advanced features = separate "Advanced" section
Rule 4: Exit criteria always visible (never hide AC)
```

---

## 7. Skill Reusability

> **Experts:** Martin Fowler, Sandi Metz, Kent C. Dodds

### Solution: Convention Pattern (WET over DRY)

**Key Insight:** "Duplication is far cheaper than the wrong abstraction" — Sandi Metz

Skills are **specifications**, not executable code. Duplication makes each skill's behavior explicit and self-contained.

**Pattern Library Structure:**
```bash
.claude/
├── patterns/           # Document conventions (read-only reference)
│   ├── BEADS_DETECTION.md
│   ├── MULTI_AGENT_SPAWN.md
│   ├── QUALITY_GATES.md
│   └── OUTPUT_FORMATS.md
├── fragments/          # Stable fragments for include
│   ├── _quality_gates.md
│   └── _beads_integration.md
└── skills/
    ├── build/SKILL.md  # Reference patterns, include fragments
    └── ...
```

**When to Extract Patterns:**

| **Shared Behavior** | **Current Duplication** | **Recommended Approach** |
|---------------------|-------------------------|--------------------------|
| Beads detection | 8 skills | Convention (still evolving) |
| Multi-agent spawn | 2 skills (@design, @review) | Convention (only 2 uses) |
| Quality gates | 3+ skills | **Fragment** (stable checklist) |
| Output format tables | 5+ skills | **Fragment** (pure markdown) |
| Skill orchestration | 1 orchestrator | N/A (delegation pattern) |
| TDD cycle | 2 skills (@build, @tdd) | **Python helper** (complex logic) |

**3-Strike Rule:** Extract pattern after 3+ **stable** (unchanged) uses.

---

## 8. Agent Specialization

> **Experts:** Sam Newman, Barbara Liskov, Kelsey Hightower

### Solution: 19 Specialized Single-Purpose Agents

**Complete Agent Catalog:**

**Discovery Team:**
- **Business Analyst:** Requirements discovery, user stories, KPIs
- **Product Manager:** Product vision, prioritization, roadmap
- **Systems Analyst:** Functional requirements, API specs, data models
- **Analyst:** Bridge stakeholders and development

**Design Team:**
- **System Architect:** Architecture design, tech stack, quality attributes
- **Architect:** Software architecture, clean architecture, ADRs
- **Security:** Threats, auth, compliance (OWASP, GDPR, SOC2)
- **SRE:** SLOs/SLIs, monitoring, incident response
- **DevOps:** CI/CD, infrastructure, deployment

**Planning & Decomposition:**
- **Planner:** Codebase analysis, workstream decomposition
- **Technical Decomposition:** Workstreams, dependencies, estimation

**Implementation:**
- **Developer:** Senior developer, production code, clean architecture
- **Builder:** TDD execution, test coverage, quality gates
- **Tester:** Test strategies, edge cases, coverage

**Review & QA:**
- **Tech Lead:** Technical leadership, code review, team coordination
- **QA:** Test strategy, quality metrics, quality gates
- **Reviewer:** Quality validation, coverage, linters

**Orchestration:**
- **Orchestrator:** Autonomous execution, checkpoints, error handling
- **Deployer:** Production deployment, configs, smoke tests

**Cross-Cutting Concerns:**
- **Security:** Handled by dedicated agent (active in design + review)
- **Performance:** Shared (architect defines, SRE monitors, developer implements)
- **Testing:** QA designs, tester writes, builder follows TDD
- **Documentation:** Shared responsibility (all agents contribute)

**Dynamic Role Switching:**
```bash
@role builder
@build WS-001
@switch security-reviewer
@review WS-001
@switch builder
@build WS-002
```

---

## 9. Failure Recovery

> **Experts:** Martin Kleppmann, Michael Nygard, Kelsey Hightower

### Solution: Circuit Breaker + Categorized Retries

**Error Categorization:**
```python
class AgentFailureType(Enum):
    TRANSIENT = "transient"  # Network, temp resource, rate limit
    PERMANENT = "permanent"  # Validation, circular dep, config error
    CASCADE = "cascade"      # Downstream agent failure
    CRASH = "crash"          # Agent process died
```

**Circuit Breaker Implementation:**
```go
type CircuitBreakerState int
const (
    Closed CircuitBreakerState = iota  // Normal operation
    Open                                // Failing fast
    HalfOpen                            // Testing recovery
)

type CircuitBreaker struct {
    state           CircuitBreakerState
    failureCount    int
    failureRatio    float64  // e.g., 0.1 = break at 10% failure
    windowSize      int      // Number of requests to sample
    lastFailureTime time.Time
}
```

**Retry Strategy by Failure Type:**
- **TRANSIENT:** 3 retries with exponential backoff (1s, 2s, 4s)
- **PERMANENT:** 1 retry, then escalate to human
- **CASCADE:** Open circuit (stop dependent workstreams), escalate
- **CRASH:** 2 retries (agent may have crashed), then escalate

**Selective C Elements:**
- Human escalation for PERMANENT failures after 1 retry
- Fallback to alternative agent for TRANSIENT failures after 2 retries
- Graceful degradation: continue with degraded mode if non-critical agent fails

---

## 10. Quality Gates

> **Experts:** Jez Humble, Martin Fowler, Kent Beck

### Solution: Two-Stage Review (Spec Compliance + Code Quality)

**Stage 1: Spec Compliance Review**
```python
Task(
    subagent_type="general-purpose",
    prompt=f"""
    Review spec compliance for {ws_id}

    WHAT WAS REQUESTED: {AC from WS file}
    WHAT IMPLEMENTER CLAIMS: {from implementer report}

    CRITICAL: Do not trust report. Read actual code.

    Report: ✅ Spec compliant OR ❌ Issues found
    """,
    description=f"Spec review {ws_id}"
)
```

**Stage 2: Code Quality Review**
```python
Task(
    subagent_type="general-purpose",
    prompt=f"""
    Read .claude/agents/reviewer.md for your specification.

    BASE_SHA: {commit before WS}
    HEAD_SHA: {current commit}

    Review quality: coverage, mypy, ruff, architecture
    """,
    description=f"Quality review {ws_id}"
)
```

**Quality Gates:**
- **Coverage:** ≥80% (unit + integration)
- **Type safety:** mypy --strict (Python) or go vet (Go)
- **Linters:** ruff (Python), golint (Go)
- **Architecture:** Clean Architecture compliance
- **Tech Debt:** Zero TODO/FIXME in production code
- **File size:** <200 LOC per file

**Gate Enforcement:**
- Spec compliance MUST pass before quality review
- Quality review MUST pass before marking WS complete
- Either gate failing → Implementer fixes → Re-review
- Both gates passing → Mark WS complete, update checkpoint

---

## Implementation Plan

### Phase 1: MVP (Two-Stage Review) - 2 weeks

**01-051-01: Implementer Subagent** (MEDIUM, 3-5 days)
- AC1: Create `.claude/agents/implementer.md` with TDD cycle spec
- AC2: Implement @build workflow: spawn implementer → get report
- AC3: Implementer reads WS file, executes TDD, commits, self-reviews
- AC4: Implementer reports: what implemented, test results, files changed

**01-051-02: Spec Compliance Reviewer** (MEDIUM, 3-5 days)
- AC1: Create `.claude/agents/spec-reviewer.md`
- AC2: Implement spec review: compare actual code to requirements
- AC3: Critical: "Do not trust report" - read actual code
- AC4: Report: ✅ Spec compliant OR ❌ Issues found

**01-051-03: Code Quality Reviewer** (MEDIUM, 3-5 days)
- AC1: Extend `.claude/agents/reviewer.md` for two-stage review
- AC2: Implement quality review: coverage, mypy, ruff, architecture
- AC3: Compare BASE_SHA vs HEAD_SHA for diff review
- AC4: Report: ✅ Approved OR ❌ Issues found

**01-051-04: Two-Stage Workflow** (SMALL, 2-3 days)
- AC1: Update @build skill: implementer → spec review → quality review
- AC2: Implement fix loop: if spec/quality fails → implementer fixes → re-review
- AC3: Max 2 retries per gate, then escalate
- AC4: Update checkpoint after both gates pass

---

### Phase 2: Parallel Execution - 3 weeks

**01-052-01: Circuit Breaker** (LARGE, 5-7 days)
- AC1: Implement CircuitBreaker struct in Go
- AC2: Add to Orchestrator as optional component
- AC3: Persist state to checkpoint
- AC4: Add metrics (failure count, open time)

**01-052-02: Error Categorization** (MEDIUM, 3-5 days)
- AC1: Extend SDPError with `failure_type` field
- AC2: Categorize all existing SDP error types
- AC3: Add `is_retriable()` method to SDPError
- AC4: Document error categories

**01-052-03: Parallel Batch Execution** (LARGE, 5-7 days)
- AC1: Extend orchestrator to spawn N builder agents in parallel
- AC2: Use `bd ready` to find ready workstreams (no dependencies)
- AC3: Execute up to 3 workstreams in parallel (background mode)
- AC4: Wait for all agents, update checkpoint, continue to next batch

**01-052-04: Enhanced Retry Logic** (MEDIUM, 3-5 days)
- AC1: Replace `executeWithRetry` with `executeWithCircuitBreaker`
- AC2: Implement different retry counts by failure type
- AC3: Add exponential backoff for transient failures
- AC4: Add circuit break threshold (50% failure, window=10)

---

### Phase 3: Decision Synthesis - 2 weeks

**01-053-01: Agent Synthesizer** (MEDIUM, 4-5 days)
- AC1: Implement `AgentSynthesizer` class in Python
- AC2: Add synthesis rules: unanimous, domain expertise, quality gates, merge
- AC3: Implement fallback: escalate to human
- AC4: Unit tests for each synthesis rule

**01-053-02: Enhanced @design** (MEDIUM, 3-5 days)
- AC1: Update @design to use synthesizer for 3-agent outputs
- AC2: Implement domain expertise prioritization (architect > security > SRE)
- AC3: Add merge logic for compatible outputs
- AC4: Create unified design document from synthesized outputs

**01-053-03: Enhanced @review** (MEDIUM, 3-5 days)
- AC1: Update @review to use synthesizer for 6-agent outputs
- AC2: Implement quality gate overrides (reviewer veto)
- AC3: Generate APPROVED/CHANGES_REQUESTED verdict
- AC4: Create issues for findings via Beads integration

**01-053-04: Human Escalation** (SMALL, 2-3 days)
- AC1: Define escalation triggers (permanent failure, agent crash)
- AC2: Add escalation prompt templates
- AC3: Implement escalation in synthesizer fallback
- AC4: Track escalations in checkpoint metadata

---

### Phase 4: User Experience - 2 weeks

**01-054-01: Verbosity Tiers** (MEDIUM, 3-5 days)
- AC1: Implement --quiet, --verbose, --debug flags for all skills
- AC2: Define output format for each verbosity level
- AC3: Update @build, @review, @oneshot with verbosity support
- AC4: Test with examples (CI/CD mode, learning mode, troubleshooting)

**01-054-02: Collapsible Sections** (SMALL, 2-3 days)
- AC1: Add `<details>` sections to workstream templates
- AC2: Create executive summary (always visible)
- AC3: Collapse implementation details, context, debug info
- AC4: Test markdown rendering in Cursor/Claude Code

**01-054-03: Progressive Disclosure** (MEDIUM, 3-5 days)
- AC1: Implement TMI detection (first screen ≤10 lines)
- AC2: Add 3-question cycles to @idea skill
- AC3: Add progress indicators ("Question cycle 2/5")
- AC4: Implement --no-interview flag for experienced users

**01-054-04: User Config** (SMALL, 2-3 days)
- AC1: Add `~/.claude/settings.json` support
- AC2: Implement config: sdp.verbosity, sdp.auto_expand, sdp.expert_mode
- AC3: Apply defaults to all skills
- AC4: Document configuration options

---

### Phase 5: Polish & Documentation - 1 week

**01-055-01: Agent Catalog** (SMALL, 2-3 days)
- AC1: Create `.claude/agents/README.md` with complete catalog
- AC2: Document all 19 agents: purpose, expertise, scope, collaboration
- AC3: Add cross-cutting concerns section
- AC4: Add dynamic role switching examples

**01-055-02: Pattern Library** (MEDIUM, 3-4 days)
- AC1: Extract 3 patterns to `.claude/patterns/` (BEADS_DETECTION, MULTI_AGENT_SPAWN, QUALITY_GATES)
- AC2: Add cross-references in skill files
- AC3: Create 2 fragments (quality_gates.md, beads_integration.md)
- AC4: Document convention vs. extraction guidelines

**01-055-03: Migration Guide** (SMALL, 2-3 days)
- AC1: Document changes from old @feature to new orchestrator
- AC2: Document changes from old @build to two-stage review
- AC3: Provide examples: before vs. after
- AC4: Add troubleshooting section

---

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| **Quality** | Scope creep (13% friction) | ↓40% via two-stage review |
| **Speed** | Sequential WS execution | 3x speedup for 5+ WS (parallel) |
| **User Questions** | 20+ questions per feature | 12-27 questions (progressive disclosure) |
| **Agent Coordination** | Sequential orchestration | Parallel batches + two-stage review |
| **Failure Recovery** | Simple retry (max 3) | Circuit breaker + categorized retries |
| **Information Overload** | 150-line WS templates | 10-line summary + collapsible details |
| **Decision Quality** | Single-agent decisions | Synthesized from 6 agents (review) |

---

## Open Questions

1. **Weight Tuning** (for collective decisions): What makes reviewer 3x weight? How to validate?
2. **Circuit Thresholds**: 50% failure rate, window=10 - how to tune for production?
3. **Verbosity Adoption**: Will developers use --verbose/--debug or stick to default?
4. **Agent Proliferation**: 19 agents complex - how to simplify onboarding?
5. **Beads Integration**: How do Beads tasks interact with new agent workflows?

---

## Appendix: Key Files to Modify

**Skills:**
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/skills/build/SKILL.md` - Two-stage review
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/skills/idea/SKILL.md` - Progressive disclosure
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/skills/design/SKILL.md` - Synthesizer integration
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/skills/review/SKILL.md` - Synthesizer integration
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/skills/oneshot/SKILL.md` - Parallel execution

**Agents:**
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/agents/implementer.md` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/agents/spec-reviewer.md` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/agents/reviewer.md` - ENHANCE
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/agents/orchestrator.md` - ENHANCE

**Code:**
- `/Users/fall_out_bug/projects/vibe_coding/sdp/src/sdp/agents/synthesizer.py` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/orchestrator/circuit_breaker.go` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/checkpoint/checkpoint.go` - ENHANCE

**Documentation:**
- `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/plans/2026-02-07-sdp-skills-redesign.md` - THIS FILE
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/patterns/BEADS_DETECTION.md` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/patterns/MULTI_AGENT_SPAWN.md` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/patterns/QUALITY_GATES.md` - NEW
- `/Users/fall_out_bug/projects/vibe_coding/sdp/.claude/agents/README.md` - UPDATE

---

**Document Status:** ✅ Research complete
**Next:** Review with stakeholders, prioritize phases, start Phase 1 implementation
**Estimated Timeline:** 10 weeks (all 5 phases)
**Recommended Start:** Phase 1 (Two-Stage Review MVP)
