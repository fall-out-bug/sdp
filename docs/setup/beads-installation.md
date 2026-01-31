# Beads CLI Installation

Beads is a git-backed issue tracker designed for AI agents. This guide shows how to install and verify the Beads CLI (`bd`).

## Prerequisites

- **Go 1.21+** (required for building Beads)
- **Git** (required for Beads data storage)

## Default Behavior (v0.6.0+)

**As of SDP v0.6.0, real Beads is the default:**

- **If `bd` CLI is installed** → Uses real Beads automatically
- **If `bd` not installed** → Falls back to mock with warning
- **Override with `BEADS_USE_MOCK=true`** for development/testing

This means you'll get real Beads integration out of the box after installation.

## Installation

### macOS

```bash
# 1. Install Go (if not already installed)
brew install go

# 2. Install Beads CLI
go install github.com/steveyegge/beads/cmd/bd@latest

# 3. Add Go bin directory to PATH
# Add this line to your ~/.zshrc or ~/.bash_profile
export PATH=$PATH:$(go env GOPATH)/bin

# 4. Reload shell configuration
source ~/.zshrc  # or source ~/.bash_profile

# 5. Verify installation
bd --version
```

### Linux

#### Ubuntu/Debian

```bash
# 1. Install Go
sudo apt update
sudo apt install golang-go

# 2. Install Beads CLI
go install github.com/steveyegge/beads/cmd/bd@latest

# 3. Add Go bin directory to PATH
# Add this line to your ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin

# 4. Reload shell configuration
source ~/.bashrc

# 5. Verify installation
bd --version
```

#### Fedora/RHEL

```bash
# 1. Install Go
sudo dnf install golang

# 2. Install Beads CLI
go install github.com/steveyegge/beads/cmd/bd@latest

# 3. Add Go bin directory to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# 4. Verify installation
bd --version
```

## Verify Installation with SDP

Once installed, verify Beads is detected by SDP:

```bash
# Run SDP health check
sdp doctor

# Expected output should include:
# ✅ Beads CLI: Beads CLI v0.1.0 at /Users/you/go/bin/bd
```

## Initialize Beads in Your Project

After installation, initialize Beads in your project directory:

```bash
# Navigate to your project
cd /path/to/your/project

# Initialize Beads (creates .beads/ directory)
bd init

# Verify initialization
ls -la .beads/
```

## Troubleshooting

### "bd: command not found"

**Cause:** The Go bin directory is not in your PATH.

**Solution:**

1. Check where Go installs binaries:
   ```bash
   go env GOPATH
   ```

2. Add the bin directory to your PATH:
   ```bash
   # Add to ~/.zshrc (macOS) or ~/.bashrc (Linux)
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

3. Reload your shell configuration:
   ```bash
   source ~/.zshrc  # or source ~/.bashrc
   ```

4. Verify:
   ```bash
   which bd
   bd --version
   ```

### "go: command not found"

**Cause:** Go is not installed.

**Solution:**

- **macOS:** `brew install go`
- **Ubuntu/Debian:** `sudo apt install golang-go`
- **Fedora/RHEL:** `sudo dnf install golang`
- **Manual:** Download from [golang.org](https://golang.org/dl/)

### Permission denied

**Cause:** The `bd` binary is not executable.

**Solution:**

```bash
chmod +x $(go env GOPATH)/bin/bd
```

### Beads fails to initialize

**Cause:** Not in a git repository.

**Solution:**

Beads requires a git repository. Initialize one if needed:

```bash
git init
bd init
```

## Environment Variables

### `BEADS_USE_MOCK`

Controls whether SDP uses real Beads or a mock client:

```bash
# Use real Beads (default if bd is installed)
export BEADS_USE_MOCK=false

# Use mock Beads (for testing without installation)
export BEADS_USE_MOCK=true
```

**Default behavior (as of SDP v0.6.0):**
- If `bd` CLI is installed → Uses real Beads
- If `bd` not installed → Falls back to mock with warning
- Override with `BEADS_USE_MOCK=true` for development

## CI/CD Configuration

For GitHub Actions or other CI systems, add Beads installation to your workflow:

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.21'

- name: Install Beads CLI
  run: |
    go install github.com/steveyegge/beads/cmd/bd@latest
    echo "$(go env GOPATH)/bin" >> $GITHUB_PATH

- name: Verify Beads
  run: bd --version
```

## Further Reading

- [Beads Documentation](https://github.com/steveyegge/beads)
- [Go Installation Guide](https://golang.org/doc/install)
- [SDP Protocol Documentation](../PROTOCOL.md)
- [SDP Doctor Command](../cli/doctor.md)

## Quick Reference

```bash
# Installation
go install github.com/steveyegge/beads/cmd/bd@latest

# Verify
bd --version
sdp doctor

# Initialize in project
bd init

# Common commands
bd list              # List all tasks
bd show <id>         # Show task details
bd create <title>    # Create new task
bd update <id>       # Update task
bd ready             # Show tasks ready to work on
```

## Support

For issues with:
- **Beads CLI:** [github.com/steveyegge/beads/issues](https://github.com/steveyegge/beads/issues)
- **SDP Integration:** [Your SDP repository issues](../../README.md)
