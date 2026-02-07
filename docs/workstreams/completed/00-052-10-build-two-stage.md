# 00-052-10: Update @build for Two-Stage Review

> **Beads ID:** sdp-syga
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 2 - Two-Stage Review (Quality Lock-in)
> **Size:** LARGE
> **Duration:** 4-5 days
> **Dependencies:**
> - 00-052-08 (Implementer Agent)
> - 00-052-09 (Spec Compliance Reviewer Agent)

## Goal

Update `.claude/skills/build/SKILL.md` to orchestrate two-stage review workflow.

## Acceptance Criteria

- **AC1:** @build orchestrates 3 stages: Implementer → Spec Reviewer → Quality Reviewer
- **AC2:** Max 2 retries per stage (implementer fixes issues)
- **AC3:** All stages pass → mark complete, move to completed/
- **AC4:** Any stage fails after retries → mark as blocked, create bug workstream
- **AC5:** Workflow documented with retry logic and escalation

## Files

**Modify:**
- `.claude/skills/build/SKILL.md` - Add two-stage review workflow

## Steps

### Step 1: Read Current @build SKILL.md

Understand existing workflow to determine what needs to be updated.

### Step 2: Add Two-Stage Review Section

Add to `.claude/skills/build/SKILL.md`:

```markdown
## Two-Stage Review Workflow

**Quality Lock-In:** All workstreams must pass 3 stages before completion.

### Stage 1: Implementer Agent

**Purpose:** Execute workstream with TDD discipline

**Retry Logic:** Up to 2 retries if quality gates fail

```python
for attempt in range(3):  # Initial + 2 retries
    agent = spawn_implementer_agent(ws_id)
    result = agent.execute()

    if result.quality_gates_pass:
        break  # Proceed to stage 2
    elif attempt < 2:
        log("Quality gates failed, retrying...")
        continue  # Retry
    else:
        return BLOCKED("Quality gates failed after 3 attempts")
```

### Stage 2: Spec Compliance Reviewer Agent

**Purpose:** Validate implementation matches specification

**Retry Logic:** Up to 2 retries if CHANGES_REQUESTED

```python
for attempt in range(3):  # Initial + 2 retries
    reviewer = spawn_spec_reviewer_agent(ws_id)
    result = reviewer.review()

    if result.verdict == "APPROVED":
        break  # Proceed to stage 3
    elif result.verdict == "CHANGES_REQUESTED" and attempt < 2:
        log("Spec reviewer requested changes, retrying...")
        # Return to implementer for fixes
        continue  # Retry from stage 1
    else:
        return BLOCKED("Spec compliance failed after 3 attempts")
```

### Stage 3: Quality Reviewer (Existing /review)

**Purpose:** Multi-agent quality check (QA, Security, DevOps, SRE, TechLead)

**Retry Logic:** Up to 2 retries if CHANGES_REQUESTED

```python
for attempt in range(3):  # Initial + 2 retries
    review = spawn_quality_reviewer(ws_id)
    result = review.execute()

    if result.verdict == "APPROVED":
        break  # All stages passed
    elif result.verdict == "CHANGES_REQUESTED" and attempt < 2:
        log("Quality reviewer requested changes, retrying...")
        # Return to implementer for fixes
        continue  # Retry from stage 1
    else:
        return BLOCKED("Quality review failed after 3 attempts")
