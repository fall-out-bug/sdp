# Completion Protocol

**Version:** 1.0  
**Last Updated:** 2026-01-22

This document describes MANDATORY completion requirements for each SDP command.  
An agent CANNOT consider a command complete without satisfying ALL requirements.

---

## /build — Exit Requirements

### Checklist

| # | Requirement | Blocking | Verification |
|---|-------------|----------|--------------|
| 1 | Tests pass | ✅ YES | `pytest tests/unit/test_*.py -v` |
| 2 | Coverage ≥ 80% | ✅ YES | `pytest --cov --cov-fail-under=80` |
| 3 | No TODO/FIXME | ✅ YES | `grep -rn "TODO\|FIXME" src/` |
| 4 | Execution Report appended | ✅ YES | `grep "Execution Report" WS-*.md` |
| 5 | Git commit with WS-ID | ✅ YES | `git log -1 \| grep "WS-"` |
| 6 | GitHub issue updated | ⚠️ If configured | `gh issue view` |
| 7 | Clean git state | ✅ YES | `git status --porcelain` |

### Quick Verification

```bash
# Run all checks
WS_ID="WS-XXX-YY"
WS_FILE="path/to/WS-*.md"

# 1-3: Post-build hook
sdp/hooks/post-build.sh $WS_ID

# 4: Execution Report
grep -q "Execution Report" "$WS_FILE" && echo "✓ Report" || echo "❌ Report"

# 5: Git commit
git log -1 --oneline | grep -q "$WS_ID" && echo "✓ Commit" || echo "❌ Commit"

# 7: Clean state
test -z "$(git status --porcelain)" && echo "✓ Clean" || echo "❌ Dirty"
```

### Output Format

```markdown
## ✅ Build Complete: {WS-ID}

**Checklist:**
- [x] Tests: X passed
- [x] Coverage: XX%
- [x] Execution Report: appended
- [x] Git commit: {hash}
- [x] Clean state: yes

**Next:** /build {next-WS} or /codereview F{XX}
```

---

## /codereview — Exit Requirements

### Checklist

| # | Requirement | Blocking | Verification |
|---|-------------|----------|--------------|
| 1 | All WS reviewed | ✅ YES | Check each WS file |
| 2 | Review Results in all WS | ✅ YES | `grep "Review Results" WS-*.md` |
| 3 | Verdict determined | ✅ YES | APPROVED or CHANGES REQUESTED |
| 4 | UAT Guide created (if APPROVED) | ✅ YES | `ls docs/uat/F*-uat-guide.md` |
| 5 | Feature Summary output | ✅ YES | Output to user |
| 6 | GitHub issues updated | ⚠️ If configured | `gh issue list` |
| 7 | Blockers listed (if CHANGES REQUESTED) | ✅ YES | Output to user |

### Quick Verification

```bash
FEATURE="F191"
FEATURE_NUM="191"

# 1-2: Review Results
for f in WS-${FEATURE_NUM}*.md; do 
    grep -q "Review Results" "$f" && echo "✓ $f" || echo "❌ $f"
done

# 4: UAT Guide (if APPROVED)
ls tools/hw_checker/docs/uat/${FEATURE}-uat-guide.md

# Full check
sdp/hooks/post-codereview.sh $FEATURE
```

### Output Format (APPROVED)

```markdown
## ✅ Review Complete: {Feature}

**Verdict:** APPROVED

**WS Results:**
| WS | Verdict | Coverage |
|----|---------|----------|
| WS-XXX-01 | ✅ APPROVED | 85% |
| WS-XXX-02 | ✅ APPROVED | 82% |

**UAT Guide:** docs/uat/F{XX}-uat-guide.md

**Next:** Human UAT → /deploy F{XX}
```

### Output Format (CHANGES REQUESTED)

```markdown
## ❌ Review Complete: {Feature}

**Verdict:** CHANGES REQUESTED

**Blockers:**
1. WS-XXX-02: AC3 not achieved
2. WS-XXX-03: Coverage 72% (< 80%)

**Next:** Fix blockers → /build affected WS → /codereview F{XX}
```

