---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation + PromptOps)
cli: sdp quality all
version: 16.0.0
changes:
  - "16.0.0: Fixed schema consistency - all 7 reviewers always spawned (F098 P1 fix)"
  - "15.0.0: Add risk-based reviewer selection"
  - "14.3.0: Add @go-modern checks for Go review surfaces"
  - "14.2.0: Handoff block when CHANGES_REQUESTED"
  - "14.1.0: Language-agnostic (platform-agnostic spawn, agents/ path)"
  - "14.0.0: Compress to ~150 lines (P2 remediation)"
---

# review

> **CLI:** `sdp quality all` | **LLM:** Spawn all 7 specialist subagents with risk-based depth allocation

Comprehensive multi-agent quality review. All 7 reviewers always spawned; risk patterns determine depth, not presence.

---

## All 7 Reviewers Always Spawned

**Base contract:** qa, security, devops, sre, techlead, docs, promptops — always present in verdict JSON.

Risk patterns determine **review depth**, not **reviewer presence**.

### Risk-Based Depth Allocation

LOC tiers set baseline depth; risk patterns override for specific files:

**LOC tiers** (baseline):
| LOC Range | Deep Reviewers |
|-----------|---------------|
| < 50 | qa, techlead |
| 50–200 | qa, security, techlead |
| > 200 | all 7 |

**Risk patterns** (additive override by file path):
| Pattern | Extra Deep Reviewers |
|---------|---------------------|
| `**/auth/**`, `**/crypto/**` | security, qa |
| `**/.github/workflows/**`, `**/ci/**` | devops, sre |
| `**/migrations/**`, `**/db/**` | sre, security |

Full config: `.sdp/config.yml` under `review` section.

### Flag Overrides

| Flag | Behavior |
|------|----------|
| `--full` | All 7 reviewers with full depth |
| `--quick` | All 7 reviewers, but only 2-3 do deep review (rest rubber-stamp) |

---

## EXECUTE THIS NOW

When user invokes `@review F{XX}`:

