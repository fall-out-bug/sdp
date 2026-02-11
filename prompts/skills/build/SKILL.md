---
name: build
description: Execute workstream with TDD and guard enforcement
tools: Read, Write, StrReplace, Shell, Skill
---

# @build - Execute Workstream

Execute a single workstream following TDD discipline with automatic guard.

## Invocation (BEADS-001)

Accepts **both** formats:

- `@build 00-001-01` — WS-ID (PP-FFF-SS), resolve beads_id from `.beads-sdp-mapping.jsonl`
- `@build sdp-xxx` — Beads task ID directly

## Verbosity Tiers

```bash
@build 00-050-01 --quiet     # Exit status only: ✅
@build 00-050-01             # Summary: ✅ 00-050-01: Workstream Parser (22m, 85%, commit:abc123)
@build 00-050-01 --verbose   # Step-by-step progress
@build 00-050-01 --debug     # Internal state + API calls
```

**Examples:**

```bash
# Quiet mode
@build 00-050-01 --quiet
# Output: ✅

# Default mode
@build 00-050-01
# Output: ✅ 00-050-01: Workstream Parser (22m, 85%, commit:abc123)

# Verbose mode
@build 00-050-01 --verbose
# Output:
# → Activating guard...
# → Reading WS spec...
# → Stage 1: Implementer (TDD cycle)
#   → Red: Writing failing test (3m)
#   → Green: Implementing minimum code (12m)
#   → Refactor: Improving code (7m)
# → Stage 2: Spec Reviewer
#   → Verifying implementation matches spec (5m)
# → Stage 3: Quality Reviewer
#   → Running quality gates (3m)
# ✅ COMPLETE

# Debug mode
@build 00-050-01 --debug
# Output:
# [DEBUG] WS ID: 00-050-01
# [DEBUG] Beads ID: sdp-abc123
# [DEBUG] Activating guard...
# [DEBUG] Guard activated: /tmp/guard-00-050-01.json
# [DEBUG] Reading WS spec: docs/workstreams/backlog/00-050-01.md
# [DEBUG] Scope files: [src/sdp/parser.py, tests/sdp/test_parser.py]
# → Activating guard...
# → Reading WS spec...
# [DEBUG] Starting Stage 1: Implementer...
# → Stage 1: Implementer (TDD cycle)
#   [DEBUG] Test file: tests/sdp/test_parser.py
#   → Red: Writing failing test (3m)
#   [DEBUG] Implementing in: src/sdp/parser.py
#   → Green: Implementing minimum code (12m)
#   → Refactor: Improving code (7m)
# [DEBUG] Coverage: 85.3%
# → Stage 2: Spec Reviewer
#   [DEBUG] Verifying: all AC implemented
#   → Verifying implementation matches spec (5m)
# → Stage 3: Quality Reviewer
#   [DEBUG] Running: pytest --cov
#   → Running quality gates (3m)
# ✅ COMPLETE
```

## Beads Integration (optional)

**When Beads is enabled** (bd installed, `.beads/` exists):

1. **Resolve ID:** ws_id → beads_id via mapping (if ws_id given)
2. **Before work:** `bd update {beads_id} --status in_progress`
3. **Execute:** TDD cycle
4. **On success:** `bd close {beads_id} --reason "WS completed" --suggest-next`
5. **On failure:** `bd update {beads_id} --status blocked`
6. **Before commit:** `bd sync`

**When Beads NOT enabled:** Skip all Beads steps. Use ws_id only.

**Detection:** Check if `bd --version` works and `.beads/` exists.

## Quick Reference