```

### Completion Criteria

**Mark complete if:**
- ✅ Implementer: Quality gates pass
- ✅ Spec reviewer: APPROVED
- ✅ Quality reviewer: APPROVED

**Mark blocked if:**
- ❌ Any stage fails after 3 attempts (initial + 2 retries)
- ❌ Create bug workstream (e.g., 00-052-99-fix-{issue})
- ❌ Move workstream to `docs/workstreams/blocked/`

### Workflow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    @build {WS_ID}                            │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Stage 1: Implementer Agent (max 3 attempts)                │
│  - Read WS spec                                              │
│  - TDD cycle (Red → Green → Refactor)                        │
│  - Quality gates (coverage, vet, etc.)                       │
│  - Self-report                                               │
└─────────────────────────────────────────────────────────────┘
                            │
                 ┌──────────┴──────────┐
                 │   Gates Pass?       │
                 └──────────┬──────────┘
                   No      │      Yes
                   │       │       │
                   ▼       │       ▼
              Retry (max 2)│    Stage 2 → Spec Reviewer
                            │       │
                            ▼       ▼
                      ┌─────────────┴──────────────┐
                      │   APPROVED?                │
                      └─────────────┬──────────────┘
                        No         │         Yes
                        │          │          │
                        ▼          │          ▼
                   Back to Stage 1 │    Stage 3 → Quality Reviewer
                   (max 2 retries) │        │
                                       │       ▼
                              ┌────────────────┴───────────────┐
                              │   APPROVED?                   │
                              └────────────────┬───────────────┘
                                 No          │           Yes
                                 │           │           │
                                 ▼           │           ▼
                            Back to Stage 1 │    ✅ COMPLETE
                            (max 2 retries) │       Move to completed/
                                               │
                                               ▼
                                        Git commit + push
```

### Example Output

**Successful completion:**
```
@build 00-052-02

→ Stage 1: Implementer Agent
  Attempt 1: Quality gates pass ✅
  Coverage: 85.3%, All tests pass

→ Stage 2: Spec Reviewer
  Attempt 1: APPROVED ✅
  All AC verified in actual code

→ Stage 3: Quality Reviewer
  Attempt 1: APPROVED ✅
  QA, Security, DevOps, SRE, TechLead all pass

→ ✅ WORKSTREAM COMPLETE
→ Moving to docs/workstreams/completed/00-052-02.md
→ Commit: feat(vision): add vision extractor implementation
```

**Blocked after retries:**
```
@build 00-052-02

→ Stage 1: Implementer Agent
  Attempt 1: Quality gates failed (coverage 72%)
  Attempt 2: Quality gates failed (coverage 76%)
  Attempt 3: Quality gates failed (coverage 79%)
  ❌ BLOCKED after 3 attempts

→ Creating bug workstream: 00-052-99-fix-coverage-gap
→ Moving to docs/workstreams/blocked/00-052-02.md

→ Required actions:
  1. Improve test coverage to ≥80%
  2. Re-run @build 00-052-02
```
```

### Step 3: Update Retry Logic

Add retry loop to main workflow:

```python
def execute_build_workflow(ws_id):
    max_attempts = 3  # Initial + 2 retries

    for stage in ["implementer", "spec_reviewer", "quality_reviewer"]:
        for attempt in range(max_attempts):
            if stage == "implementer":
                result = run_implementer(ws_id)
                if result.quality_gates_pass:
                    break
                elif attempt < max_attempts - 1:
                    continue  # Retry
                else:
                    return block_workstream(ws_id, "Quality gates failed")

            elif stage == "spec_reviewer":
                result = run_spec_reviewer(ws_id)
                if result.verdict == "APPROVED":
                    break
                elif result.verdict == "CHANGES_REQUESTED" and attempt < max_attempts - 1:
                    continue  # Back to implementer
                else:
                    return block_workstream(ws_id, "Spec compliance failed")

            elif stage == "quality_reviewer":
                result = run_quality_reviewer(ws_id)
                if result.verdict == "APPROVED":
                    break
                elif result.verdict == "CHANGES_REQUESTED" and attempt < max_attempts - 1:
                    continue  # Back to implementer
                else:
                    return block_workstream(ws_id, "Quality review failed")

    # All stages passed
    return complete_workstream(ws_id)
```

### Step 4: Commit Changes

```bash
git add .claude/skills/build/SKILL.md
git commit -m "feat(build): add two-stage review workflow

- Implementer → Spec Reviewer → Quality Reviewer
- Max 2 retries per stage
- Block if fails after 3 attempts
- Create bug workstream for blocked items"
```

## Quality Gates

- Workflow is clear and unambiguous
- Retry logic is well-specified
- Escalation path documented (blocked → bug workstream)
- Example outputs provided
- Consistent with existing @build workflow

## Success Metrics

- Two-stage review is enforceable
- Retry logic prevents infinite loops
- Blocked workstreams are tracked
- Bug workstreams are created automatically
