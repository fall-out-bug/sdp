---
name: review
description: Quality review with AI-based validation
tools: Read, Shell, Grep
---

# @review - Quality Review

Review feature by validating workstreams against quality gates using AI-based validators.

**Language-Agnostic:** Uses AI validators that work with Python, Java, Go, and other languages.

## Invocation

```bash
@review F01       # Feature ID (markdown workflow)
@review sdp-xxx   # Beads task ID
```

## Workflow Summary

| Step | Action | Gate |
|------|--------|------|
| 1 | Detect project type | Auto-detect Python/Java/Go |
| 2 | List workstreams | All WS found |
| 3 | Check traceability | All ACs have tests |
| 4 | Run AI validators | All checks pass |
| 5 | Verify goals | All ACs achieved |
| 6 | Verdict | APPROVED or CHANGES_REQUESTED |
| 7 | Post-review (if needed) | Track all findings |

## Step 1: Detect Project Type

```bash
# Python
if [ -f "pyproject.toml" ] || [ -f "setup.py" ]; then
    PROJECT_TYPE="python"

# Java
elif [ -f "pom.xml" ] || [ -f "build.gradle" ]; then
    PROJECT_TYPE="java"

# Go
elif [ -f "go.mod" ]; then
    PROJECT_TYPE="go"

else
    PROJECT_TYPE="agnostic"
fi
```

## Step 2: List & Check Traceability

```bash
# List workstreams
bd list --parent {feature-id}  # Beads
ls docs/workstreams/completed/{feature-id}-*.md  # Markdown

# Check traceability
sdp trace check {WS-ID}
```

**Gate:** 100% AC coverage (all ACs have mapped tests).

## Step 3: Quality Gates (AI Validators)

**Run AI validators** (language-agnostic, works for any language):

### 3.1 Coverage Validator

```
Ask Claude: "Analyze test coverage by:
1. Reading all test files and source files
2. Mapping tests to functions/classes
3. Calculating (tested_functions / total_functions) * 100

Report:
- Total coverage percentage
- Untested functions
- Missing branch coverage
- Verdict: ✅ PASS (≥80%) or ❌ FAIL (<80%)"
```

### 3.2 Architecture Validator

```
Ask Claude: "Validate Clean Architecture by:
1. Parsing imports from all source files
2. Mapping files to layers (domain/, application/, infrastructure/, presentation/)
3. Checking for violations:
   - ❌ domain/ imports infrastructure/
   - ❌ domain/ imports application/
   - ❌ application/ imports presentation/

Report:
- Number of violations
- File paths and line numbers
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

### 3.3 Error Handling Validator

```
Ask Claude: "Find unsafe error handling by searching for:
1. Bare except clauses (Python: `except:`)
2. Empty catch blocks (Java: `catch(Exception e) { }`)
3. Ignored errors (Go: `func(), _`)

For each match, check:
- Is error logged?
- Is error re-raised?
- Is specific exception type?

Report:
- Number of violations
- File paths and line numbers
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

### 3.4 Complexity Validator

```
Ask Claude: "Identify overly complex code by:
1. Counting lines per file
2. Calculating cyclomatic complexity per function
3. Checking nesting depth

Thresholds:
- Complexity <10: ✅ OK
- Complexity 10-20: ⚠️ Warning
- Complexity >20: ❌ FAIL
- Lines >200: ❌ FAIL
- Nesting >4: ❌ FAIL

Report:
- Files exceeding thresholds
- Complexity scores
- Recommendations for refactoring
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

### 3.5 Optional Tool-Based Validation

**If language tools are available**, run them as additional checks:

**Python:**
```bash
pytest tests/ -v
pytest --cov=src --cov-fail-under=80
mypy src/ --strict
ruff check src/
```

**Java:**
```bash
mvn test
mvn verify
mvn checkstyle:check
```

**Go:**
```bash
go test ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
go vet ./...
golint ./...
```

**Note:** AI validators are PRIMARY. Tool-based validation is OPTIONAL supplementary check.

## Step 3.5: Retrieve Architectural Decisions

**Fetch all decisions related to this feature/workstream:**

```bash
# Get all decisions for the feature
sdp decisions list | grep "{feature-id}"

