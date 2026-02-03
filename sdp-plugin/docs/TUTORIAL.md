# SDP Plugin Tutorial

Learn how to use the Spec-Driven Protocol (SDP) Claude Plugin with your project, regardless of programming language.

## Quick Start (5 minutes)

### 1. Install Plugin

```bash
# Clone or download SDP plugin
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp

# Copy prompts to your project
cp -r ~/.claude/sdp/prompts/* .claude/

# Or use the Go binary (optional)
cd sdp-plugin
make build
./bin/sdp init --project-type=python  # Auto-detects if omitted
```

### 2. Start Development

Open Claude Code and start building:

```
@feature "Add REST API for user management"
```

SDP will:
1. Ask deep questions about requirements
2. Design workstreams with acceptance criteria
3. Execute each workstream with TDD discipline
4. Validate quality with AI-based checks

## Language Examples

### Python Workflow

**Project Detection:** `pyproject.toml` found → Python

```bash
# 1. Install plugin
cp -r sdp-plugin/prompts/* .claude/

# 2. Create feature
@feature "Add user authentication"

# SDP asks:
# - What auth method? (JWT, sessions, OAuth)
# - Where should auth be applied? (specific routes, global middleware)
# - Error handling for failed auth? (401, redirect, custom)
# - Security concerns? (rate limiting, password hashing)

# 3. Design workstreams
@design feature-user-auth
# Creates: 00-001-01, 00-001-02, 00-001-03, 00-001-04

# 4. Execute workstream
@build 00-001-01
# Process:
# - Reads tests/ and src/ structure
# - Runs: pytest tests/ -v
# - AI validates:
#   * Coverage: Reads test files, maps to functions, calculates %
#   * Architecture: Checks src/domain/ doesn't import infrastructure/
#   * Errors: Finds bare except, except: pass patterns
#   * Complexity: Counts lines per file, checks <200 LOC
# - Output: ✅ PASS (all gates met)

# 5. Review feature
@review F01
# Runs all validators, produces quality report
```

**Expected Output:**
```
✅ Coverage: 87% (≥80% required)
✅ Architecture: No violations
✅ Error Handling: No unsafe patterns
✅ Complexity: All files <200 LOC

Verdict: ✅ PASS
```

### Java Workflow

**Project Detection:** `pom.xml` or `build.gradle` found → Java

```bash
# 1. Install plugin
cp -r sdp-plugin/prompts/* .claude/

# 2. Create feature
@feature "Add user authentication"

# SDP asks same deep questions (language-agnostic)

# 3. Design workstreams
@design feature-user-auth
# Creates: 00-001-01 through 00-001-05

# 4. Execute workstream
@build 00-001-01
# Process:
# - Reads src/main/java/ and src/test/java/ structure
# - Runs: mvn test
# - AI validates:
#   * Coverage: Reads JaCoCo reports or calculates from test/method mapping
#   * Architecture: Checks package separation (com.example.domain, .application, .infrastructure)
#   * Errors: Finds empty catch blocks, catch(Exception e) without logging
#   * Complexity: Counts lines per .java file, checks <200 LOC
# - Output: ✅ PASS (all gates met)

# 5. Review feature
@review F01
# Runs all validators, produces quality report
```

**Expected Output:**
```
✅ Coverage: 85% (JaCoCo report)
✅ Architecture: Clean package separation
✅ Error Handling: No empty catch blocks
✅ Complexity: All classes <200 LOC

Verdict: ✅ PASS
```

### Go Workflow

**Project Detection:** `go.mod` found → Go

```bash
# 1. Install plugin
cp -r sdp-plugin/prompts/* .claude/

# 2. Create feature
@feature "Add user authentication"

# SDP asks same deep questions (language-agnostic)

# 3. Design workstreams
@design feature-user-auth
# Creates: 00-001-01 through 00-001-04

# 4. Execute workstream
@build 00-001-01
# Process:
# - Reads .go files and *_test.go files
# - Runs: go test ./...
# - AI validates:
#   * Coverage: Reads go test -cover output or maps tests to functions
#   * Architecture: Checks import paths for layer violations
#   * Errors: Finds ignored errors (func(), _), empty recover()
#   * Complexity: Counts lines per .go file, checks <200 LOC
# - Output: ✅ PASS (all gates met)

# 5. Review feature
@review F01
# Runs all validators, produces quality report
```

