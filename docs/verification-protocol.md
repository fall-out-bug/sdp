# Verification Protocol

**Purpose:** Enforce evidence-based claims in workstream completion. Prevent false positives from uncertain language.

**Iron Law:** No completion without fresh verification.

---

## Overview

The verification protocol ensures that workstream completion claims are backed by actual command output evidence, not assumptions or uncertain language.

### Red Flag Phrases

The following phrases indicate uncertainty and are **forbidden** in Execution Reports:

- "should" — implies assumption, not fact
- "probably" — indicates uncertainty
- "seems to" / "seems like" — suggests guesswork
- "might" — uncertain possibility
- "may" — uncertain possibility
- "appears to" — suggests assumption

### Evidence Requirement

All completion claims must be backed by **command output** showing actual results, not just commands.

**Required format:**
```bash
$ pytest tests/unit/test_service.py -v
===== 15 passed in 0.5s =====

$ pytest --cov=hw_checker/module --cov-fail-under=80
===== Coverage: 85% =====
```

**Not acceptable:**
- Just commands without output
- Text descriptions without evidence
- Empty code blocks
- Uncertain claims without proof

---

## Hook: `verification-completion.sh`

### Usage

```bash
sdp/hooks/verification-completion.sh <ws-file-path>
```

### What It Checks

1. **Execution Report exists** — WS file must have "### Execution Report" section
2. **No red flag phrases** — Scans for uncertain language
3. **Command output evidence** — Requires bash/shell code blocks with actual output

### Integration with `/build`

The hook is automatically called by `post-build.sh` after Execution Report is appended:

```bash
# After appending Execution Report to WS file
sdp/hooks/verification-completion.sh tools/hw_checker/docs/workstreams/backlog/WS-XXX-YY.md
```

### Exit Codes

- `0` — Verification passed
- `1` — Verification failed (red flags or missing evidence)

---

## Examples

### ✅ Correct Execution Report

```markdown
### Execution Report

**Executed by:** Auto
**Date:** 2026-01-15

#### 🎯 Goal Status

- [x] AC1: Hook created — ✅
- [x] AC2: Red flag detection — ✅

**Goal Achieved:** ✅ YES

#### Self-Check Results

```bash
$ pytest sdp/tests/unit/hooks/test_verification_completion.py -v
===== 12 passed in 0.8s =====

$ ruff check sdp/hooks/verification-completion.sh
All checks passed!

$ grep -rn "TODO\|FIXME" sdp/hooks/
(empty - OK)
```
```

### ❌ Incorrect Execution Report (Red Flag)

```markdown
### Execution Report

**Goal Achieved:** ✅ YES

#### Self-Check Results

Tests should pass. Everything probably works.
```

**Error:** Red flag phrases detected: should, probably

### ❌ Incorrect Execution Report (No Evidence)

```markdown
### Execution Report

**Goal Achieved:** ✅ YES

#### Self-Check Results

All tests passed. Coverage is good.
```

**Error:** No command output evidence found

### ❌ Incorrect Execution Report (Empty Code Block)

```markdown
### Execution Report

**Goal Achieved:** ✅ YES

#### Self-Check Results

```bash
$ pytest tests/unit/test_service.py -v
```
```

**Error:** No command output evidence found (empty code block)

---

## Workflow Integration

### `/build` Command Flow

1. Pre-build checks (`pre-build.sh`)
2. Execute workstream (TDD)
3. Post-build checks (`post-build.sh`)
4. **Append Execution Report to WS file**
5. **Run verification hook** (`verification-completion.sh`)
6. Commit changes

### Manual Verification

If verification fails during automated check, run manually:

```bash
# Find WS file
WS_FILE=$(find tools/hw_checker/docs/workstreams -name "WS-XXX-YY-*.md" | head -1)

# Run verification
sdp/hooks/verification-completion.sh "$WS_FILE"
```

---

## Rationale

### Why Red Flags Matter

Uncertain language in completion claims leads to:
- False positives (claiming success without proof)
- Technical debt (untested code marked as complete)
- Broken builds (assumptions that don't hold)

### Why Evidence Is Required

Command output provides:
- **Proof** — Actual test results, not assumptions
- **Reproducibility** — Others can verify the same results
- **Transparency** — Clear indication of what was tested

### Iron Law

> **No completion without fresh verification**

This means:
- Every completion claim must have evidence
- Evidence must be fresh (from actual execution)
- Evidence must be verifiable (command output)

---

## Troubleshooting

### "Red flag phrases detected"

**Solution:** Replace uncertain language with command output:

```markdown
# ❌ Wrong
Tests should pass.

# ✅ Correct
```bash
$ pytest tests/unit/test_service.py -v
===== 15 passed in 0.5s =====
```
```

### "No command output evidence found"

**Solution:** Add code blocks with actual command output:

```markdown
#### Self-Check Results

```bash
$ pytest tests/unit/test_service.py -v
===== 15 passed in 0.5s =====

$ ruff check src/hw_checker/module/
All checks passed!
```
```

### "Execution Report section not found"

**Solution:** Add Execution Report section to WS file:

```markdown
---

### Execution Report

**Executed by:** {agent}
**Date:** {YYYY-MM-DD}

...
```

---

## Related Documents

- `sdp/PROTOCOL.md` — Full SDP specification
- `sdp/prompts/commands/build.md` — Build command instructions
- `sdp/hooks/post-build.sh` — Post-build hook (calls verification)

---

**Version:** 1.0.0  
**Last Updated:** 2026-01-15  
**Status:** Active
