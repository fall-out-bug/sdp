---
name: review
description: Quality review with traceability check
tools: Read, Shell, Grep
---

# @review - Quality Review

Review feature by validating workstreams against quality gates and traceability.

## Invocation

```bash
@review F01       # Feature ID (markdown workflow)
@review sdp-xxx   # Beads task ID
```

## Workflow Summary

| Step | Action | Gate |
|------|--------|------|
| 1 | List workstreams | All WS found |
| 2 | Check traceability | All ACs have tests |
| 3 | Run quality gates | All checks pass |
| 4 | Verify goals | All ACs achieved |
| 5 | Verdict | APPROVED or CHANGES_REQUESTED |
| 6 | Post-review (if needed) | Track all findings |

## Step 1-2: List & Check Traceability

```bash
# List workstreams
bd list --parent {feature-id}  # Beads
ls docs/workstreams/completed/{feature-id}-*.md  # Markdown

# Check traceability
sdp trace check {WS-ID}
```

**Gate:** 100% AC coverage (all ACs have mapped tests).

## Step 3: Quality Gates

```bash
pytest tests/ -v                    # All tests pass
pytest --cov=src --cov-fail-under=80  # Coverage ≥80%
mypy src/ --strict                  # Type checking
ruff check src/                     # Linting
grep -r "except:" src/ | grep "pass"  # No except:pass
```

## Step 4: Goal Achievement

For each WS verify:
- [ ] All ACs have passing tests
- [ ] Implementation matches description
- [ ] No TODO/FIXME in code

## Step 5: Verdict

**APPROVED** — All gates pass, all ACs traceable  
**CHANGES_REQUESTED** — Any failure

No middle ground. No "approved with notes."

## Step 6: Post-Review (when CHANGES_REQUESTED)

**⚠️ MANDATORY when verdict is CHANGES_REQUESTED**

| Finding type | Action | Output |
|--------------|--------|--------|
| **Bugs** | @issue | `docs/issues/` → /bugfix |
| **Planned work** | Add WS to **same feature** | `docs/workstreams/backlog/` |
| **Tech debt** | @issue for triage | Backlog |

**Rules:**
- Never create new feature for review follow-up
- Every finding must have Issue or WS link
- "Deferred" without tracking = protocol violation

### Completion Checklist

```markdown
- [ ] Verdict recorded
- [ ] Report saved to docs/reports/
- [ ] All bugs → Issue created
- [ ] All planned work → WS created
- [ ] No "deferred" without tracking
```

## Errors

| Error | Fix |
|-------|-----|
| Missing trace | Add test for AC |
| Coverage <80% | Add more tests |
| Goal not met | Fix implementation |

## See Also

- [@issue skill](../issue/SKILL.md)
- [Traceability Guide](../../docs/reference/traceability.md)
