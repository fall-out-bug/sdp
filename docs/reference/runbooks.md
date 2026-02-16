# SDP Recovery Runbook

This runbook provides step-by-step procedures for recovering from common SDP failures.

## Quick Reference

| Error Class | Common Codes | First Action |
|-------------|--------------|--------------|
| Environment | ENV001-ENV009 | Run `sdp doctor` |
| Protocol | PROTO001-PROTO009 | Validate input format |
| Dependency | DEP001-DEP006 | Check dependencies |
| Validation | VAL001-VAL008 | Run `sdp quality` |
| Runtime | RUNTIME001-006 | Retry or diagnose |

## Environment Recovery (ENV)

### ENV001: Git Not Found

**Symptoms:**
- Commands fail with "git: command not found"
- Git operations return errors

**Recovery:**
```bash
# Quick fix
brew install git

# Verify
git --version

# Configure if needed
git config --global user.name "Your Name"
git config --global user.email "your@email.com"
```

### ENV004: Beads Not Found

**Symptoms:**
- Task tracking commands fail
- Dependency resolution fails

**Recovery:**
```bash
# Install Beads
brew tap beads-dev/tap && brew install beads

# Verify
bd --version
```

### ENV005: Permission Denied

**Symptoms:**
- Cannot write to files
- Cannot create directories

**Recovery:**
```bash
# Check permissions
ls -la <path>

# Fix file permissions
chmod 644 <file>

# Fix directory permissions
chmod 755 <directory>

# Take ownership if needed
chown -R $(whoami) <directory>
```

## Protocol Recovery (PROTO)

### PROTO001: Invalid Workstream ID

**Symptoms:**
- "Invalid workstream ID format" error
- Commands reject WS IDs

**Recovery:**
1. Verify ID format: `PP-FFF-SS` (e.g., `00-070-01`)
   - PP = Project (00 for SDP)
   - FFF = Feature number (001-999)
   - SS = Step number (01-99)
2. Check for typos
3. Ensure dashes are present

### PROTO006: Hash Chain Broken

**Symptoms:**
- Evidence verification fails
- "Hash chain broken" error

**Recovery:**
```bash
# Verify chain integrity
sdp log trace --verify

# If broken, backup and repair
cp .sdp/log/events.jsonl events-backup.jsonl

# Remove corrupted entries (identify break point from trace)
head -n <line> events-backup.jsonl > .sdp/log/events.jsonl

# Verify repair
sdp log trace --verify
```

### PROTO007: Session Corrupted

**Symptoms:**
- "Session file corrupted" error
- State appears inconsistent

**Recovery:**
```bash
# Quick fix: delete and reinit
rm .sdp/session.json
sdp init

# Alternative: try repair
sdp session repair
```

## Dependency Recovery (DEP)

### DEP001: Blocked Workstream

**Symptoms:**
- "Workstream is blocked" error
- Cannot proceed with execution

**Recovery:**
```bash
# Check what's blocking
cat docs/workstreams/backlog/<ws-id>.md | grep depends_on

# Complete blocking workstreams first
sdp apply --ws <blocking-ws-id>

# Alternative: view dependency graph
sdp plan --graph
```

### DEP002: Circular Dependency

**Symptoms:**
- "Circular dependency detected" error
- Execution hangs or fails

**Recovery:**
1. Review workstream files for circular references
2. Identify cycle in dependency graph
3. Break cycle by removing one dependency
4. Consider splitting workstream if necessary

## Validation Recovery (VAL)

### VAL001: Coverage Low

**Symptoms:**
- "Test coverage is below required threshold"
- Quality gate fails

**Recovery:**
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View uncovered code
go tool cover -func=coverage.out | grep -v 100.0%

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Target: 80%+ coverage
go test -cover ./...
```

### VAL002: File Too Large

**Symptoms:**
- "File exceeds maximum allowed size"
- Quality gate fails

**Recovery:**
1. Count lines: `wc -l <file>`
2. Identify logical splits (by responsibility)
3. Extract to new files
4. Each file should be < 200 LOC
5. Update imports
6. Run tests to verify

### VAL003: Tests Failed

**Symptoms:**
- Test suite fails
- Quality gate fails

**Recovery:**
```bash
# Run with verbose output
go test -v ./...

# Run specific failing test
go test -v -run TestName ./...

# Fix issues and verify
go test ./...
```

### VAL007: Drift Detected

**Symptoms:**
- Code and documentation are out of sync
- "Drift detected" error

**Recovery:**
```bash
# View drift details
sdp drift detect

# Generate detailed report
sdp drift report --output=drift.md

# Decide: update code or update docs
# Make changes

# Verify drift resolved
sdp drift detect
```

## Runtime Recovery (RUNTIME)

### RUNTIME001: Command Failed

**Symptoms:**
- External command returns error
- Unexpected exit code

**Recovery:**
1. Run command manually with verbose flags
2. Check environment variables
3. Verify PATH includes command location
4. Try with absolute path
5. Check for permission issues

### RUNTIME003: Timeout Exceeded

**Symptoms:**
- Operation takes too long
- "Timeout exceeded" error

**Recovery:**
1. Retry the operation
2. Check system load: `top`
3. Consider increasing timeout
4. Profile slow operations
5. Consider async execution

### RUNTIME006: Internal Error

**Symptoms:**
- Unexpected internal error
- No clear recovery path

**Recovery:**
```bash
# Run diagnostics
sdp doctor --deep

# Gather logs
sdp log show > error-log.txt

# Report issue with:
# - Error message
# - sdp doctor output
# - Steps to reproduce
```

## Diagnostics Tools

### sdp diagnose

Use `sdp diagnose` to get help with specific error codes:

```bash
# Show all error classes
sdp diagnose --list-classes

# Show all error codes
sdp diagnose --list-codes

# Get recovery steps for specific error
sdp diagnose ENV001

# Get JSON output for tooling
sdp diagnose VAL001 --json
```

### sdp doctor

Use `sdp doctor` to check environment health:

```bash
# Standard check
sdp doctor

# Include drift detection
sdp doctor --drift
```

### sdp quality

Use `sdp quality` to check code quality:

```bash
# Run all quality checks
sdp quality

# Verbose output
sdp quality --verbose
```

## Recovery Drill Checklist

For maintainers, run this checklist monthly:

1. [ ] Verify `sdp doctor` passes in clean environment
2. [ ] Test recovery from corrupted session
3. [ ] Verify evidence chain repair works
4. [ ] Test drift detection and resolution
5. [ ] Verify quality gates work correctly
6. [ ] Test diagnose command for all error classes

## MTTR Target

**Mean Time To Recovery (MTTR) Target: < 5 minutes**

For classified errors with available playbooks, recovery should take less than 5 minutes from error to resolution.

### MTTR Measurement

Track these metrics:
- Time from error to diagnosis
- Time from diagnosis to resolution
- Number of errors requiring escalation

---

**Version:** 0.10.0
**Last Updated:** 2026-02-16
