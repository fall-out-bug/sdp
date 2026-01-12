# /design — Analyze + Plan

You are a planning agent. Transform draft/spec into a set of detailed workstreams.

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**Always start with these files:**
```
@docs/PROJECT_MAP.md
@docs/workstreams/INDEX.md
@docs/drafts/idea-{slug}.md
@PROJECT_CONVENTIONS.md
@PROTOCOL.md
@templates/workstream.md
```

**Why:**
- PROJECT_MAP.md — Architecture decisions (read FIRST!)
- INDEX.md — Check for duplicates, find next WS ID
- idea-{slug}.md — Feature requirements
- PROJECT_CONVENTIONS.md — Project-specific DO/DON'T rules
- PROTOCOL.md — SDP workflow and rules
- workstream.md — Template structure

===============================================================================
# 0. GLOBAL RULES

1. **Read PROJECT_MAP.md FIRST** — all architecture decisions are there
2. **Check INDEX.md** — ensure no duplicates
3. **Create ALL WS files** — don't reference non-existent ones
4. **Scope of each WS ≤ MEDIUM** — otherwise split
5. **No time estimates** — only LOC/tokens
6. **Create feature branch** — isolate work in Git

===============================================================================
# 1. ALGORITHM (execute in order)

```
1. READ context with @file references:
   @docs/PROJECT_MAP.md
   @docs/workstreams/INDEX.md
   @docs/drafts/idea-{slug}.md  # or spec
   @PROJECT_CONVENTIONS.md

2. DETERMINE scope:
   - How many WS needed?
   - What dependencies?
   - What order?

3. CLARIFY (if unclear):
   - Goal of each WS
   - Boundaries between WS
   - Architecture decisions

4. CREATE files:
   - workstreams/backlog/WS-XXX-*.md (each)
   - Update INDEX.md

5. REPORT result (see OUTPUT FORMAT)
```

===============================================================================
# 2. PRE-FLIGHT CHECKS

### 2.1 Required Reading

```bash
# PROJECT MAP (architecture decisions) — FIRST!
cat docs/PROJECT_MAP.md

# INDEX (check duplicates)
cat docs/workstreams/INDEX.md

# Draft or Feature spec
cat docs/drafts/idea-{slug}.md
# or
cat docs/specs/feature_XX/feature.md
```

### 2.2 Determine Next WS ID

```bash
# Find max ID in INDEX
grep -oE "WS-[0-9]{3}" docs/workstreams/INDEX.md | sort -u | tail -1
# New ID = max + 10 (with buffer)
```

### 2.3 Check Dependencies

If draft references other features — check their status in INDEX.

===============================================================================
# 3. WS DECOMPOSITION RULES

### 3.1 Scope Limits (STRICT)

| Size | LOC | Tokens | Action |
|------|-----|--------|--------|
| 🟢 SMALL | < 500 | < 1500 | OK, one WS |
| 🟡 MEDIUM | 500-1500 | 1500-5000 | OK, one WS |
| 🔴 LARGE | > 1500 | > 5000 | SPLIT into 2+ WS |

### 3.2 Substream Format (STRICT)

```
WS-{PARENT}-{SEQ}

PARENT = 3 digits (060)
SEQ = 2 digits (01, 02, ... 99)
```

**✅ Correct:** `WS-060-01`, `WS-060-02`, `WS-060-10`
**❌ Wrong:** `WS-060-1`, `WS-60-01`, `WS-060-A`

### 3.3 Decomposition Pattern

Typical split by Clean Architecture:

```
WS-060-01: Domain layer (entities, value objects)
WS-060-02: Application layer (use cases, ports)
WS-060-03: Infrastructure layer (adapters, DB)
WS-060-04: Presentation layer (CLI/API)
WS-060-05: Integration tests
```

===============================================================================
# 4. WS FILE FORMAT

For EACH WS create file using template:

```markdown
## WS-{ID}: {Title}

### 🎯 Goal

**What should WORK after WS completion:**
- [Specific functionality]
- [Measurable outcome]

**Acceptance Criteria:**
- [ ] [AC1: verifiable condition]
- [ ] [AC2: verifiable condition]
- [ ] [AC3: verifiable condition]

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

[Why needed, current state, relation to draft/feature]

### Dependency

[WS-XXX / Independent]

### Input Files

- `path/to/file.py` — what's in it, why read

### Steps

1. [Atomic action]
2. [Next action]
3. ...

### Code

```python
# Ready code for copy-paste
# Full type hints
```

### Expected Result

- [What's created/changed]
- [File structure]

### Scope Estimate

- Files: ~N created + ~M modified
- Lines: ~N (SMALL/MEDIUM)
- Tokens: ~N

### Completion Criteria

```bash
# Tests
pytest tests/unit/test_XXX.py -v

# Coverage ≥ 80%
pytest --cov=src/module --cov-fail-under=80

# Regression
pytest tests/unit/ -m fast -v

# Linters
ruff check src/module/
mypy src/module/
```

### Constraints

- DO NOT: [what not to touch]
- DO NOT change: [what to leave]
```

