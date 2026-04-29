---
name: delivery-loop
description: Autonomous delivery cycle — bootstrap → preflight → build/review/fix loop (bounded) → PR → codex review (bounded, stable-N) → closeout. Replaces manual "build ws → review → fix → PR → codex → fix".
version: 2.0.0
tags:
  - delivery
  - orchestration
  - loop
requires_cli:
  - bd
  - git
  - go
  - gh
  - codex
compatibility:
  - claude-code
  - opencode
  - cursor
  - codex
---

# Delivery Loop

## Purpose

Run the full feature delivery cycle autonomously without user intervention.

**Entry condition:** operator invoked `@delivery-loop` (no arguments) or `@delivery-loop --resume` after compaction.
**Exit condition:** PR merged + children closed + worktree removed, OR explicit `--abort`.

## Invocation

```
@delivery-loop           # full run, picks feature from bd ready
@delivery-loop --resume  # read .sdp/checkpoints/${FEATURE}.json, continue
@delivery-loop --abort   # cleanup: kill subagents, stash work, archive checkpoint
```

Per-harness dispatch is delegated to `scripts/sdp-dispatch.sh`. Adding a harness = one `case` branch in that script.

## Phases (declarative)

Phases below are **the source of truth**; the prose "Loop structure" section renders each of them for humans. Per-feature overrides live at `.sdp/delivery.yaml` (same schema; any key present there wins; phases can be disabled with `enabled: false`).

```yaml
phases:
  - name: bootstrap
    required: true
  - name: preflight
    required: true
  - name: build
    required: true
    max_cycles: 5
    wallclock_budget_hours: 4
    backoff_seconds: [30, 120, 300, 600]
  - name: design_gap          # triggered only from build review
    required: false
    max_scope_delta: 2
  - name: impact_review
    required: true
    subagent_timeout_minutes: 10
  - name: traceability        # advisory — see "Traceability gate"
    required: false
    mode: advisory            # warn-only; promote by setting mode: gate
  - name: pr
    required: true
  - name: codex
    required: true
    max_cycles: 4
    wallclock_budget_hours: 2
    stable_n: 2
    enabled: true             # set false for offline mode (.sdp/delivery.yaml override)
  - name: closeout
    required: true
whole_loop_budget_hours: 72   # runaway detector only, not a delivery SLO
```

**Override example** — offline dogfood that skips codex:

```yaml
# .sdp/delivery.yaml
phases:
  - name: codex
    enabled: false
```

**Override example** — feature-specific contract-gen stage:

```yaml
# .sdp/delivery.yaml
phases:
  - name: contract_gen        # new phase, must be defined in skill logic
    required: true
    before: build
```

Phases not listed in an override inherit the skill defaults. Unknown phase names in an override are errors (operator gets a message; loop refuses to start).

## Loop structure