1. **Run CLI:** `sdp quality all`
2. **Determine depth:** Match risk patterns to identify which reviewers get deep focus (rest get rubber-stamp).
3. **Spawn all 7 subagents IN PARALLEL** (use your platform's subagent spawn). **DO NOT skip.**

**All 7 roles always spawned:** qa, security, devops, sre, techlead, docs, promptops

**Per-subagent task template** (replace F{XX}, round-N, {role}):

**5-step evaluation structure:**

1. **SCOPE:** What files/packages does this feature touch?
2. **RISK MAP:** Top 3 risk areas for your domain ({role}) in this scope?
3. **EVIDENCE:** For each risk, what did you find? (file:line, test name, config entry)
4. **SEVERITY:** P0 = exploitable in production. P1 = breaks on edge case. P2 = maintenance debt. P3 = style.
5. **VERDICT:** PASS if max severity ≤ P2. FAIL if any P0/P1.

For Go files, also check modern stdlib usage (`slices`, `maps`, `strings.Cut`, `strings.CutPrefix`, `any`) instead of legacy patterns.

For each finding: `bd create --silent --labels "review-finding,F{XX},round-1,{role}" --priority={0-3} --type=bug`. Output: `FINDINGS_CREATED: id1 id2` or `FINDINGS_CREATED: (none)`. Output verdict: `PASS` or `FAIL`.

**Role files:** `prompts/agents/qa.md`, `prompts/agents/security.md`, `prompts/agents/devops.md`, `prompts/agents/sre.md`, `prompts/agents/tech-lead.md`. Docs and PromptOps: inline.

**Docs expert:** Check drift (`sdp drift detect`), AC coverage (jq `.ac_evidence|length` vs WS file). Labels: `review-finding,F{XX},round-1,docs`

**PromptOps expert:** Review prompts/skills, prompts/agents, prompts/commands. Check: language-agnostic, no phantom CLI, no handoff lists, skill size ≤200 LOC. Labels: `review-finding,F{XX},round-1,promptops`. Output `checks` array per schema/review-verdict.schema.json.

## Write Plan (F101)

Before writing review output files (verdict, findings), emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason. Covers review_verdict.json and any findings files created during review.
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"review"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `sdp/schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @review <target>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Write-plan flags:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

> **Note:** `--dry-run` and `--yes` are orthogonal to skill mode flags (`--full`, `--quick`). They can be combined with any mode (e.g. `@review F098 --quick --dry-run`).

## After All Complete — Synthesis Phase

**MUST include all 7 reviewers in `reviewers` object.** All 7 roles always spawned; missing any = FAIL (set verdict=CHANGES_REQUESTED).

1. **CONFLICT CHECK:** Do any two reviewers contradict? If yes, create escalation finding with both perspectives; add to `synthesis.conflicts`.
2. **COVERAGE CHECK:** Did any reviewer report 0 findings? Note role in `synthesis.rubber_stamps`.
3. **ADVERSARIAL SYNTHESIS:** Before final verdict, ask "What if we're wrong?" For each PASS, note one plausible blind spot. Add to synthesis.
4. **VERDICT:** **APPROVED** only if **all 7 reviewers** have an entry in `reviewers` and verdict=PASS. **Missing reviewer = FAIL** (set verdict=CHANGES_REQUESTED). **CHANGES_REQUESTED** if any FAIL or escalation.

**Before final verdict:** Verify `reviewers` contains all 7 roles: qa, security, devops, sre, techlead, docs, promptops. If any role is missing, set verdict=CHANGES_REQUESTED and add a note.

**Synthesize:** `## Feature Review: F{XX}` with `### {ROLE}: PASS/FAIL` for all 7 reviewers.

**Save verdict** to `.sdp/review_verdict.json` (required for @deploy, @oneshot). **Output must validate against** `schema/review-verdict.schema.json` before saving.

```json
{
  "feature": "F{XX}",
  "verdict": "APPROVED|CHANGES_REQUESTED",
  "timestamp": "...",
  "round": 1,
  "reviewers": {
    "qa": {"verdict": "PASS", "findings": []},
    "security": {"verdict": "PASS", "findings": []},
    "devops": {"verdict": "PASS", "findings": []},
    "sre": {"verdict": "PASS", "findings": []},
    "techlead": {"verdict": "PASS", "findings": []},
    "docs": {"verdict": "PASS", "findings": []},
    "promptops": {"verdict": "PASS", "findings": []}
  },
  "reviewer_selection": {
    "deep_reviewers": ["qa", "security"],
    "risk_patterns_matched": ["**/auth/**"],
    "flag": null
  },
  "finding_ids": [],
  "blocking_ids": [],
  "synthesis": {
    "conflicts": [],
    "rubber_stamps": ["devops", "sre", "techlead", "docs", "promptops"]
  },
  "summary": "..."
}
```

**Priority:** P0/P1 block; P2/P3 track only.

**When verdict=CHANGES_REQUESTED** — output this handoff block prominently:

```
---
## Next Step
Run `@design phase4-remediation` with findings to create workstreams.
---
```

---

## Beads

`bd create --title "{AREA}: {desc}" --priority {0-3} --labels "review-finding,F{XX},round-{N},{role}" --type bug --silent`

Replace `F{NNN}` with feature ID, `round-{N}` with iteration (e.g. round-1), `{role}` with qa/security/devops/sre/techlead/docs/promptops.

After creating findings, include in subagent output: `FINDINGS_CREATED: id1 id2 id3`

---

## Few-Shot Examples

**Good P0 finding (Security):**
```
bd create --title "Security: auth bypass via missing role check in API handler" --priority 0 --labels "review-finding,<feature-id>,round-1,security" --type bug --silent
```

**Good P2 finding (Docs):**
```
bd create --title "Docs: typo in README deployment section" --priority 2 --labels "review-finding,<feature-id>,round-1,docs" --type bug --silent
```

**Bad — vague finding (no file:line):**
```
bd create --title "Security: possible vulnerability" --priority 0 ...
```
Reason: P0 requires evidence. Add file:line or downgrade to P2.

**Good — no findings (explicit output):**
```
SCOPE: internal/auth/*.go (3 files). RISK MAP: token validation, rate limit. EVIDENCE: All checks present. VERDICT: PASS
FINDINGS_CREATED: (none)
PASS
```

## See Also
@oneshot — review-fix loop | @deploy — requires APPROVED verdict | @go-modern — Go modernization checklist