**Expected Output:**
```
✅ Coverage: 82% (go test -cover)
✅ Architecture: Clean import paths
✅ Error Handling: No ignored errors
✅ Complexity: All files <200 LOC

Verdict: ✅ PASS
```

## Quality Gates

SDP applies the same quality standards across all languages:

| Gate | Threshold | Python Check | Java Check | Go Check |
|------|-----------|--------------|------------|----------|
| **Test Coverage** | ≥80% | pytest --cov | mvn verify (JaCoCo) | go test -cover |
| **Type Safety** | Complete | mypy --strict | javac -Xlint:all | go vet ./... |
| **Error Handling** | No unsafe patterns | No bare except | No empty catch | No ignored errors |
| **File Size** | <200 LOC | wc -l *.py | wc -l *.java | wc -l *.go |
| **Architecture** | Clean layers | Import checks | Package checks | Import path checks |

**AI-Based Validation:**
When tools aren't available, Claude AI analyzes code directly:
- **Coverage:** Reads all test files and source files, maps tests to functions, calculates percentage
- **Architecture:** Parses imports, checks for layer violations
- **Errors:** Searches for unsafe patterns (bare except, empty catch, ignored errors)
- **Complexity:** Counts lines, checks cyclomatic complexity, nesting depth

## Language-Specific Patterns

### Error Handling

**Python (❌ BAD):**
```python
try:
    risky()
except:
    pass  # Swallows all errors!
```

**Python (✅ GOOD):**
```python
try:
    risky()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise
```

**Java (❌ BAD):**
```java
try {
    risky();
} catch (Exception e) {
    // Empty - error lost
}
```

**Java (✅ GOOD):**
```java
try {
    risky();
} catch (ValueException e) {
    logger.error("Invalid value", e);
    throw e;
}
```

**Go (❌ BAD):**
```go
data, _ := fetchData()  // Error lost
```

**Go (✅ GOOD):**
```go
data, err := fetchData()
if err != nil {
    log.Printf("Fetch failed: %v", err)
    return err
}
```

### Clean Architecture

**Python:**
```python
# ✅ CORRECT: Infrastructure imports Domain
from src.domain.entities import User

# ❌ VIOLATION: Domain imports Infrastructure
from src.infrastructure.persistence import Database
```

**Java:**
```java
// ✅ CORRECT: Service imports Domain
import com.example.domain.entities.User;

// ❌ VIOLATION: Domain imports Infrastructure
import com.example.infrastructure.persistence.Database;
```

**Go:**
```go
// ✅ CORRECT: Infrastructure imports Domain
import "github.com/user/project/domain"

// ❌ VIOLATION: Domain imports Infrastructure
import "github.com/user/project/infrastructure"
```

## Workflow Comparison

### Python SDP (Old) vs Claude Plugin (New)

| Aspect | Python SDP | Claude Plugin |
|--------|-----------|---------------|
| **Installation** | pip install sdp | Copy prompts to .claude/ |
| **Dependencies** | Python 3.10+, Click, PyYAML | None (prompts work in any language) |
| **Validation** | pytest, mypy, ruff | AI analysis (or tools if available) |
| **Languages** | Python only | Python, Java, Go, any language |
| **Workflow** | sdp commands | Claude Code skills (@feature, @build, etc.) |

### Migration from Python SDP

If you're migrating from the Python SDP CLI:

1. **Install plugin:**
   ```bash
   cp -r sdp-plugin/prompts/* .claude/
   ```

2. **Existing workstreams:** No changes needed - format is identical

3. **Quality checks:** Now use AI validators (slower but language-agnostic)