```
PHASE 0: BOOTSTRAP
  1. **Pick feature (deterministic):** `pick="$(scripts/deliver-pick.sh)"; EPIC_ID="${pick%%$'\t'*}"; TITLE="${pick#*$'\t'}"`
       - exit 0 → EPIC_ID + title printed; derive `FEATURE` from EPIC_ID metadata or title (e.g. F129 from "F129: ...")
       - exit 4 → no deliverable feature in ready queue → exit Phase 0 cleanly (no work)
       - exit 1 → bd error → escalate
       - **Do NOT improvise a picker. Do NOT pick a workstream-leaf task. Do NOT pick a coordination/meta epic.** Tag any program/coordination epic with `bd label add <id> coordination` so the picker skips it.
  2. Identify workstreams: cross-reference `docs/workstreams/backlog/${FEATURE}-*.md` with `bd list -n 200 | grep "^${FEATURE}-"`
  3. **Acquire delivery slot (HARD GATE):** `scripts/deliver-acquire.sh ${FEATURE} ${EPIC_ID}`
       - exit 0 → claim + lock acquired, proceed
       - exit 2 → foreign claim (another operator owns this epic) → exit Phase 0 cleanly
       - exit 3 → lock held by live PID on this host (parallel shell of same user) → exit Phase 0 cleanly
       - exit 1 → bd error → escalate
       - **No state-mutating step (worktree create, checkpoint write, @build dispatch) may run before this returns 0.**
  4. Create worktree: `git worktree add .worktrees/${FEATURE}`
  5. Write initial `.sdp/checkpoints/${FEATURE}.json` (schema v2)

PHASE 0.5: PREFLIGHT
  # Fail-fast before any build work
  - command -v bd gh go codex       # all required CLIs present
  - gh auth status                  # gh authenticated
  - ws_count == bd_count            # workstream files match bead children
  - ws_count > 0                    # feature has at least one workstream
  - bd show ${EPIC_ID} assignee == ${USER}   # epic claimed to this operator
  - lock file pid == $$              # we still hold the delivery slot (no hijack)
  - disk free > 2GB                 # room for worktree + build artifacts
  On any failure: emit actionable error, release lock, unclaim epic, exit 2.

PHASE 1: BUILD LOOP (max 5 cycles, max 4h wallclock)
  repeat:
    1. Dispatch @build subagents (one per WS, parallel, haiku; per-subagent timeout 20m)
    2. Dispatch @review subagent (fresh context, sonnet; timeout 10m)
    3. If APPROVED (zero findings) → break
    4. Else:
       - P1/P2/P3 findings → dispatch @fix subagents per finding (haiku; timeout 15m)
       - Design gaps → dispatch @design subagent; cap: max 2 new WS per feature (scope_delta ≤ 2)
    5. Exponential backoff between cycles: 30s / 2m / 5m / 10m
  until: APPROVED, OR cycle=5 hit, OR 4h wallclock hit

  On cap hit with residual P3:
    Operator MUST paste the deferred-P3 list into a new bead description
    (file:line: description for each P3). NOT a plain y/N prompt.
    Create: `bd create --parent ${EPIC_ID} --type=task --priority=3 \
            --title="Deferred P3 from ${FEATURE}" --description=<pasted list>`
    Continue to Phase 2.

PHASE 2: PR CREATION
  1. @review --dimension impact (check blast radius; timeout 10m)
  2. Traceability gate — `scripts/traceability-gate.sh ${FEATURE}`
     - For each WS: grep AC[0-9]+ tokens in `docs/workstreams/backlog/${FEATURE}-*.md`
       confirm they appear in the touched `*_test.go` files (or marked NONE).
     - For schema changes (schema/*.schema.json): require an adjacent `_test.go`
       invoking `jsonschema.Validate` OR a `testdata/<schema>` directory.
     - Mode `advisory` (default): emit warnings, do NOT block.
     - Mode `gate` (promoted via `.sdp/delivery.yaml`): missing coverage → fail phase 2.
  3. If impact + traceability OK: run `./scripts/run_go_quality_gates.sh` LOCALLY — must be green
  4. `gh pr create` — MUST NOT run before gates green
  5. Record `pr_number` in checkpoint

PHASE 3: CODEX REVIEW LOOP (max 4 cycles, max 2h wallclock, stable-N=2)

  **HARD RULE — DO NOT IMPROVISE A SKIP.** Phase 3 is required by default.
  The ONLY supported way to skip codex review is the operator override
  `.sdp/delivery.yaml { phases: [{ name: codex, enabled: false }] }`. If
  the override is not present, you MUST run the loop. You MUST NOT mark
  the phase SKIPPED for any of these reasons:
    - "codex model availability"
    - "spark not supported" / "gpt-5.x too slow" / "no model works"
    - "codex command failed" (that means retry, not skip)
    - "running it would take too long" (the budget is 2h wallclock —
      the loop self-terminates without inventing a reason)
  If the codex command exits non-zero or returns invalid JSON: log the
  exact stderr, sleep with the backoff in step 5, retry. After 4 cycles
  or 2h, exit the loop honestly with phase_status=exhausted and emit
  the last error. Do not pretend the phase succeeded; do not pretend
  the phase was skippable.

  Default invocation passes no `--model` flag; the rescue command picks
  the best available model. Override with `--model spark` explicitly
  only when the operator asks. Do not pre-select a model and then
  declare it unavailable.

  repeat:
    1. scripts/sdp-dispatch.sh codex_review "Review PR #${PR}. Steps: (1) read all changed files, (2) run ./scripts/run_go_quality_gates.sh, (3) emit JSON {tests_passed: bool, findings: [{file, line, rule, symbol, severity, msg}]}. Do not skip tests."
    2. Parse codex JSON output. If parse fails or stderr non-empty, log both verbatim and proceed to step 5 retry.
    3. Dedupe findings against prior cycle by hash(rule + symbol_path + normalized_snippet). NOT file:line:rule — line shifts invalidate naive hashes.
    4. Mark findings absent for ≥2 consecutive cycles as "non-reproducible candidate" — these are NEVER auto-closed at v1. They enter **manual Phase-4 triage**: the operator sees the list in Phase 4 and decides close vs re-raise. (Technician minority: rename this to "auto-close" only if/when an AST-unchanged check lands — see §7.3 of the design doc.)
    5. If zero NEW findings + tests pass → consecutive_clean_cycles++. Break when consecutive_clean_cycles == 2.
    6. Else: dispatch @fix per finding (haiku/sonnet; timeout 15m), run gates locally, `git push`.
       Backoff between cycles: 30s / 2m / 5m / 10m.
  until: 2 consecutive clean cycles, OR cycle=4 hit, OR 2h wallclock hit

  **On exit, write to checkpoint:**
    phase_status: "done"      → loop converged (2 clean cycles)
    phase_status: "exhausted" → cycle/wallclock budget hit
    phase_status: "skipped"   → ONLY if .sdp/delivery.yaml override disables phase
    phase_status: "error"     → something else; capture exact error, do not invent

PHASE 4: CLOSEOUT
  1. Confirm merge: `gh pr view ${PR} --json state -q .state == "MERGED"`
  2. Manual Phase-4 triage: operator reviews "non-reproducible candidates" from Phase 3 (list shown with file:line). Either close-as-resolved or re-raise as follow-up bead. No silent auto-close.
  3. Batch close children: `bd close ${WS1} ${WS2} ... --reason "merged via PR#${PR}"`
  4. Close epic: `bd close ${EPIC_ID}`
  5. `scripts/beads_transport.sh export && git push`
  6. `git worktree remove .worktrees/${FEATURE}`
  7. `git push origin --delete ${BRANCH}` (remote branch cleanup)
  8. Archive: `mv .sdp/checkpoints/${FEATURE}.json .sdp/archive/delivered/${FEATURE}-$(date -u +%Y%m%dT%H%M%SZ).json`
  9. Release lock

WHOLE-LOOP BUDGET
  72h hard ceiling (runaway detector only — not a delivery SLO).
  On ceiling: write `phase_status: "exhausted"`, post PR comment summarizing state, exit non-zero.
```

