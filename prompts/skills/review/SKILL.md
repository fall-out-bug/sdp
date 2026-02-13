---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation)
cli: sdp quality all
version: 11.1.0
---

# review

> **CLI:** `sdp quality all` (quality checks only)
> **LLM:** Spawn 6 specialist subagents for full review

Comprehensive multi-agent quality review.

---

## EXECUTE THIS NOW

When user invokes `@review F067`, you MUST:

1. First run CLI quality checks:
```bash
sdp quality all
```

2. Then spawn 6 specialist subagents IN PARALLEL for review:
   - QA expert
   - Security expert
   - DevOps expert
   - SRE expert
   - TechLead expert
   - Documentation expert

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
4. For each finding: create beads issue with `bd create`

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
4. For each finding: create beads issue with `bd create --priority=0` for critical

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
4. For each finding: create beads issue with `bd create`

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
4. For each finding: create beads issue with `bd create`

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
4. For each finding: create beads issue with `bd create`

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
2. Run `sdp drift detect` for each workstream
3. Verify all AC documented
4. For each finding: create beads issue with `bd create`

Output verdict: PASS or FAIL
```

---

## After All Subagents Complete

**Synthesize verdict:**

```
## Feature Review: F067

### QA: {PASS/FAIL} - {summary}
### Security: {PASS/FAIL} - {summary}
### DevOps: {PASS/FAIL} - {summary}
### SRE: {PASS/FAIL} - {summary}
### TechLead: {PASS/FAIL} - {summary}
### Documentation: {PASS/FAIL} - {summary}

## Overall Verdict

**APPROVED** if all 6 PASS
**CHANGES_REQUESTED** if any FAIL
```

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

For each finding, create issue:
```bash
bd create --title="{AREA}: {description}" --priority={0-3} --type=task
```

---

## See Also

- `@oneshot` - Execution with review-fix loop
- `.claude/patterns/quality-gates.md` - Quality gates
- `.claude/agents/*.md` - Agent specifications
