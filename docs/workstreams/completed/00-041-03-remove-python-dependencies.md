# 00-041-03: Remove Python Dependencies from Skills

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** LARGE
> **Created:** 2026-02-02

## Goal

Update all 18 skills to work without Python tools, add project type detection, integrate AI validators.

## Acceptance Criteria

- AC1: `@build` skill detects project type (Python/Java/Go) from build files
- AC2: `@tdd` skill has language-agnostic examples (not pytest-only)
- AC3: `@review` skill uses AI-based validation (not pytest/mypy/ruff)
- AC4: All 18 skills updated and tested
- AC5: Skills work in pure Java/Go projects (no Python required)

## Scope

### Input Files
- `.claude/skills/build/SKILL.md` (148 lines, MAJOR rewrite needed)
- `.claude/skills/review/SKILL.md` (103 lines, MAJOR rewrite needed)
- `.claude/skills/tdd/SKILL.md` (142 lines, minor rewrite)
- All other skills (15 files, minor updates)

### Output Files
- `sdp-plugin/prompts/skills/build.md` (MAJOR changes)
- `sdp-plugin/prompts/skills/review.md` (MAJOR changes)
- `sdp-plugin/prompts/skills/tdd.md` (minor changes)
- All other skills (adapted)

### Out of Scope
- Creating AI validators (WS-00-041-04 does this)
- Go binary development (WS-00-041-05)

## Implementation Steps

### Step 1: Update @build Skill

**File: sdp-plugin/prompts/skills/build.md**

Replace tool-based quality check (lines 96-103):

```markdown
# OLD
sdp quality check --module {module}
Must pass:
- Coverage ≥80%
- mypy --strict
- ruff (no errors)
- Files <200 LOC

# NEW
## Step 5: Quality Check

### 1. Detect Project Type

Check for build files in project root:
- `pyproject.toml` or `setup.py` → **Python**
- `pom.xml` or `build.gradle` → **Java**
- `go.mod` → **Go**
- `package.json` → **Node.js**
- `Cargo.toml` → **Rust**

If multiple found, ask user to specify.

### 2. Run Tests

| Language | Test Command |
|----------|--------------|
| Python | `pytest tests/ -v` |
| Java | `mvn test` or `gradle test` |
| Go | `go test ./...` |
| Node.js | `npm test` |
| Rust | `cargo test` |

### 3. AI Validation

Ask Claude to analyze:

**Coverage:**
"Read all test files and source files. Map tests to functions. Calculate: (tested_functions / total_functions) * 100. Is ≥80%? Report missing coverage."

**Type Safety:**
"Review all function signatures. Are types complete? Any missing annotations? Report violations."

**File Size:**
"Count lines in each file. List files >200 LOC. Report violations."

**Error Handling:**
"Search for bare except clauses or catch-all exception handlers. Report violations."

**Gate:** All AI validators must PASS.
```

### Step 2: Update @review Skill

**File: sdp-plugin/prompts/skills/review.md**

Replace tool-based quality gates (lines 44-49):

```markdown
# OLD
pytest tests/ -v
pytest --cov=src --cov-fail-under=80
mypy src/ --strict
ruff check src/
grep -r "except:" src/ | grep "pass"

# NEW
## Step 3: Quality Gates (AI Validators)

Run all AI validators in sequence:

1. **Coverage Validator** (`/coverage-validator`)
   - Analyzes test coverage
   - Reports percentage ≥80%

2. **Architecture Validator** (`/architecture-validator`)
   - Checks Clean Architecture layer violations
   - Reports dependency issues

3. **Error Validator** (`/error-validator`)
   - Finds bare except clauses
   - Reports unsafe error handling

4. **Complexity Validator** (`/complexity-validator`)
   - Counts lines per file
   - Reports files >200 LOC

**Output Format:**
```markdown
## Quality Gates Report

✅ Coverage: 85% (≥80% required)
✅ Architecture: No violations
❌ Error Handling: 3 bare except clauses found
✅ Complexity: All files <200 LOC

**Verdict:** CHANGES_REQUESTED
```

**Gate:** All validators must PASS for APPROVED verdict.
```

### Step 3: Update @tdd Skill

**File: sdp-plugin/prompts/skills/tdd.md**

Add language-agnostic examples:

```markdown
# RED Phase: Write Failing Test

**Python Example:**
```bash
# Create test file
cat > tests/test_user.py <<EOF
def test_user_creation():
    user = User("Alice")
    assert user.name == "Alice"
EOF

# Run test (should fail)
pytest tests/test_user.py -v
# Expected: FAILED (User not implemented)
```

**Java Example:**
```bash
# Create test class
cat > src/test/java/UserTest.java <<EOF
@Test
public void testUserCreation() {
    User user = new User("Alice");
    assertEquals("Alice", user.getName());
}
EOF

# Run test (should fail)
mvn test -Dtest=UserTest
# Expected: FAILED (User not implemented)
```

**Go Example:**
```bash
# Create test file
cat > user_test.go <<EOF
func TestUserCreation(t *testing.T) {
    user := User{Name: "Alice"}
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
}
EOF

# Run test (should fail)
go test -v
# Expected: FAILED (User not implemented)
```
```

### Step 4: Update All Other Skills

For each skill:
- Replace "pytest" with generic "run tests"
- Replace "mypy" with generic "type checking"
- Replace "ruff" with generic "linting"
- Remove Python-specific paths (src/**/*.py)
- Keep language-agnostic logic

Skills to update:
- `feature` (minor)
- `design` (minor)
- `deploy` (minor - remove uv run pytest)
- `bugfix` (minor - remove pytest coverage assumptions)
- `oneshot` (minor - simplify Python integration)
- Others (check for Python assumptions)

## Verification

```bash
# Test 1: Python project (SDP itself)
cd /Users/fall_out_bug/projects/vibe_coding/sdp
cp -r sdp-plugin/prompts/* .claude/
claude "@build 00-041-01"
# Expected: Detects Python (pyproject.toml), runs pytest, quality gates pass

# Test 2: Java project (Spring Petclinic)
git clone https://github.com/spring-projects/spring-petclinic.git tests/test-java
cd tests/test-java
cp -r ../../sdp-plugin/prompts/* .claude/
claude "@build TEST-001-01"
# Expected: Detects Java (pom.xml), runs mvn test, quality gates pass

# Test 3: Go project (Gin)
git clone https://github.com/gin-gonic/gin.git tests/test-go
cd tests/test-go
cp -r ../../sdp-plugin/prompts/* .claude/
claude "@build TEST-002-01"
# Expected: Detects Go (go.mod), runs go test, quality gates pass
```

## Quality Gates

- @build detects project type correctly
- @review uses AI validators (no pytest/mypy/ruff)
- @tdd has examples for all 3 languages
- All skills work without Python installed
- Language detection logic documented

## Dependencies

- 00-041-01 (Plugin Package Structure) - provides plugin structure
- 00-041-02 (Language-Agnostic Protocol Docs) - provides language examples

## Blocks

- 00-041-04 (AI-Based Validation Prompts) - depends on updated skills
- 00-041-06 (Cross-Language Validation) - needs updated skills for testing
