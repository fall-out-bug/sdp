---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation)
cli: sdp quality review
version: 11.0.0
---

# review

> **CLI:** `sdp quality review <feature-id>`

Comprehensive multi-agent quality review.

---

## EXECUTE THIS NOW

When user invokes `@review F067`, run:

```bash
sdp quality review F067
```

This is not an example. Execute this command.

---

## Quick Start

```bash
sdp quality review F067          # Review feature
sdp quality review F067 --quiet  # Exit status only
sdp quality review sdp-xxx       # Review by beads ID
```

---

## Review Areas

| Agent | Checks | Blocks? |
|-------|--------|---------|
| QA | Test coverage, quality metrics | YES (if <80%) |
| Security | Vulnerabilities, OWASP | YES (P0/P1) |
| DevOps | CI/CD, deployment | YES (P0) |
| SRE | SLOs, monitoring | YES (P1) |
| TechLead | Code quality, architecture | YES (P1) |
| Documentation | Drift, missing docs | NO (P2) |

---

## Verdict

**APPROVED:** All 6 areas pass
**CHANGES_REQUESTED:** Any area fails

No middle ground.

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

When Beads is enabled, agents create issues automatically:

```bash
# For each finding
bd create --title="{AREA}: {description}" --priority={0-3}
sdp guard finding add --feature={FID} --area={AREA} --priority={PRI}
```

---

## Output

**Success:**
```
✅ APPROVED
📊 QA: PASS (82% coverage)
🔒 Security: PASS
⚙️ DevOps: PASS
⏱️ SRE: PASS
👨‍💻 TechLead: PASS
📚 Documentation: PASS
```

**Failure:**
```
❌ CHANGES_REQUESTED
📊 QA: FAIL (65% coverage)
🔒 Security: PASS
⚙️ DevOps: FAIL (no rollback)
...

Findings tracked: 3 issues
```

---

## Contract Validation

If contract exists (`.contracts/{feature}.yaml`):

```bash
sdp contract validate --contracts .contracts/{feature}.yaml
```

Blocks review if validation fails.

---

## See Also

- `@oneshot` - Execution with review-fix loop
- `.claude/patterns/quality-gates.md` - Quality gates
- `@bugfix` - Fix review findings

**Implementation:** `sdp-plugin/cmd/sdp/quality.go`
