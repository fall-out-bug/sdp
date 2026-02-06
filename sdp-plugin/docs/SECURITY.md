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

## Security Checklist

- ✅ All exec.Command calls use whitelisted commands
- ✅ All arguments validated for injection patterns
- ✅ All subprocess calls have context with timeout
- ✅ No shell execution (sh -c, cmd /c)
- ✅ Environment variables sanitized if passed

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
```

## Testing

Security tests verify:

1. Whitelist enforcement
2. Injection pattern detection
3. Timeout application
4. Context propagation

Run tests:
```bash
go test ./internal/security/... -v
```

## References

- [OWASP Command Injection](https://owasp.org/www-community/attacks/Command_Injection)
- [Go exec package documentation](https://pkg.go.dev/os/exec)
