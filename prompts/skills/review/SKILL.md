---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation + PromptOps)
cli: sdp quality all
version: 14.0.0
changes:
  - "14.0.0: Compress to ~150 lines (P2 remediation)"
  - Subagent tasks consolidated into template
---

# review

> **CLI:** `sdp quality all` | **LLM:** Spawn 7 specialist subagents

Comprehensive multi-agent quality review.

---

## EXECUTE THIS NOW

When user invokes `@review F{XX}`:

1. **Run CLI:** `sdp quality all`
2. **Spawn 7 subagents IN PARALLEL** (Task tool, agent panel). **DO NOT skip.** CLI is basic; full review needs subagents.

**Roles:** qa, security, devops, sre, techlead, docs, promptops

**Per-subagent task template** (replace F{XX}, round-N, {role}):

**5-step evaluation structure** — complete in order:

1. **SCOPE:** What files/packages does this feature touch? (list from checkpoint or scope files)
2. **RISK MAP:** For your domain ({role}), what are the top 3 risk areas in this scope?
3. **EVIDENCE:** For each risk area, what did you find? (file:line, test name, config entry)
4. **SEVERITY:** Classify each finding. P0 = exploitable in production. P1 = breaks on edge case. P2 = maintenance debt. P3 = style.
5. **VERDICT:** PASS if max severity ≤ P2. FAIL if any P0/P1.

For each finding: `bd create --silent --labels "review-finding,F{XX},round-1,{role}" --priority={0-3} --type=bug`. Output: `FINDINGS_CREATED: id1 id2` or `FINDINGS_CREATED: (none)`. Output verdict: `PASS` or `FAIL`.

**Role files:** `.claude/agents/qa.md`, `security.md`, `devops.md`, `sre.md`, `tech-lead.md`. Docs and PromptOps: inline (see below).

**Docs expert:** Check drift (`sdp drift detect`), AC coverage (jq `.ac_evidence|length` vs WS file). Labels: `review-finding,F{XX},round-1,docs`

**PromptOps expert:** Review sdp/prompts/skills, agents, commands. Check: language-agnostic, no phantom CLI, no handoff lists, skill size ≤200 LOC. Labels: `review-finding,F{XX},round-1,promptops`. Output `checks` array: `[{"check_id":"language-agnostic","status":"pass","note":"..."},{"check_id":"no-phantom-cli","status":"pass"},...]` per schema/review-verdict.schema.json.

---

## After All Complete — Synthesis Phase

**MUST include all 7 roles in `reviewers`:** qa, security, devops, sre, techlead, docs, promptops. **Missing role = FAIL** (set verdict=CHANGES_REQUESTED).

1. **CONFLICT CHECK:** Do any two reviewers contradict? (e.g. Security says "add auth" but SRE says "remove auth middleware for latency"). If yes, create one escalation finding with both perspectives; add to `synthesis.conflicts`.
2. **COVERAGE CHECK:** Did any reviewer report 0 findings? Note role in `synthesis.rubber_stamps` for transparency (does not by itself change verdict).
3. **ADVERSARIAL SYNTHESIS:** Before final verdict, ask: "What if we're wrong?" For each PASS, note one plausible blind spot (e.g. "QA passed but load tests not run"). Add to synthesis. Does not change verdict unless evidence supports it.
4. **VERDICT:** **APPROVED** only if **ALL 7 roles** have an entry in `reviewers` and verdict=PASS. **Missing role = FAIL** (set verdict=CHANGES_REQUESTED). **CHANGES_REQUESTED** if any FAIL or escalation.

**Before final verdict:** Verify `reviewers` contains exactly these keys: **qa, security, devops, sre, techlead, docs, promptops**. If any is missing, set verdict=CHANGES_REQUESTED and add a note.

**Synthesize:** `## Feature Review: F{XX}` with `### QA: PASS/FAIL`, etc.

**Save verdict** to `.sdp/review_verdict.json` (required for @deploy, @oneshot). **Output must validate against** `schema/review-verdict.schema.json` before saving.

```json
{"feature":"F{XX}","verdict":"APPROVED|CHANGES_REQUESTED","timestamp":"...","round":1,"reviewers":{"qa":{"verdict":"PASS","findings":[]},"security":{"verdict":"PASS","findings":[]},"devops":{"verdict":"PASS","findings":[]},"sre":{"verdict":"PASS","findings":[]},"techlead":{"verdict":"PASS","findings":[]},"docs":{"verdict":"PASS","findings":[]},"promptops":{"verdict":"PASS","findings":[],"checks":[{"check_id":"language-agnostic","status":"pass"},{"check_id":"no-phantom-cli","status":"pass"},{"check_id":"skill-size","status":"pass"}]}},"finding_ids":[],"blocking_ids":[],"synthesis":{"conflicts":[],"rubber_stamps":[]},"summary":"..."}
```

**Priority:** P0/P1 block; P2/P3 track only.

---

## Beads

`bd create --title "{AREA}: {desc}" --priority {0-3} --labels "review-finding,F{XX},round-{N},{role}" --type bug --silent`

Replace `F{NNN}` with feature ID (e.g. F067), `round-{N}` with iteration (e.g. round-1), `{role}` with qa/security/devops/sre/techlead/docs.

After creating findings, include in subagent output: `FINDINGS_CREATED: id1 id2 id3`

---

## Few-Shot Examples

**Good P0 finding (Security):**
```
bd create --title "Security: auth bypass via missing role check in API handler" --priority 0 --labels "review-finding,F001,round-1,security" --type bug --silent
```
Reason: Exploitable in production; attacker can bypass auth. Include file:line.

**Good P2 finding (style):**
```
bd create --title "Docs: typo in README deployment section" --priority 2 --labels "review-finding,F001,round-1,docs" --type bug --silent
```
Reason: Maintenance debt, not blocking.

**Bad — vague finding (no file:line):**
```
bd create --title "Security: possible vulnerability" --priority 0 ...
```
Reason: P0 requires evidence. Add file:line or downgrade to P2.

**Bad — missing FINDINGS_CREATED:**
```
SCOPE: internal/auth/*.go. RISK MAP: token validation. EVIDENCE: All checks present. VERDICT: PASS
```
Reason: Must output `FINDINGS_CREATED: (none)` explicitly; never leave blank.

**Good — no findings (explicit output):**
```
SCOPE: internal/auth/*.go (3 files). RISK MAP: token validation, rate limit. EVIDENCE: All checks present. VERDICT: PASS
FINDINGS_CREATED: (none)
PASS
```

---

## See Also

- `@oneshot` — review-fix loop
- `@deploy` — requires APPROVED verdict
