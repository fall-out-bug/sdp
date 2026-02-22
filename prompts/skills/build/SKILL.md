---
name: build
description: Execute workstream with TDD, guard enforcement, evidence lifecycle, and ws-verdict output
cli: sdp apply --ws
llm: Spawn subagents for 3-stage review
version: 7.0.0
changes:
  - Evidence file creation at WS start (intent.acceptance from AC)
  - Real coverage measurement, ws-verdict with ac_evidence
  - Checkpoint-based defensive branch check
  - Provenance and trace.commits in evidence
---

# build

> **CLI:** `sdp apply --ws <workstream-id>` (file operations only)
> **LLM:** Spawn subagents for TDD cycle + review

Execute a single workstream following TDD discipline.

---

## 🚨 CRITICAL RULES

1. **CHECK EXISTING CODE FIRST** - Run `@reality --quick` or grep for existing implementations before starting new features.
2. **NEVER STOP** - Continue to next workstream after commit. No summaries. No pauses.
3. **USE SPAWN OR DO IT YOURSELF** - If spawn available, use it. If not, implement manually.
4. **AUTO-CONTINUE** - After commit, immediately start next WS in dependency order.
5. **POST-COMPACTION RECOVERY** - After context compaction, run `bd ready` to find your task. Never drift to side tasks.

---

## 🔄 POST-COMPACTION PROTOCOL

**After any context compaction, you MUST:**

1. **Check active task:**
```bash
bd list --status=in_progress
bd ready
```

2. **Resume PRIMARY TASK, not side task:**
   - If you were fixing a bug as side task → return to main feature
   - If you were improving coverage → return to main feature
   - Side tasks are distractions from roadmap

3. **Ask yourself: "What was I doing BEFORE the side task?"**
   - Roadmap execution? → Back to roadmap
   - Feature implementation? → Back to feature
   - Review? → Back to review

---

## Git Safety

**CRITICAL:** Before ANY git operation, verify context.

**MANDATORY before starting work:**

```bash
# Step 1: Verify context
pwd
git branch --show-current

# Step 2: Defensive branch check via checkpoint (from @oneshot)
FEATURE_ID=$(grep "^feature_id:" docs/workstreams/backlog/${WS_ID}.md 2>/dev/null | awk '{print $2}')
EXPECTED=$(jq -r .branch .sdp/checkpoints/${FEATURE_ID}.json 2>/dev/null)
CURRENT=$(git branch --show-current)
if [ -n "$EXPECTED" ] && [ "$CURRENT" != "$EXPECTED" ]; then
  echo "ERROR: Wrong branch. Expected $EXPECTED, got $CURRENT."
  echo "Run: git checkout $EXPECTED"
  exit 1
fi

# Step 3: If sdp guard available, use it as secondary check
sdp guard context check 2>/dev/null || true
sdp guard branch check --feature=$FEATURE_ID 2>/dev/null || true
```

**NOTE:** Features MUST be implemented in feature branches. @oneshot creates the branch; @build only verifies.

---

## Evidence Lifecycle

**BEFORE any code** (at WS start):

1. Resolve beads_id from `.beads-sdp-mapping.jsonl` (sdp_id = WS_ID) or from "Feature: F{NNN} (beads_id)" line in WS file
2. Extract AC list from WS file (lines under `## Acceptance Criteria`, e.g. `- [ ] ...` or `- AC1: ...`)
3. Create `.sdp/evidence/{beads_id}.json`:

```bash
mkdir -p .sdp/evidence .sdp/ws-verdicts
RUN_ID=$(ls .sdp/runs/oneshot-F*.json 2>/dev/null | head -1 | xargs basename .json)
# Extract ACs from WS file
ACCEPTANCES=$(grep -A 20 "## Acceptance Criteria" docs/workstreams/backlog/${WS_ID}.md | grep "^- " | sed 's/^- \[.\] //' | sed 's/^- //' | head -20)
# Create evidence with intent.acceptance populated (NOT empty [])
# Include provenance: run_id, orchestrator: cursor-oneshot, captured_at
```

Schema: `intent` (acceptance, issue_id, risk_class, trigger), `plan` (workstreams, ordering_rationale), `provenance` (artifact_id, run_id, orchestrator, captured_at).

**DURING execution:** Patch `execution.branch`, `execution.changed_files` as files change.

**AFTER go test:** Patch `verification.tests`, `verification.coverage.value` (real value from `go test -coverprofile`), `verification.lint`. Add `review.self_review` with per-AC evidence: `"AC1: {text} -> satisfied by TestXxx in file.go:NN"`.

**AFTER git commit:** Patch `trace.commits = [$(git rev-parse HEAD)]`, `execution.claimed_issue_ids = [beads_id]`.

**AFTER commit:** Write `.sdp/ws-verdicts/{ws-id}.json` with `verdict`, `commit`, `quality_gates`, `ac_evidence[]` (per-AC proof).

---

## EXECUTE THIS NOW

When user invokes `@build 00-067-01`:

1. **Create evidence file** (see Evidence Lifecycle above) — BEFORE any code
2. Run CLI to setup and validate:
```bash
# Git safety verification (F065)
sdp guard context check
sdp guard branch check --feature=F067

# Guard activation
sdp guard activate 00-067-01
sdp apply --ws 00-067-01 --dry-run  # Preview first
```

