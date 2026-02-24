---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation + PromptOps)
cli: sdp quality all
version: 13.0.0
changes:
  - "13.0.0: Add Subagent 7 — Prompt Engineer for skill/agent/command file review"
  - P2/P3 PASS rule: Output PASS if all findings are P2 or P3
  - Beads: --silent --labels "review-finding,F{NNN},round-N,{role}"
  - Verdict: finding_ids, blocking_ids for @oneshot integration
  - Documentation Expert: AC coverage check via ws-verdicts
---

# review

> **CLI:** `sdp quality all` (quality checks only)
> **LLM:** Spawn 7 specialist subagents for full review

Comprehensive multi-agent quality review.

---

## EXECUTE THIS NOW

When user invokes `@review F067`, you MUST:

1. First run CLI quality checks:
```bash
sdp quality all
```

2. Then spawn 7 specialist subagents IN PARALLEL for review:
   - QA expert
   - Security expert
   - DevOps expert
   - SRE expert
   - TechLead expert
   - Documentation expert
   - Prompt Engineer expert

**DO NOT skip step 2.** The CLI only runs basic checks. Full review requires spawning subagents.

---

## How to Spawn Subagents

Use your tool's subagent capability. For example:
- Claude Code: Use Task tool with `subagent_type="general-purpose"`
- Cursor: Use agent panel
- Windsurf: Use agent spawning

Each subagent should read its specification from `.claude/agents/{role}.md`:

---

## Subagent 1: QA Expert

**Role file:** `.claude/agents/qa.md`

**Task:**
```
You are the QA expert for feature F067.

Your task:
1. Review test coverage (target: 80%+)
2. Check test quality
3. Verify quality metrics
4. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,qa" --priority={0-3} --type=bug`
5. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 2: Security Expert

**Role file:** `.claude/agents/security.md`

**Task:**
```
You are the SECURITY expert for feature F067.

Your task:
1. Review security controls
2. Check for vulnerabilities (OWASP Top 10)
3. Verify compliance
4. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,security" --priority={0-3} --type=bug`
5. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 3: DevOps Expert

**Role file:** `.claude/agents/devops.md`

**Task:**
```
You are the DEVOPS expert for feature F067.

Your task:
1. Review CI/CD pipeline
2. Check infrastructure
3. Verify deployment strategy
4. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,devops" --priority={0-3} --type=bug`
5. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 4: SRE Expert

**Role file:** `.claude/agents/sre.md`

**Task:**
```
You are the SRE expert for feature F067.

Your task:
1. Review SLOs/SLIs
2. Check monitoring
3. Verify incident response procedures
4. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,sre" --priority={0-3} --type=bug`
5. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 5: TechLead Expert

**Role file:** `.claude/agents/tech-lead.md`

**Task:**
```
You are the TECH LEAD expert for feature F067.

Your task:
1. Review code quality (SOLID, clean code)
2. Check architecture decisions
3. Verify LOC compliance (max 200 per file)
4. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,techlead" --priority={0-3} --type=bug`
5. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 6: Documentation Expert

**Role file:** None (inline task)

**Task:**
```
You are the DOCUMENTATION expert for feature F067.

Your task:
1. Check drift: Vision → Specs → Code
2. Run `sdp drift detect {ws-id}` for each workstream in the feature
3. Verify AC coverage: for each ws-id, check jq '.ac_evidence | length' .sdp/ws-verdicts/{ws-id}.json matches AC count in WS file. If gap, create P1 finding.
4. Verify all AC documented
5. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,docs" --priority={0-3} --type=bug`
6. Include in your output: FINDINGS_CREATED: {space-separated ids}

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## Subagent 7: Prompt Engineer Expert

**Role file:** None (inline task)