| Step | Action | Gate | Beads? | Retries |
|------|--------|------|--------|---------|
| 0 | Detect Beads | Check `bd --version` + `.beads/` | Detection | - |
| 0a | Resolve beads_id | ws_id → mapping (if Beads) | Optional | - |
| 0b | Beads IN_PROGRESS | `bd update --status in_progress` (if Beads) | Optional | - |
| 1 | Activate guard | `sdp guard activate {ws_id}` | Always | - |
| 2 | Read WS spec | AC present and clear | Always | - |
| 3 | **Stage 1: Implementer** | TDD cycle | Always | 2 |
| 4 | **Stage 2: Spec Reviewer** | Spec compliance | Always | 2 |
| 5 | **Stage 3: Quality Reviewer** | Final quality check | Always | 2 |
| 6 | All stages pass? | If yes, proceed | Always | - |
| 7 | Beads CLOSED/blocked | `bd close` or `bd update --status blocked` (if Beads) | Optional | - |
| 8 | Beads sync + commit | `bd sync` then commit (if Beads) | Optional | - |
| 9 | Commit only | `git commit` (if no Beads) | Fallback | - |

## Workflow

### Step 0: Resolve Task ID

```bash
# Input: ws_id (00-001-01) OR beads_id (sdp-xxx)
# If ws_id: beads_id = grep mapping for sdp_id
# If beads_id: ws_id = grep mapping for beads_id (reverse lookup)
# Guard needs ws_id; Beads needs beads_id
beads_id=$(grep -m1 "\"sdp_id\": \"{WS-ID}\"" .beads-sdp-mapping.jsonl 2>/dev/null | grep -o '"beads_id": "[^"]*"' | cut -d'"' -f4)
ws_id=$(grep -m1 "\"beads_id\": \"{beads_id}\"" .beads-sdp-mapping.jsonl 2>/dev/null | grep -o '"sdp_id": "[^"]*"' | cut -d'"' -f4)
```

### Step 1: Beads IN_PROGRESS (when Beads enabled)

```bash
[ -n "$beads_id" ] && bd update "$beads_id" --status in_progress
```

### Step 2: Activate Guard

```bash
sdp guard activate {WS-ID}
```

**Gate:** Must succeed. If fails, WS not ready.

### Step 3: Read Workstream

```bash
Read("docs/workstreams/backlog/{WS-ID}-*.md")
```

Extract:
- Goal and Acceptance Criteria
- Input/Output files
- Steps to execute

### Step 4: Two-Stage Review (Three Stages)

**IMPORTANT:** Execute all three stages with retry logic.

#### Stage 1: Implementer Agent

**Purpose:** Execute TDD cycle (Red → Green → Refactor)

```python
Task(
    subagent_type="general-purpose",
    prompt="""You are the IMPLEMENTER agent.

Read .claude/agents/implementer.md for your specification.

WORKSTREAM: {WS-ID}
SPEC: docs/workstreams/backlog/{WS-ID}.md

Execute TDD cycle for each AC:
1. RED: Write failing test
2. GREEN: Implement minimum code
3. REFACTOR: Improve code

Generate self-report with metrics.
Run quality gates.

Return verdict: PASS/FAIL
""",
    description="Implementer agent - Stage 1"
)
```

**Retry Logic:**
- Max 2 retries
- If FAIL: Fix issues, retry
- If FAIL after 2 retries: Stop, report failure

**Success Criteria:**
- ✅ All AC implemented
- ✅ All tests passing
- ✅ Quality gates passed
- ✅ Self-report generated

#### Stage 2: Spec Compliance Reviewer Agent

**Purpose:** Verify implementation matches specification

```python
Task(
    subagent_type="general-purpose",
    prompt="""You are the SPEC COMPLIANCE REVIEWER agent.

Read .claude/agents/spec-reviewer.md for your specification.

WORKSTREAM: {WS-ID}
SPEC: docs/workstreams/backlog/{WS-ID}.md
IMPLEMENTER REPORT: {from stage 1}

CRITICAL: DO NOT TRUST implementer report.
Verify everything yourself:
1. Read actual code
2. Run tests yourself
3. Check coverage yourself
4. Verify each AC manually

Generate evidence-based verdict.
Return verdict: PASS/FAIL
""",
    description="Spec compliance reviewer - Stage 2"
)
```

**Retry Logic:**
- Max 2 retries
- If FAIL: Implementer fixes issues, retry
- If FAIL after 2 retries: Stop, report failure

