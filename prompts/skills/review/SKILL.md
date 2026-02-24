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

```
You are the {ROLE} expert for feature F{XX}. Review your domain. For each finding: bd create --silent --labels "review-finding,F{XX},round-1,{role}" --priority={0-3} --type=bug. Output: FINDINGS_CREATED: id1 id2. Rule: PASS if all P2/P3; FAIL if any P0/P1. Output verdict: PASS or FAIL
```

**Role files:** `.claude/agents/qa.md`, `security.md`, `devops.md`, `sre.md`, `tech-lead.md`. Docs and PromptOps: inline (see below).

**Docs expert:** Check drift (`sdp drift detect`), AC coverage (jq `.ac_evidence|length` vs WS file). Labels: `review-finding,F{XX},round-1,docs`

**PromptOps expert:** Review sdp/prompts/skills, agents, commands. Check: language-agnostic, no phantom CLI, no handoff lists, skill size ≤200 LOC. Labels: `review-finding,F{XX},round-1,promptops`

---

## After All Complete

**Synthesize:** `## Feature Review: F{XX}` with `### QA: PASS/FAIL`, etc. **APPROVED** if all 7 PASS; **CHANGES_REQUESTED** if any FAIL.

**Save verdict** to `.sdp/review_verdict.json` (required for @deploy, @oneshot):

```json
{"feature":"F{XX}","verdict":"APPROVED|CHANGES_REQUESTED","timestamp":"...","round":1,"reviewers":{...},"finding_ids":[...],"blocking_ids":[...],"summary":"..."}
```

**Priority:** P0/P1 block; P2/P3 track only.

---

## Beads

`bd create --title "{AREA}: {desc}" --priority {0-3} --labels "review-finding,F{XX},round-{N},{role}" --type bug --silent`

Replace `F{NNN}` with feature ID (e.g. F067), `round-{N}` with iteration (e.g. round-1), `{role}` with qa/security/devops/sre/techlead/docs.

After creating findings, include in subagent output: `FINDINGS_CREATED: id1 id2 id3`

---

## See Also

- `@oneshot` — review-fix loop
- `@deploy` — requires APPROVED verdict
