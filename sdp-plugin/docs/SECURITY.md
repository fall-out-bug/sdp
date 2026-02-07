# Subprocess Security Model

## Overview

This document describes the security model for subprocess execution in the SDP Go binary.

## Threat Model

### Attack Vectors

1. **Command Injection**: User input concatenated into command strings
2. **Argument Injection**: Malicious arguments to safe commands
3. **Path Traversal**: Accessing files outside project directory
4. **Resource Exhaustion**: Commands that hang indefinitely

## Mitigation Strategies

### 1. Command Whitelisting

Only explicitly allowed commands may be executed:

```go
// Whitelisted commands
- Test runners: pytest, go test, mvn test, gradle test, npm test
- Safe tools: git, claude, gh
```

### 2. Argument Validation

All arguments are checked for injection patterns:

```
Blocked patterns:
- ; (command separator)
- | (pipe)
- & (background)
- ` (backtick substitution)
- $( (dollar substitution)
- \n, \r (newlines)
- ../ (path traversal)
- Absolute paths to /etc/, /usr/, /bin/, /sbin/
```

### 3. Timeout Enforcement

All subprocess calls have timeouts:

```go
const (
    DefaultTimeout = 30 * time.Second   // Standard operations
    ShortTimeout   = 5 * time.Second    // Version checks
    LongTimeout    = 5 * time.Minute    // Full test suites
)
```

### 4. Context Propagation

All execution uses `context.Context` for cancellation:

```go
cmd := exec.CommandContext(ctx, command, args...)
```

## Usage

### Safe Command Creation

```go
import "github.com/ai-masters/sdp/internal/security"

ctx := context.Background()
cmd, err := security.SafeCommand(ctx, "pytest", []string{"tests/"}...)
if err != nil {
    return err
}
output, err := cmd.CombinedOutput()
```

### Custom Test Commands

When accepting custom test commands from users:

```go
testCmd := "pytest" // From user input or config
if err := security.ValidateTestCommand(testCmd); err != nil {
    return fmt.Errorf("invalid test command: %w", err)
}
```

### 6. File Permissions

Sensitive data files must have restrictive permissions:

```go
// Sensitive files (owner read/write only)
os.WriteFile(filename, data, 0600)  // Checkpoint files
os.WriteFile(filename, data, 0600)  // Telemetry files

// Executable files (owner execute)
os.WriteFile(filename, data, 0755)  // Git hooks
```

**Protected files:**
- `.beads/beads.db` - Issue tracker data
- `.oneshot` - Execution state
- `~/.sdp/telemetry.jsonl` - Usage telemetry
- Checkpoint files - Feature execution state

**Doctor check:**
```bash
sdp doctor
# Checks file permissions on sensitive data
# Warns if permissions > 0600
```

### 8. YAML Injection Protection

All YAML parsing uses secure decoder with limits:

```go
import "github.com/ai-masters/sdp/internal/parser"

// Safe YAML unmarshaling with limits
content, _ := os.ReadFile("workstream.md")
var fm frontmatter
if err := SafeYAMLUnmarshal(content, &fm); err != nil {
    return err
}
```

**Security limits:**
- **Max file size:** 1MB for workstream files
- **Max field length:** 10KB for YAML string fields
- **Max content length:** 1MB for markdown content
- **Max YAML depth:** 100 nesting levels

**Protected against:**
- YAML bombs (exponential expansion aliases)
- Recursive anchors (stack overflow)
- Oversized files (DoS via memory exhaustion)
- Excessively long field values

### 9. Path Traversal Protection

All user-provided file paths are validated:

```go
import "github.com/ai-masters/sdp/internal/security"

// Sanitize user input
cleanPath, err := security.SanitizePath(userPath)
if err != nil {
    return fmt.Errorf("invalid path: %w", err)
}

// Safe path joining
safePath, err := security.SafeJoinPath(baseDir, userPath)
if err != nil {
    return fmt.Errorf("path traversal attempt: %w", err)
}

// Validate path is within directory
err = security.ValidatePathInDirectory(baseDir, targetPath)
if err != nil {
    return fmt.Errorf("path outside base: %w", err)
}
```

**Blocked patterns:**
- `../` and `..\\` (parent directory traversal)
- Absolute paths (`/etc/passwd`, `C:\Windows\System32`)
- Paths that escape base directory

## Security Checklist

- ✅ All exec.Command calls use whitelisted commands
- ✅ All arguments validated for injection patterns
- ✅ All subprocess calls have context with timeout
- ✅ No shell execution (sh -c, cmd /c)
- ✅ Environment variables sanitized if passed
- ✅ All user-provided paths validated for traversal

## Examples

### Safe (✅)

```go
// Whitelisted command, safe arguments
exec.Command("pytest", "tests/", "-v")

// Git version check
exec.Command("git", "--version")

// Go test with context
exec.CommandContext(ctx, "go", "test", "./...")
```

### Unsafe (❌)

```go
// User-controlled command
testCmd := getUserInput()
exec.Command(testCmd, args...) // ❌ No validation

// Shell execution
exec.Command("sh", "-c", userString) // ❌ Arbitrary code

// No timeout
exec.Command("go", "test").Run() // ❌ Can hang

// Unsafe path handling
userPath := "../../../etc/passwd"
f, _ := os.Open(userPath) // ❌ Path traversal
```

### Safe Path Handling (✅)

```go
// Validate and sanitize user paths
userPath := "docs/workstreams/00-001-01.md"
safePath, err := security.SafeJoinPath(baseDir, userPath)
if err != nil {
    return err
}
f, err := os.Open(safePath) // ✅ Safe

// Block traversal attempts
userPath = "../../../etc/passwd"
_, err = security.SafeJoinPath(baseDir, userPath)
// Returns error: "path contains traversal pattern '../'"
```

## Testing

Security tests verify:

1. Whitelist enforcement
2. Injection pattern detection
3. Timeout application
4. Context propagation
5. Path traversal prevention
6. Directory containment validation

Run tests:
```bash
go test ./internal/security/... -v

# With coverage
go test ./internal/security/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## References

- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [Go exec package documentation](https://pkg.go.dev/os/exec)
