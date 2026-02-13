---
name: build
description: Execute workstream with TDD and guard enforcement
cli: sdp apply --ws
llm: Spawn subagents for 3-stage review
version: 6.1.0
---

# build

> **CLI:** `sdp apply --ws <workstream-id>` (file operations only)
> **LLM:** Spawn 3 subagents for TDD cycle + review

Execute a single workstream following TDD discipline.

---

## EXECUTE THIS NOW

When user invokes `@build 00-067-01`:

1. Run CLI to setup and validate:
```bash
sdp guard activate 00-067-01
sdp apply --ws 00-067-01 --dry-run  # Preview first
```

2. Then spawn 3 subagents for TDD cycle:
   - **Implementer** - Write tests and code
   - **Spec Reviewer** - Verify matches spec
   - **Quality Reviewer** - Run quality gates

**DO NOT skip step 2.** The CLI only handles file operations. TDD cycle requires spawning subagents.

---

## How to Spawn Subagents

Use your tool's subagent capability. For example:
- Claude Code: Use Task tool with `subagent_type="general-purpose"`
- Cursor: Use agent panel
- Windsurf: Use agent spawning

---

## Subagent 1: Implementer

**Role file:** `.claude/agents/implementer.md`

**Task:**
```
You are the IMPLEMENTER for workstream 00-067-01.

Read the spec: docs/workstreams/backlog/00-067-01.md

Execute TDD cycle for each Acceptance Criteria:
1. RED: Write failing test first
2. GREEN: Write minimum code to pass
3. REFACTOR: Clean up while keeping tests green

Quality gates:
- Test coverage >= 80%
- All tests passing
- No lint errors

Output: Verdict PASS or FAIL with evidence
```

---

## Subagent 2: Spec Reviewer

**Role file:** `.claude/agents/spec-reviewer.md`

**Task:**
```
You are the SPEC COMPLIANCE REVIEWER for workstream 00-067-01.

CRITICAL: Do NOT trust the implementer's report. Verify yourself.

1. Read the actual code
2. Run tests yourself
3. Check coverage yourself
4. Verify each AC is implemented

Output: Verdict PASS or FAIL with evidence
```

---

## Subagent 3: Quality Reviewer

**Task:**
```
You are the QUALITY REVIEWER for workstream 00-067-01.

Run comprehensive quality check:
1. Test coverage (>=80%)
2. LOC check (<=200 lines per file) - MANDATORY
3. Code quality (complexity, duplication)
4. Security check
5. Lint passes

LOC Gate (MANDATORY):
```bash
for file in *.go; do
  loc=$(wc -l < "$file")
  if [ "$loc" -gt 200 ]; then
    echo "ERROR: $file is $loc LOC (max: 200)"
    exit 1
  fi
done
```

Output: Verdict PASS or FAIL with evidence
```

---

## After All Subagents Complete

**If all 3 PASS:**
```bash
sdp guard complete 00-067-01
git add .
git commit -m "feat(F067): 00-067-01 - {title}"
```

**If any FAIL:** Report failure, do not commit.

---

## Identifier Formats

```bash
@build 00-067-01      # Workstream ID (PP-FFF-SS)
@build 99-F064-01     # Fix workstream (99-{FEATURE}-{SEQ})
@build sdp-xxx        # Beads task ID (resolved)
```

---

## Quality Gates

| Gate | Threshold | Check |
|------|-----------|-------|
| Tests | 100% pass | `go test ./...` |
| Coverage | >= 80% | `go test -cover ./...` |
| Lint | 0 errors | `golangci-lint run` |
| File Size | <= 200 LOC | `wc -l *.go` |

---

## Beads Integration

When Beads enabled:
1. **Before:** `bd update {beads_id} --status in_progress`
2. **Success:** `bd close {beads_id} --reason "WS completed"`
3. **Failure:** `bd update {beads_id} --status blocked`

---

## See Also

- `.claude/patterns/tdd.md` - TDD pattern
- `.claude/patterns/quality-gates.md` - Quality gates
- `@oneshot` - Execute all workstreams

**Implementation:** `sdp-plugin/cmd/sdp/apply.go`
