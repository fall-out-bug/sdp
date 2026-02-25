---
name: build
description: Execute ONE workstream with TDD, guard enforcement, and ws-verdict output
cli: sdp guard activate
llm: Spawn subagents for TDD cycle
version: 8.0.0
changes:
  - F020: Remove auto-continue rules; @build does ONE WS then STOPS
  - F020: Strip evidence boilerplate to orchestrator/CLI
  - Single subagent strategy (no Option A/B ambiguity)
---

# build

> **CLI:** `sdp guard activate <workstream-id>` (scope enforcement)
> **LLM:** Execute one workstream following TDD discipline

Execute **this ONE workstream**. After commit, **STOP**. Continuation is the orchestrator's job (@oneshot / sdp orchestrate).

---

## CRITICAL RULES

1. **CHECK EXISTING CODE FIRST** — Run `@reality --quick` or grep before starting new features. Output `existing_work_summary` in ws-verdict — **required**. Short summary: files/functions/risks found before implementation.
2. **ONE WORKSTREAM** — Execute this workstream only. After commit, STOP. Do not start the next WS.
3. **USE SPAWN OR DO IT YOURSELF** — If spawn available, use it. If not, implement manually.
4. **POST-COMPACTION RECOVERY** — After context compaction, run `bd ready` to find your task. Never drift to side tasks.

---

## Git Safety

Before ANY git operation:

```bash
pwd
git branch --show-current

FEATURE_ID=$(grep "^feature_id:" docs/workstreams/backlog/${WS_ID}.md 2>/dev/null | awk '{print $2}')
EXPECTED=$(jq -r .branch .sdp/checkpoints/${FEATURE_ID}.json 2>/dev/null)
CURRENT=$(git branch --show-current)
if [ -n "$EXPECTED" ] && [ "$CURRENT" != "$EXPECTED" ]; then
  echo "ERROR: Wrong branch. Expected $EXPECTED, got $CURRENT."
  exit 1
fi
```

**NOTE:** Features MUST be implemented in feature branches. @oneshot creates the branch; @build only verifies.

---

## EXECUTE THIS NOW

When user invokes `@build 00-067-01`:

1. **Setup:**
```bash
sdp guard activate 00-067-01
```

2. **TDD cycle** (spawn subagents if available, else do yourself):
   - Implementer: RED → GREEN → REFACTOR per AC
   - Spec Reviewer: Verify each AC with evidence
   - Quality Reviewer: Coverage >= 80%, LOC <= 200, lint pass

3. **Commit and STOP:**
```bash
sdp guard deactivate 2>/dev/null || true
git add .
git commit -m "feat(F067): 00-067-01 - {title}"
# STOP. Orchestrator continues to next WS if any.
```

4. **Write ws-verdict** (required):
```bash
mkdir -p .sdp/ws-verdicts
# Populate: ws_id, feature_id, verdict, commit, quality_gates, ac_evidence[], existing_work_summary (required)
```
**Required fields:** `existing_work_summary` — one line summary of pre-existing code/tests found before implementation. **Output must validate against** `schema/ws-verdict.schema.json` ([ws-verdict.schema.json](../../schema/ws-verdict.schema.json)).

Evidence lifecycle (create/patch `.sdp/evidence/*.json`) is orchestrator or post-build CLI responsibility.

---

## Subagent Tasks (if spawning)

**Implementer:** TDD per AC. Output verdict + evidence.

**Spec Reviewer:** Verify code matches spec. Output ac_evidence per [ws-verdict.schema.json](../../schema/ws-verdict.schema.json): `{"ac":"AC text","met": true|false,"evidence":"file:line or test name"}`.

**Quality Reviewer:** Coverage >= 80%, LOC <= 200, lint. Output verdict.

---

## Quality Gates

| Gate | Threshold |
|------|-----------|
| Tests | 100% pass |
| Coverage | >= 80% |
| Lint | 0 errors |
| File Size | <= 200 LOC |

---

## Beads Integration

- **Before:** `bd update {beads_id} --status in_progress`
- **Success:** `bd close {beads_id} --reason "WS completed"`
- **Failure:** `bd update {beads_id} --status blocked`

---

## Few-Shot Examples

**Good ws-verdict (all gates green):**
```json
{
  "ws_id": "00-067-01",
  "feature_id": "F067",
  "verdict": "PASS",
  "commit": "a1b2c3d",
  "quality_gates": {"tests_pass": true, "lint_clean": true, "coverage_pct": 85.2, "coverage_threshold": 80, "max_file_loc": 142, "build_ok": true, "vet_ok": true},
  "ac_evidence": [
    {"ac": "User can reset password via email", "met": true, "evidence": "TestResetPassword in internal/auth/reset_test.go:42"},
    {"ac": "Rate limit 5/min per email", "met": true, "evidence": "TestRateLimit in internal/auth/reset_test.go:89"}
  ],
  "existing_work_summary": "Found ResetToken in pkg/auth; extended with email flow."
}
```

**Bad — missing existing_work_summary:**
```json
{"ws_id": "00-067-01", "feature_id": "F067", "verdict": "PASS", "ac_evidence": [...]}
```
Reason: `existing_work_summary` is required. Add one-line summary of pre-existing code/tests found before implementation.

**Bad — ac_evidence without evidence field:**
```json
{"ac": "User can reset password", "met": true}
```
Reason: Each ac_evidence entry must include `evidence` (file:line or test name).

Schema: `schema/ws-verdict.schema.json` (from sdp root; project: `sdp/schema/`)

---

## See Also

- `@oneshot` — Orchestrator that invokes @build per WS
- `@tdd` — TDD pattern
