# 00-041-06: Cross-Language Validation

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** MEDIUM
> **Created:** 2026-02-02

## Goal

Test plugin on Python, Java, and Go projects to validate language-agnostic functionality.

## Acceptance Criteria

- AC1: Plugin works on Python project (existing SDP repo)
- AC2: Plugin works on Java project (Spring Boot Petclinic)
- AC3: Plugin works on Go project (Gin web framework)
- AC4: All quality gates pass in each language
- AC5: Documentation updated with language examples

## Scope

### Input Files
- `sdp-plugin/prompts/` (from WS-00-041-01, 00-041-03, 00-041-04)
- Existing test projects

### Output Files
- `tests/test-python/` (symlink to existing SDP repo)
- `tests/test-java/` (clone spring-petclinic)
- `tests/test-go/` (clone gin-gonic/gin)
- `sdp-plugin/docs/TUTORIAL.md` (NEW - language-specific workflows)
- `sdp-plugin/docs/examples/*/QUICKSTART.md` (from WS-00-041-02, validated here)

### Out of Scope
- Creating new workstreams (only testing existing)
- Go binary development (WS-00-041-05)

## Implementation Steps

### Step 1: Python Test (Existing SDP)

```bash
# Navigate to existing SDP repo
cd /Users/fall_out_bug/projects/vibe_coding/sdp

# Install plugin (copy prompts to .claude/)
cp -r sdp-plugin/prompts/* .claude/

# Test @build skill with existing workstream
claude "@build 00-001-01"

# Expected Behavior:
# 1. Detects Python project (pyproject.toml found)
# 2. Runs pytest tests/
# 3. AI validators run:
#    - Coverage: Analyzes pytest output
#    - Architecture: Checks src/ imports
#    - Errors: Finds bare except clauses
#    - Complexity: Counts lines per file
# 4. All quality gates PASS

# Verify Outputs:
# - Project type detected: Python
# - Test command: pytest tests/ -v
# - Coverage report: ≥80%
# - No architecture violations
# - No unsafe error handling
# - All files <200 LOC

# Save test results
echo "✓ Python test passed" > tests/test-python/results.txt
```

### Step 2: Java Test (Spring Petclinic)

```bash
# Clone Java project
git clone https://github.com/spring-projects/spring-petclinic.git tests/test-java
cd tests/test-java

# Initialize SDP prompts
cp -r ../../sdp-plugin/prompts/* .claude/

# Create test workstream
cat > docs/workstreams/TEST-001-01.md <<'EOF'
---
ws_id: TEST-001-01
feature: Test
status: backlog
size: SMALL
goal: Add validation to Owner class
AC:
- AC1: Add @NotNull to firstName field
- AC2: Add @Size(min=2, max=30) to lastName field
dependencies: []
scope:
  inputs:
    - src/main/java/org/springframework/samples/petclinic/model/Owner.java
  outputs:
    - src/main/java/org/springframework/samples/petclinic/model/Owner.java
EOF

# Test @build skill
claude "@build TEST-001-01"

# Expected Behavior:
# 1. Detects Java project (pom.xml found)
# 2. Runs mvn test
# 3. AI validators run:
#    - Coverage: Analyzes JaCoCo report
#    - Architecture: Checks layer separation
#    - Errors: Finds empty catch blocks
#    - Complexity: Counts lines per Java file
# 4. All quality gates PASS

# Verify Outputs:
# - Project type detected: Java
# - Test command: mvn test
# - Coverage: JaCoCo report shows ≥80%
# - No architecture violations (model doesn't import infrastructure)
# - No unsafe error handling
# - All files <200 LOC

# Save test results
echo "✓ Java test passed" > tests/test-java/results.txt
```

### Step 3: Go Test (Gin Framework)

```bash
# Clone Go project
git clone https://github.com/gin-gonic/gin.git tests/test-go
cd tests/test-go

# Initialize SDP prompts
cp -r ../../sdp-plugin/prompts/* .claude/

# Create test workstream
cat > docs/workstreams/TEST-002-01.md <<'EOF'
---
ws_id: TEST-002-01
feature: Test
status: backlog
size: SMALL
goal: Add validation middleware
AC:
- AC1: Create validator middleware
- AC2: Add test cases for validator
dependencies: []
scope:
  inputs:
    - examples/
  outputs:
    - middleware/validator.go
    - middleware/validator_test.go
EOF

# Test @build skill
claude "@build TEST-002-01"

# Expected Behavior:
# 1. Detects Go project (go.mod found)
# 2. Runs go test ./...
# 3. AI validators run:
#    - Coverage: go test -coverprofile output
#    - Architecture: Checks import paths
#    - Errors: Finds ignored errors
#    - Complexity: Counts lines per Go file
# 4. All quality gates PASS

# Verify Outputs:
# - Project type detected: Go
# - Test command: go test ./...
# - Coverage: go tool cover shows ≥80%
# - No architecture violations
# - No ignored errors
# - All files <200 LOC

# Save test results
echo "✓ Go test passed" > tests/test-go/results.txt
```