# Or export to markdown for full report
sdp decisions export docs/decisions/{feature-id}-decisions.md
```

**Include in review report:**
1. Count of related decisions
2. Key decisions with dates, types, and rationale
3. Link to full decision export

## Step 4: Goal Achievement

For each WS verify:
- [ ] All ACs have passing tests
- [ ] Implementation matches description
- [ ] No TODO/FIXME in code
- [ ] AI validators pass

## Step 5: Verdict

**APPROVED** — All AI validators pass, all ACs traceable
**CHANGES_REQUESTED** — Any AI validator fails

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
- [ ] Report saved to appropriate documentation location
- [ ] All bugs → Issue created
- [ ] All planned work → WS created
- [ ] No "deferred" without tracking
```

## Output Format

### APPROVED Example

```markdown
## Quality Review Report for F01

**Project Type:** Python

### AI Validators Results

✅ **Coverage:** 85% (≥80% required)
- All functions tested
- Missing: src/service.py:UserService.delete_user() (45-50)

✅ **Architecture:** No violations
- Clean layers maintained
- No circular dependencies

✅ **Error Handling:** No violations
- No bare except clauses
- All errors logged or re-raised

✅ **Complexity:** All files <200 LOC
- Max complexity: 8 (within threshold)
- Max file size: 150 LOC

### Tool-Based Validation (Optional)

✅ pytest tests/ -v: All 42 tests pass
✅ pytest --cov=src/: 85% coverage
✅ mypy src/ --strict: 0 errors
✅ ruff check src/: No issues

### Traceability

✅ 100% AC coverage (42/42 ACs have tests)

### Architectural Decisions

**Related decisions:** 3

1. **[2026-02-06] Use JSONL for decision logging**
   - Type: technical
   - Question: Should we use JSONL for decision logging?
   - Rationale: JSONL is append-only, easy to parse, and supports streaming

2. **[2026-02-05] Go for binary implementation**
   - Type: technical
   - Question: Which language for CLI binary?
   - Rationale: Single binary, cross-compilation, Windows support

3. **[2026-02-04] Quality gates use AI validation**
   - Type: tradeoff
   - Question: Tool-based or AI-based validation?
   - Rationale: Language-agnostic, no tool dependencies

### Verdict

**✅ APPROVED**

All quality gates pass. Feature ready for deployment.
```

### CHANGES_REQUESTED Example

```markdown
## Quality Review Report for F01

**Project Type:** Java

### AI Validators Results

❌ **Coverage:** 72% (<80% required)
- Untested classes:
  - src/main/java/service/UserService.java
  - src/main/java/repository/UserRepository.java

❌ **Architecture:** 2 violations
1. src/main/java/domain/User.java (line 5)
   Imports: com.example.infrastructure.Database
   Domain must not import infrastructure

2. src/main/java/service/AuthService.java (line 12)
   Imports: com.example.presentation.UserController
   Service must not import presentation

❌ **Error Handling:** 3 violations
1. src/main/java/service/Service.java (line 45)
   Empty catch block - add logging or re-throw

✅ **Complexity:** All files <200 LOC

### Architectural Decisions

**Related decisions:** 1

1. **[2026-02-05] Use Spring Boot**
   - Type: technical
   - Question: Framework choice for Java backend?
   - Rationale: Enterprise support, auto-configuration, ecosystem

### Verdict

**❌ CHANGES_REQUESTED**

**Required Actions:**
1. Fix coverage: Add tests for UserService, UserRepository
2. Fix architecture: Remove infrastructure imports from domain
3. Fix error handling: Log exceptions in catch blocks

**Issue Tracking:**
- Create issue for coverage: docs/issues/ISSUE-001.md
- Create WS for architecture fix: 00-001-05
- Fix error handling in same WS

**Follow-up:** Re-run @review after fixes
```

## Errors

| Error | Fix |
|-------|-----|
| Missing trace | Add test for AC |
| Coverage <80% | Add more tests |
| Architecture violations | Fix import structure |
| Error handling violations | Add logging or re-raise |
| Goal not met | Fix implementation |

## See Also

- [@issue skill](../issue/SKILL.md)
- [Traceability Guide](../../docs/reference/traceability.md)
- [Quality Gates Reference](../../docs/quality-gates.md)
- [Python Quick Start](../../docs/examples/python/QUICKSTART.md)
- [Java Quick Start](../../docs/examples/java/QUICKSTART.md)
- [Go Quick Start](../../docs/examples/go/QUICKSTART.md)