---

## /deploy — Exit Requirements

### Pre-merge Checklist

| # | Requirement | Blocking | Verification |
|---|-------------|----------|--------------|
| 1 | All WS APPROVED | ✅ YES | Check Review Results |
| 2 | UAT passed (human sign-off) | ✅ YES | Check UAT Guide |
| 3 | CI/CD green | ✅ YES | Check GitHub Actions |
| 4 | No uncommitted changes | ✅ YES | `git status` |

### Post-merge Checklist

| # | Requirement | Blocking | Verification |
|---|-------------|----------|--------------|
| 5 | Merged to main | ✅ YES | `git branch` |
| 6 | Tag created | ✅ YES | `git tag -l` |
| 7 | Feature branch deleted | ⚠️ Recommended | `git branch -a` |
| 8 | WS moved to completed/ | ⚠️ Recommended | `ls workstreams/completed/` |
| 9 | INDEX.md updated | ✅ YES | Check file |
| 10 | CHANGELOG.md updated | ✅ YES | Check file |
| 11 | Release notes created | ✅ YES | `ls docs/releases/` |

### Quick Verification

```bash
VERSION="v2.1.0"
FEATURE="F191"

# Pre-merge
sdp/hooks/pre-deploy.sh $FEATURE

# Post-merge
git branch --show-current | grep -q "main" && echo "✓ On main"
git tag -l | grep -q "$VERSION" && echo "✓ Tag exists"
ls tools/hw_checker/docs/workstreams/completed/WS-191*.md
```

---

## /oneshot — Exit Requirements

Combines /design + /build (all WS) + /codereview + /deploy

### Checklist

| Phase | Requirements | Verification |
|-------|--------------|--------------|
| Design | WS files created | `ls workstreams/backlog/WS-*.md` |
| Build (each WS) | All /build requirements | `post-build.sh` per WS |
| Codereview | All /codereview requirements | `post-codereview.sh` |
| Deploy | All /deploy requirements | `pre-deploy.sh` |

---

## Anti-Patterns (What NOT to Do)

❌ **Finish without commit:**
```
"Tests pass, done!" → forgot git commit
```

❌ **Skip UAT Guide:**
```
"All WS approved, ready for deploy!" → no UAT Guide created
```

❌ **Leave WS in backlog:**
```
"Feature deployed!" → completed WS still in backlog/
```

❌ **Partial completion:**
```
"95% done, remaining work tracked in WS-XXX-02" → incomplete WS
```

---

## Recovery Procedures

### Forgot Git Commit

```bash
# Check current state
git status
git log -1 --oneline

# Create commit now
git add .
git commit -m "feat({feature}): {WS-ID} - {title}

- {what was done}
- Tests: X passed, coverage XX%"
```

### Forgot Execution Report

```bash
# Append to WS file
cat >> WS-XXX-YY-*.md << 'EOF'

---

## Execution Report

**Executed:** $(date +%Y-%m-%d)
**Agent:** {name}

### What Was Done
- {changes}

### Goal Status
- [x] AC1: ✅
- [x] AC2: ✅

**Goal:** ✅ ACHIEVED
EOF
```

### Forgot UAT Guide

```bash
# Create from template
cp sdp/templates/uat-guide.md tools/hw_checker/docs/uat/F{XX}-uat-guide.md
# Edit with feature-specific content
```

---

## Quick Reference Card

```
/build WS-XXX-YY
  ✓ Tests pass
  ✓ Coverage ≥ 80%
  ✓ Execution Report → WS file
  ✓ Git commit with WS-ID
  ✓ Clean git state

/codereview F{XX}
  ✓ All WS reviewed
  ✓ Review Results → each WS
  ✓ UAT Guide created (if APPROVED)
  ✓ Feature Summary output

/deploy F{XX}
  ✓ UAT passed (human)
  ✓ Merge to main
  ✓ Tag version
  ✓ Update CHANGELOG
  ✓ Release notes
```