**Task:**
```
You are the PROMPT ENGINEER expert for feature F067.

Review scope: all files in sdp/prompts/skills/, sdp/prompts/agents/, sdp/prompts/commands/ that were changed or are referenced by the feature's workstreams.

Your task:
1. Language-agnostic compliance: skills must NOT hardcode language-specific commands (go test, npm test, cargo test, etc.) in CRITICAL execution paths. Skills reference "quality gates (see AGENTS.md)" instead. AGENTS.md is the project-specific toolchain config.
2. No handoff lists: skills must NOT end with "Next steps:", "Follow-up:", or delegation lists. The agent completes its work or reports failure — no handoff.
3. No phantom CLI: every CLI command referenced in a skill must actually exist. Check with `which {cmd}` or verify in cmd/ directory. Flag references to commands that don't exist.
4. Symlink integrity: verify .cursor/commands, .opencode/commands, .claude/skills, .claude/agents are symlinks to sdp/ canonical source (not duplicated files).
5. Skill size: skills should be 50-150 lines. Flag skills over 200 lines as bloated.
6. Agent liveness: every agent in sdp/prompts/agents/ must be referenced by at least one skill or command. Flag unreferenced agents as dead code.
7. Command consistency: every command in sdp/prompts/commands/ should have correct YAML frontmatter (description, agent fields).
8. For each finding: create beads issue with `bd create --silent --labels "review-finding,F067,round-1,promptops" --priority={0-3} --type=bug`
9. Include in your output: FINDINGS_CREATED: {space-separated ids}

Priority guide:
- P0: phantom CLI command (agent will fail at runtime)
- P1: hardcoded language-specific command in CRITICAL path, handoff list in skill output
- P2: bloated skill (>200 lines), missing frontmatter field
- P3: style issue, minor inconsistency

Rule: Output PASS if ALL your findings are P2 or P3 priority. Output FAIL only if you have P0 or P1 findings.

Output verdict: PASS or FAIL
```

---

## After All Subagents Complete

**Step 1: Synthesize verdict:**

```
## Feature Review: F067

### QA: {PASS/FAIL} - {summary}
### Security: {PASS/FAIL} - {summary}
### DevOps: {PASS/FAIL} - {summary}
### SRE: {PASS/FAIL} - {summary}
### TechLead: {PASS/FAIL} - {summary}
### Documentation: {PASS/FAIL} - {summary}
### PromptOps: {PASS/FAIL} - {summary}

## Overall Verdict

**APPROVED** if all 7 PASS
**CHANGES_REQUESTED** if any FAIL
```

**Step 2: Aggregate finding IDs from all subagents**

Parse `FINDINGS_CREATED: id1 id2 ...` from each subagent output. Collect all IDs into `finding_ids`. Filter P0/P1 into `blocking_ids` (query beads for priority, or infer from FAIL reviewers).

**Step 3: Save verdict to file (CRITICAL):**

After synthesizing, write the verdict to `.sdp/review_verdict.json`:

```bash
cat > .sdp/review_verdict.json << EOF
{
  "feature": "F067",
  "verdict": "APPROVED" or "CHANGES_REQUESTED",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "round": 1,
  "reviewers": {
    "qa": "PASS/FAIL",
    "security": "PASS/FAIL",
    "devops": "PASS/FAIL",
    "sre": "PASS/FAIL",
    "techlead": "PASS/FAIL",
    "docs": "PASS/FAIL",
    "promptops": "PASS/FAIL"
  },
  "finding_ids": ["sdp_dev-abc", "sdp_dev-xyz"],
  "blocking_ids": ["sdp_dev-abc"],
  "summary": "Brief summary of review findings"
}
EOF
```

This file is required for @deploy and @oneshot. If missing or verdict is not APPROVED, @deploy must block. @oneshot uses `blocking_ids` for the review-fix loop.

---

## Finding Priority

| Priority | Action | Blocks? |
|----------|--------|---------|
| P0 | Fix immediately | YES |
| P1 | Create bugfix | YES |
| P2 | Track only | NO |
| P3 | Track only | NO |

---

## Beads Integration

For each finding, create issue and capture ID:

```bash
FINDING_ID=$(bd create \
  --title "{AREA}: {description}" \
  --priority {0-3} \
  --labels "review-finding,F{NNN},round-{N},{role}" \
  --type bug \
  --silent)
echo "FINDING:$FINDING_ID"
```

Replace `F{NNN}` with feature ID (e.g. F067), `round-{N}` with iteration (e.g. round-1), `{role}` with qa/security/devops/sre/techlead/docs/promptops.

After creating findings, include in subagent output: `FINDINGS_CREATED: id1 id2 id3`

---

## See Also

- `@oneshot` - Execution with review-fix loop
- `.claude/patterns/quality-gates.md` - Quality gates
- `.claude/agents/*.md` - Agent specifications
