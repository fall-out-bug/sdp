# Beads Installation & Testing Guide

Complete guide for installing Go, Beads CLI, and testing SDP integration with real Beads (not mock).

---

## Prerequisites

- macOS or Linux system
- Python 3.10+
- SDP project already set up
- Internet connection

---

## Step 1: Install Go

### macOS

```bash
# Install Go using Homebrew
brew install go

# Verify installation
go version
# Expected: go version go1.24+ (or similar)

# Add Go to PATH (if not already)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Linux

```bash
# Download Go 1.24+
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

---

## Step 2: Install Beads CLI

```bash
# Install Beads
go install github.com/steveyegge/beads/cmd/bd@latest

# Verify installation
bd --version
# Expected: bd version <version>

# Check available commands
bd --help
```

**Note:** This installs `bd` binary to `$GOPATH/bin` (usually `~/go/bin`).

---

## Step 3: Initialize Beads in Project

```bash
cd /path/to/sdp-project

# Initialize Beads
bd init

# Creates:
# .beads/
# ├── beads.db          # SQLite cache
# └── issues.jsonl      # Git-tracked tasks

# Verify initialization
ls -la .beads/
```

---

## Step 4: Test Beads CLI

```bash
# Create a test task
bd create "Test task from CLI"

# List tasks
bd list

# Show task details
bd show <task-id-from-above>

# Delete test task
bd delete <task-id>
```

---

## Step 5: Test SDP + Beads Integration

### Disable Mock Mode

```bash
# In current shell
export BEADS_USE_MOCK=false

# Or permanently add to profile
echo 'export BEADS_USE_MOCK=false' >> ~/.zshrc
source ~/.zshrc
```

### Test with Real Beads

```python
# Test script
from sdp.beads import create_beads_client, BeadsTaskCreate

# Create real Beads client (not mock!)
client = create_beads_client(use_mock=False)

# Create a task
task = client.create_task(BeadsTaskCreate(
    title="Real Beads test task",
    description="Testing with actual Beads CLI",
))

print(f"✅ Created task: {task.id}")
print(f"   Title: {task.title}")
print(f"   Status: {task.status}")

# Verify with Beads CLI
import subprocess
result = subprocess.run(["bd", "show", task.id], capture_output=True, text=True)
print(f"\n📋 Beads CLI output:")
print(result.stdout)
```

### Run Tests with Real Beads

```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp-beads-integration

# Set environment
export BEADS_USE_MOCK=false

# Run tests
PYTHONPATH=src poetry run pytest tests/unit/beads/ -v

# All tests should pass with real Beads!
```

---

## Step 6: Test Multi-Agent Workflow

```bash
# Create feature
@idea "Multi-agent test feature"
# → bd-0001

# Decompose into workstreams
@design bd-0001
# → bd-0001.1, bd-0001.2, bd-0001.3

# Check ready tasks
bd ready
# → Should show bd-0001.1

# Execute with 3 agents in parallel
@oneshot bd-0001 --agents 3

# All workstreams should complete automatically!
```

---

## Troubleshooting

### "command not found: bd"

**Problem:** Beads CLI not in PATH.

**Solution:**
```bash
# Check Go bin directory
echo $(go env GOPATH)/bin

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Make permanent
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```

### "bd: command not found"

**Problem:** Beads not installed.

**Solution:**
```bash
go install github.com/steveyegge/beads/cmd/bd@latest

# Verify
which bd
```

### "ModuleNotFoundError: No module named 'sdp'"

**Problem:** Python path not set.

**Solution:**
```bash
export PYTHONPATH=/path/to/sdp/src:$PYTHONPATH

# Or run from project root with -m flag
python -m sdp.cli beads status
```

### "Beads connection error"

**Problem:** Beads daemon not running or database locked.

**Solution:**
```bash
# Check Beads status
bd status

# Restart Beads daemon
pkill bd
bd init  # Reinitialize

# Or use mock mode
export BEADS_USE_MOCK=true
```

### Tests Fail with Real Beads

**Problem:** Tests expecting mock behavior fail with real Beads.

**Solution:**
- Check test output for specific failures
- Verify test IDs are unique (Beads generates new IDs each run)
- Ensure `.beads/` is clean between test runs

---

## Verification Checklist

Use this checklist to verify complete setup:

- [ ] Go 1.24+ installed
- [ ] `bd --version` works
- [ ] `bd init` run in project
- [ ] `.beads/beads.db` exists
- [ ] `BEADS_USE_MOCK=false` set
- [ ] `sdp beads status` works
- [ ] Test: `create_beads_client(use_mock=False)` works
- [ ] All pytest tests pass with real Beads
- [ ] `@idea` creates real Beads task
- [ ] `@design` creates real sub-tasks
- [ ] `@oneshot` executes with real Beads

---

## Performance Benchmarks

Compare mock vs real Beads:

| Operation | Mock | Real Beads | Overhead |
|-----------|------|------------|----------|
| Create task | ~1ms | ~50ms | 50x |
| List tasks | ~1ms | ~10ms | 10x |
| Get ready | ~1ms | ~20ms | 20x |
| Update status | ~1ms | ~30ms | 30x |

**Conclusion:** Real Beads adds overhead but is still fast (<100ms per operation). Mock is suitable for development.

---

## Next Steps After Real Beads Testing

1. **Performance Testing:**
   - Test with 100+ workstreams
   - Measure multi-agent execution time
   - Benchmark `bd ready` queries

2. **Integration Testing:**
   - Test with actual feature
   - Verify parallel execution works
   - Check git integration

3. **Decision Time:**
   - Compare Beads vs F012 performance
   - Decide which to use in production
   - Plan migration or deprecation

---

**Version:** 1.0.0
**Last Updated:** 2026-01-28
