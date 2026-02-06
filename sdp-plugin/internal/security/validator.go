package security

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	// DefaultTimeout is the default timeout for subprocess execution
	DefaultTimeout = 30 * time.Second

	// ShortTimeout is for quick checks like version commands
	ShortTimeout = 5 * time.Second

	// LongTimeout is for long-running operations like full test suites
	LongTimeout = 5 * time.Minute
)

// CommandValidator validates command safety before execution
type CommandValidator struct {
	whitelist map[string]bool
}

// NewCommandValidator creates a validator with default whitelist
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		whitelist: map[string]bool{
			// Test runners
			"pytest":    true,
			"pytest-3":  true,
			"python":    true, // When used with -m pytest
			"go":        true, // When used with test, vet, build
			"mvn":       true,
			"mvnw":      true,
			"gradle":    true,
			"gradlew":   true,
			"./gradlew": true,
			"npm":       true,
			"yarn":      true,
			"pnpm":      true,
			"jest":      true,
			"mocha":     true,
			"jasmine":   true,
			"cargo":     true, // Rust
			"dart":      true,
			"flutter":   true,

			// Safe system tools
			"git":    true,
			"claude": true,
			"gh":     true, // GitHub CLI
		},
	}
}

// ValidateCommand checks if a command is allowed to execute
func (v *CommandValidator) ValidateCommand(command string) error {
	if !v.whitelist[command] {
		return fmt.Errorf("command '%s' is not whitelisted for execution", command)
	}
	return nil
}

// ValidateArgs checks arguments for injection patterns
func (v *CommandValidator) ValidateArgs(args []string) error {
	for _, arg := range args {
		if err := v.validateArg(arg); err != nil {
			return err
		}
	}
	return nil
}

// validateArg checks a single argument for injection patterns
func (v *CommandValidator) validateArg(arg string) error {
	injectionPatterns := []string{
		";",      // Command separator
		"|",      // Pipe
		"&",      // Background command
		"`",      // Command substitution (backtick)
		"$(",     // Command substitution (dollar)
		"\n",     // Newline
		"\r",     // Carriage return
		"\\",     // Escape character
		"../",    // Path traversal (partial)
		"../../", // Path traversal (double)
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(arg, pattern) {
			return fmt.Errorf("argument '%s' contains injection pattern '%s'", arg, pattern)
		}
	}

	// Check for command substitution
	if strings.Contains(arg, "$(") || strings.Contains(arg, "`") {
		return fmt.Errorf("argument '%s' contains command substitution", arg)
	}

	// Check for absolute paths to sensitive system files
	if strings.HasPrefix(arg, "/etc/") || strings.HasPrefix(arg, "/usr/") || strings.HasPrefix(arg, "/bin/") || strings.HasPrefix(arg, "/sbin/") {
		return fmt.Errorf("argument '%s' contains absolute path to system directory", arg)
	}

	return nil
}

// SafeCommand creates a safe exec.Command with timeout and validation
func SafeCommand(ctx context.Context, command string, args ...string) (*exec.Cmd, error) {
	validator := NewCommandValidator()

	// Validate command
	if err := validator.ValidateCommand(command); err != nil {
		return nil, err
	}

	// Validate arguments
	if err := validator.ValidateArgs(args); err != nil {
		return nil, err
	}

	// Set timeout if context doesn't have one
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}

	// Create command with context
	cmd := exec.CommandContext(ctx, command, args...)

	return cmd, nil
}

// MustSafeCommand creates a safe command or panics
// Use this for commands that MUST be safe (e.g., hardcoded system checks)
func MustSafeCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd, err := SafeCommand(ctx, command, args...)
	if err != nil {
		panic(fmt.Sprintf("safe command: %v", err))
	}
	return cmd
}

// ValidateTestCommand validates a custom test command string
func ValidateTestCommand(testCmd string) error {
	// Split into command and args
	parts := strings.Fields(testCmd)
	if len(parts) == 0 {
		return fmt.Errorf("test command is empty")
	}

	validator := NewCommandValidator()

	// Validate base command
	if err := validator.ValidateCommand(parts[0]); err != nil {
		return err
	}

	// Validate arguments
	if len(parts) > 1 {
		if err := validator.ValidateArgs(parts[1:]); err != nil {
			return err
		}
	}

	return nil
}