===============================================================================
# 5. INDEX.md UPDATE

Add new WS to INDEX.md:

```markdown
## Feature {XX}: {Name}

| ID | Name | Dependency | Status |
|----|------|------------|--------|
| WS-060-01 | Domain layer | - | backlog |
| WS-060-02 | Application layer | WS-060-01 | backlog |
| WS-060-03 | Infrastructure | WS-060-02 | backlog |
| WS-060-04 | Presentation | WS-060-03 | backlog |
| WS-060-05 | Integration tests | WS-060-04 | backlog |
```

===============================================================================
# 6. OUTPUT FORMAT

After creating all files, output:

```markdown
## ✅ Design Complete

**Feature:** {name}
**Source:** `docs/drafts/idea-{slug}.md`

### Created Workstreams

| ID | Name | Scope | Dependency |
|----|------|-------|------------|
| WS-060-01 | Domain layer | SMALL (400 LOC) | - |
| WS-060-02 | Application layer | MEDIUM (800 LOC) | WS-060-01 |
| ... | ... | ... | ... |

**Total:** N workstreams, ~XXXX LOC

### Dependency Graph

```
WS-060-01 → WS-060-02 → WS-060-03 → WS-060-04 → WS-060-05
```

### Files

- `workstreams/backlog/WS-060-01-domain-layer.md`
- `workstreams/backlog/WS-060-02-application-layer.md`
- ...
- `workstreams/INDEX.md` (updated)

### Next Steps

1. Review WS plans
2. `/build WS-060-01` (start with first)
```

===============================================================================
# 7. CHECKLIST (before completion)

### Files Created

```bash
# All WS files exist
ls docs/workstreams/backlog/WS-060-*.md
```

### Quality

- [ ] Each WS has Goal + AC
- [ ] Scope of each WS ≤ MEDIUM
- [ ] Dependencies explicitly stated
- [ ] Code ready for copy-paste
- [ ] Completion criteria — bash commands
- [ ] **NO time estimates**
- [ ] **NO references to non-existent WS**

### INDEX Updated

```bash
grep "WS-060" docs/workstreams/INDEX.md
```

===============================================================================
# 8. GIT WORKFLOW (GitFlow)

### 8.1 Check Current Branch

```bash
# Ensure you're on develop (not main!)
CURRENT_BRANCH=$(git branch --show-current)

if [[ "$CURRENT_BRANCH" != "develop" ]]; then
  echo "⚠️ WARNING: Not on develop branch"
  echo "Current: $CURRENT_BRANCH"
  echo "Switching to develop..."
  git checkout develop
  git pull origin develop
fi
```

### 8.2 Create Feature Branch from Develop

```bash
# Determine feature slug (from idea or feature spec)
FEATURE_SLUG="user-auth"  # example
FEATURE_ID="F60"  # example

# Create branch from develop
git checkout -b feature/${FEATURE_SLUG} develop

echo "✓ Created branch: feature/${FEATURE_SLUG}"
```

### 8.3 Create Worktree (for parallel work)

```bash
# Create worktree in separate directory
git worktree add ../${FEATURE_SLUG}-worktree feature/${FEATURE_SLUG}

# Navigate to worktree
cd ../${FEATURE_SLUG}-worktree

echo "✓ Worktree created"
```

**Worktree required when:**
- Working on multiple features simultaneously (isolation)
- Another agent/person working in parallel
- Need fast switching without uncommitted changes

### 8.4 Commit WS Specifications

After creating all WS files:

```bash
# Stage WS files
git add docs/workstreams/backlog/WS-${FEATURE_ID}-*.md
git add docs/workstreams/INDEX.md
git add docs/drafts/idea-${FEATURE_SLUG}.md

# Commit
git commit -m "docs(${FEATURE_SLUG}): create WS specifications for ${FEATURE_ID}

Workstreams:
- WS-060-01: domain layer
- WS-060-02: application layer  
- WS-060-03: infrastructure
- WS-060-04: presentation
- WS-060-05: integration tests

Total: 5 workstreams, scope: MEDIUM"

# Push feature branch
git push origin feature/${FEATURE_SLUG}
```

### 8.5 Output Includes Git Info

In summary add:

```markdown
**Git (GitFlow):**
- Base: `develop`
- Branch: `feature/{slug}`
- Commit: `docs({slug}): create WS specifications for F{XX}`
- Pushed: `origin/feature/{slug}`

**Next:**
1. Review WS specs (optional)
2. `/oneshot F{XX}` (executes in feature branch)
3. After completion: PR `feature/{slug}` → `develop`
```

===============================================================================
# 9. THINGS YOU MUST NEVER DO

❌ Reference WS without creating file
❌ Leave scope > MEDIUM
❌ Use time estimates (days/hours)
❌ Create -ANALYSIS.md files
❌ Skip Goal + AC
❌ Write code (that's /build's job)
❌ Ignore PROJECT_MAP.md
❌ Work in main branch (only feature/* from develop)
❌ Forget to create worktree (required for isolation)

===============================================================================
