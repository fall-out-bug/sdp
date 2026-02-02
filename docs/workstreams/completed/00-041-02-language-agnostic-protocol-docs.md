# 00-041-02: Language-Agnostic Protocol Docs

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** MEDIUM
> **Created:** 2026-02-02

## Goal

Rewrite PROTOCOL.md and quality-gates.md removing Python assumptions, adding multi-language examples.

## Acceptance Criteria

- AC1: PROTOCOL.md updated with generic test commands (no pytest/mypy/ruff hardcoded)
- AC2: quality-gates.md uses language-specific tables (Python, Java, Go)
- AC3: Example directories created for Python, Java, Go
- AC4: All "pytest/mypy/ruff" references generalized to language-agnostic equivalents
- AC5: Documentation examples verified valid in respective languages

## Scope

### Input Files
- `docs/PROTOCOL.md` (existing, Python-centric)
- `docs/reference/quality-gates.md` (existing, tool-specific)
- Existing language examples for reference

### Output Files
- `sdp-plugin/docs/PROTOCOL.md` (adapted)
- `sdp-plugin/docs/quality-gates.md` (adapted)
- `sdp-plugin/docs/examples/python/QUICKSTART.md` (NEW)
- `sdp-plugin/docs/examples/java/QUICKSTART.md` (NEW)
- `sdp-plugin/docs/examples/go/QUICKSTART.md` (NEW)

### Out of Scope
- Modifying skills (WS-00-041-03)
- Creating AI validators (WS-00-041-04)

## Implementation Steps

1. **Create Language Example Directories**
   ```
   sdp-plugin/docs/examples/python/
   sdp-plugin/docs/examples/java/
   sdp-plugin/docs/examples/go/
   ```

2. **Adapt PROTOCOL.md**
   - Rewrite quality gates section
   - Replace hardcoded pytest with language-specific table
   - Remove "pip install sdp" from quick start
   - Add language-agnostic installation

3. **Adapt quality-gates.md**
   - Create language-specific tables
   - Python column: pytest, mypy, ruff
   - Java column: Maven/Gradle, JaCoCo, javac
   - Go column: go test, go tool cover, go vet

4. **Create QUICKSTART Guides**
   - Python: pytest workflow
   - Java: Maven workflow
   - Go: go test workflow
   - Each: prerequisites, commands, quality gates

## Key Changes Example

```markdown
# OLD (Python-specific)
## Quality Gates
pytest --cov=src/ --cov-fail-under=80
mypy src/ --strict
ruff check src/

# NEW (Language-agnostic)
## Quality Gates

### Test Coverage ≥80%

| Language | Command |
|----------|---------|
| Python | pytest --cov=src/ --cov-fail-under=80 |
| Java | mvn verify (JaCoCo report) |
| Go | go test -coverprofile=coverage.out && go tool cover -func=coverage.out |

### Type Checking

| Language | Command |
|----------|---------|
| Python | mypy src/ --strict |
| Java | javac -Xlint:all |
| Go | go vet ./... |
```

## Verification

```bash
# Check all language examples present
ls sdp-plugin/docs/examples/python/
ls sdp-plugin/docs/examples/java/
ls sdp-plugin/docs/examples/go/

# Verify no hardcoded pytest in PROTOCOL.md
grep -n "pytest" sdp-plugin/docs/PROTOCOL.md
# Should only appear in Python column of tables

# Verify language-specific tables exist
grep -A 10 "Test Coverage" sdp-plugin/docs/quality-gates.md
# Should show table with 3 languages
```

## Quality Gates

- No "pytest", "mypy", "ruff" outside Python context
- All 3 languages have examples
- Tables properly formatted
- Examples copy-paste valid (verify by manual test)

## Dependencies

None

## Blocks

- 00-041-03 (Remove Python Dependencies from Skills)
