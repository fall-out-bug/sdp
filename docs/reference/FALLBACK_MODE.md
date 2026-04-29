# Fallback Mode

Fallback mode is the manual path for harnesses that cannot spawn subagents.

The goal is not “approximate the workflow.” The goal is to preserve the same
outputs and review discipline, just sequentially instead of in parallel.

## When To Use It

Use fallback mode when the harness cannot reliably dispatch subagents for:

- `@review`
- `@vision`
- `@reality`
- `@build`
- `@feature`

If native subagent spawning works, do not use fallback mode.

## General Workflow

1. Confirm that subagent spawning is unavailable or unreliable.
2. Pick the skill you need.
3. Execute the matching manual checklist below.
4. Produce the same artifacts you would expect from the normal path.
5. Verify outputs against the relevant schemas and acceptance criteria.

## `@review` Fallback

**Normal mode:** parallel specialist review.  
**Fallback expectation:** 7 sections, in this order:

1. QA
2. Security
3. DevOps
4. SRE
5. Tech Lead
6. Docs
7. Verdict

### Checklist

1. Read the workstream or PR scope.
2. For each role, identify the top risks, collect evidence, and assign severity.
3. Record findings with concrete file and line references where possible.
4. Produce a verdict section that summarizes blocking vs non-blocking findings.
5. If your repo uses review verdict artifacts, update `.sdp/review_verdict.json`.

### Minimum output

- 6 specialist sections plus 1 verdict section
- explicit blocking vs non-blocking call
- evidence for each serious finding

## `@vision` Fallback

**Normal mode:** analyst + architect + product-manager.  
**Fallback expectation:** one structured strategy pass.

### Checklist

1. Ask the minimum clarifying questions: problem, user, success metric, MVP scope.
2. Write short sections for:
   - product
   - market
   - technical feasibility
   - UX constraints
   - business value
   - delivery risk
3. Produce the intended planning artifacts for the repo:
   - product vision
   - PRD or feature brief
   - roadmap delta if needed

## `@reality` Fallback

**Normal mode:** analyst + architect.  
**Fallback expectation:** one evidence-backed repo audit.

### Checklist

1. Detect project type and main tech stack.
2. Check architecture shape, tests, docs, TODO/HACK debt, and obvious risk areas.
3. Report:
   - current state
   - top issues
   - quick wins
   - mismatch between docs and code
4. If a vision or roadmap exists, compare intended state vs actual state.

## `@build` Fallback

**Normal mode:** implementer + spec-reviewer + quality-reviewer.  
**Fallback expectation:** sequential TDD plus self-review.

### Checklist

1. Confirm the scope is executable. If the task is not scoped, stop and clarify.
2. Perform TDD:
   - RED: add or update a failing test
   - GREEN: implement the smallest fix
   - REFACTOR: clean up without changing behavior
3. Review acceptance criteria one by one and attach evidence.
4. Run the relevant quality gates.
5. If the repo uses workstream verdict artifacts, update `.sdp/ws-verdicts/<ws-id>.json`.

### Required outputs

- code change
- test evidence
- AC evidence
- gate results

## `@feature` Fallback

**Normal mode:** analyst + architect + planner.  
**Fallback expectation:** sequential discovery and decomposition.

### Checklist

1. Confirm the problem, outcome, and non-goals.
2. Decide whether the work needs:
   - full discovery
   - direct workstream design
   - roadmap extraction
3. Produce scoped workstreams with:
   - goal
   - scope files
   - acceptance criteria
   - dependencies
4. If the repo uses issue tracking, create or link the executable units.

## Command Mapping

The canonical command mapping for fallback-capable harnesses lives in:

- [`prompts/commands.yml`](../../prompts/commands.yml)

## Rule Of Thumb

Fallback mode is slower, not lower quality.

If the manual checklist cannot produce the same artifact quality as the normal
path, stop and say so instead of pretending the harness supports the workflow.
