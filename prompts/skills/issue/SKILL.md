---
name: issue
description: Analyze bugs, classify severity (P0-P3), route to appropriate fix command (@hotfix, @bugfix, or backlog).
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @issue - Analyze & Route Issues

Systematic bug analysis with severity classification and routing.

---

## EXECUTE THIS NOW

When user invokes `@issue "description"`:

### Step 1: Systematic Debugging

**Phase 1: Symptom Documentation**
- Record exact error messages
- Note reproduction steps
- Document environment

**Phase 2: Hypothesis Formation**
- List all possible causes
- Rank by likelihood
- Select top theory to test

**Phase 3: Systematic Elimination**
- Test hypotheses one at a time
- Record results objectively

**Phase 4: Root Cause Isolation**
- Confirm root cause
- Document findings

### Step 2: Severity Classification

| Severity | Keyword Signals | Route |
|----------|----------------|-------|
| **P0** | "production down", "crash", "blocked", "security" | @hotfix |
| **P1** | "doesn't work", "failing", "error", "broken" | @bugfix |
| **P2** | "edge case", "sometimes", "inconsistently" | backlog |
| **P3** | "cosmetic", "typo", "minor" | defer |

### Step 3: Create Issue & Route

```bash
# Create issue
bd create --title="Bug: {description}" --type=bug --priority={0-3}

# Route to appropriate fix
@hotfix {issue}  # P0 - emergency
@bugfix {issue}  # P1 - quality fix
# P2/P3 - schedule as workstream
```

---

## Auto-Classification Rules

- **P0 (CRITICAL)**: Production down -> @hotfix
- **P1 (HIGH)**: Feature broken -> @bugfix
- **P2 (MEDIUM)**: Edge case -> New WS
- **P3 (LOW)**: Cosmetic -> Defer

---

## Output

- Issue file: `docs/issues/{ID}-{slug}.md`
- GitHub issue (if gh available)
- Routing recommendation

---

## Quick Reference

| Input | Output | Next |
|-------|--------|------|
| Bug description | Issue file + Routing | @hotfix or @bugfix or schedule WS |

---

## See Also

- `@debug` - Systematic debugging workflow
- `@hotfix` - Emergency P0 fixes
- `@bugfix` - Quality P1/P2 fixes