## Subagent model policy

| Task | Model | Timeout |
|------|-------|---------|
| @build per WS | haiku | 20m |
| @fix per finding | haiku | 15m |
| @review | sonnet | 10m |
| @review --dimension impact | sonnet | 10m |
| @design (new WS) | sonnet | 15m |
| codex:rescue | default (Codex CLI) | 30m |

Whole-loop wallclock ceiling: 72h. Phase budgets: build=4h, codex=2h.

## Rules

**Resolve or explicitly defer P3 with operator signoff.** (Renamed from "Never skip P3".) All findings from @review block the loop inside cycles 1–4. At cycle 5 the operator may defer remaining P3 — but only by manually pasting the `file:line: description` list into a new bead. Plain confirmation prompts are disallowed.

**Never create PR with red tests.** `./scripts/run_go_quality_gates.sh` must be green **locally** before `gh pr create`. Not after.

**Codex must run tests.** The codex prompt MUST include "run ./scripts/run_go_quality_gates.sh and report failures". Never send codex a code-only review prompt.

**Codex output must be structured.** Codex is invoked with a JSON output contract so findings can be deduped across non-deterministic runs.

**Independent WS → parallel subagents.** Check for file overlap before dispatching parallel @build agents. Overlapping WS → sequential. (Enforcement mechanism: WS frontmatter `touches:` — tracked as deferred work.)

**Max parallel subagents: 5.** Orchestrator enforces via a semaphore; queued agents start after first batch completes.

## Checkpoint schema v2

Location: `.sdp/checkpoints/${FEATURE}.json` (per-feature, not single-slot).
Atomic writes: `scripts/sdp-checkpoint-write.sh` (tmp+rename; only orchestrator writes, subagents return structured output which orchestrator merges).

```json
{
  "schema_version": 2,
  "skill": "delivery-loop",
  "feature_id": "F134",
  "epic_bead_id": "sdplab-xxx",
  "worktree_path": ".worktrees/F134",
  "branch": "feature/F134-...",
  "pr_number": null,
  "phase": 1,
  "step": "build",
  "phase_status": "running",
  "cycle_number": 2,
  "max_cycles": 5,
  "consecutive_clean_cycles": 0,
  "ws_done": ["00-134-01"],
  "ws_in_progress": ["00-134-02", "00-134-03"],
  "findings": [
    {
      "id": "F-001",
      "hash": "sha1(rule+symbol_path+normalized_snippet)",
      "file": "x.go",
      "line": 42,
      "rule": "gofmt",
      "severity": "P2",
      "status": "fixing"
    }
  ],
  "scope_delta_count": 0,
  "started_at": "2026-04-22T09:00:00Z",
  "deadline": "2026-04-25T09:00:00Z",
  "last_heartbeat": "2026-04-22T11:00:00Z",
  "subagent_pids": [12345, 12346]
}
```

