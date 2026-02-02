---
name: build
description: Execute workstream with TDD and guard enforcement
tools: Read, Write, StrReplace, Shell, Skill
---

# @build - Execute Workstream

Execute a single workstream following TDD discipline with automatic guard.

**Language-Agnostic:** Detects project type (Python, Java, Go) and uses appropriate quality checks.

## Invocation

Accepts **both** formats:

- `@build 00-001-01` — WS-ID (PP-FFF-SS)
- `@build sdp-xxx` — Beads task ID directly (when Beads enabled)

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 0 | Detect project type | Auto-detect Python/Java/Go |
| 1 | Resolve beads_id | ws_id → mapping (if Beads enabled) |
| 2 | Beads IN_PROGRESS | `bd update {beads_id} --status in_progress` |
| 3 | Activate guard | `sdp guard activate {ws_id}` |
| 4 | Read WS spec | AC present and clear |
| 5 | TDD cycle | `@tdd` for each AC |
| 6 | Quality check | Language-specific tests + AI validators |
| 7 | Beads CLOSED/blocked | `bd close` or `bd update --status blocked` |
| 8 | Beads sync + Complete | `bd sync` then commit |

## Workflow

### Step 0: Detect Project Type

**Check for build files in project root:**

```bash
# Python
if [ -f "pyproject.toml" ] || [ -f "setup.py" ]; then
    PROJECT_TYPE="python"
    TEST_CMD="pytest tests/ -v"
    COVERAGE_CMD="pytest --cov=src/ --cov-fail-under=80"
    TYPE_CHECK_CMD="mypy src/ --strict"
    LINT_CMD="ruff check src/"

# Java (Maven)
elif [ -f "pom.xml" ]; then
    PROJECT_TYPE="java"
    TEST_CMD="mvn test"
    COVERAGE_CMD="mvn verify"
    TYPE_CHECK_CMD="javac -Xlint:all"
    LINT_CMD="mvn checkstyle:check"

# Java (Gradle)
elif [ -f "build.gradle" ]; then
    PROJECT_TYPE="java"
    TEST_CMD="gradle test"
    COVERAGE_CMD="gradle test jacocoTestReport"
    TYPE_CHECK_CMD="javac -Xlint:all"
    LINT_CMD="gradle checkstyleMain"

# Go
elif [ -f "go.mod" ]; then
    PROJECT_TYPE="go"
    TEST_CMD="go test ./..."
    COVERAGE_CMD="go test -coverprofile=coverage.out ./..."
    TYPE_CHECK_CMD="go vet ./..."
    LINT_CMD="golint ./..."

# Node.js
elif [ -f "package.json" ]; then
    PROJECT_TYPE="node"
    TEST_CMD="npm test"
    COVERAGE_CMD="npm run test:coverage"
    TYPE_CHECK_CMD=""
    LINT_CMD="npm run lint"

# Unknown/agnostic
else
    PROJECT_TYPE="agnostic"
    # Ask Claude to detect structure
fi
```

**If multiple build files exist, ask user to specify.**

### Step 1: Beads IN_PROGRESS (when Beads enabled)

**When Beads is enabled** (bd installed, `.beads/` exists):

```bash
# Resolve ID: ws_id → beads_id
beads_id=$(grep -m1 "\"sdp_id\": \"{WS-ID}\"" .beads-sdp-mapping.jsonl 2>/dev/null | grep -o '"beads_id": "[^"]*"' | cut -d'"' -f4)

# Update status
bd update "$beads_id" --status in_progress
```

**When Beads NOT enabled:** Skip Beads steps. Use ws_id only.

### Step 2: Activate Guard

```bash
sdp guard activate {WS-ID}
```

**Gate:** Must succeed. If fails, WS not ready.

### Step 3: Read Workstream

```bash
Read("docs/workstreams/backlog/{WS-ID}-*.md")
```

Extract:
- Goal and Acceptance Criteria
- Input/Output files
- Steps to execute

### Step 4: TDD Cycle

For each AC, call internal TDD skill:

```
@tdd "AC1: {description}"
```

Cycle: Red → Green → Refactor

**@tdd will automatically use language-specific test commands based on PROJECT_TYPE.**

### Step 5: Quality Check

#### 5.1: Run Tests

Use detected TEST_CMD from Step 0:

**Python:**
```bash
pytest tests/ -v
```

**Java:**
```bash
mvn test
# OR
gradle test
```

**Go:**
```bash
go test ./...
```

**Gate:** All tests must PASS.

#### 5.2: Check Coverage

Use detected COVERAGE_CMD from Step 0:

**Python:**
```bash
pytest --cov=src/ --cov-fail-under=80
# Expected: Coverage ≥80%
```