4. **Commands:**
   - Old: `sdp feature create "Add auth"`
   - New: `@feature "Add auth"` (interactive, deeper questions)

5. **Validation:**
   - Old: `sdp quality check --module src/` (pytest, mypy, ruff)
   - New: `@review` (AI analysis, works for any language)

## Common Commands

### Create Feature
```
@feature "Add user authentication"
```
- Deep requirements gathering
- Progressive disclosure (asks questions as needed)
- Explores tradeoffs and concerns

### Plan Workstreams
```
@design feature-user-auth
```
- Explores codebase
- Designs WS decomposition
- Asks architecture questions

### Execute Workstream
```
@build 00-001-01
```
- TDD discipline (Red → Green → Refactor)
- AI quality validation
- Real-time progress tracking

### Quality Review
```
@review F01
```
- Runs all validators
- Checks acceptance criteria
- Produces verdict (PASS/FAIL)

### Deploy
```
@deploy F01
```
- Generates artifacts (CHANGELOG, release notes)
- Executes GitFlow merge
- Creates and pushes git tag

## Troubleshooting

### Plugin not loading

**Problem:** Claude Code doesn't show SDP skills

**Solution:**
```bash
# Check prompts exist
ls -la .claude/skills/
# Should show: feature.md, design.md, build.md, etc.

# Check settings.json
cat .claude/settings.json
# Should list skills
```

### Language detection fails

**Problem:** Wrong language detected

**Solution:**
```bash
# For Python: Ensure pyproject.toml exists
touch pyproject.toml

# For Java: Ensure pom.xml or build.gradle exists
touch pom.xml

# For Go: Ensure go.mod exists
touch go.mod

# Or specify manually:
@build 00-001-01 --project-type=python
```

### Validators not working

**Problem:** AI validators give wrong results

**Solution:**
```bash
# Check validator prompts exist
ls -la .claude/validators/
# Should show: coverage.md, architecture.md, errors.md, complexity.md, all.md

# Run validators manually:
claude "/coverage-validator"
claude "/architecture-validator"
claude "/error-validator"
claude "/complexity-validator"
```

### Tests not running

**Problem:** Build command doesn't run tests

**Solution:**
```bash
# Check for language-specific test files:
# Python: tests/test_*.py
# Java: src/test/java/**/*Test.java
# Go: *_test.go

# Run tests manually to verify:
# Python: pytest tests/ -v
# Java: mvn test
# Go: go test ./...
```

## Advanced Usage

### Multi-Language Projects

If your project uses multiple languages:

```bash
# SDP detects primary language from build files
# Specify explicitly if needed:
@build 00-001-01 --project-type=python

# Or use language-agnostic mode:
@build 00-001-01 --project-type=agnostic
```

### Custom Validators

Add your own validators:

```bash
# Create custom validator
cat > .claude/validators/custom.md <<'EOF'
# Custom Validator

Check for project-specific rules:
- No TODO comments in production code
- All public APIs have documentation
- Error messages follow company style guide

Run with: claude "/custom-validator"
EOF
```

### Integrating with CI/CD

GitHub Actions example:

```yaml
name: SDP Quality Check
on: [push, pull_request]

jobs:
  sdp-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run SDP validators
        run: |
          claude "/coverage-validator"
          claude "/architecture-validator"
          claude "/error-validator"
          claude "/complexity-validator"
```

## Next Steps

1. **Install plugin:** `cp -r sdp-plugin/prompts/* .claude/`
2. **Read examples:** `sdp-plugin/docs/examples/*/QUICKSTART.md`
3. **Start building:** `@feature "Your first feature"`
4. **Join community:** https://github.com/ai-masters/sdp/discussions

## Support

- **Documentation:** `sdp-plugin/docs/`
- **Examples:** `sdp-plugin/docs/examples/`
- **Issues:** https://github.com/ai-masters/sdp/issues
- **Discussions:** https://github.com/ai-masters/sdp/discussions
