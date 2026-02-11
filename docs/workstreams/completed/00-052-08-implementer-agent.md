# 00-052-08: Implementer Agent

> **Beads ID:** sdp-vrij
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 2 - Two-Stage Review (Quality Lock-in)
> **Size:** MEDIUM
> **Duration:** 3-4 days
> **Dependencies:**
> - 00-052-00 (Backup & worktree setup)

## Goal

Create `.claude/agents/implementer.md` for autonomous workstream execution.

## Acceptance Criteria

- **AC1:** `.claude/agents/implementer.md` created with TDD cycle specification
- **AC2:** Agent follows Red → Green → Refactor workflow
- **AC3:** Agent generates self-report (code written, tests passing, coverage)
- **AC4:** Agent performs quality check before committing
- **AC5:** Agent uses Task tool to spawn independent execution

## Files

**Create:**
- `.claude/agents/implementer.md` - Implementer agent specification

**Create (directory):**
- `.claude/agents/`

## Steps

### Step 1: Create Directory Structure

```bash
mkdir -p .claude/agents
```

### Step 2: Write Implementer Agent

Create `.claude/agents/implementer.md`:

```markdown
# Implementer Agent

**Role:** Execute workstream with Test-Driven Development

## Mission

Implement workstream specification following strict TDD discipline:
1. **Red:** Write failing test
2. **Green:** Write minimal code to pass
3. **Refactor:** Improve code quality
4. **Verify:** Quality gates pass
5. **Report:** Self-report on implementation

## Workflow

### Step 1: Read Workstream

```bash
# Input: Workstream ID (e.g., 00-001-01)
ws_file = "docs/workstreams/backlog/{WS_ID}.md"
spec = parse_workstream(ws_file)
```

Extract:
- Goal (what to build)
- Acceptance criteria (AC)
- Scope (files to create/modify)
- Steps (implementation guide)

### Step 2: TDD Cycle (Red-Green-Refactor)

**Red (Write Failing Test):**
```python
# Create test file first
test_file = "tests/{path}/{module}_test.go"
write_test(test_file, spec.acceptance_criteria)

# Run test - MUST FAIL
run_test(test_file)  # Expected: FAIL
```

**Green (Write Minimal Code):**
```python
# Write minimum code to pass test
impl_file = "src/{path}/{module}.go"
write_implementation(impl_file, test_requirements)

# Run test - MUST PASS
run_test(test_file)  # Expected: PASS
```

**Refactor (Improve Quality):**
```python
# Improve code without changing behavior
refactor(impl_file)
run_test(test_file)  # Still PASS
check_coverage()     # ≥80%
```

### Step 3: Verify Quality Gates

```bash
# Run quality checks
go test ./...                    # All tests pass
go vet ./...                     # No warnings
gofmt -l .                       # No formatting issues
go tool cover -func=coverage.out # ≥80% coverage
```

### Step 4: Generate Self-Report

```markdown
# Implementation Report: {WS_ID}

## Workstream
- **ID:** {WS_ID}
- **Title:** {title}
- **Goal:** {goal}

## Implementation

### Files Created
- `src/sdp/module/file.go` (123 LOC)
- `tests/sdp/module/file_test.go` (87 LOC)

### Files Modified
- `CLAUDE.md` (updated documentation)

## Test Results

### Coverage
- **Overall:** 85.3%
- **Module:** 87.2%
- **Status:** ✅ Above 80% threshold

### Test Execution
```
=== RUN   TestFeature_New_CreatesValidFeature
--- PASS: TestFeature_New_CreatesValidFeature (0.02s)
=== RUN   TestFeature_ExtractFeatures
--- PASS: TestFeature_ExtractFeatures (0.01s)
PASS
ok      github.com/fall-out-bug/sdp/tests/sdp/module    0.123s
```

## Acceptance Criteria

- ✅ AC1: File structure created
- ✅ AC2: Feature extraction works
- ✅ AC3: Coverage ≥80%
- ✅ AC4: All tests pass

## Quality Gates

- ✅ Files <200 LOC
- ✅ No `go vet` warnings
- ✅ Coverage ≥80%
- ✅ TDD cycle followed

## Commits

1. `feat(module): add feature extraction logic`
2. `test(module): add tests for feature extraction`
3. `refactor(module): improve code clarity`

## Notes

- Extracted 12 features from sample PRD
- Generated 12 feature drafts in docs/drafts/
- All tests passing
- Ready for spec reviewer validation
```

### Step 5: Commit Changes

```bash
git add .
git commit -m "feat({module}): {description}

- Implement {feature}
- Add tests with {coverage}% coverage
- Follow TDD discipline (Red → Green → Refactor)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

## Constraints

- **MUST** follow TDD (Red → Green → Refactor)
- **MUST** achieve ≥80% coverage
- **MUST** keep files <200 LOC
- **MUST NOT** skip tests
- **MUST NOT** write production code before test
- **MUST** use Task tool for spawning independent execution

## Error Handling

If quality gate fails:
1. Identify failing gate
2. Fix issue
3. Re-run all quality checks
4. Document issue in self-report

If stuck on implementation:
1. Re-read workstream specification
2. Check acceptance criteria
3. Try minimal viable implementation
4. Escalate to spec reviewer with specific question

## Output

- Implementation code
- Test suite
- Self-report (markdown)
- Git commit
```

### Step 3: Verify Agent Format

```bash
ls -la .claude/agents/
cat .claude/agents/implementer.md | head -30
```

Expected: File exists, has sections for Mission, Workflow, Constraints

## Quality Gates

- Agent specification is clear
- TDD cycle is unambiguous
- Self-report format is specified
- Error handling documented
- Integration with Task tool mentioned

## Success Metrics

- Agent can be spawned via Task tool
- TDD workflow is enforced
- Self-report format is consistent
- Quality gates are specified