**Success Criteria:**
- ✅ All AC verified (evidence provided)
- ✅ Implementation matches spec
- ✅ Tests are real (not mocked)
- ✅ Quality gates passed (verified)

#### Stage 3: Quality Reviewer Agent

**Purpose:** Final quality check (comprehensive review)

```python
Task(
    subagent_type="general-purpose",
    prompt="""You are the QUALITY REVIEWER agent.

WORKSTREAM: {WS-ID}
SPEC: docs/workstreams/backlog/{WS-ID}.md

Run comprehensive quality check:
1. Test coverage (≥80%)
2. Code quality (LOC, complexity)
3. Security check
4. Performance check
5. Documentation check

Generate quality report.
Return verdict: PASS/FAIL
""",
    description="Quality reviewer - Stage 3"
)
```

**Retry Logic:**
- Max 2 retries
- If FAIL: Fix quality issues, retry
- If FAIL after 2 retries: Stop, report failure

**Success Criteria:**
- ✅ All quality gates passed
- ✅ No security issues
- ✅ No performance issues
- ✅ Documentation complete

#### Retry Logic Summary

```python
def execute_stage(stage_name, agent, max_retries=2):
    for attempt in range(1, max_retries + 1):
        print(f"Stage: {stage_name}, Attempt: {attempt}/{max_retries}")

        result = agent.execute()

        if result.verdict == "PASS":
            print(f"✅ {stage_name} PASSED")
            return True
        else:
            print(f"❌ {stage_name} FAILED")
            if attempt < max_retries:
                print(f"Retrying...")
                # Fix issues and retry
            else:
                print(f"Failed after {max_retries} retries")
                return False

    return False

# Execute all stages
stages = [
    ("Implementer", implementer_agent),
    ("Spec Reviewer", spec_reviewer_agent),
    ("Quality Reviewer", quality_reviewer_agent)
]

all_passed = True
for stage_name, agent in stages:
    if not execute_stage(stage_name, agent):
        all_passed = False
        break  # Stop if any stage fails

if all_passed:
    print("✅ All stages passed - proceeding to commit")
else:
    print("❌ Stage failed - workstream blocked")
```

### Step 5: Verify All Stages Passed

**If all three stages passed:**
```markdown
✅ All Stages Passed:
- Stage 1 (Implementer): ✅ PASS
- Stage 2 (Spec Reviewer): ✅ PASS
- Stage 3 (Quality Reviewer): ✅ PASS

Proceeding to commit...
```

**If any stage failed:**
```markdown
❌ Stage Failed:
- Stage 1 (Implementer): ❌ FAIL (after 2 retries)
  Reason: {details}
  Action: Fix issues, retry

Workstream blocked.
```

### Step 6: Beads CLOSED or blocked

**On success (all 3 stages passed):**
```bash
[ -n "$beads_id" ] && bd close "$beads_id" --reason "WS completed (3-stage review passed)" --suggest-next
```

**On failure (any stage failed after retries):**
```bash
[ -n "$beads_id" ] && bd update "$beads_id" --status blocked --notes="Failed at {stage}: {reason}"
```

### Step 7: Complete

```bash
# When Beads enabled: sync before commit
[ -d .beads ] && bd sync

sdp guard complete {WS-ID}
git add .
git commit -m "feat({scope}): {WS-ID} - {title}"
```

## Quality Gates

See [Quality Gates Reference](../../docs/reference/quality-gates.md)

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| No active WS | Guard not activated | `sdp guard activate` |
| File not in scope | Editing wrong file | Check WS scope |
| Coverage <80% | Missing tests | Add tests |

## See Also

- [BEADS-001 Phase 2.3](../../docs/workstreams/backlog/BEADS-001-skills-integration.md) — Beads @build spec
- [WorkstreamExecutor](../../src/sdp/beads/skills_build.py) — Python implementation
- [Full Build Spec](../../docs/reference/build-spec.md)
- [TDD Skill](../tdd/SKILL.md)
- [Guard Skill](../guard/SKILL.md)