2. **CHOOSE ONE:**

   **Option A (Preferred):** Spawn 3 subagents for TDD cycle:
   - **Implementer** - Write tests and code
   - **Spec Reviewer** - Verify matches spec
   - **Quality Reviewer** - Run quality gates

   **Option B (Fallback):** If subagent spawning not available, implement yourself:
   - Write test first (RED)
   - Write minimal code (GREEN)
   - Refactor while keeping tests green
   - Verify coverage >= 80%, LOC <= 200

3. **COMMIT AND CONTINUE:**
```bash
git commit -m "feat(F067): 00-067-01 - {title}"
# IMMEDIATELY start next workstream - NO PAUSE, NO SUMMARY
```

---

## How to Spawn Subagents

Use your tool's subagent capability. For example:
- Claude Code: Use Task tool with `subagent_type="general-purpose"`
- Cursor: Use agent panel
- Windsurf: Use agent spawning

---

## Subagent 1: Implementer

**Role file:** `.claude/agents/implementer.md`

**Task:**
```
You are the IMPLEMENTER for workstream 00-067-01.

Read the spec: docs/workstreams/backlog/00-067-01.md

Execute TDD cycle for each Acceptance Criteria:
1. RED: Write failing test first
2. GREEN: Write minimum code to pass
3. REFACTOR: Clean up while keeping tests green

Quality gates:
- Test coverage >= 80%
- All tests passing
- No lint errors

Output: Verdict PASS or FAIL with evidence
```

---

## Subagent 2: Spec Reviewer

**Role file:** `.claude/agents/spec-reviewer.md`

**Task:**
```
You are the SPEC COMPLIANCE REVIEWER for workstream 00-067-01.

CRITICAL: Do NOT trust the implementer's report. Verify yourself.

1. Read the actual code
2. Run tests yourself
3. Check coverage yourself
4. Verify each AC is implemented
5. Output ac_evidence mapping: for each AC, list {"ac_id": "AC1", "ac_text": "...", "evidence": "TestName in file.go:line", "status": "SATISFIED"}

Output: Verdict PASS or FAIL with evidence. Include AC_EVIDENCE: [array of ac_evidence objects]
```

---

## Subagent 3: Quality Reviewer

**Task:**
```
You are the QUALITY REVIEWER for workstream 00-067-01.

Run comprehensive quality check:
1. Test coverage (>=80%)
2. LOC check (<=200 lines per file) - MANDATORY
3. Code quality (complexity, duplication)
4. Security check
5. Lint passes

LOC Gate (MANDATORY):
```bash
for file in *.go; do
  loc=$(wc -l < "$file")
  if [ "$loc" -gt 200 ]; then
    echo "ERROR: $file is $loc LOC (max: 200)"
    exit 1
  fi
done
```

Output: Verdict PASS or FAIL with evidence
```

---

## After All Subagents Complete

**If all 3 PASS:**

```bash
# 1. Measure real coverage
go test -coverprofile=/tmp/cover.out ./... 2>/dev/null
COVERAGE=$(go tool cover -func=/tmp/cover.out 2>/dev/null | tail -1 | awk '{print $3}' | tr -d '%')
[ -z "$COVERAGE" ] && COVERAGE=0

# 2. Build ac_evidence array (per-AC proof)
# For each AC in WS file, map to test/evidence: "AC1: {text} -> satisfied by TestXxx in file:line"

# 3. Write ws-verdict file
mkdir -p .sdp/ws-verdicts
cat > .sdp/ws-verdicts/${WS_ID}.json << EOF
{
  "ws_id": "${WS_ID}",
  "feature_id": "${FEATURE_ID}",
  "verdict": "PASS",
  "commit": "$(git rev-parse HEAD)",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "quality_gates": {
    "tests": "PASS",
    "coverage": ${COVERAGE},
    "lint": "PASS",
    "loc_ok": true
  },
  "ac_evidence": [
    {"ac_id": "AC1", "ac_text": "...", "evidence": "TestXxx in file.go:NN", "status": "SATISFIED"}
  ]
}
EOF

# 4. Patch evidence file: verification.coverage.value, trace.commits, review.self_review

# 5. Commit
sdp guard complete 00-067-01 2>/dev/null || true
git add .
git commit -m "feat(${FEATURE_ID}): ${WS_ID} - {title}"
```

**If any FAIL:** Report failure, do not commit. Do not write ws-verdict.

**ac_evidence:** Populate from Implementer/Spec Reviewer output. Each AC in the WS file must have one entry: `{"ac_id": "AC1", "ac_text": "...", "evidence": "TestName in file.go:line", "status": "SATISFIED"}`.

---

## Identifier Formats

```bash
@build 00-067-01      # Workstream ID (PP-FFF-SS)
@build 99-F064-01     # Fix workstream (99-{FEATURE}-{SEQ})
@build sdp-xxx        # Beads task ID (resolved)
```

---

## Quality Gates

| Gate | Threshold | Check |
|------|-----------|-------|
| Tests | 100% pass | `go test ./...` |
| Coverage | >= 80% | `go test -cover ./...` |
| Lint | 0 errors | `golangci-lint run` |
| File Size | <= 200 LOC | `wc -l *.go` |

---

## Beads Integration

When Beads enabled:
1. **Before:** `bd update {beads_id} --status in_progress`
2. **Success:** `bd close {beads_id} --reason "WS completed"`
3. **Failure:** `bd update {beads_id} --status blocked`

---

## See Also

- `.claude/patterns/tdd.md` - TDD pattern
- `.claude/patterns/quality-gates.md` - Quality gates
- `@oneshot` - Execute all workstreams

**Implementation:** `sdp-plugin/cmd/sdp/apply.go`