### Step 4: AI Validation Test

```bash
# Test AI validators work correctly across languages

cd tests/test-python
claude "@review"
# Expected:
# - /coverage-validator reads Python code
# - /architecture-validator checks src/ imports
# - /error-validator finds bare except clauses
# - /complexity-validator counts Python file lines
# Output: Structured report with PASS/FAIL

cd tests/test-java
claude "@review"
# Expected:
# - /coverage-validator reads Java code
# - /architecture-validator checks package structure
# - /error-validator finds empty catch blocks
# - /complexity-validator counts Java file lines
# Output: Structured report with PASS/FAIL

cd tests/test-go
claude "@review"
# Expected:
# - /coverage-validator reads Go code
# - /architecture-validator checks import paths
# - /error-validator finds ignored errors
# - /complexity-validator counts Go file lines
# Output: Structured report with PASS/FAIL
```

### Step 5: Documentation Validation

**File: sdp-plugin/docs/TUTORIAL.md** (NEW)

```markdown
# SDP Plugin Tutorial

## Quick Start (5 minutes)

1. **Install Plugin**
   ```bash
   git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
   cp -r ~/.claude/sdp/prompts/* .claude/
   ```

2. **Start Development**
   ```bash
   @feature "Add REST API"
   @design feature-rest-api
   @build 00-001-01
   ```

## Language Examples

### Python Workflow

```bash
# 1. Detect project type
cat pyproject.toml  # → Python

# 2. Create workstream
@feature "Add user authentication"

# 3. Plan workstreams
@design feature-auth

# 4. Execute
@build 00-001-01
# Runs: pytest tests/ -v
# AI validates coverage, architecture, errors, complexity

# 5. Review
@review feature-auth
# All quality gates PASS
```

### Java Workflow

```bash
# 1. Detect project type
cat pom.xml  # → Java

# 2. Create workstream
@feature "Add user authentication"

# 3. Plan workstreams
@design feature-auth

# 4. Execute
@build 00-001-01
# Runs: mvn test
# AI validates JaCoCo coverage, layer separation, catch blocks

# 5. Review
@review feature-auth
# All quality gates PASS
```

### Go Workflow

```bash
# 1. Detect project type
cat go.mod  # → Go

# 2. Create workstream
@feature "Add user authentication"

# 3. Plan workstreams
@design feature-auth

# 4. Execute
@build 00-001-01
# Runs: go test ./...
# AI validates go tool cover, import paths, error handling

# 5. Review
@review feature-auth
# All quality gates PASS
```

## Quality Gates

All languages use same quality gates:

| Gate | Threshold |
|------|-----------|
| Test Coverage | ≥80% |
| Type Safety | Complete signatures |
| Error Handling | No unsafe patterns |
| File Size | <200 LOC |
| Architecture | Clean layers |

Language-specific tools:
- **Python**: pytest, mypy, ruff
- **Java**: Maven, JaCoCo, javac
- **Go**: go test, go tool cover, go vet
```

**Validate:**
```bash
# Follow Python quick start in SDP repo
# Expected: All commands work

# Follow Java quick start in test-java
# Expected: All commands work

# Follow Go quick start in test-go
# Expected: All commands work
```

## Verification Summary

```bash
# Check all tests passed
cat tests/test-python/results.txt  # ✓ Python test passed
cat tests/test-java/results.txt    # ✓ Java test passed
cat tests/test-go/results.txt      # ✓ Go test passed

# Verify language detection
grep "Project type" tests/*/results.txt
# Expected: All 3 show correct language detected

# Verify quality gates
grep "quality gates" tests/*/results.txt
# Expected: All 3 show all gates PASS

# Verify documentation
ls -la sdp-plugin/docs/examples/
# Expected: python/, java/, go/ directories with QUICKSTART.md
```

## Quality Gates

- All 3 languages pass @build successfully
- Language detection works (pyproject.toml, pom.xml, go.mod)
- AI validators work correctly in all languages
- Documentation examples validated
- No Python dependencies required for Java/Go tests

## Dependencies

- 00-041-03 (Remove Python Dependencies from Skills)
- 00-041-04 (AI-Based Validation Prompts)
- 00-041-05 (Go Binary CLI) - optional, for init/doctor testing

## Blocks

- 00-041-07 (Marketplace Release) - needs validated plugin for release
