# SDP Error Codes Reference

This document defines the standard error taxonomy and codes for SDP operations.

## Error Classes

SDP errors are organized into five classes:

| Class | Code | Description | Recovery Pattern |
|-------|------|-------------|------------------|
| Environment | `ENV` | Runtime environment issues (tools, permissions, filesystem) | Install/fix environment |
| Protocol | `PROTO` | SDP protocol violations (invalid IDs, malformed files) | Fix input format |
| Dependency | `DEP` | Workstream dependency issues (blocked, circular) | Resolve dependencies |
| Validation | `VAL` | Quality gate failures (coverage, size, tests) | Fix code quality |
| Runtime | `RUNTIME` | Unexpected runtime failures (network, timeout) | Retry or investigate |

## Error Code Format

Error codes follow the pattern: `{CLASS}{NUMBER}` (e.g., `ENV001`, `PROTO002`)

## Environment Errors (ENV001-ENV099)

| Code | Message | Recovery Hint |
|------|---------|---------------|
| ENV001 | Git is not installed or not found in PATH | Install Git from https://git-scm.com |
| ENV002 | Go is not installed or not found in PATH | Install Go from https://go.dev/dl/ |
| ENV003 | Claude Code CLI is not installed | Install Claude Code CLI from Anthropic |
| ENV004 | Beads CLI is not installed (required for task tracking) | Install Beads: `brew tap beads-dev/tap && brew install beads` |
| ENV005 | Permission denied | Check file permissions or run with appropriate privileges |
| ENV006 | Git worktree not found | Verify you are in a valid git worktree |
| ENV007 | SDP configuration file not found | Run `sdp init` to create configuration |
| ENV008 | Required directory not found | Ensure required directories exist |
| ENV009 | File is not writable | Check file permissions |

## Protocol Errors (PROTO001-PROTO099)

| Code | Message | Recovery Hint |
|------|---------|---------------|
| PROTO001 | Invalid workstream ID format (expected PP-FFF-SS) | Verify the ID format matches PP-FFF-SS or FNNN |
| PROTO002 | Invalid feature ID format (expected FNNN) | Verify the ID format matches PP-FFF-SS or FNNN |
| PROTO003 | YAML parsing error | Check YAML syntax and structure |
| PROTO004 | Required field is missing | Provide all required fields |
| PROTO005 | Invalid status value | Use valid status: pending, in_progress, completed, failed |
| PROTO006 | Evidence hash chain is broken | Run `sdp log trace --verify` to diagnose |
| PROTO007 | Session file is corrupted or tampered | Run `sdp session repair` or delete .sdp/session.json |
| PROTO008 | Invalid event type | Use valid event types: plan, generation, verification |
| PROTO009 | Schema validation failed | Verify file matches expected schema |

## Dependency Errors (DEP001-DEP099)

| Code | Message | Recovery Hint |
|------|---------|---------------|
| DEP001 | Workstream is blocked by unresolved dependencies | Complete blocking workstreams first or use `--force` |
| DEP002 | Circular dependency detected | Review workstream dependencies for cycles |
| DEP003 | Required prerequisite is not satisfied | Ensure all prerequisites are satisfied |
| DEP004 | Feature not found | Verify the ID exists in docs/workstreams/ |
| DEP005 | Workstream not found | Verify the ID exists in docs/workstreams/ |
| DEP006 | File scope collision detected between workstreams | Review workstream scope files for overlaps |

## Validation Errors (VAL001-VAL099)

| Code | Message | Recovery Hint |
|------|---------|---------------|
| VAL001 | Test coverage is below required threshold | Add tests to increase coverage to >= 80% |
| VAL002 | File exceeds maximum allowed size | Split file into smaller modules (< 200 LOC) |
| VAL003 | Tests failed | Run tests with verbose output to diagnose failures |
| VAL004 | Linting failed | Fix linting errors reported by linter |
| VAL005 | Type mismatch | Verify types match expected signatures |
| VAL006 | Quality gate failed | Review quality gate output for specific failures |
| VAL007 | Drift detected between code and documentation | Run `sdp drift detect` for details and sync |
| VAL008 | Edit scope violation | Stay within workstream scope or use `sdp guard deactivate` |

## Runtime Errors (RUNTIME001-RUNTIME099)

| Code | Message | Recovery Hint |
|------|---------|---------------|
| RUNTIME001 | External command failed | Check command output for details |
| RUNTIME002 | Network error | Check network connectivity and retry |
| RUNTIME003 | Operation timed out | Increase timeout or optimize operation |
| RUNTIME004 | Resource exhausted | Free up resources and retry |
| RUNTIME005 | Unexpected state encountered | Run `sdp doctor` to diagnose environment |
| RUNTIME006 | Internal error | Report this issue with full error context |

## Error JSON Format

All SDP errors can be serialized to JSON for programmatic consumption:

```json
{
  "code": "ENV001",
  "class": "ENV",
  "message": "Git is not installed or not found in PATH",
  "recovery_hint": "Install Git from https://git-scm.com",
  "context": {
    "file": "config.yml",
    "operation": "read"
  },
  "cause": "executable not found in $PATH"
}
```

## Using Errors Programmatically

```go
import "github.com/fall-out-bug/sdp/internal/errors"

// Create a new error
err := errors.New(errors.ErrGitNotFound, nil)

// With context
err.WithContext("file", "config.yml")

// Check error class
if errors.GetClass(err) == errors.ClassEnvironment {
    // Handle environment error
}

// Get error code
code := errors.GetCode(err) // returns ErrGitNotFound

// Serialize to JSON
jsonStr, _ := errors.ToJSON(err)
```

## Recovery Patterns by Class

### Environment Errors
1. Run `sdp doctor` to diagnose missing tools
2. Install missing dependencies
3. Fix file/directory permissions
4. Reinitialize configuration with `sdp init`

### Protocol Errors
1. Verify input format matches specification
2. Check YAML frontmatter syntax
3. Run `sdp parse` to validate workstream files
4. Review PROTOCOL.md for correct format

### Dependency Errors
1. Check workstream dependencies in backlog files
2. Complete blocking workstreams first
3. Use `sdp apply --ws` for single workstream
4. Review collision detection output

### Validation Errors
1. Run `sdp quality` for detailed report
2. Increase test coverage
3. Split large files
4. Fix linting errors
5. Run `sdp drift detect` to sync code and docs

### Runtime Errors
1. Retry with exponential backoff
2. Check network connectivity
3. Free up system resources
4. Report persistent internal errors with logs

---

**Version:** 0.10.0
**Last Updated:** 2026-02-16