**Java:**
```bash
mvn verify
# Expected: JaCoCo report shows ≥80%

# OR
gradle test jacocoTestReport
# Expected: JaCoCo report shows ≥80%
```

**Go:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
# Expected: Total coverage ≥80%
```

**Gate:** Coverage ≥80%

#### 5.3: Type Checking

Use detected TYPE_CHECK_CMD from Step 0:

**Python:**
```bash
mypy src/ --strict
# Expected: 0 errors
```

**Java:**
```bash
javac -Xlint:all
# Expected: 0 warnings
```

**Go:**
```bash
go vet ./...
# Expected: 0 warnings
```

**Gate:** No type errors/warnings

#### 5.4: Linting

Use detected LINT_CMD from Step 0:

**Python:**
```bash
ruff check src/
# Expected: No errors
```

**Java:**
```bash
mvn checkstyle:check
# OR
gradle checkstyleMain
# Expected: No violations
```

**Go:**
```bash
golint ./...
# Expected: No warnings
```

**Gate:** No linting violations

#### 5.5: AI Validators

**Always run AI validators** (language-agnostic):

1. **Coverage Validator**
   ```
   Ask Claude: "Read all test and source files. Calculate tested/total ratio. Is coverage ≥80%? Report missing coverage."
   ```

2. **Architecture Validator**
   ```
   Ask Claude: "Parse imports. Check Clean Architecture layers. Any violations? Report."
   ```

3. **Error Validator**
   ```
   Ask Claude: "Find bare except clauses (Python), empty catch blocks (Java), ignored errors (Go). Report violations."
   ```

4. **Complexity Validator**
   ```
   Ask Claude: "Count lines per file. List files >200 LOC. Report violations."
   ```

**Gate:** All AI validators must PASS.

### Step 6: Beads CLOSED or blocked

**On success (all gates pass):**
```bash
bd close "$beads_id" --reason "WS completed" --suggest-next
```

**On failure (any gate fails):**
```bash
bd update "$beads_id" --status blocked
```

### Step 7: Complete

```bash
# When Beads enabled: sync before commit
[ -d .beads ] && bd sync

# Complete guard
sdp guard complete {WS-ID}

# Commit changes
git add .
git commit -m "feat({scope}): {WS-ID} - {title}"
```

## Quality Gates Summary

| Gate | Python | Java | Go | AI Validator |
|------|--------|------|-----|-------------|
| Tests | pytest | mvn test | go test | ✓ |
| Coverage | pytest-cov ≥80% | JaCoCo ≥80% | go tool cover ≥80% | ✓ |
| Type Check | mypy --strict | javac -Xlint:all | go vet | ✓ |
| Linting | ruff | checkstyle | golint | - |
| Architecture | - | - | - | ✓ |
| Error Handling | - | - | - | ✓ |
| File Size | - | - | - | ✓ |

**All gates must PASS for workstream completion.**

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| No active WS | Guard not activated | `sdp guard activate` |
| File not in scope | Editing wrong file | Check WS scope |
| Project type unknown | No build file found | Specify `--project-type=python/java/go` |
| Coverage <80% | Missing tests | Add tests |
| Type errors | Missing type hints | Add types |
| Lint errors | Code style issues | Fix style |

## Examples

### Python Project

```bash
@build 00-001-01

# Output:
# ✓ Project type detected: Python (pyproject.toml)
# ✓ Running tests: pytest tests/ -v
# ✓ Coverage: 85% (≥80% required)
# ✓ Type checking: mypy src/ --strict
# ✓ Linting: ruff check src/
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!
```

### Java Project

```bash
@build 00-001-01

# Output:
# ✓ Project type detected: Java (pom.xml)
# ✓ Running tests: mvn test
# ✓ Coverage: 87% (≥80% required)
# ✓ Type checking: javac -Xlint:all
# ✓ Linting: mvn checkstyle:check
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!
```

### Go Project

```bash
@build 00-001-01

# Output:
# ✓ Project type detected: Go (go.mod)
# ✓ Running tests: go test ./...
# ✓ Coverage: 91% (≥80% required)
# ✓ Type checking: go vet ./...
# ✓ Linting: golint ./...
# ✓ AI validators: PASS
#
# Workstream 00-001-01 complete!
```

## See Also

- [Quality Gates Reference](../../docs/quality-gates.md) - Language-specific quality criteria
- [Python Quick Start](../../docs/examples/python/QUICKSTART.md)
- [Java Quick Start](../../docs/examples/java/QUICKSTART.md)
- [Go Quick Start](../../docs/examples/go/QUICKSTART.md)
- [TDD Skill](../tdd/SKILL.md)
- [Guard Skill](../guard/SKILL.md)