Heartbeat: every 2 min orchestrator updates `last_heartbeat`. External tooling (e.g., `bd status` extension) can distinguish stuck from working.

## Abort & rollback

`@delivery-loop --abort` performs:

1. Kill tracked subagent PIDs from `checkpoint.subagent_pids`
2. `bd update ${EPIC_ID} --status blocked --notes "aborted at phase=${P} cycle=${C}"`
   (NOT `--release` / unclaim — preserve claim history)
3. `git stash push -u -m "delivery-loop abort ${FEATURE}"` in worktree
4. `mv .sdp/checkpoints/${FEATURE}.json .sdp/archive/aborted/${FEATURE}-$(date -u +%Y%m%dT%H%M%SZ).json`
5. Release the lock: `rm .sdp/locks/deliver-${FEATURE}.lock`
6. Print recovery steps: stash ref, checkpoint archive path, how to resume later

## Compaction recovery

If loop is interrupted (compaction, crash):

1. `cat .sdp/checkpoints/${FEATURE}.json` → read phase + step + cycle_number + last_heartbeat
2. `bd show ${EPIC_ID}` → verify claim still held by ${USER}
3. `git diff main --name-only` → confirm which WS are implemented
4. Resume from the current phase's step — do NOT restart from WS 1
5. Stale heartbeat (>10 min old + no live subagent PIDs) → treat as interrupted, resume from the last completed step

## Output (per phase)

```
## Phase 0: Bootstrap
Feature picked: F134 (sdplab-xxx) — "Feature title"
Workstreams: 00-134-01, 00-134-02, 00-134-03 (3 files, 3 beads, drift=0)
Lock: .sdp/locks/deliver-F134.lock acquired (pid=12340)
Worktree: .worktrees/F134 created
Checkpoint: .sdp/checkpoints/F134.json initialized

## Phase 0.5: Preflight
bd✓ gh✓ go✓ codex✓  gh-auth✓  ws-count=3  bd-count=3  claim=${USER}✓  disk=45GB✓

## Phase 1: Build Loop — Cycle 1/5 [00:02 elapsed]
WS dispatched: 00-134-01, 00-134-02, 00-134-03 (parallel, haiku)
WS completed: 3/3 (4m avg)
@review verdict: FINDINGS (3: 1×P2, 2×P3)
  P2: pkg/x/foo.go:42 missing error wrap → dispatching @fix
  P3: pkg/x/bar.go:91 inconsistent naming → dispatching @fix
  P3: pkg/y/baz.go:12 missing doc comment → dispatching @fix

## Phase 1: Build Loop — Cycle 2/5 [00:08 elapsed]
@review verdict: APPROVED — zero findings

## Phase 2: PR
Impact review: OK (blast radius: 2 packages, 0 exported symbols changed)
Quality gates (local): 47/47 passed, 0 failed
PR: #123 created

## Phase 3: Codex Review — Cycle 1/4 [00:15 elapsed]
Codex findings: 2 (1 test failure, 1 code issue)
  test: TestFoo panics → dispatching @fix
  code: pkg/y/baz.go:91 unused import → dispatching @fix
Tests after fix: 47/47 passed
Push: done
consecutive_clean_cycles: 0

## Phase 3: Codex Review — Cycle 2/4 [00:22 elapsed]
Codex: zero findings, tests pass
consecutive_clean_cycles: 1

## Phase 3: Codex Review — Cycle 3/4 [00:28 elapsed]
Codex: zero findings, tests pass
consecutive_clean_cycles: 2 → EXIT

## Phase 4: Closeout
Non-reproducible candidates: (none)
Beads closed: 00-134-01, 00-134-02, 00-134-03, sdplab-xxx (4 issues)
Beads export: pushed
Worktree: removed
Branch (remote): deleted
Checkpoint archived: .sdp/archive/delivered/F134-20260422T113000Z.json
Lock released.

Done. PR #123 merged, F134 closed.
```

## References

- `docs/reference/go-patterns.md` — Go conventions for @build subagents
- `AGENTS.md` — beads workflow, quality gates
- `.agents/skills/build.md` — build skill (worktree bootstrap anchor)
- `.agents/skills/review.md` — review dimensions
- `scripts/run_go_quality_gates.sh` — quality gate script
- `scripts/sdp-dispatch.sh` — per-harness subagent + codex invocation
- `scripts/sdp-checkpoint-write.sh` — atomic checkpoint writer
- `docs/plans/2026-04-22-deliver-skill-review-design.md` — design rationale and consensus record
